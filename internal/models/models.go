package models

import "time"

type PointTypeID int

const (
	PurchaseEarning PointTypeID = iota + 1
	ReferralBonus
	Cashback
	Redemption
)

type Member struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	MemberID  int       `json:"member_id"`
	CreatedAt time.Time `json:"created_at"`
}

type MemberInfo struct {
	Member
	PointBalance int `json:"point_balance"`
}

type Reward struct {
	RewardID    int         `json:"reward_id"`
	MemberID    int         `json:"member_id"`
	PointTypeID PointTypeID `json:"point_type_id"`
	Points      int         `json:"points"`
	Description string      `json:"description"`
	EventDate   time.Time   `json:"event_date"`
}
