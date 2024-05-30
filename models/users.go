package models

import (
	"errors"
	"time"

	"github.com/upper/db/v4"
	"golang.org/x/crypto/bcrypt"
)

const cost = 12

type User struct {
	ID        int       `db:"id,omitempty"`
	Name      string    `db:"name"`
	Password  string    `db:"password_hash"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
	Activated bool      `db:"activated"`
}

type UserModel struct {
	db db.Session
}

func (m UserModel) Table() string {
	return "users"
}

func (m UserModel) GetById(id int) (*User, error) {

	var u User
	err := m.db.Collection("users").Find(db.Cond{"id": id}).One(&u)
	if err != nil {

		if errors.Is(err, db.ErrNoMoreRows) {
			return nil, ErrorNoMoreRows
		}

		return nil, err

	}

	return &u, err

}

func (m UserModel) FindByEmail(email string) (*User, error) {

	var u User
	err := m.db.Collection("users").Find(db.Cond{"email": email}).One(&u)
	if err != nil {

		if errors.Is(err, db.ErrNoMoreRows) {
			return nil, ErrorNoMoreRows
		}

		return nil, err

	}

	return &u, err

}

func (m UserModel) Insert(u *User) error {

	newHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), cost)

	if err != nil {
		return err
	}

	u.Password = string(newHash)
	u.CreatedAt = time.Now()
	col := m.db.Collection("users")
	res, err := col.Insert(u)
	if err != nil {
		switch {
		case errHasDuplicate(err, "user_email_key"):
			return ErrorDuplicateEmail
		default:
			return err
		}
	}

	u.ID = convertUpperIDToInt(res.ID())

	return nil
}

func (u *User) ComparePassword(plainPassword string) (bool, error) {

	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

func (m UserModel) Authenticate(email string, password string) (*User, error) {

	user, err := m.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	if !user.Activated {
		return nil, ErrorAccountNotActivated
	}

	isCorrectPassword, err := user.ComparePassword(password)
	if err != nil {
		return nil, err
	}

	if !isCorrectPassword {
		return nil, ErrorInvalidLogin
	}

	return user, nil

}
