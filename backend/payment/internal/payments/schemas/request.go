package schemas

import "time"

type PaymentRequest struct {
	Id 					string `json:"id" db:"id"`
	UserId 				string `json:"userId" db:"user_id"`
	SubscriptionType 	int `json:"subscriptionType" db:"subscription_type"`
	ActiveTillTo 		time.Duration `json:"activeTillTo" db:"active_till_to"`
	Status 				int `json:"status" db:"status"`
}