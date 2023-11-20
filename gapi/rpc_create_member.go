package gapi

import (
	"context"

	"github.com/lib/pq"
	db "github.com/YuanData/allegro-trade/db/sqlc"
	"github.com/YuanData/allegro-trade/pb"
	"github.com/YuanData/allegro-trade/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateMember(ctx context.Context, req *pb.CreateMemberRequest) (*pb.CreateMemberResponse, error) {
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
