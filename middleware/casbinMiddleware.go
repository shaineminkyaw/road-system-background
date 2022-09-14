package middleware

import (
	"casbin/ds"
	"casbin/dto"
	"casbin/model"
	"fmt"
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

type authHeaders struct {
	Auth string `header:"Authorization"`
}

func CasbinMiddleware(e *casbin.Enforcer, permission, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp := dto.RespObj{}
		fmt.Println("Headers .....", c.Request.Header)
		fmt.Println("BODY...", c.Request.Body)
		fmt.Println("Request....", c.Request)
		// now existing admin
		// user, bol := c.Get("admin")
		// admin := c.MustGet("admin").(*model.Admin)
		// fmt.Println("User....", user)

		h := authHeaders{}
		err := c.ShouldBindHeader(&h)
		if err != nil {
			resp.ErrCode = 506
			resp.ErrMsg = err.Error()
			c.JSON(http.StatusOK, resp)
			c.Abort()
			return
		}
		// if !bol {
		// 	resp.ErrCode = 401
		// 	resp.ErrMsg = "Unauthorized"
		// 	c.JSON(http.StatusOK, resp)
		// 	c.Abort()
		// 	return
		// }
		// admin := user.(*model.Admin) //bind admin data
		user := &model.Admin{}
		err = ds.DB.Model(&model.Admin{}).Where("username = ?", "admin").First(&user).Error
		if err != nil {
			resp.ErrCode = 900
			resp.ErrMsg = err.Error()
			c.JSON(http.StatusOK, resp)
			c.Abort()
			return
		}
		fmt.Println("ADMIN ...", user)

		//find admin ids
		rolesID := &model.AdminRole{}
		err = ds.DB.Model(&model.AdminRole{}).Where("uid = ?", user.ID).First(&rolesID).Error
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
		err = ds.DB.Model(&model.Roles{}).Where("id = ?", rolesID.RoleID).Find(&userRoles).Error
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

		// // var userName string
		// userName := c.GetHeader("admin")
		// if userName == "" {
		// 	fmt.Println("headers invalid")
		// 	c.JSON(200, gin.H{
		// 		"code":    401,
		// 		"message": "Unauthorized",
		// 		"data":    "",
		// 	})
		// 	c.Abort()
		// 	return
		// }
		// // 请求的path
		// p := c.Request.URL.Path
		// // 请求的方法
		// m := c.Request.Method

		// res, err := e.Enforce(userName, p, m)
		// if err != nil {
		// 	fmt.Println("no permission ")
		// 	fmt.Println(err)
		// 	c.JSON(200, gin.H{
		// 		"code":    401,
		// 		"message": "Unauthorized",
		// 		"data":    "",
		// 	})
		// 	c.Abort()
		// 	return
		// }
		// if !res {
		// 	fmt.Println("permission check failed")
		// 	c.JSON(200, gin.H{
		// 		"code":    401,
		// 		"message": "Unauthorized",
		// 		"data":    "",
		// 	})
		// 	c.Abort()
		// 	return
		// }
		c.Next()
	}
}
