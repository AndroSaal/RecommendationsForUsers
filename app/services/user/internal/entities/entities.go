package entities

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
	Usrname      Username        `json:"username" `
	Email        Email           `json:"email" binding:"required"`
	Password     Password        `json:"password" binding:"required"`
	UsrDesc      UserDiscription `json:"discription" binding:"required"`
	UserInterest UserInterests   `json:"interests" binding:"required"`
	UsrAge       UserAge         `json:"age" binding:"required"`
}

func (ui *UserId) ValidateUserId() error {
	return nil
}

func (u *Username) ValidateUsername() error {
	return nil
}

func (e *Email) ValidateEmail() error {
	return nil
}

func (p *Password) ValidatePassword() error {
	return nil
}

func (ud *UserDiscription) ValidateUserDiscription() error {
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
