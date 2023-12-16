package gapi

import (
	"context"
	"time"

	"github.com/hibiken/asynq"
	db "github.com/YuanData/allegro-trade/db/sqlc"
	"github.com/YuanData/allegro-trade/pb"
	"github.com/YuanData/allegro-trade/util"
	"github.com/YuanData/allegro-trade/vld"
	"github.com/YuanData/allegro-trade/worker"
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

	arg := db.CreateMemberTxParams{
		CreateMemberParams: db.CreateMemberParams{
			Membername:       req.GetMembername(),
			PasswordHash: hashedPassword,
			NameEntire:       req.GetNameEntire(),
			Email:          req.GetEmail(),
		},
		AfterCreate: func(member db.Member) error {
			taskPayload := &worker.PayloadSendVerifyEmail{
				Membername: member.Membername,
			}
			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(10 * time.Second),
				asynq.Queue(worker.QueueCritical),
			}

			return server.taskDistributor.DistributeTaskSendVerifyEmail(ctx, taskPayload, opts...)
		},
	}

	txResult, err := server.store.CreateMemberTx(ctx, arg)
	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			return nil, status.Errorf(codes.AlreadyExists, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "create member err: %s", err)
	}

	rsp := &pb.CreateMemberResponse{
		Member: convertMember(txResult.Member),
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
