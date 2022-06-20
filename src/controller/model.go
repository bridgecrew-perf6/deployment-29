package controller

import "github.com/gin-gonic/gin"

type Ping interface {
	ExecPing(c *gin.Context)
	UpdateTP(c *gin.Context)
	ShowLog(*gin.Context)
}

type Wireguard interface {
	ShowConfig(c *gin.Context)
	ShowStatus(c *gin.Context)
	UpdateWgConfig(c *gin.Context)
}

type ICategoryController interface {
	Wireguard
	Ping
	Update(c *gin.Context)
}

type CategoryController struct {
}
