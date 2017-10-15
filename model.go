package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"
)

type Dialog struct {
	Id               int       `json:"id"`
	OwnerId          int       `json:"owner_id"`
	ContactId        int       `json:"contact_id"`
	LastMessageId    int       `json:"last_message_id"`
	IsRead           int       `json:"is_read"`
	LastReadDate     time.Time `json:"last_read_date"`
	CountNewMessages int       `json:"count_new_messages"`
	Status           int       `json:"status"`
}

type Room struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Type    int    `json:"type"`
	OwnerId int    `json:"owner_id"`
	Created string `json:"created"`
}

type User struct {
	Id        int       `json:"id"`
	UserName  string    `json:"user_name"`
	Password  string    `json:"-"`
	Token     string    `json:"token"`
	Sex       string    `json:"sex"`
	DateBirth time.Time `json:"date_birth"`
}

type DialogMessage struct {
	Id        int       `json:"id"`
	Text      string    `json:"text"`
	OwnerId   int       `json:"owner_id"`
	ContactId int       `json:"contact_id"`
	Created   time.Time `json:"created"`
	DialogId  int       `json:"dialog_id"`
}

type RoomMessage struct {
	Id      int       `json:"id"`
	Text    string    `json:"text"`
	OwnerId int       `json:"owner_id"`
	Created time.Time `json:"created"`
	RoomId  int       `json:"room_id"`
}

type RoomUser struct {
	RoomId              int       `json:"room_id"`
	UserId              int       `json:"user_id"`
	Created             time.Time `json:"created"`
	LastDateVisitToRoom time.Time `json:"last_date_visit_to_room"`
}

type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type Token struct {
	Token string `json:"token"`
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
func (u *User) login(db *sql.DB) {
	db.QueryRow("SELECT id, sex, token FROM users WHERE user_name=$1 AND password=$2",
		u.UserName, u.Password).Scan(&u.Id, &u.Sex, &u.Token)
}

func (u *User) registration(db *sql.DB) error {
	return errors.New("No implements")
}

//get user info
func (u *User) me(a *App, key string) error {
	val, err := a.Redis.Get(key).Result()
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(val), &u)
	return err
}

/**
Message methods
*/
func (m *DialogMessage) send(db *sql.DB) error {
	dialog := &Dialog{
		OwnerId:   m.OwnerId,
		ContactId: m.ContactId,
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

func (m *DialogMessage) save(db *sql.DB, dialogId int) error {
	_, err := db.Exec("INSERT INTO dialog_messages(owner_id, contact_id, dialog_id, text)"+
		"VALUES($1,$2,$3,$4,$5)", m.OwnerId, m.ContactId, dialogId, m.Text)
	if err != nil {
		return err
	}

	return nil
}

func (m *DialogMessage) remove(db *sql.DB) error {
	return errors.New("No implement")
}
func (m *DialogMessage) getMessages(db *sql.DB, dialog_id int) error {
	return errors.New("No implement")
}

//rooms methods
func (room *Room) create(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO rooms(name, owner_id) VALUES($1, $2) RETURNING id,created,type",
		room.Name, room.OwnerId).Scan(&room.Id, &room.Created, &room.Type)

	if err != nil {
		return err
	}

	return nil
}

func (room_user *RoomUser) join(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO room_users(room_id, user_id, last_date_visit_to_room) VALUES($1, $2, $3) RETURNING last_date_visit_to_room",
		room_user.RoomId, room_user.UserId, time.Now()).Scan(&room_user.LastDateVisitToRoom)

	if err != nil {
		return err
	}

	return nil
}

func (m *RoomMessage) Send(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO room_messages(room_id, owner_id, text) VALUES($1, $2, $3) RETURNING id",
		m.RoomId, m.OwnerId, m.Text).Scan(&m.Id)

	if err != nil {
		return err
	}

	return nil
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
