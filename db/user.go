package db

type User struct {
	Id       string `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
}
