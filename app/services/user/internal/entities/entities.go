package entities

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

const (
	emailMaxLenth = 64
	emailMinLenth = 10

	usernameMaxLenth = 20
	usernameMinLenth = 3

	passwordMaxLenth = 32
	passwordMinLenth = 8

	userDiscriptionMaxLenth = 1024
)

type UserId int

type Username string

type Email string

type Password string

type UserDiscription string

type UserInterest string

type UserInterests []string

type UserAge int

type ErrorResponse struct {
	Reason string `json:"reason"`
}

type UserInfo struct {
	UsrId        UserId          `json:"userId"`
	Usrname      Username        `json:"username"`
	Email        Email           `json:"email" binding:"required"`
	Password     Password        `json:"password" binding:"required"`
	UsrDesc      UserDiscription `json:"discription" binding:"required"`
	UserInterest UserInterests   `json:"interests" binding:"required"`
	UsrAge       UserAge         `json:"age" binding:"required"`
}

func (ui *UserId) ValidateUserId() error {

	if *ui < 0 {
		return errors.New("invalid user id: can`t be less 0")
	}

	return nil
}

func (u *Username) ValidateUsername() error {
	re := regexp.MustCompile(`a-zA-Z0-9_`)

	if !re.MatchString(string(*u)) {
		return errors.New("invalid username: does not math regexp")
	}
	if len(*u) > usernameMaxLenth {
		return errors.New("invalid username: too long, max length is " + strconv.Itoa(usernameMaxLenth))
	}

	if len(*u) < usernameMinLenth {
		return errors.New("invalid username: too short, min length is " + strconv.Itoa(usernameMinLenth))
	}

	return nil
}

func (e *Email) ValidateEmail() error {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-z]{2,}$`)

	if !re.MatchString(string(*e)) {
		return errors.New("invalid email: does not math regexp")
	}

	if len(*e) > emailMaxLenth {
		return fmt.Errorf("%s %s", "invalid email: too long, max length is",
			strconv.Itoa(emailMaxLenth))
	}

	if len(*e) < emailMinLenth {
		return fmt.Errorf("%s %s", "invalid email: too short, min length is",
			strconv.Itoa(emailMinLenth))
	}

	return nil
}

func (p *Password) ValidatePassword() error {
	re := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9!@#$%^&()*]+$`)

	if !re.MatchString(string(*p)) {
		return errors.New("invalid password: does not math regexp")
	}

	if len(*p) > passwordMaxLenth {
		return fmt.Errorf("%s %s", "invalid password: too long, max length is",
			strconv.Itoa(passwordMaxLenth))
	}

	if len(*p) < passwordMinLenth {
		return fmt.Errorf("%s %s", "invalid password: too short, min length is",
			strconv.Itoa(passwordMinLenth))
	}

	return nil
}

func (ud *UserDiscription) ValidateUserDiscription() error {

	if len(*ud) > userDiscriptionMaxLenth {
		return fmt.Errorf("%s %s", "invalid user discription: too long, max length is",
			strconv.Itoa(userDiscriptionMaxLenth))
	}
	return nil
}

func (ui *UserInterest) ValidateUserInterest() error {
	return nil
}

func (ui *UserInterests) ValidateUserInterests() error {
	return nil
}

func (a *UserAge) ValidateUserAge() error {
	return nil
}

func (inf *UserInfo) ValidateUserInfo() error {
	return nil
}
