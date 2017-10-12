package main

import (
	"database/sql"
	"errors"
	"crypto/rand"
	"encoding/base64"
)

type Dialog struct {
	Id    				int     `json:"id"`
	OwnerId 			int  	`json:"owner_id"`
	ContactId			int 	`json:"contact_id"`
	LastMessageId		int		`json:"last_message_id"`
	IsRead				int		`json:"is_read"`
	LastReadDate		int		`json:"last_read_date"`
	CountNewMessages 	int		`json:"count_new_messages"`
	Status	 			int     `json:"status"`
}

type Room struct {
	Id 					int 	`json:"id"`
	Name 				string 	`json:"name"`
	Type 				int 	`json:"type"`
	OwnerId				int 	`json:"owner_id"`
	Created				int		`json:"created"`
}

type User struct {
	Id 					int 	`json:"id"`
	UserName 			string 	`json:"user_name"`
	Password 			string  `json:"-"`
	Token 				string	`json:"token"`
	Sex					string	`json:"sex"`
	DateBirth 			int 	`json:"date_birth"`

}

type DialogMessage struct {
	Id 					int 	`json:"id"`
	Text 				string	`json:"text"`
	OwnerId 			int		`json:"owner_id"`
	ContactId 			int 	`json:"contact_id"`
	Created				int		`json:"created"`
	DialogId			int		`json:"dialog_id"`
}

type RoomMessage struct {
	Id 					int 	`json:"id"`
	Text 				string	`json:"text"`
	OwnerId 			int		`json:"owner_id"`
	Created				int		`json:"created"`
	RoomId 				int 	`json:"room_id"`
}

type RoomUser struct {
	RoomId 				int 	`json:"room_id"`
	UserId 				int 	`json:"user_id"`
	Created 			int 	`json:"created"`
	DateLastUserMessage int 	`json:"date_last_user_message"`
}

type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}


func (d *Dialog) read(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (d *Dialog) getExistOrCreate(db *sql.DB) error {
	return errors.New("Not implemented")
}



func getDialogs(db *sql.DB, id int) ([]Dialog, error) {
	rows, err := db.Query("SELECT * FROM dialogs WHERE owner_id=$1", id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	dialogs := []Dialog{}

	for rows.Next() {
		var d Dialog
		if err := rows.Scan(&d.Id, &d.OwnerId, &d.CountNewMessages); err != nil {
			return nil, err
		}
		dialogs = append(dialogs, d)
	}

	return dialogs, nil

}

/**
 User methods
 */
func (u *User) login(db *sql.DB){
	 db.QueryRow("SELECT id, sex, token FROM users WHERE user_name=$1 AND password=$2",
		u.UserName,u.Password).Scan(&u.Id, &u.Sex, &u.Token)
}

func (u *User) registration(db *sql.DB) error {
	return errors.New("No implements")
}

/**
 Message methods
 */
func (m * DialogMessage) send(db *sql.DB)  error  {
	dialog := &Dialog{
		OwnerId : m.OwnerId,
		ContactId : m.ContactId,
	}

	var err error
	if err = dialog.getExistOrCreate(db); err != nil {
		return err
	}

	err = m.save(db, dialog.Id)
	if err != nil {
		return err
	}

	return nil

}

func (m * DialogMessage) save(db *sql.DB, dialogId int) error  {
	_, err := db.Exec("INSERT INTO dialog_messages(owner_id, contact_id, dialog_id, text)" +
		"VALUES($1,$2,$3,$4,$5)", m.OwnerId, m.ContactId, dialogId, m.Text)
	if err != nil {
		return err
	}

	return nil
}

func (m * DialogMessage) remove(db *sql.DB) error  {
	return errors.New("No implement")
}
func (m *DialogMessage) getMessages(db *sql.DB, dialog_id int) (error)  {
	return errors.New("No implement")
}


// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}
