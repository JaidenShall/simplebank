package gapi

import (
	"context"
	"database/sql"

	"github.com/JaidenShall/simplebank/pb"
	"github.com/JaidenShall/simplebank/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) GetAccount(ctx context.Context, req *pb.GetAccountRequest) (*pb.GetAccountResponse, error) {
	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	violations := validateGetAccountRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	account, err := server.store.GetAccount(ctx, req.GetId())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "account not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get account: %s", err)
	}

	if account.Owner != authPayload.Username {
		return nil, status.Errorf(codes.PermissionDenied, "account doesn't belong to the authenticated user")
	}

	rsp := &pb.GetAccountResponse{
		Account: convertAccount(account),
	}
	return rsp, nil
}

func validateGetAccountRequest(req *pb.GetAccountRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateAccountID(req.GetId()); err != nil {
		violations = append(violations, fieldViolation("id", err))
	}

	return violations
}
