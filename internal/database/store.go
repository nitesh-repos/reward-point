package database

import (
	"reward-points/internal/models"
	"sync"
)

type Store interface {
	CreateMember(m models.Member) (models.Member, error)
	GetMember(id int) (models.MemberInfo, error)
	CreateMemberReward(r models.Reward) (models.Reward, error)
	GetMemberReward(id int) ([]models.Reward, error)
}

type LocalStorage struct {
	members      map[string]models.MemberInfo
	rewards      map[string][]models.Reward
	mu           sync.RWMutex
	lastMemberID int
	lastRewardID int
}

func NewLocalStorage() *LocalStorage {
	return &LocalStorage{
		members: make(map[string]models.MemberInfo),
		rewards: make(map[string][]models.Reward),
	}
}

func (conn *LocalStorage) CreateMember(m models.Member) (models.Member, error) {
	return models.Member{}, nil
}

func (conn *LocalStorage) GetMember(id int) (models.MemberInfo, error) {
	return models.MemberInfo{}, nil
}

func (conn *LocalStorage) CreateMemberReward(r models.Reward) (models.Reward, error) {
	return models.Reward{}, nil
}

func (conn *LocalStorage) GetMemberReward(id int) ([]models.Reward, error) {
	return nil, nil
}

var _ Store = (*LocalStorage)(nil)
