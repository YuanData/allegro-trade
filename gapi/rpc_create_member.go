package gapi

import (
	"context"

	"github.com/lib/pq"
	db "github.com/YuanData/allegro-trade/db/sqlc"
	"github.com/YuanData/allegro-trade/pb"
	"github.com/YuanData/allegro-trade/util"
	"github.com/YuanData/allegro-trade/vld"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateMember(ctx context.Context, req *pb.CreateMemberRequest) (*pb.CreateMemberResponse, error) {
	violations := validateCreateMemberRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "gen hash err: %s", err)
	}

	arg := db.CreateMemberParams{
		Membername:       req.GetMembername(),
		PasswordHash: hashedPassword,
		NameEntire:       req.GetNameEntire(),
		Email:          req.GetEmail(),
	}

	member, err := server.store.CreateMember(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "membername AlreadyExists err: %s", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "create member err: %s", err)
	}

	rsp := &pb.CreateMemberResponse{
		Member: convertMember(member),
	}
	return rsp, nil
}

func validateCreateMemberRequest(req *pb.CreateMemberRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := vld.ValidateMembername(req.GetMembername()); err != nil {
		violations = append(violations, fieldViolation("membername", err))
	}

	if err := vld.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	if err := vld.ValidateNameEntire(req.GetNameEntire()); err != nil {
		violations = append(violations, fieldViolation("name_entire", err))
	}

	if err := vld.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	return violations
}
