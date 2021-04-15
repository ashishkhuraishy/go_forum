package api

import "github.com/gin-gonic/gin"

func createRoutes(router *gin.Engine, server *Server) {

	// TODO: add auth middleware
	// authMiddleWare := authMiddleWare(server.tokenMaker)

	// LoginApi
	// SignupApi
	router.POST("/user", server.createUser)

	// User Apis
	router.GET("/user", server.listUsers)
	router.GET("/user/:id", server.getUser)
	// router.PUT("/user/:id", server.updateUser)
	// router.DELETE("/user/:id", server.deleteUser)

	// Feed apis

	// Like apis
}
