package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shaineminkyaw/road-system-background/ds"
	"github.com/shaineminkyaw/road-system-background/dto"
	"github.com/shaineminkyaw/road-system-background/factory"
	"github.com/shaineminkyaw/road-system-background/middleware"
	"github.com/shaineminkyaw/road-system-background/model"
	"github.com/shaineminkyaw/road-system-background/utils"

	"gorm.io/gorm"
)

type adminController struct {
	H *Handler
}

func NewAdminController(h *Handler) *adminController {
	return &adminController{
		H: h,
	}
}

func (ctr *adminController) Register() {
	//
	h := ctr.H
	h.R.POST("refresh", ctr.refresh)
	h.R.POST("resetPermission", ctr.setAdminPermission)
	group := ctr.H.R.Group("/api/v0")
	group.GET("list", ctr.list)
	group.POST("create", ctr.create)
	group.POST("ping", middleware.CasbinMiddleware(h.Enforcer, "/api/v0/ping", "GET"), ctr.ping)
	group.POST("login", ctr.login)
	group.POST("delete", ctr.delete)
	group.POST("updatePassword", ctr.editPassword)

}

func (ctr *adminController) ping(c *gin.Context) {
	//
	var response dto.Response
	response.Code = 0
	response.Message = "success"
	response.Data = "How are you...."
	c.JSON(http.StatusOK, response)

}

type ReqUpdatePassword struct {
	//
	Uid         uint64 `json:"uid" form:"uid"`
	OldPassword string `json:"old_password" form:"old_password"`
	NewPassword string `json:"new_password" form:"new_password"`
}

//@@@ edit admin password
func (ctr *adminController) editPassword(c *gin.Context) {
	//
	resp := dto.RespObj{}
	req := ReqUpdatePassword{}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		resp.ErrCode = 403
		resp.ErrMsg = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}

	eAdmin := &model.Admin{}
	//validate admin user
	err = ds.DB.Model(&model.Admin{}).Where("id = ?", req.Uid).First(&eAdmin).Error
	if err != nil {
		resp.ErrCode = 9020
		resp.ErrMsg = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}

	//validate old password
	ok, err := utils.ValidateHashedPassword(eAdmin.Password, req.OldPassword)
	if err != nil {
		resp.ErrCode = 9021
		resp.ErrMsg = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}
	if !ok {
		resp.ErrCode = 9022
		resp.ErrMsg = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}

	//update new password
	nPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		resp.ErrCode = 9023
		resp.ErrMsg = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}

	err = ds.DB.Model(&model.Admin{}).Where("id =?", req.Uid).Update("password", nPassword).Error
	if err != nil {
		resp.ErrCode = 9024
		resp.ErrMsg = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}

	resp.ErrCode = 0
	resp.ErrMsg = "success"
	c.JSON(http.StatusOK, resp)
}

type RespAdminList struct {
	ID          uint64 `json:"id"`
	UserName    string `json:"user_name"`
	LoginName   string `json:"login_name"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	Avatar      string `json:"avatar"`
	IsOnline    bool   `json:"is_online"`
	Gender      int8   `json:"gender"`
	LoginIP     string `json:"login_ip"`
	LastLoginIP string `json:"last_login_ip"`
	CreatedBy   string `json:"created_by"`
	UpdatedBy   string `json:"updated_by"`
	CreatedAt   string `json:"created_at"`
	DeletedAt   string `json:"deleted_at"`
}

//@@ admin list
func (ctr *adminController) list(c *gin.Context) {
	//
	ds.DB.Scopes()

	resp := dto.RespObj{}
	req := dto.ReqAdminList{}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		resp.ErrCode = 403
		resp.ErrMsg = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}
	respAdmin := make([]*RespAdminList, 0)

	fAdmin := make([]*model.Admin, 0)
	total := int64(0)
	db := ds.DB.Model(&model.Admin{})
	if len(req.LoginName) > 0 {
		db = db.Where("login_name LIKE ", fmt.Sprintf("%s%s%s", "%", req.LoginName, "%"))
	}
	if len(req.UserName) > 0 {
		db = db.Where("username LIKE ", fmt.Sprintf("%s%s%s", "%", req.UserName, "%"))
	}
	if req.Email != "" {
		db = db.Where("email LIKE ", fmt.Sprintf("%s%s%s", "%", req.Email, "%"))
	}
	db = db.Count(&total)
	db = db.Order("id DESC")
	db = db.Scopes(utils.BetweenDate(req.StartDate, req.EndDate))
	db = db.Scopes(utils.Paginate(req.Page, req.PageSize))
	err = db.Find(&fAdmin).Error
	if err != nil {
		resp.ErrCode = 9030
		resp.ErrMsg = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}
	for _, repAdmin := range fAdmin {
		respAdmin = append(respAdmin, &RespAdminList{
			ID:        repAdmin.ID,
			UserName:  repAdmin.UserName,
			LoginName: repAdmin.LoginName,
			Email:     repAdmin.Email,
			Password:  repAdmin.Password,
			Avatar:    repAdmin.Avatar,
			IsOnline:  repAdmin.IsOnline,
			Gender:    repAdmin.Gender,
			CreatedBy: repAdmin.CreadtedBy,
			UpdatedBy: repAdmin.UpdatedBy,
			CreatedAt: repAdmin.CreatedAt.Format("2006-01-02 15:04:05"),
			DeletedAt: repAdmin.DeletedAt.Time.Format("2006-01-02 15:04:05"),
		})
	}

	resp.ErrCode = 0
	resp.ErrMsg = "success"
	resp.Data = respAdmin
	c.JSON(http.StatusOK, resp)
}

//@@@ delete admin
func (ctr *adminController) delete(c *gin.Context) {
	//
	resp := dto.RespObj{}
	req := dto.ReqDeleteAdminAccount{}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		resp.ErrCode = 403
		resp.ErrMsg = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}

	tx := ds.DB.Begin()
	fAdmin := &model.Admin{}
	err = tx.Model(&model.Admin{}).Where("id = ?", req.AdminID).First(&fAdmin).Error
	if err == gorm.ErrRecordNotFound && err != nil {
		resp.ErrCode = 9010
		resp.ErrMsg = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	} else {
		err = tx.Model(&model.Roles{}).Where("id = ?", fAdmin.ID).Delete(&model.Roles{
			ID: fAdmin.ID,
		}).Error
		if err != nil {
			resp.ErrCode = 9012
			resp.ErrMsg = err.Error()
			c.JSON(http.StatusOK, resp)
			tx.Rollback()
			return
		}
		err = tx.Model(&model.AdminRole{}).Where("uid = ?", fAdmin.ID).Delete(&model.AdminRole{
			UserID: fAdmin.ID,
		}).Error
		if err != nil {
			resp.ErrCode = 9013
			resp.ErrMsg = err.Error()
			c.JSON(http.StatusOK, resp)
			tx.Rollback()
			return
		}
		err = tx.Model(&model.Admin{}).Delete(&fAdmin).Error
		if err != nil {
			resp.ErrCode = 9011
			resp.ErrMsg = err.Error()
			c.JSON(http.StatusOK, resp)
			tx.Rollback()
			return
		}

	}
	resp.ErrCode = 0
	resp.ErrMsg = "success"
	c.JSON(http.StatusOK, resp)
}

//@@@login
func (ctr *adminController) login(c *gin.Context) {

	resp := dto.RespObj{}
	req := dto.AdminLogin{}

	fmt.Println("Whole ....", c.Request)
	err := c.ShouldBindJSON(&req)
	if err != nil {
		resp.ErrCode = 403
		resp.ErrMsg = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}
	admin := &model.Admin{}
	tx := ds.DB.Begin()
	err = tx.Model(&model.Admin{}).Where("username = ?", req.Username).First(&admin).Error
	if err == gorm.ErrRecordNotFound && err != nil {
		resp.ErrCode = 508
		resp.ErrMsg = "user not found"
		c.JSON(http.StatusOK, resp)
		return
	}
	eAdmin := &model.Admin{
		ID:          admin.ID,
		UserName:    admin.UserName,
		Email:       admin.Email,
		Password:    admin.Password,
		Avatar:      admin.Avatar,
		Gender:      admin.Gender,
		LoginIP:     c.ClientIP(),
		LastLoginIP: c.ClientIP(),
		CreadtedBy:  admin.CreadtedBy,
		UpdatedBy:   admin.UpdatedBy,
		IsOnline:    true,
		LoginName:   req.Username,
	}

	err = tx.Model(&model.Admin{}).Where("id = ?", admin.ID).Updates(&eAdmin).Error
	if err != nil {
		resp.ErrCode = 509
		resp.ErrMsg = err.Error()
		c.JSON(http.StatusOK, resp)
		tx.Rollback()
		return
	}

	ok, _ := utils.ValidateHashedPassword(admin.Password, req.Password)
	if !ok {
		resp.ErrCode = 601
		resp.ErrMsg = "wrong password"
		c.JSON(http.StatusOK, resp)
		return
	}

	token, err := utils.GetAccessToken(admin.ID)
	if err != nil {
		resp.ErrCode = 602
		resp.ErrMsg = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}

	resp.ErrCode = 0
	resp.ErrMsg = "success"
	resp.Data = gin.H{
		"token": token,
	}
	c.JSON(http.StatusOK, resp)

}

//@@@ create admin
func (ctr *adminController) create(c *gin.Context) {
	//

	resp := dto.RespObj{}
	req := dto.ReqAdminCreate{}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		resp.ErrCode = 403
		resp.ErrMsg = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}
	eAdmin := c.MustGet("admin").(*model.Admin)
	cAdmin := &model.Admin{
		UserName:    req.Username,
		LoginName:   "",
		Email:       req.Email,
		Password:    req.Password,
		Avatar:      req.Avatar,
		IsOnline:    false,
		Gender:      req.Gender,
		LoginIP:     c.ClientIP(),
		LastLoginIP: c.ClientIP(),
		CreadtedBy:  eAdmin.UserName,
	}
	cRole := &model.Roles{
		Slug:      req.Slug,
		Name:      req.Username,
		Status:    1,
		CreatedBy: eAdmin.UserName,
	}
	admin := &model.Admin{}
	tx := ds.DB.Begin()
	err = tx.Model(&model.Admin{}).Where("email = ?", req.Email).First(&admin).Error
	if err == gorm.ErrRecordNotFound && err != nil {
		err = tx.Model(&model.Admin{}).Create(&cAdmin).Error
		if err != nil {
			resp.ErrCode = 9001
			resp.ErrMsg = err.Error()
			c.JSON(http.StatusOK, resp)
			tx.Rollback()
			return
		}
		err = tx.Model(&model.Roles{}).Create(&cRole).Error
		if err != nil {
			resp.ErrCode = 9002
			resp.ErrMsg = err.Error()
			c.JSON(http.StatusOK, resp)
			tx.Rollback()
			return
		}
		err = tx.Model(&model.AdminRole{}).Create(&model.AdminRole{
			UserID: cAdmin.ID,
			RoleID: cRole.ID,
		}).Error
		if err != nil {
			resp.ErrCode = 9003
			resp.ErrMsg = err.Error()
			c.JSON(http.StatusOK, resp)
			tx.Rollback()
			return
		}
	} else {
		resp.ErrCode = 9004
		resp.ErrMsg = "create admin already exists ...."
		c.JSON(http.StatusOK, resp)
		return
	}

	resp.ErrCode = 0
	resp.ErrMsg = "success"
	c.JSON(http.StatusOK, resp)

}

//@@ set route
func (ctr *adminController) refresh(c *gin.Context) {
	//

	//add policy route
	type Data struct {
		Route  string
		Method string
	}

	resp := dto.RespObj{}
	menu := []*Data{
		{
			Route:  "/api/v2/ping",
			Method: "GET",
		}, {
			Route:  "/api/v1/ping",
			Method: "GET",
		},
		{
			Route:  "/api/v0/ping",
			Method: "GET",
		},
		{
			Route:  "/api/v3/ping",
			Method: "GET",
		},
		{
			Route:  "/api/v4/ping",
			Method: "POST",
		},
	}

	for _, d := range menu {
		rule := &model.CasbinPloicy{
			Permission: d.Route,
			Action:     d.Method,
		}

		casbin := &model.CasbinPloicy{}
		err := ds.DB.Model(&model.CasbinPloicy{}).Where("permission = ?", d.Route).First(&casbin).Error
		if err == gorm.ErrRecordNotFound && err != nil {
			err := ds.DB.Model(&model.CasbinPloicy{}).Create(&rule).Error
			if err != nil {
				resp.ErrCode = http.StatusBadRequest
				resp.ErrMsg = err.Error()
				c.JSON(http.StatusOK, resp)
				return

			}
		} else {
			resp.ErrCode = 408
			resp.ErrMsg = "route already exists..."
			c.JSON(http.StatusOK, resp)
			continue
		}

	}

	resp.ErrCode = 0
	resp.ErrMsg = "success"
	c.JSON(http.StatusOK, resp)

}

//@@set permission
func (ctr *adminController) setAdminPermission(c *gin.Context) {
	//
	resp := dto.RespObj{}

	//setup admin
	err := factory.FactoryAdmin()
	if err != nil {
		resp.ErrCode = 801
		resp.ErrMsg = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}

	//setup admin role
	err = factory.FactoryRole()
	if err != nil {
		resp.ErrCode = 802
		resp.ErrMsg = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}

	var rules []*model.CasbinPloicy

	err = ds.DB.Model(&model.CasbinPloicy{}).Find(&rules).Error
	if err != nil {
		resp.ErrCode = http.StatusBadRequest
		resp.ErrMsg = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}

	//role
	adminRole := &model.Roles{}
	err = ds.DB.Model(&model.Roles{}).Where("slug = ?", "admin").First(&adminRole).Error
	if err != nil {
		resp.ErrCode = 509
		resp.ErrMsg = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}

	//assign role
	for _, principle := range rules {
		bol, err := principle.AssignPermisssion(ctr.H.Enforcer, principle, adminRole.Slug)
		if !bol || err != nil {
			resp.ErrCode = 700
			resp.ErrMsg = err.Error()
			c.JSON(http.StatusOK, resp)
			return
		}
	}

	resp.ErrCode = 0
	resp.ErrMsg = "success"
	c.JSON(http.StatusOK, resp)
}
