package gapi

import (
	"context"
	"errors"

	db "github.com/YuanData/allegro-trade/db/sqlc"
	"github.com/YuanData/allegro-trade/pb"
	"github.com/YuanData/allegro-trade/util"
	"github.com/YuanData/allegro-trade/vld"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginMember(ctx context.Context, req *pb.LoginMemberRequest) (*pb.LoginMemberResponse, error) {
	violations := validateLoginMemberRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	member, err := server.store.GetMember(ctx, req.GetMembername())
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "member NotFound err")
		}
		return nil, status.Errorf(codes.Internal, "member Internal err")
	}

	err = util.VerifyPassword(req.Password, member.PasswordHash)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "wrong password")
	}

	accessToken, accessPayload, err := server.tokenAuthzr.CreateToken(
		member.Membername,
		member.Role,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create access token err")
	}

	refreshToken, refreshPayload, err := server.tokenAuthzr.CreateToken(
		member.Membername,
		member.Role,
		server.config.RefreshTokenDuration,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create refresh token err")
	}

	mtdata := server.extractMetadata(ctx)
	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Membername:     member.Membername,
		RefreshToken: refreshToken,
		UserAgent:    mtdata.UserAgent,
		ClientIp:     mtdata.ClientIP,
		IsBlocked:    false,
		ExpiredTime:    refreshPayload.ExpiredTime,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create session err")
	}

	rsp := &pb.LoginMemberResponse{
		Member:                  convertMember(member),
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiredTime:  timestamppb.New(accessPayload.ExpiredTime),
		RefreshTokenExpiredTime: timestamppb.New(refreshPayload.ExpiredTime),
	}
	return rsp, nil
}

func validateLoginMemberRequest(req *pb.LoginMemberRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := vld.ValidateMembername(req.GetMembername()); err != nil {
		violations = append(violations, fieldViolation("membername", err))
	}

	if err := vld.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	return violations
}
