package entities

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

type UserDiscription string

type UserInterest string

type UserInterests []UserInterest

type UserAge int

type ErrorResponse struct {
	Reason string `json:"reason"`
}

const (
	emailMaxLenth = 64
	emailMinLenth = 10

	usernameMaxLenth = 20
	usernameMinLenth = 3

	passwordMaxLenth = 32
	passwordMinLenth = 8

	userDiscriptionMaxLenth = 1024

	userInterestMaxLenth = 32
	userInterestMinLenth = 3

	maxUserAge = 150
	minUserAge = 5
)

type UserInfo struct {
	UsrId           int             `json:"userId" db:"id"`
	Usrname         string          `json:"username" binding:"required" db:"username"`
	Email           string          `json:"email" binding:"required"`
	Password        string          `json:"password" binding:"required"`
	UsrDesc         UserDiscription `json:"description" binding:"required"`
	UserInterests   UserInterests   `json:"interests" binding:"required"`
	UsrAge          UserAge         `json:"age" binding:"required"`
	IsEmailVerified bool            `json:"isVerified"`
}

func ValidateUserId(usrId int) error {

	if usrId <= 0 {
		return errors.New("invalid user id: can`t be less or equel 0")
	}

	return nil
}

func ValidateUsername(username string) error {
	re := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

	if len(username) > usernameMaxLenth {
		return errors.New("invalid username: too long, max length is " + strconv.Itoa(usernameMaxLenth))
	}

	if len(username) < usernameMinLenth {
		return errors.New("invalid username: too short, min length is " + strconv.Itoa(usernameMinLenth))
	}

	if !re.MatchString(username) {
		return fmt.Errorf("invalid username: %s does not match regexp", username)
	}

	return nil
}

func ValidateEmail(email string) error {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-z]{2,}$`)

	if len(email) > emailMaxLenth {
		return fmt.Errorf("%s %s", "invalid email: too long, max length is",
			strconv.Itoa(emailMaxLenth))
	}

	if len(email) < emailMinLenth {
		return fmt.Errorf("%s %s", "invalid email: too short, min length is",
			strconv.Itoa(emailMinLenth))
	}

	if !re.MatchString(email) {
		return errors.New("invalid email: does not math regexp")
	}

	return nil
}

func ValidatePassword(password string) error {
	re := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9!@#$%^&()*]+$`)

	if len(password) > passwordMaxLenth {
		return fmt.Errorf("%s %s", "invalid password: too long, max length is",
			strconv.Itoa(passwordMaxLenth))
	}

	if len(password) < passwordMinLenth {
		return fmt.Errorf("%s %s", "invalid password: too short, min length is",
			strconv.Itoa(passwordMinLenth))
	}

	if !re.MatchString(string(password)) {
		return errors.New("invalid password: does not math regexp")
	}

	return nil
}

func ValidateCode(code string) error {
	re := regexp.MustCompile(`^[0-9]{5}$`)

	if !re.MatchString(code) {
		return errors.New("invalid code: does not match regexp")
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

	if len(*ui) > userInterestMaxLenth {
		return fmt.Errorf("%s %s", "invalid user intersest: too long, max length is",
			strconv.Itoa(userInterestMaxLenth))
	}

	if len(*ui) < userInterestMinLenth {
		return fmt.Errorf("%s %s", "invalid user intersest: too short, min length is",
			strconv.Itoa(userInterestMinLenth))

	}
	return nil
}

func (ui *UserInterests) ValidateUserInterests() error {

	if len(*ui) == 0 {
		return errors.New("invalid user intersest: can`t be empty")
	}

	for elem, userInterest := range *ui {
		if err := userInterest.ValidateUserInterest(); err != nil {
			return fmt.Errorf("user interests[%d]: %s", elem, err.Error())
		}
	}

	return nil
}

func (a *UserAge) ValidateUserAge() error {

	if *a > maxUserAge || *a < minUserAge {
		return fmt.Errorf("%s %s and %s", "invalid user age: must be between",
			strconv.Itoa(minUserAge), strconv.Itoa(maxUserAge))
	}

	return nil

}

func (inf *UserInfo) ValidateUserInfo() error {

	if err := ValidateUsername(inf.Usrname); err != nil {
		return err
	}

	if err := ValidateEmail(inf.Email); err != nil {
		return err
	}

	if err := ValidatePassword(inf.Password); err != nil {
		return err
	}

	if err := inf.UsrDesc.ValidateUserDiscription(); err != nil {
		return err
	}

	if err := inf.UserInterests.ValidateUserInterests(); err != nil {
		return err
	}

	if err := inf.UsrAge.ValidateUserAge(); err != nil {
		return err
	}

	return nil
}
