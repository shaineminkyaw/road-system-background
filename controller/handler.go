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

func Inject(r *gin.Engine, e *casbin.Enforcer, s *gorm.DB) *Handler {

	h := &Handler{
		R:        r,
		Enforcer: e,
		Source:   s,
	}
	//admin controller
	adminController := NewAdminController(h)
	adminController.Register()

	//bet controller
	betController := NewBetController(h)
	betController.Register()
	return h
}
