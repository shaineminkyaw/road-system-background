package controller

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	R        *gin.Engine
	Enforcer *casbin.Enforcer
	Source   *gorm.DB
}
