package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/YuanData/allegro-trade/db/sqlc"
	"github.com/YuanData/allegro-trade/util"
)

type createMemberRequest struct {
	Membername string `json:"membername" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	NameEntire string `json:"name_entire" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type createMemberResponse struct {
	Membername          string    `json:"membername"`
	NameEntire          string    `json:"name_entire"`
	Email             string    `json:"email"`
	PasswordChangedTime time.Time `json:"password_changed_time"`
	CreatedTime         time.Time `json:"created_time"`
}

func (server *Server) createMember(ctx *gin.Context) {
	var req createMemberRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateMemberParams{
		Membername:       req.Membername,
		PasswordHash: hashedPassword,
		NameEntire:       req.NameEntire,
		Email:          req.Email,
	}

	member, err := server.store.CreateMember(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := createMemberResponse{
		Membername:          member.Membername,
		NameEntire:          member.NameEntire,
		Email:             member.Email,
		PasswordChangedTime: member.PasswordChangedTime,
		CreatedTime:         member.CreatedTime,
	}
	ctx.JSON(http.StatusOK, rsp)
}
