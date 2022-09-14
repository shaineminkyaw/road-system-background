package controller

import "github.com/gin-gonic/gin"

type betController struct {
	H *Handler
}

func NewBetController(h *Handler) *betController {
	return &betController{
		H: h,
	}
}

func (ctr *betController) Register() {
	h := ctr.H
	group := h.R.Group("/api/table")
	group.GET("list", ctr.list)
	group.POST("create", ctr.create)
	group.POST("update", ctr.edit)
	group.POST("delete", ctr.erase)
}

func (ctr *betController) list(c *gin.Context) {
	//
}

func (ctr *betController) create(c *gin.Context) {
	//

}

func (ctr *betController) edit(c *gin.Context) {
	//
}

func (ctr *betController) erase(c *gin.Context) {
	//
}
