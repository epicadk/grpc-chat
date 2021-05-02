package db

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DBconn *gorm.DB

func init() {
	var err error
	// TODO use env vars
	host, check := os.LookupEnv("DB_HOST")
	if !check {
		host = "db"
	}
	user, check := os.LookupEnv("DB_USER")
	if !check {
		user = "postgres"
	}
	pass, check := os.LookupEnv("DB_PASSWORD")
	if !check {
		pass = "postgres"
	}
	dbname, check := os.LookupEnv("DB_NAME")
	if !check {
		dbname = "chats"
	}
	port, check := os.LookupEnv("DB_PORT")
	if !check {
		port = "5432"
	}
	tz, check := os.LookupEnv("DB_TIMEZONE")
	if !check {
		// TODO use hosts tz
		tz, _ = time.Now().Zone()
		fmt.Println(tz)
	}
	sslmode, check := os.LookupEnv("DB_SSLMODE")
	if !check {
		sslmode = "disable"
	}

	dns := fmt.Sprintf("host=%s user=%s password =%s dbname=%s port=%s sslmode=%s TimeZone=%s", host, user, pass, dbname, port, sslmode, tz)
	//"host=db user=postgres password=postgres dbname=chats port=5432 sslmode=disable TimeZone=Asia/Kolkata"

	DBconn, err = gorm.Open(postgres.Open(dns), &gorm.Config{})

	if err != nil {
		panic("error connecting to database")
	}

	err = DBconn.AutoMigrate(Chat{})
	if err != nil {
		panic("error in auto migration")
	}
	DBconn.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	err = DBconn.AutoMigrate(&User{})
	if err != nil {
		panic("error in auto migration")
	}
}
