package db

type Chat struct {
	ID   uint64 // ID of the chat
	From User   `gorm:"not null;foreignKey:Phonenumberconstraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`            // Sender of the chat
	Body string `gorm:"not null"`                                                                                // Body of the chat
	To   User   `gorm:"index;not null;foreignKey:Phonenumber;constraint : OnUpdate:CASCADE,OnDelete : CASCADE;"` // Reciever of the chat
	Time uint64 `gorm:"autoCreateTime"`                                                                          // Time at which server recieved the chat
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
