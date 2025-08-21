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
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) Deposit(ctx context.Context, req *pb.DepositRequest) (*pb.DepositResponse, error) {
	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	// Validate input
	violations := validateDepositRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	// Check if account exists and belongs to the authenticated user
	account, err := server.store.GetAccount(ctx, req.GetAccountId())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "account not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to find account")
	}

	if account.Owner != authPayload.Username {
		return nil, status.Errorf(codes.PermissionDenied, "account doesn't belong to the authenticated user")
	}

	// Perform deposit transaction
	arg := db.DepositTxParams{
		AccountID: req.GetAccountId(),
		Amount:    req.GetAmount(),
	}

	result, err := server.store.DepositTx(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "deposit transaction failed: %s", err)
	}

	rsp := &pb.DepositResponse{
		Id:        result.Account.ID,
		Owner:     result.Account.Owner,
		Balance:   result.Account.Balance,
		Currency:  result.Account.Currency,
		CreatedAt: timestamppb.New(result.Account.CreatedAt),
		UpdatedAt: timestamppb.New(result.Account.CreatedAt), // You might want to add updated_at to your schema
	}

	return rsp, nil
}

func validateDepositRequest(req *pb.DepositRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateAccountID(req.GetAccountId()); err != nil {
		violations = append(violations, fieldViolation("account_id", err))
	}

	if err := val.ValidateAmount(req.GetAmount()); err != nil {
		violations = append(violations, fieldViolation("amount", err))
	}

	return violations
}
