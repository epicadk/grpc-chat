package db

import "gorm.io/gorm"

type Chat struct {
	gorm.Model
	Sender   string
	Body     string
	Reciever string `gorm:"index"`
	Sent     uint64 `gorm:"autoCreateTime"`
}
