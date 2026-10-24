package database

import (
	"reward-points/internal/models"
	"strings"
	"sync"
	"time"
)

type Store interface {
	CreateMember(m models.Member) (models.Member, error)
	GetMember(id int) (models.MemberInfo, error)
	CreateMemberReward(r models.Reward) (models.Reward, error)
	GetMemberReward(id int) ([]models.Reward, error)
}

type LocalStorage struct {
	members      map[int]models.MemberInfo
	rewards      map[int][]models.Reward
	emailIndex   map[string]int
	mu           sync.RWMutex
	lastMemberID int
	lastRewardID int
}

func NewLocalStorage() *LocalStorage {
	return &LocalStorage{
		members:    make(map[int]models.MemberInfo),
		rewards:    make(map[int][]models.Reward),
		emailIndex: make(map[string]int),
	}
}

func (conn *LocalStorage) CreateMember(m models.Member) (models.Member, error) {
	conn.mu.Lock()
	defer conn.mu.Unlock()

	emailID := strings.ToLower(strings.TrimSpace(m.Email))
	if _, exists := conn.emailIndex[emailID]; exists {
		return models.Member{}, ErrDuplicateEmail
	}

	m.MemberID = conn.lastMemberID + 1
	m.CreatedAt = time.Now().UTC()

	conn.lastMemberID = m.MemberID
	conn.members[m.MemberID] = models.MemberInfo{Member: m}
	conn.rewards[m.MemberID] = make([]models.Reward, 0)
	conn.emailIndex[emailID] = m.MemberID

	return m, nil
}

func (conn *LocalStorage) GetMember(id int) (models.MemberInfo, error) {
	conn.mu.RLock()
	defer conn.mu.RUnlock()

	m, ok := conn.members[id]
	if !ok {
		return models.MemberInfo{}, ErrMemberNotFound
	}

	return m, nil
}

func (conn *LocalStorage) CreateMemberReward(r models.Reward) (models.Reward, error) {
	conn.mu.Lock()
	defer conn.mu.Unlock()

	m, ok := conn.members[r.MemberID]
	if !ok {
		return models.Reward{}, ErrMemberNotFound
	}

	r.RewardID = conn.lastRewardID + 1
	r.EventDate = time.Now().UTC()

	if r.Points <= 0 {
		return models.Reward{}, ErrNegativePoints
	}
	switch r.PointTypeID {
	case models.PurchaseEarning, models.ReferralBonus, models.Cashback:
		m.PointBalance += r.Points
	case models.Redemption:
		if m.PointBalance < r.Points {
			return models.Reward{}, ErrInsufficientPoints
		}
		m.PointBalance -= r.Points
		r.Points = r.Points * -1
	default:
		return models.Reward{}, ErrInvalidPointType
	}
	conn.members[r.MemberID] = m
	conn.lastRewardID = r.RewardID
	conn.rewards[r.MemberID] = append(conn.rewards[r.MemberID], r)

	return r, nil
}

func (conn *LocalStorage) GetMemberReward(id int) ([]models.Reward, error) {
	conn.mu.RLock()
	defer conn.mu.RUnlock()

	rewards, ok := conn.rewards[id]
	if !ok {
		return nil, ErrMemberNotFound
	}

	return rewards, nil
}

var _ Store = (*LocalStorage)(nil)
