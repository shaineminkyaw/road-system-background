package model

import "time"

type UserLog struct {
	Id        uint64    `gorm:"column:id" json:"id"`
	UID       uint64    `gorm:"column:uid" json:"uid"`
	Type      int8      `gorm:"column:type" json:"type"`
	Ip        string    `gorm:"column:ip" json:"ip"`
	Info      string    `gorm:"column:log-info" json:"log-info"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (ul *UserLog) TableName() string {
	return "user-log"
}
