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
	userTimestampsTable = "user_timestamps"
	//её поля
	//idField = "id" PK
	userIdField    = "user_id"
	timestampField = "timestamp"
)

const (
	//таблица
	userUpdatesTable = "user_updates"
	//её поля
	//timestampField = "timestamp" PK
	userInterestsField = "interests"
	//kwIdField      = "kw_id"
)

const (
	//таблица
	productsTable = "products"
	//её поля
	//idField = "id" PK
)

const (
	//таблица
	productsTimestampsTable = "products_timestamps"
	//её поля
	//idField = "id" PK
	productIdField = "product_id"
	// timestampField = "timestamp"

)

const (
	//таблица
	productUpsetesTable = "product_updates"
	//её поля
	//timestampField = "timestamp" PK
	kwField = "keyWords"
	//kwIdField      = "kw_id"
)
