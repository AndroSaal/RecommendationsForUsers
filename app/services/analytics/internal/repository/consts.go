package repository

// const (
// 	all = "*"
// )

const (
	//таблица
	usersTable = "users"
	//её поля
	idField = "id" //PK
)

const (
	//таблица
	userUpdatesTable = "user_updates"
	//её поля
	//id PK
	userIdField        = "user_id"
	timestampField     = "timestamp_column"
	userInterestsField = "interests"
)

const (
	//таблица
	productsTable = "products"
	//её поля
	//idField = "id" PK
)

const (
	//таблица
	productUpdatesTable = "product_updates"
	//её поля
	//idField = "id" PK
	productIdField = "product_id"
	//timestampField = "timestamp_column"
	kwField = "keywords"
)
