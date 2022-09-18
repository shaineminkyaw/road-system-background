package middleware

import (
	"fmt"
	"net/http"

	"github.com/shaineminkyaw/road-system-background/config"
	"github.com/shaineminkyaw/road-system-background/ds"
	"github.com/shaineminkyaw/road-system-background/dto"
	"github.com/shaineminkyaw/road-system-background/model"
	"github.com/shaineminkyaw/road-system-background/utils"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

type HeaderAuth struct {
	Header string `header:"Authorization"`
}

func AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		res := &dto.RespObj{}
		h := &HeaderAuth{}
		err := c.ShouldBindHeader(&h)
		if err != nil {
			res.ErrCode = 409
			res.ErrMsg = "Error on header binding"
			c.JSON(http.StatusOK, res)
			return
		}
		if h.Header == "" {
			res.ErrCode = 403
			res.ErrMsg = "Unprivileged to vie this page"
			c.Abort()
			c.JSON(http.StatusOK, res)
			return
		}
		token, err := utils.ValidateAccessToken(h.Header, config.PublicKey)
		if err != nil {
			res.ErrCode = 403
			res.ErrMsg = "Unauthorized"
			c.Abort()
			c.JSON(http.StatusOK, res)
			return
		}

		var user *model.Admin
		db := ds.DB.Model(&model.Admin{})
		db = db.Where("id = ?", token.UserID)
		err = db.First(&user).Error
		if err != nil {
			res.ErrCode = 400
			res.ErrMsg = "Unauthorized"
			c.Abort()
			c.JSON(http.StatusOK, res)
			return
		}
		c.Set("admin", user)
		c.Next()
	}
}

func Authorize(e *casbin.Enforcer, permission, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp := dto.RespObj{}
		// now existing admin
		user, bol := c.Get("admin")

		if !bol {
			resp.ErrCode = 501
			resp.ErrMsg = "Unauthorized"
			c.JSON(http.StatusOK, resp)
			c.Abort()
			return
		}

		admin := user.(*model.Admin) //bind admin data
		fmt.Printf("Data .....%v\n", admin)
		admin1 := &model.Admin{}
		err := ds.DB.Model(&model.Admin{}).Where("id = ?", admin.ID).First(&admin1).Error
		if err != nil {
			resp.ErrCode = 900
			resp.ErrMsg = "Unauthorized"
			c.JSON(http.StatusOK, resp)
			c.Abort()
			return
		}

		//find admin ids
		rolesID := &model.AdminRole{}
		err = ds.DB.Model(&model.AdminRole{}).Where("uid = ?", admin1.ID).First(&rolesID).Error
		if err != nil {
			resp.ErrCode = 402
			resp.ErrMsg = "Unauthorized"
			c.JSON(http.StatusOK, resp)
			c.Abort()
			return
		}

		fmt.Println("IDS ....", rolesID)

		//find admins with admin ids
		userRoles := &model.Roles{}
		err = ds.DB.Model(&model.Roles{}).Where("id = ?", rolesID.RoleID).First(&userRoles).Error
		if err != nil {
			resp.ErrCode = 403
			resp.ErrMsg = "Unauthorized"
			c.JSON(http.StatusOK, resp)
			c.Abort()
			return
		}

		err = e.LoadPolicy() //reload policy
		if err != nil {
			resp.ErrCode = 405
			resp.ErrMsg = "Unauthorized"
			c.JSON(http.StatusOK, resp)
			c.Abort()
			return
		}

		ok, err := e.Enforce(fmt.Sprintf(userRoles.Slug), permission, action) //check permission each admin
		if !ok || err != nil {
			resp.ErrCode = 406
			resp.ErrMsg = "Unauthorized"
			c.JSON(http.StatusOK, resp)
			c.Abort()
			return
		}

		c.Next()
	}
}
