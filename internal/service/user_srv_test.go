package service

import (
	"context"
	"testing"

	"inventory-system/internal/dto/request"
	"inventory-system/internal/model"
	"inventory-system/internal/repository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestUpdateUser_AdminCannotUpdateSuperAdmin(t *testing.T) {
	// 1. SETUP: Siapkan Stuntman dan Logger dummy
	mockUserRepo := new(repository.MockUserRepository)
	logger := zap.NewNop() // Nop = No Operation (Logger yang gak nge-print apa-apa biar testnya bersih)

	// Inject stuntman ke dalam struktur Repository utama
	mockRepos := &repository.Repository{
		User: mockUserRepo,
	}

	// Bikin Service-nya menggunakan Mock Repository
	userService := NewUserService(mockRepos, logger)

	// 2. DATA DUMMY
	targetUserID := uuid.New()
	targetUserRole := model.RoleSuperAdmin // Targetnya adalah Bos Besar!

	requesterRole := string(model.RoleAdmin) // Yang nge-request adalah Admin biasa
	reqPayload := request.UpdateUserRequest{
		Name: "Hacker Jahat",
		Role: string(model.RoleAdmin),
	}

	// 3. ATUR SKENARIO MOCK ("Kalau Service manggil FindByID, kasih targetUserRole ini ya!")
	mockUserRepo.On("FindByID", mock.Anything, targetUserID).Return(&model.User{
		BaseModel: model.BaseModel{ // <-- Panggil nama Base struct-nya di sini
			ID: targetUserID,
		},
		Role: targetUserRole,
	}, nil)

	// Perhatikan: Kita TIDAK nge-mock Update(), karena logika kita HARUSNYA
	// menolak request sebelum nyentuh db.Update(). Kalau Update() terpanggil, berarti test GAGAL.

	// 4. EKSEKUSI: Panggil fungsi aslinya
	res, err := userService.UpdateUser(context.Background(), targetUserID, reqPayload, requesterRole)

	// 5. VALIDASI (Assert): Pastikan hasilnya sesuai harapan kita!
	assert.Error(t, err)                                                         // Harus ada error
	assert.Nil(t, res)                                                           // Respon harus kosong (nil)
	assert.Equal(t, "forbidden: admin cannot modify a super_admin", err.Error()) // Pesan error harus presisi

	// Pastikan Stuntman bekerja sesuai skenario (FindByID dipanggil 1x)
	mockUserRepo.AssertExpectations(t)
}
