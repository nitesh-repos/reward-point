package database

import (
	"errors"
	"reward-points/internal/models"
	"testing"
)

func TestGetMember(t *testing.T) {
	s := NewLocalStorage()
	m, err := s.CreateMember(models.Member{Name: "abc", Email: "abc@xyz.com"})
	if err != nil {
		t.Fatalf("initial setup error with create member %v", err)
	}

	t.Run("known member", func(t *testing.T) {
		_, err := s.GetMember(m.MemberID)
		if err != nil {
			t.Errorf("err = %v, wanted nil", err)
		}
	})

	t.Run("unknown member", func(t *testing.T) {
		_, err := s.GetMember(0)
		if !errors.Is(err, ErrMemberNotFound) {
			t.Errorf("err = %v, wanted %v", err, ErrMemberNotFound)
		}
	})
}

func TestCreateMemberReward(t *testing.T) {
	tests := []struct {
		name        string
		seedPoints  int // starting balance
		reward      models.Reward
		wantErr     error
		wantBalance int
	}{
		{"earn purchase", 0, models.Reward{PointTypeID: models.PurchaseEarning, Points: 100}, nil, 100},
		{"redeem ok", 100, models.Reward{PointTypeID: models.Redemption, Points: 40}, nil, 60},
		{"redeem too much", 50, models.Reward{PointTypeID: models.Redemption, Points: 80}, ErrInsufficientPoints, 50},
		{"zero points", 0, models.Reward{PointTypeID: models.Cashback, Points: 0}, ErrNegativePoints, 0},
		{"invalid type", 0, models.Reward{PointTypeID: 99, Points: 10}, ErrInvalidPointType, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewLocalStorage()
			m, err := s.CreateMember(models.Member{Name: "abc", Email: "abc@xyz.com"})
			if err != nil {
				t.Fatalf("initial setup error with create member %v", err)
			}

			// set initial point +ve balance
			if tt.seedPoints > 0 {
				s.CreateMemberReward(models.Reward{MemberID: m.MemberID, PointTypeID: models.PurchaseEarning, Points: tt.seedPoints})
			}

			tt.reward.MemberID = m.MemberID
			_, err = s.CreateMemberReward(tt.reward)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, wanted %v", err, tt.wantErr)
			}
			got, _ := s.GetMember(m.MemberID)
			if got.PointBalance != tt.wantBalance {
				t.Errorf("balance = %d, wanted %d", got.PointBalance, tt.wantBalance)
			}
		})
	}
}

func TestGetMemberReward(t *testing.T) {
	s := NewLocalStorage()
	m, err := s.CreateMember(models.Member{Name: "abc", Email: "abc@xyz.com"})
	if err != nil {
		t.Fatalf("initial setup error with create member %v", err)
	}

	t.Run("known member", func(t *testing.T) {
		_, err := s.GetMemberReward(m.MemberID)
		if err != nil {
			t.Errorf("err = %v, wanted nil", err)
		}
	})

	t.Run("known member 0 rewards", func(t *testing.T) {
		rewards, _ := s.GetMemberReward(m.MemberID)
		if len(rewards) > 0 {
			t.Errorf("found %d reward, wanted 0", len(rewards))
		}
	})

	t.Run("known member 0 rewards", func(t *testing.T) {
		rewards, _ := s.GetMemberReward(m.MemberID)
		if len(rewards) > 0 {
			t.Errorf("found %d reward, wanted 0", len(rewards))
		}
	})

	t.Run("known member 3 rewards", func(t *testing.T) {
		s.CreateMemberReward(models.Reward{MemberID: m.MemberID, PointTypeID: models.PurchaseEarning, Points: 100})
		s.CreateMemberReward(models.Reward{MemberID: m.MemberID, PointTypeID: models.ReferralBonus, Points: 0})
		s.CreateMemberReward(models.Reward{MemberID: m.MemberID, PointTypeID: models.Redemption, Points: 50})
		s.CreateMemberReward(models.Reward{MemberID: m.MemberID, PointTypeID: models.Cashback, Points: 25})

		rewards, _ := s.GetMemberReward(m.MemberID)
		if len(rewards) != 3 {
			t.Errorf("found %d reward, wanted 3", len(rewards))
		}
	})

}
