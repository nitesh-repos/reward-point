package server

import "reward-points/internal/models"

type MemberRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type RewardRequest struct {
	MemberID    int                `json:"member_id"`
	PointTypeID models.PointTypeID `json:"point_type_id"`
	Points      int                `json:"points"`
	Description string             `json:"description"`
}
