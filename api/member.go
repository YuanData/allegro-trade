package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/YuanData/allegro-trade/db/sqlc"
	"github.com/YuanData/allegro-trade/util"
)

type createMemberRequest struct {
	Membername string `json:"membername" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	NameEntire string `json:"name_entire" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type memberResponse struct {
	Membername          string    `json:"membername"`
	NameEntire          string    `json:"name_entire"`
	Email             string    `json:"email"`
	PasswordChangedTime time.Time `json:"password_changed_time"`
	CreatedTime         time.Time `json:"created_time"`
}

func newMemberResponse(member db.Member) memberResponse {
	return memberResponse{
		Membername:          member.Membername,
		NameEntire:          member.NameEntire,
		Email:             member.Email,
		PasswordChangedTime: member.PasswordChangedTime,
		CreatedTime:         member.CreatedTime,
	}
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
		if db.ErrorCode(err) == db.UniqueViolation {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := newMemberResponse(member)
	ctx.JSON(http.StatusOK, rsp)
}

type loginMemberRequest struct {
	Membername string `json:"membername" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginMemberResponse struct {
	SessionID             uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiredTime  time.Time    `json:"access_token_expired_time"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiredTime time.Time    `json:"refresh_token_expired_time"`
	Member                  memberResponse `json:"member"`
}

func (server *Server) loginMember(ctx *gin.Context) {
	var req loginMemberRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	member, err := server.store.GetMember(ctx, req.Membername)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.VerifyPassword(req.Password, member.PasswordHash)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, accessPayload, err := server.tokenAuthzr.CreateToken(
		member.Membername,
		member.Role,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshPayload, err := server.tokenAuthzr.CreateToken(
		member.Membername,
		member.Role,
		server.config.RefreshTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Membername:     member.Membername,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiredTime:    refreshPayload.ExpiredTime,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := loginMemberResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiredTime:  accessPayload.ExpiredTime,
		RefreshToken:          refreshToken,
		RefreshTokenExpiredTime: refreshPayload.ExpiredTime,
		Member:                  newMemberResponse(member),
	}
	ctx.JSON(http.StatusOK, rsp)
}
