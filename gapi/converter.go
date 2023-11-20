package gapi

import (
	db "github.com/YuanData/allegro-trade/db/sqlc"
	"github.com/YuanData/allegro-trade/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertMember(member db.Member) *pb.Member {
	return &pb.Member{
		Membername:          member.Membername,
		NameEntire:          member.NameEntire,
		Email:             member.Email,
		PasswordChangedTime: timestamppb.New(member.PasswordChangedTime),
		CreatedTime:         timestamppb.New(member.CreatedTime),
	}
}
