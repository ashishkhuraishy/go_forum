package api

import (
	db "github.com/ashishkhuraishy/go_forum/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/user", server.CreateUser)

	server.router = router
	return server
}

func (s *Server) Start(addess string) error {
	return s.router.Run(addess)
}

func errorResponse(err error) *gin.H {
	return &gin.H{
		"error": err.Error(),
	}
}
