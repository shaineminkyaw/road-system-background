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

	//type1 controller
	betController := NewBetController(h)
	betController.Register()

	//type2 controller
	type2Controller := NewType2Controller(h)
	type2Controller.Register()

	//type3 controller
	type3Controller := NewType3Controller(h)
	type3Controller.Register()

	//type4 controller
	type4Controller := NewType4Controller(h)
	type4Controller.Register()

	return h
}
