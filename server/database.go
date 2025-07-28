package main

import (
	"os"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const dbFileName string = "blabber.db"

var db *gorm.DB

func Initialize() error {
	var err error
	db, err = gorm.Open(sqlite.Open(GetPlayPath(dbFileName)), &gorm.Config{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&Message{})
	if err != nil {
		return err
	}
	return nil
}

func Destruct() error {
	err := os.Truncate(GetPlayPath(dbFileName), 0)
	if err != nil {
		return err
	}
	return nil
}

type Message struct {
	gorm.Model
	Username string
	Message  string
}

func NewMessage(m Message) (*Message, error) {
	err := db.Create(&m).Error
	if err != nil {
		return nil, err
	}

	return &m, nil
}

func ListMessages() ([]Message, error) {
	var cas []Message
	err := db.Find(&cas).Error
	if err != nil {
		return nil, err
	}
	return cas, nil
}
