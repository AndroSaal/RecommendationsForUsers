package repository

// const (
// 	all = "*"
// )

const (
	//таблица
	productsTable = "products"
	//её поля
	id               = "id" //PK
	categoryField    = "category"
	describtionField = "prd_description"
	statusField      = "prd_status"
)

const (
	//таблица
	productKwTable = "product_keyWord"
	//её поля
	//id = "id" PK
	productIdField = "product_id"
	kwIdField      = "kw_id"
)

const (
	//таблица
	kwTable = "keyWords"
	//её поля
	//id = "id" PK
	kwNameField = "kw_name"
)
