package gapi

import (
	"context"
	// "fmt"

	db "github.com/JaidenShall/simplebank/db/sqlc"
	"github.com/JaidenShall/simplebank/pb"

	// "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) ListAccounts(ctx context.Context, req *pb.ListAccountsRequest) (*pb.ListAccountsResponse, error) {
	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	// violations := validateListAccountsRequest(req)
	// if violations != nil {
	// 	return nil, invalidArgumentError(violations)
	// }

	arg := db.ListAccountsParams{
		Owner:  authPayload.Username,
		Limit:  3,
		Offset: 0,
	}

	accounts, err := server.store.ListAccounts(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list accounts: %s", err)
	}

	rsp := &pb.ListAccountsResponse{
		Accounts: convertAccounts(accounts),
	}
	return rsp, nil
}

// func validateListAccountsRequest(req *pb.ListAccountsRequest) (violations []*errdetails.BadRequest_FieldViolation) {
// 	if req.GetPageId() < 1 {
// 		violations = append(violations, fieldViolation("page_id", fmt.Errorf("must be greater than 0")))
// 	}

// 	if req.GetPageSize() < 5 {
// 		violations = append(violations, fieldViolation("page_size", fmt.Errorf("must be at least 5")))
// 	}

// 	if req.GetPageSize() > 10 {
// 		violations = append(violations, fieldViolation("page_size", fmt.Errorf("must not exceed 10")))
// 	}

// 	return violations
// }

func convertAccounts(accounts []db.Account) []*pb.Account {
	var result []*pb.Account
	for _, account := range accounts {
		result = append(result, convertAccount(account))
	}
	return result
}
