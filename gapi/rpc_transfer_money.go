package gapi

import (
	"context"
	"database/sql"

	db "github.com/JaidenShall/simplebank/db/sqlc"
	"github.com/JaidenShall/simplebank/pb"
	"github.com/JaidenShall/simplebank/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) TransferMoney(ctx context.Context, req *pb.TransferMoneyRequest) (*pb.TransferMoneyResponse, error) {
	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	violations := validateTransferMoneyRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	// Get and validate from account
	fromAccount, err := server.store.GetAccount(ctx, req.GetFromAccountId())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "from account not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get from account: %s", err)
	}

	// Check if the authenticated user owns the from account
	if fromAccount.Owner != authPayload.Username {
		return nil, status.Errorf(codes.PermissionDenied, "from account doesn't belong to the authenticated user")
	}

	// Validate from account currency
	if fromAccount.Currency != req.GetCurrency() {
		return nil, status.Errorf(codes.InvalidArgument, "from account currency mismatch: %s vs %s", fromAccount.Currency, req.GetCurrency())
	}

	// Get and validate to account
	toAccount, err := server.store.GetAccount(ctx, req.GetToAccountId())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "to account not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get to account: %s", err)
	}

	// Validate to account currency
	if toAccount.Currency != req.GetCurrency() {
		return nil, status.Errorf(codes.InvalidArgument, "to account currency mismatch: %s vs %s", toAccount.Currency, req.GetCurrency())
	}

	// Perform the transfer transaction
	arg := db.TransferTxParams{
		FromAccountID: req.GetFromAccountId(),
		ToAccountID:   req.GetToAccountId(),
		Amount:        req.GetAmount(),
	}

	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "transfer transaction failed: %s", err)
	}

	rsp := &pb.TransferMoneyResponse{
		Transfer:    convertTransfer(result.Transfer),
		FromAccount: convertAccount(result.FromAccount),
		ToAccount:   convertAccount(result.ToAccount),
		FromEntry:   convertEntry(result.FromEntry),
		ToEntry:     convertEntry(result.ToEntry),
	}

	return rsp, nil
}

func validateTransferMoneyRequest(req *pb.TransferMoneyRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateAccountID(req.GetFromAccountId()); err != nil {
		violations = append(violations, fieldViolation("from_account_id", err))
	}

	if err := val.ValidateAccountID(req.GetToAccountId()); err != nil {
		violations = append(violations, fieldViolation("to_account_id", err))
	}

	if err := val.ValidateAmount(req.GetAmount()); err != nil {
		violations = append(violations, fieldViolation("amount", err))
	}

	if err := val.ValidateCurrency(req.GetCurrency()); err != nil {
		violations = append(violations, fieldViolation("currency", err))
	}

	return violations
}
