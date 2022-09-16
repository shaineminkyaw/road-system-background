package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shaineminkyaw/road-system-background/ds"
	"github.com/shaineminkyaw/road-system-background/dto"
	"github.com/shaineminkyaw/road-system-background/middleware"
	"github.com/shaineminkyaw/road-system-background/model"
	"github.com/shaineminkyaw/road-system-background/utils"
	"gorm.io/gorm"
)

//dragon tiger controller

type type2Controller struct {
	H *Handler
}

func NewType2Controller(h *Handler) *type2Controller {
	return &type2Controller{
		H: h,
	}
}

func (ctr *type2Controller) Register() {
	//
	h := ctr.H
	group := h.R.Group("/api/type2_table/", middleware.Cors(), middleware.AuthMiddleware())
	group.GET("list", ctr.list)
	group.POST("create", middleware.Authorize(h.Enforcer, "/api/type2_table/create", "POST"), ctr.create)
	group.POST("update", middleware.Authorize(h.Enforcer, "/api/type2_table/update", "POST"), ctr.edit)
	group.POST("delete", middleware.Authorize(h.Enforcer, "/api/type2_table/delete", "POST"), ctr.erase)
}

type ReqType2List struct {
	RoadType    int8   `json:"road_type" form:"road_type" binding:"required"`
	TableNumber uint64 `json:"table_no,omitempty" form:"table_no"`
	Title       string `json:"title,omitempty" form:"title"`
	Status      int8   `json:"status,omitempty" form:"omitempty"`
	StartDate   string `json:"start_date,omitempty" form:"start_date"`
	EndDate     string `json:"end_date,omitempty" form:"end_date"`
	Page        int    `json:"page,omitempty" form:"page"`
	PageSize    int    `json:"page_size,omitempty" form:"page_size"`
}

type RespTableType2List struct {
	TableNumber      uint64 `json:"table_no"`
	RoadType         int8   `json:"road_type"`
	Title            string `json:"title"`
	Password         string `json:"password"`
	Type             int8   `json:"type"`
	Cover            string `json:"cover"`
	Placard          string `json:"placard"`
	IpLimit          string `json:"ip_limit"`
	AskTime          string `json:"ask_time"`
	Status           int8   `json:"status"`
	OnlineUserNumber int    `json:"online_user"`
	BetTime          int    `json:"bet_time"`
	TableRound       int64  `json:"bs_round"`
	MatchRound       int64  `json:"oe_round"`
	FirstLimit       string `json:"player_limit"`
	SecondLimit      string `json:"banker_limit"`
	ThirdLimit       string `json:"pair_limit"`
	FourthLimit      string `json:"tie_limit"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}

//type 2 list
func (ctr *type2Controller) list(c *gin.Context) {
	//
	resp := &dto.RespObj{}
	req := ReqType2List{}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		resp.ErrCode = 403
		resp.ErrMsg = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}

	respBet := make([]*RespTableType2List, 0)
	bets := make([]*model.BetTable, 0)
	total := int64(0)
	db := ds.DB.Model(&model.BetTable{})
	if req.TableNumber > 0 {
		db = db.Where("table_no = ?", req.TableNumber)
	}
	if len(req.Title) > 0 {
		db = db.Where("title = ?", req.Title)
	}
	if req.Status > 0 {
		db = db.Where("status = ?", req.Status)
	}
	db = db.Where("road_type = ?", req.RoadType)
	db = db.Count(&total)
	db = db.Order("id DESC")
	db = db.Scopes(utils.BetweenDate(req.StartDate, req.EndDate))
	db = db.Scopes(utils.Paginate(req.Page, req.PageSize))
	err = db.Find(&bets).Error
	if err != nil {
		resp.ErrCode = 2006
		resp.ErrMsg = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}

	for _, bet := range bets {
		respBet = append(respBet, &RespTableType2List{
			TableNumber:      bet.TableNumber,
			RoadType:         bet.RoadType,
			Password:         bet.Password,
			Type:             bet.Type,
			Title:            bet.Title,
			Cover:            bet.Cover,
			Placard:          bet.Placard,
			IpLimit:          bet.IPLimit,
			AskTime:          bet.AskTime,
			Status:           bet.Status,
			OnlineUserNumber: bet.OnlineUserNumber,
			BetTime:          bet.BetTime,
			TableRound:       bet.TableRound,
			MatchRound:       bet.MatchRound,
			FirstLimit:       bet.FirstLimit,
			SecondLimit:      bet.SecondLimit,
			ThirdLimit:       bet.ThirdLimit,
			FourthLimit:      bet.FourthLimit,
			CreatedAt:        bet.CreatedAt.Format("2006-04-05 15:04:05"),
			UpdatedAt:        bet.UpdatedAt.Format("2006-04-05 15:04:05"),
		})
	}
	resp.ErrCode = 0
	resp.ErrMsg = "success"
	resp.Data = respBet
	c.JSON(http.StatusOK, resp)

}

//type 2 request
type ReqType2Create struct {
	TableType   int8   `json:"table_type" form:"table_type" binding:"required"`
	TableNumber uint64 `json:"table_no" form:"table_no" binding:"required"`
	Title       string `json:"title,omitempty" form:"title"`
	Password    string `json:"password,omitempty" form:"password"`
	PlayerLimit string `json:"player_limit,omitempty" form:"player_limit"`
	BankerLimit string `json:"banker_limit,omitempty" form:"banker_limit"`
	TieLimit    string `json:"tie_limit,omitempty" form:"tie_limit"`
	AskTime     string `json:"ask_time,omitempty" form:"ask_time"`
	BetTime     int    `json:"bet_time,omitempty" form:"bet_time"`
	MatchRound  int64  `json:"oe_round,omitempty" form:"oe_round"`
	Status      int8   `json:"status,omitempty" form:"status"`
	Placard     string `json:"placard,omitempty" form:"placard"`
	IpLimit     string `json:"ip_limit,omitempty" form:"ip_limit"`
	Type        int8   `json:"type,omitempty" form:"type"`
}

//create type2
func (ctr *type2Controller) create(c *gin.Context) {
	//
	resp := &dto.RespObj{}
	req := ReqType2Create{}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		resp.ErrCode = 403
		resp.ErrMsg = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}
	cBet := &model.BetTable{
		RoadType:    req.TableType,
		TableNumber: req.TableNumber,
		Title:       req.Title,
		Password:    req.Password,
		FirstLimit:  req.PlayerLimit,
		SecondLimit: req.BankerLimit,
		FourthLimit: req.TieLimit,
		AskTime:     req.AskTime,
		BetTime:     req.BetTime,
		MatchRound:  req.MatchRound,
		Status:      req.Status,
		Placard:     req.Placard,
		IPLimit:     req.IpLimit,
		Type:        req.Type,
	}
	fBet := &model.BetTable{}
	err = ds.DB.Model(&model.BetTable{}).Where("road_type = ? and table_no = ?", req.TableType, req.TableNumber).First(&fBet).Error
	if err == gorm.ErrRecordNotFound && err != nil {
		err = ds.DB.Model(&model.BetTable{}).Create(&cBet).Error
		if err != nil {
			resp.ErrCode = 2007
			resp.ErrMsg = err.Error()
			c.JSON(http.StatusOK, resp)
			return
		}
	} else {
		resp.ErrCode = 2008
		resp.ErrMsg = "table already exists ..."
		c.JSON(http.StatusOK, resp)
		return
	}
	resp.ErrCode = 0
	resp.ErrMsg = "success"
	c.JSON(http.StatusOK, resp)
}

//type2 edit
func (ctr *type2Controller) edit(c *gin.Context) {
	//

	resp := &dto.RespObj{}
	req := ReqType2Create{}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		resp.ErrCode = 403
		resp.ErrMsg = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}
	fBet := &model.BetTable{}
	eBet := &model.BetTable{
		RoadType:    req.TableType,
		TableNumber: req.TableNumber,
		Title:       req.Title,
		Password:    req.Password,
		FirstLimit:  req.PlayerLimit,
		SecondLimit: req.BankerLimit,
		FourthLimit: req.TieLimit,
		AskTime:     req.AskTime,
		BetTime:     req.BetTime,
		MatchRound:  req.MatchRound,
		Status:      req.Status,
		Placard:     req.Placard,
		IPLimit:     req.IpLimit,
		Type:        req.Type,
	}

	err = ds.DB.Model(&model.BetTable{}).Where("road_type = ? and table_no = ?", req.TableType, req.TableNumber).First(&fBet).Error
	if err == gorm.ErrRecordNotFound && err != nil {
		resp.ErrCode = 2009
		resp.ErrMsg = "not found ..."
		c.JSON(http.StatusOK, resp)
		return
	} else {
		err = ds.DB.Model(&model.BetTable{}).Where("road_type = ? and table_no = ?", req.TableType, req.TableNumber).Updates(&eBet).Error
		if err != nil {
			resp.ErrCode = 2010
			resp.ErrMsg = err.Error()
			c.JSON(http.StatusOK, resp)
			return
		}
	}

	resp.ErrCode = 0
	resp.ErrMsg = "success"
	c.JSON(http.StatusOK, resp)
}

//delete request for bacarat
type ReqDeleteType2 struct {
	TableType uint64 `json:"road_type"`
	ID        uint64 `json:"tid"`
}

//type 2 delete
func (ctr *type2Controller) erase(c *gin.Context) {
	//
	resp := &dto.RespObj{}
	req := ReqDeleteType2{}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		resp.ErrCode = 403
		resp.ErrMsg = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}

	err = ds.DB.Model(&model.BetTable{}).Delete(&model.BetTable{
		ID:       req.ID,
		RoadType: int8(req.TableType),
	}).Error
	if err != nil {
		resp.ErrCode = 2012
		resp.ErrMsg = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}

	resp.ErrCode = 0
	resp.ErrMsg = "success"
	c.JSON(http.StatusOK, resp)
}
