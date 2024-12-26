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
	describtionField = "describtion"
	statusField      = "status"
)

const (
	//таблица
	productTagsTable = "product_tags"
	//её поля
	//id = "id" PK
	productIdField = "product_id"
	tagIdField     = "tag_id"
)

const (
	//таблица
	tagsTable = "tags"
	//её поля
	//id = "id" PK
	tagNameField = "tag_name"
)
