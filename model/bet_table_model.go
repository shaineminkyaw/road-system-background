package model

import "time"

type BetTable struct {
	ID                uint64    `gorm:"column:id" json:"id"`
	TableNumber       uint64    `gorm:"column:table_no" json:"table_no"`
	Password          string    `gorm:"column:password" json:"password"`
	Title             string    `gorm:"column:title" json:"title"`
	Type              int8      `gorm:"column:type" json:"type"`
	Cover             string    `gorm:"column:cover;comment:dealer profile or game rule photo" json:"cover"` //dealer profile or game rule photo
	Placard           string    `gorm:"column:placard" json:"placard"`                                       //table's marquee text
	IPLimit           string    `gorm:"column:ip_limit" json:"ip_limit"`
	FirstLimit        string    `gorm:"column:first_limit" json:"first_limit"`
	SecondLimit       string    `gorm:"column:second_limit" json:"second_limit"`
	ThirdLimit        string    `gorm:"column:third_limit" json:"third_limit"`
	FourthLimit       string    `gorm:"column:fourth_limit" json:"fourth_limit"`
	FifthLimit        string    `gorm:"column:fifth_limit" json:"fifth_limit"`
	AskTime           string    `gorm:"column:ask_time;default:5" json:"ask_time"`
	Status            int8      `gorm:"column:status;default:1;comment:table status" json:"status"`
	OnlineUserNumber  int       `gorm:"column:online_user_number" json:"online_user_number"`
	BetTime           int       `gorm:"column:bet_time;default:30;comment:countdown time" json:"bet_time"`
	FrontType         int8      `gorm:"column:front_type;comment:terminal type" json:"front_type"`
	TableRound        int64     `gorm:"column:bs_round;default:0;comment:table's round" json:"bs_round"`
	MatchRound        int64     `gorm:"column:oe_round;default:0;comment:match 's round" json:"oe_round"`
	TableStatusCode   string    `gorm:"column:status_code;default:0;comment:table status code" json:"status code"`
	RoadResultHistory string    `gorm:"column:road;comment:result history" json:"road"`
	CreatedAt         time.Time `gorm:"created_at" json:"created_at"`
	UpdatedAt         time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (bt *BetTable) TableName() string {
	return "bet_table_info"
}
