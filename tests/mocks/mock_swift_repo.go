package mocks

import (
	"github.com/mroczekDNF/swift-api/internal/models"
	"github.com/mroczekDNF/swift-api/internal/repositories"
	"github.com/stretchr/testify/mock"
)

// MockSwiftCodeRepository to mock repozytorium dla testów.
type MockSwiftCodeRepository struct {
	mock.Mock
}

// GetBySwiftCode mockuje metodę repozytorium.
func (m *MockSwiftCodeRepository) GetBySwiftCode(code string) (*models.SwiftCode, error) {
	args := m.Called(code)
	if args.Get(0) != nil {
		return args.Get(0).(*models.SwiftCode), args.Error(1)
	}
	return nil, args.Error(1)
}

// GetBranchesByHeadquarter mockuje metodę repozytorium.
func (m *MockSwiftCodeRepository) GetBranchesByHeadquarter(code string) ([]models.SwiftCode, error) {
	args := m.Called(code)
	return args.Get(0).([]models.SwiftCode), args.Error(1)
}

// DeleteSwiftCode mockuje metodę repozytorium.
func (m *MockSwiftCodeRepository) DeleteSwiftCode(code string) error {
	args := m.Called(code)
	return args.Error(0)
}

// DetachBranchesFromHeadquarter mockuje metodę repozytorium.
func (m *MockSwiftCodeRepository) DetachBranchesFromHeadquarter(headquarterID int64) error {
	args := m.Called(headquarterID)
	return args.Error(0)
}

// ✅ **TUTAJ DODAJ NOWĄ METODĘ**
// GetByCountryISO2 mockuje metodę repozytorium.
func (m *MockSwiftCodeRepository) GetByCountryISO2(countryISO2 string) ([]models.SwiftCode, error) {
	args := m.Called(countryISO2)
	if args.Get(0) != nil {
		return args.Get(0).([]models.SwiftCode), args.Error(1)
	}
	return nil, args.Error(1)
}

// InsertSwiftCode mockuje metodę repozytorium.
func (m *MockSwiftCodeRepository) InsertSwiftCode(swift *models.SwiftCode) error {
	args := m.Called(swift)
	return args.Error(0)
}

// ✅ **Wymuszenie implementacji interfejsu**
var _ repositories.SwiftCodeRepositoryInterface = (*MockSwiftCodeRepository)(nil)
