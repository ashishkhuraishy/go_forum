package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type createUserReq struct {
	Username string `json:"username" binding:"required"`
}

func (s *Server) CreateUser(ctx *gin.Context) {
	var req createUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := s.store.CreateUser(ctx, req.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}
