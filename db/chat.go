package db

type Chat struct {
	ID   uint64 // ID of the chat
	From string `gorm:"not null;"`      // Sender of the chat
	Body string `gorm:"not null;"`      // Body of the chat
	To   string `gorm:"not null;index"` // Reciever of the chat
	Time uint64 `gorm:"not null"`       // Time at which server recieved the chat
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
