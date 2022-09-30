package models

import "time"

type Payment struct {
	Id 					string `json:"id" db:"id"`
	UserId 				string `json:"userId" db:"user_id"`
	SubscriptionType 	int `json:"subscriptionType" db:"subscription_type"`
	ActiveTillTo 		time.Time `json:"activeTillTo" db:"active_till_to"`
	Status 				string `json:"status" db:"status"`
}


