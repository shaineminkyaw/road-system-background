package ds

import (
	"log"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

func NewRBAC(db *gorm.DB) (*casbin.Enforcer, error) {
	//

	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}

	e, err := casbin.NewEnforcer("./config/rbac_model.conf", adapter)
	if err != nil {
		return nil, err
	}
	e.EnableLog(true)
	err = e.LoadPolicy()
	if err != nil {
		log.Println(err.Error())
	}

	return e, nil
}

func AddPolicy(e *casbin.Enforcer, sub, obj, act string) (bool, string) {
	ans := "policy is exist..."
	ans2 := "policy is not exist , adding policy ...."
	e.HasPolicy()
	bol, _ := e.AddPolicy(sub, obj, act)
	if !bol {
		return false, ans
	}
	return bol, ans2
}

func AddRoleForUser(e *casbin.Enforcer, user, role string) (bool, string) {
	//
	ans := "user role exist ..."
	ans2 := "user role is not exist , adding user role ..."
	bol, _ := e.AddRoleForUser(user, role)
	if !bol {
		return false, ans
	}
	return bol, ans2
}
