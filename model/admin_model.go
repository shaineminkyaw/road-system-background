package model

import (
	"time"

	"github.com/casbin/casbin/v2"
	"gorm.io/gorm"
)

//for each admin data
type Admin struct {
	ID          uint64         `gorm:"column:id" json:"id"`
	UserName    string         `gorm:"column:username" json:"username"`
	LoginName   string         `gorm:"column:login_name" json:"login_name"`
	Email       string         `gorm:"column:email"  json:"email"`
	Password    string         `gorm:"column:password" json:"password"`
	Avatar      string         `gorm:"column:avatar" json:"avatar"`
	IsOnline    bool           `gorm:"column:isOnline" json:"isOnline"`
	Gender      int8           `gorm:"column:gender" json:"gender"`
	LoginIP     string         `gorm:"column:login_ip" json:"login_ip"`
	LastLoginIP string         `gorm:"column:last_loginIP" json:"last_loginIP"`
	CreadtedBy  string         `gorm:"created_by" json:"created_by"`
	UpdatedBy   string         `gorm:"updated_by" json:"updated_by"`
	Remark      string         `gorm:"remark" json:"remark"`
	CreatedAt   time.Time      `gorm:"created_at" json:"created_at"`
	DeletedAt   gorm.DeletedAt `gorm:"deleted_at" json:"deleted_at"`
}

type AdminRole struct {
	ID     uint64 `gorm:"column:id" json:"id"`
	UserID uint64 `gorm:"column:uid" json:"uid"`
	RoleID uint64 `gorm:"column:rid" json:"rid"`
}

//for admin or user role define table
type Roles struct {
	ID        uint64    `gorm:"column:id" json:"id"`
	Slug      string    `gorm:"column:slug" json:"slug"`
	Name      string    `gorm:"column:name" json:"name"`
	Status    int8      `gorm:"column:status" json:"status"`
	CreatedBy string    `gorm:"column:created_by" json:"created_by"`
	UpdatedBy string    `gorm:"column:updated_by" json:"updated_by"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type CasbinPloicy struct {
	ID         uint64    `gorm:"column:id" json:"id"`
	Username   string    `gorm:"column:username" json:"username"`
	Slug       string    `gorm:"column:slug" json:"slug"`             //role
	Permission string    `gorm:"column:permission" json:"permission"` //route
	Action     string    `gorm:"column:action" json:"action"`         //http method
	CreatedBy  string    `gorm:"created_by" json:"created_by"`
	UpdatedBy  string    `gorm:"updated_by" json:"updated_by"`
	Remark     string    `gorm:"column:remark" json:"remark"`
	CreatedAt  time.Time `gorm:"created_at" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (cp *CasbinPloicy) TableName() string {
	return "Casbin_Policy"
}

func (cp *CasbinPloicy) AssignPermisssion(e *casbin.Enforcer, permission *CasbinPloicy, role string) (bool, error) {
	//

	if ok := e.HasPolicy(role, permission.Permission, permission.Action); !ok {
		ok, err := e.AddPolicy(role, permission.Permission, permission.Action)
		if err != nil {
			return false, err
		} else if !ok {
			return false, err
		}
		bol, err := e.AddRoleForUser("admin", role)
		if !bol || err != nil {
			return false, err
		}
	}
	return true, nil

}
