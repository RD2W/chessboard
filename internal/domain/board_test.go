package domain

import (
	"testing"
)

func TestConstants(t *testing.T) {
	testCases := []struct {
		name     string
		actual   int
		expected int
	}{
		{"DefaultBoardSize", DefaultBoardSize, 8},
		{"MinBoardSize", MinBoardSize, 4},
		{"MaxBoardSize", MaxBoardSize, 100},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.actual != tc.expected {
				t.Errorf("%s: ожидалось %d, получено %d",
					tc.name, tc.expected, tc.actual)
			}
		})
	}
}

func TestBoard_Initialization(t *testing.T) {
	t.Run("создание доски с положительным размером", func(t *testing.T) {
		board := &Board{Size: 10}

		if board.Size != 10 {
			t.Errorf("ожидался размер 10, получен %d", board.Size)
		}
	})

	t.Run("создание доски с минимальным размером", func(t *testing.T) {
		board := &Board{Size: MinBoardSize}

		if board.Size != MinBoardSize {
			t.Errorf("ожидался размер %d, получен %d", MinBoardSize, board.Size)
		}
	})

	t.Run("создание доски с максимальным размером", func(t *testing.T) {
		board := &Board{Size: MaxBoardSize}

		if board.Size != MaxBoardSize {
			t.Errorf("ожидался размер %d, получен %d", MaxBoardSize, board.Size)
		}
	})
}

func TestConstants_Validation(t *testing.T) {
	// Проверяем, что константы имеют логичные значения
	if MinBoardSize >= MaxBoardSize {
		t.Errorf("MinBoardSize (%d) должен быть меньше MaxBoardSize (%d)",
			MinBoardSize, MaxBoardSize)
	}

	if DefaultBoardSize < MinBoardSize || DefaultBoardSize > MaxBoardSize {
		t.Errorf("DefaultBoardSize (%d) должен быть между MinBoardSize (%d) и MaxBoardSize (%d)",
			DefaultBoardSize, MinBoardSize, MaxBoardSize)
	}

	if MinBoardSize <= 0 {
		t.Errorf("MinBoardSize (%d) должен быть положительным числом", MinBoardSize)
	}
}

// Тесты на соответствие интерфейсам (compile-time проверки)
func TestInterfaceImplementation(t *testing.T) {
	// Эти тесты проверяют, что интерфейсы правильно объявлены
	// Компилятор сам проверит реализацию

	var _ BoardRepository = (*mockRepository)(nil)
	var _ BoardService = (*mockService)(nil)
}

// Mock реализации для проверки интерфейсов
type mockRepository struct{}

func (m *mockRepository) GenerateBoard(size int) *Board {
	return &Board{Size: size}
}

type mockService struct{}

func (m *mockService) CreateBoard(size int) *Board {
	return &Board{Size: size}
}

func (m *mockService) ValidateSize(size int) error {
	return nil
}

func (m *mockService) GeneratePattern() string {
	return ""
}
