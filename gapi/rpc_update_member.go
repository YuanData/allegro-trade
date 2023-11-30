package gapi

import (
	"context"
	"database/sql"
	"time"

	db "github.com/YuanData/allegro-trade/db/sqlc"
	"github.com/YuanData/allegro-trade/pb"
	"github.com/YuanData/allegro-trade/util"
	"github.com/YuanData/allegro-trade/vld"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) UpdateMember(ctx context.Context, req *pb.UpdateMemberRequest) (*pb.UpdateMemberResponse, error) {
	authPayload, err := server.authorizeMember(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	violations := validateUpdateMemberRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	if authPayload.Membername != req.GetMembername() {
		return nil, status.Errorf(codes.PermissionDenied, "update member failed")
	}

	arg := db.UpdateMemberParams{
		Membername: req.GetMembername(),
		NameEntire: sql.NullString{
			String: req.GetNameEntire(),
			Valid:  req.NameEntire != nil,
		},
		Email: sql.NullString{
			String: req.GetEmail(),
			Valid:  req.Email != nil,
		},
	}

	if req.Password != nil {
		hashedPassword, err := util.HashPassword(req.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "gen hash err: %s", err)
		}

		arg.PasswordHash = sql.NullString{
			String: hashedPassword,
			Valid:  true,
		}

		arg.PasswordChangedTime = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
	}

	member, err := server.store.UpdateMember(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "member NotFound err")
		}
		return nil, status.Errorf(codes.Internal, "failed to update member: %s", err)
	}

	rsp := &pb.UpdateMemberResponse{
		Member: convertMember(member),
	}
	return rsp, nil
}

func validateUpdateMemberRequest(req *pb.UpdateMemberRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := vld.ValidateMembername(req.GetMembername()); err != nil {
		violations = append(violations, fieldViolation("membername", err))
	}

	if req.Password != nil {
		if err := vld.ValidatePassword(req.GetPassword()); err != nil {
			violations = append(violations, fieldViolation("password", err))
		}
	}

	if req.NameEntire != nil {
		if err := vld.ValidateNameEntire(req.GetNameEntire()); err != nil {
			violations = append(violations, fieldViolation("name_entire", err))
		}
	}

	if req.Email != nil {
		if err := vld.ValidateEmail(req.GetEmail()); err != nil {
			violations = append(violations, fieldViolation("email", err))
		}
	}

	return violations
}
