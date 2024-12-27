package repository

type KeyValueDatabse interface {
	GetRecom(userId int) ([]int, error)
	SetRecom(userId int, productIds []int) error
}
