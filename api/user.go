package api

import (
	"net/http"

	db "github.com/ashishkhuraishy/go_forum/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createUserReq struct {
	Username string `json:"username" binding:"required"`
}

func (s *Server) createUser(ctx *gin.Context) {
	var req createUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.CreateUserParams{}

	user, err := s.store.CreateUser(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

type getUserReq struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (s *Server) getUser(ctx *gin.Context) {
	var req getUserReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := s.store.GetUser(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

type listUsersReq struct {
	Limit  int32 `form:"limit,default=10" binding:"min=5,max=50"`
	Offset int32 `form:"page_no,default=1" binding:"min=1"`
}

func (s *Server) listUsers(ctx *gin.Context) {
	var req listUsersReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.ListUsersParams{
		Limit:  req.Limit,
		Offset: (req.Offset - 1) * req.Limit,
	}

	users, err := s.store.ListUsers(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, users)
}
