package main

import (
	"database/sql"
	"errors"
)

type Dialog struct {
	Id    				int     `json:"id"`
	OwnerId 			int  	`json:"name"`
	Status	 			int     `json:"status"`
}

type User struct {
	Id 					int 	`json:"id"`
	UserName 			string 	`json:"user_name"`
	Password 			string  `json:"password"`
	Token 				string	`json:"token"`
	Sex					string	`json:"sex"`

}

type DialogUser struct {
	UserId 	   			int 	`json:"user_id"`
	DialogId 			int 	`json:"dialog_id"`
	LastMessageId		int		`json:"last_message_id"`
	IsRead				int		`json:"is_read"`
	LastReadDate		int		`json:"last_read_date"`
	CountNewMessages 	int		`json:"count_new_messages"`

}

type Message struct {
	Id 					int 	`json:"id"`
	Text 				string	`json:"text"`
	UserId 				int		`json:"user_id"`
	Created				int		`json:"created"`
	DialogId			int		`json:"dialog_id"`
	DeletedForMe		int		`json:"deleted_for_me"`
	DeletedForAll		int		`json:"deleted_for_all"`
}


func (d *Dialog) read(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (d *Dialog) create(db *sql.DB) error {
	return errors.New("Not implemented")
}

func getDialogs(db *sql.DB) ([]DialogUser, error) {
	rows, err := db.Query("SELECT * FROM dialog_users")

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	dialogs := []DialogUser{}

	for rows.Next() {
		var d DialogUser
		if err := rows.Scan(&d.DialogId, &d.UserId, &d.CountNewMessages); err != nil {
			return nil, err
		}
		dialogs = append(dialogs, d)
	}

	return dialogs, nil

}

/**
 User methods
 */
func (u *User) login(db *sql.DB) error {
	return errors.New("No implements")
}

func (u *User) registration(db *sql.DB) error {
	return errors.New("No implements")
}

func (u *User) exist(db *sql.DB) error {
	return errors.New("No implements")
}


/**
 Message methods
 */
func (m * Message) send(db *sql.DB) error  {
	return errors.New("No implement")
}

func (m * Message) remove(db *sql.DB) error  {
	return errors.New("No implement")
}

