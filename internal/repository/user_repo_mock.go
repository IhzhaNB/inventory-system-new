package repository

import (
	"context"

	"inventory-system/internal/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository adalah "Stuntman" untuk UserRepository asli kita
type MockUserRepository struct {
	mock.Mock
}

// 1. Tiruan untuk FindByEmail
func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) != nil {
		return args.Get(0).(*model.User), args.Error(1)
	}
	return nil, args.Error(1)
}

// 2. Tiruan untuk Create
func (m *MockUserRepository) Create(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// 3. Tiruan untuk Count (Penting: Return int64!)
func (m *MockUserRepository) Count(ctx context.Context, search string) (int64, error) {
	args := m.Called(ctx, search)
	// Kita cast jadi int64 biar Golang gak ngamuk
	return args.Get(0).(int64), args.Error(1)
}

// 4. Tiruan untuk FindAll (Penting: Return []*model.User)
func (m *MockUserRepository) FindAll(ctx context.Context, limit, offset int, search string) ([]*model.User, error) {
	args := m.Called(ctx, limit, offset, search)
	if args.Get(0) != nil {
		return args.Get(0).([]*model.User), args.Error(1)
	}
	return nil, args.Error(1)
}

// 5. Tiruan untuk FindByID
func (m *MockUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*model.User), args.Error(1)
	}
	return nil, args.Error(1)
}

// 6. Tiruan untuk Update
func (m *MockUserRepository) Update(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// 7. Tiruan untuk Delete
func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
