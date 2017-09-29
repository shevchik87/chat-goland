package main

import (
	"database/sql"
	"errors"
)

type Dialog struct {
	Id    		int     `json:"id"`
	OwnerId 	int  	`json:"name"`
	Status	 	int     `json:"status"`
}

type User struct {
	Id 		int 	`json:"id"`
	UserName 	string 	`json:"user_name"`
	Password 	string  `json:"password"`
	Token 		string	`json:"token"`

}

type DialogUser struct {
	UserId 	   	int 	`json:"user_id"`
	DialogId 	int 	`json:"dialog_id"`
	LastMessageId	int	`json:"last_message_id"`
	IsRead		int	`json:"is_read"`
	LastReadDate	int	`json:"last_read_date"`
	CountNewMessages int	`json:"count_new_messages"`

}

type Message struct {
	Id 		int 	`json:"id"`
	Text 		string	`json:"text"`
	UserId 		int	`json:"user_id"`
	Created		int	`json:"created"`
	DialogId	int	`json:"dialog_id"`
	DeletedForMe	int	`json:"deleted_for_me"`
	DeletedForAll	int	`json:"deleted_for_all"`
}


func (p *Dialog) getDialog(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (p *dialog) updateProduct(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (p *dialog) deleteProduct(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (p *dialog) createProduct(db *sql.DB) error {
	return errors.New("Not implemented")
}

func getdialogs(db *sql.DB, start, count int) ([]dialog, error) {
	return nil, errors.New("Not implemented")
}
