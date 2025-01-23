package repository

const (
	all = "*"
)

const (
	//таблица
	usersTable = "users"
	//её поля
	id                  = "id" //PK
	emailPole           = "email"
	usernamePole        = "username"
	passwordPole        = "password_hash"
	describtionPole     = "usr_description"
	agePole             = "age"
	isEmailVerifiedPole = "is_email_verified"
)

const (
	//таблица
	codesTable = "codes"
	//её поля
	userIdPole = "user_id"
	codePole   = "email_code"
)

const (
	//таблица
	userInterestsTable = "user_interests"
	//её поля
	interestIdPole = "interest_id"
)

const (
	//таблица
	interestsTable = "interests"
	//её поля
	intersestPole = "interest"
)

type UserInfoForDB struct {
	UsrId        int    `db:"id"`
	Usrname      string `db:"username"`
	Email        string `db:"email"`
	Password     string `db:"password_hash"`
	UsrDesc      string `db:"usr_description"`
	IsEmailValid bool   `db:"is_email_verified"`
	UsrAge       int    `db:"age"`
}
