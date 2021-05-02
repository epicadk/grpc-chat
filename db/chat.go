package db

type Chat struct {
	ID       uint64
	Sender   string `gorm:"not null"`
	Body     string `gorm:"not null"`
	Receiver string `gorm:"index;not null"`
	Sent     uint64 `gorm:"autoCreateTime"`
}

func (chat *Chat) SaveToDB() error {
	return DBconn.Create(chat).Error
}

func (chat *Chat) DeleteChat() error {
	return DBconn.Delete(chat, chat.ID).Error
}

func (chat *Chat) FindChat() ([]Chat, error) {
	var chats []Chat
	tx := DBconn.Where(chat).Find(&chats)
	return chats, tx.Error
}
