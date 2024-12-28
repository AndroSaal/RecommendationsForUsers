package repository

import myproto "github.com/AndroSaal/RecommendationsForUsers/app/services/analytics/internal/transport/kafka/pb"

type KeyValueDatabse interface {
	SetProductUpdate(product *myproto.ProductAction) error
	SetUserUpdate(user *myproto.UserUpdate) error
}


