package database

import "time"

type Coupon struct {
	Id        string    `json:"id" bson:"_id"`
	Name      string    `json:"name"`
	Brand     string    `json:"brand"`
	Value     int       `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	Expiry    time.Time `json:"expiry"`
}
