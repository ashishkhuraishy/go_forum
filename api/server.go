package api

import (
	db "github.com/ashishkhuraishy/go_forum/db/sqlc"
	"github.com/ashishkhuraishy/go_forum/utils"
	"github.com/ashishkhuraishy/go_forum/utils/token"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store      *db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(store *db.Store) (*Server, error) {
	server := &Server{store: store}
	router := gin.Default()

	tokenMaker, err := token.NewPasetoMaker(utils.RandomString(32))
	if err != nil {
		return nil, err
	}

	server.tokenMaker = tokenMaker
	server.router = router

	createRoutes(router, server)

	return server, nil
}

func (s *Server) Start(addess string) error {
	return s.router.Run(addess)
}

func errorResponse(err error) *gin.H {
	return &gin.H{
		"error": err.Error(),
	}
}
