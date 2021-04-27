package db

type Chat struct {
	ID       uint64
	Sender   string
	Body     string
	Reciever string `gorm:"index"`
	Sent     uint64 `gorm:"autoCreateTime"`
}
