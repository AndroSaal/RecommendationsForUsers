package repository

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
	userIdField        = "user_id"
	timestampField     = "timestamp_column"
	userInterestsField = "interests"
)

const (
	//таблица
	productsTable = "products"
)

const (
	//таблица
	productUpdatesTable = "product_updates"
	//её поля
	productIdField = "product_id"
	kwField        = "keywords"
)
