package dto

//dummy data
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

//delete adminReq
type ReqDeleteAdminAccount struct {
	AdminID uint64 `json:"aid" form:"aid" binding:"required"`
}

//resp data
type RespObj struct {
	ErrCode uint64      `json:"err_code"`
	ErrMsg  string      `json:"err_msg"`
	Data    interface{} `json:"data,omitempty"`
}

//req admin login
type AdminLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//reqAdmin
type ReqAdminCreate struct {
	Email    string `json:"email" form:"email" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
	Username string `json:"username,omitempty" form:"username"`
	Avatar   string `json:"avatar" form:"avatar"`
	Gender   int8   `json:"gender" form:"gender"`
	Slug     string `json:"slug" form:"slug" binding:"required"`
}

//ReqAdminList
type ReqAdminList struct {
	LoginName string `json:"login_name,omitempty" form:"login_name"`
	UserName  string `json:"username,omitempty" form:"username"`
	Email     string `json:"email,omitempty" form:"email"`
	StartDate string `json:"start_date,omitempty" form:"start_date"`
	EndDate   string `json:"end_date,omitempty" form:"end_date"`
	Page      int    `json:"page,omitempty" form:"page"`
	PageSize  int    `json:"page_size,omitempty" form:"page_size"`
}
