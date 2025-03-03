package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Session struct {
	SessionID   string    `json:"session_id"`
	UserID      int       `json:"user_id"`
	LastEntered time.Time `json:"last_entered"`
}
type Users struct {
	User_id      int    `json:"user_id"`
	Username     string `json:"username"`
	PasswordHash string `json:"passwordhash"`
	Online       bool   `json:"online"`
}

var PathToDB = "user=postgres dbname=persondata password=sosuhui247 sslmode=disable"

func ConnectToDB() *sql.DB {
	db, _ := sql.Open("postgres", PathToDB)
	return db
}

type PersonRepository struct {
	Dbrepos *sql.DB
}

func NewPersonRepository(db *sql.DB) *PersonRepository {
	return &PersonRepository{Dbrepos: db}
}
func (pc *PersonRepository) CheckUniqueReg(v Users, ctx context.Context) bool {

	var id int
	err := pc.Dbrepos.QueryRowContext(ctx, "SELECT user_id from users where name = $1", v.Username).Scan(&id)
	if err == sql.ErrNoRows {
		return true
	} else {
		return false
	}
}
func (pc *PersonRepository) CheckPassword(v Users, ctx context.Context) error {
	var password string
	err := pc.Dbrepos.QueryRowContext(ctx, "SELECT password from users where name = $1", v.Username).Scan(&password)
	if err != nil {

		return fmt.Errorf("Ошибка!")
	}
	return bcrypt.CompareHashAndPassword([]byte(password), []byte(v.PasswordHash))
}
func (pc *PersonRepository) Delete(v Users, ctx context.Context) error {
	_, err := pc.Dbrepos.ExecContext(ctx, "DELETE FROM users where user_id = $1", v.User_id)
	if err != nil {
		return err
	}
	return nil
}
func (pc *PersonRepository) InsertReg(v Users, ctx context.Context) (bool, int) {
	var id int
	bytes, err := bcrypt.GenerateFromPassword([]byte(v.PasswordHash), bcrypt.DefaultCost)
	if err != nil {

		return false, 0
	}
	err = pc.Dbrepos.QueryRowContext(ctx, "INSERT INTO users (name, password, online) values ($1, $2, $3) returning user_id", v.Username, string(bytes), v.Online).Scan(&id)
	if err != nil {

		return false, 0
	}
	return true, id
}
func (pc *PersonRepository) CreateSession(v Users, ctx context.Context) (string, error) {
	sessionID := uuid.New().String()
	_, err := pc.Dbrepos.ExecContext(ctx,
		"INSERT INTO sessions (session_id, user_id) VALUES ($1, $2)",
		sessionID, v.User_id)
	if err != nil {
		return "", fmt.Errorf("не удалось создать сессию: %w", err)
	}

	return sessionID, nil
}
func (pc *PersonRepository) SearchPerson(v Users, ctx context.Context) (int, error) {
	var id int
	err := pc.Dbrepos.QueryRowContext(ctx, "select user_id from users where name = $1", v.Username).Scan(&id)
	_, err = pc.Dbrepos.ExecContext(ctx, "update users set online = true where user_id = $1", id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
func (pc *PersonRepository) GetSession(v Users, ctx context.Context) (*Session, error) {
	var id int
	err := pc.Dbrepos.QueryRowContext(ctx, "SELECT user_id from users where name = $1", v.Username).Scan(&id)
	session := &Session{}
	err = pc.Dbrepos.QueryRowContext(ctx, "SELECT session_id, user_id FROM sessions WHERE user_id = $1", id).Scan(
		&session.SessionID, &session.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return session, nil
}
func (pc *PersonRepository) Exit(id int, cook string, ctx context.Context) error {
	_, err := pc.Dbrepos.ExecContext(ctx, "UPDATE users SET online = false where user_id = $1", id)
	if err != nil {

		return err
	}
	_, err = pc.Dbrepos.ExecContext(ctx, "DELETE from sessions where user_id = $1 and session_id = $2", id, cook)
	if err != nil {

		return err
	}
	return nil
}
func (pc *PersonRepository) UpdateOnline(id int, flag bool, ctx context.Context) error {
	switch flag {
	case true:
		_, err := pc.Dbrepos.ExecContext(ctx, "UPDATE users SET online = true where user_id = $1", id)
		if err != nil {

			return err
		}
	case false:
		_, err := pc.Dbrepos.ExecContext(ctx, "UPDATE users SET online = false where user_id = $1", id)
		if err != nil {

			return err
		}
	}
	return nil
}
func (pc *PersonRepository) UpdateLastEntered(id int, currenttime time.Time, sess string, ctx context.Context) error {
	log.Println(currenttime)
	_, err := pc.Dbrepos.ExecContext(ctx, "UPDATE sessions SET lastentered = $1 where user_id = $2 and session_id = $3", currenttime, id, sess)
	return err
}
func (pc *PersonRepository) GetLastEntered(ctx context.Context) (map[int]time.Time, error) {
	var id int
	var data time.Time
	var datalast = make(map[int]time.Time)
	rows, err := pc.Dbrepos.QueryContext(ctx, "select user_id, max(lastentered) from sessions group by user_id")

	for rows.Next() {
		rows.Scan(&id, &data)
		datalast[id] = data
	}
	datalast[id] = data
	return datalast, err
}
func (pc *PersonRepository) GetUsers(ctx context.Context) ([]Users, error) {
	var newUser Users
	var datalast = make([]Users, 0)
	rows, err := pc.Dbrepos.QueryContext(ctx, "select user_id, name, online from users ")
	for rows.Next() {
		rows.Scan(&newUser.User_id, &newUser.Username, &newUser.Online)
		datalast = append(datalast, newUser)
	}
	return datalast, err
}
