package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExchangeRate struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id"`
	Crypto    string             `json:"crypto,omitempty" bson:"crypto"`
	Fiat      string             `json:"fiat,omitempty" bson:"fiat"`
	Rate      float64            `json:"rate,omitempty" bson:"rate"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}
