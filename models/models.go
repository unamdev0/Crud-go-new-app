package models

import (
	"errors"
	"fmt"
	"strings"

	"github.com/upper/db/v4"
)

var (
	ErrorNoMoreRows          = errors.New("No Record Found")
	ErrorDuplicateEmail      = errors.New("Email id already exists")
	ErrorAccountNotActivated = errors.New("Account isn't activated")
	ErrorInvalidLogin        = errors.New("Invalid Login attempt")
	ErrDuplicateTitle        = errors.New("title already exist in database")
	ErrDuplicateVotes        = errors.New("you already voted")
)

type Models struct {
	Users    UserModel
	Posts    PostModel
	Comments CommentsModel
}

func New(db db.Session) Models {
	return Models{
		Users:    UserModel{db: db},
		Posts:    PostModel{db: db},
		Comments: CommentsModel{db: db},
	}
}

func errHasDuplicate(err error, key string) bool {

	str := fmt.Sprintf(`ERROR: duplicate key value violates unique constraint "%s"`, key)
	return strings.Contains(err.Error(), str)
}

func convertUpperIDToInt(id db.ID) int {
	idType := fmt.Sprintf(`%T`, id)
	if idType == "int64" {
		return int(id.(int64))
	}

	return id.(int)

}
