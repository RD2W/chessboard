package usecase

import (
	"strings"
	"testing"

	"chessboard/internal/domain"
)

// MockBoardRepository для тестирования
type MockBoardRepository struct {
	generateError error
}

func (m *MockBoardRepository) GenerateBoard(size int) *domain.Board {
	if m.generateError != nil {
		return nil
	}
	return &domain.Board{Size: size}
}

func TestNewBoardUsecase(t *testing.T) {
	repo := &MockBoardRepository{}
	usecase := NewBoardUsecase(repo)

	if usecase == nil {
		t.Error("NewBoardUsecase вернул nil")
	}
}

func TestBoardUsecase_ValidateSize(t *testing.T) {
	repo := &MockBoardRepository{}
	usecase := NewBoardUsecase(repo)

	testCases := []struct {
		name        string
		size        int
		expectError bool
		errorMsg    string
	}{
		{
			name:        "валидный минимальный размер",
			size:        domain.MinBoardSize,
			expectError: false,
		},
		{
			name:        "валидный средний размер",
			size:        8,
			expectError: false,
		},
		{
			name:        "валидный максимальный размер",
			size:        domain.MaxBoardSize,
			expectError: false,
		},
		{
			name:        "размер меньше минимального",
			size:        0,
			expectError: true,
			errorMsg:    "размер доски не может быть меньше",
		},
		{
			name:        "отрицательный размер",
			size:        -5,
			expectError: true,
			errorMsg:    "размер доски не может быть меньше",
		},
		{
			name:        "размер больше максимального",
			size:        domain.MaxBoardSize + 1,
			expectError: true,
			errorMsg:    "размер доски не может превышать",
		},
		{
			name:        "очень большой размер",
			size:        1000,
			expectError: true,
			errorMsg:    "размер доски не может превышать",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := usecase.ValidateSize(tc.size)

			if tc.expectError {
				if err == nil {
					t.Errorf("ожидалась ошибка для размера %d, но ошибки нет", tc.size)
				}
				if tc.errorMsg != "" && err != nil {
					if !strings.Contains(err.Error(), tc.errorMsg) {
						t.Errorf("ожидалась ошибка с текстом '%s', получено: '%s'", tc.errorMsg, err.Error())
					}
				}
			} else {
				if err != nil {
					t.Errorf("неожиданная ошибка для размера %d: %v", tc.size, err)
				}
			}
		})
	}
}

func TestBoardUsecase_CreateBoard(t *testing.T) {
	t.Run("создание доски с валидным размером", func(t *testing.T) {
		repo := &MockBoardRepository{}
		usecase := NewBoardUsecase(repo)

		board := usecase.CreateBoard(8)

		if board == nil {
			t.Error("CreateBoard вернул nil")
			return
		}
		if board.Size != 8 {
			t.Errorf("ожидался размер доски 8, получен %d", board.Size)
		}
	})

	t.Run("создание доски с невалидным размером - использует размер по умолчанию", func(t *testing.T) {
		repo := &MockBoardRepository{}
		usecase := NewBoardUsecase(repo)

		board := usecase.CreateBoard(0) // Невалидный размер

		if board == nil {
			t.Error("CreateBoard вернул nil")
			return
		}
		if board.Size != domain.DefaultBoardSize {
			t.Errorf("при невалидном размере ожидался размер по умолчанию %d, получен %d",
				domain.DefaultBoardSize, board.Size)
		}
	})

	t.Run("создание доски с слишком большим размером - использует размер по умолчанию", func(t *testing.T) {
		repo := &MockBoardRepository{}
		usecase := NewBoardUsecase(repo)

		board := usecase.CreateBoard(domain.MaxBoardSize + 10)

		if board == nil {
			t.Error("CreateBoard вернул nil")
			return
		}
		if board.Size != domain.DefaultBoardSize {
			t.Errorf("при невалидном размере ожидался размер по умолчанию %d, получен %d",
				domain.DefaultBoardSize, board.Size)
		}
	})
}

func TestGenerateChessboard(t *testing.T) {
	testCases := []struct {
		name     string
		board    *domain.Board
		expected string
		desc     string
	}{
		{
			name:     "доска 1x1",
			board:    &domain.Board{Size: 1},
			expected: " ",
			desc:     "одна белая клетка",
		},
		{
			name:     "доска 2x2",
			board:    &domain.Board{Size: 2},
			expected: " #\n# ",
			desc:     "правильное чередование 2x2",
		},
		{
			name:     "доска 3x3",
			board:    &domain.Board{Size: 3},
			expected: " # \n# #\n # ",
			desc:     "правильное чередование 3x3",
		},
		{
			name:     "доска 4x4",
			board:    &domain.Board{Size: 4},
			expected: " # #\n# # \n # #\n# # ",
			desc:     "правильное чередование 4x4",
		},
		{
			name:     "доска 8x8",
			board:    &domain.Board{Size: 8},
			expected: " # # # #\n# # # # \n # # # #\n# # # # \n # # # #\n# # # # \n # # # #\n# # # # ",
			desc:     "стандартная шахматная доска",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := GenerateChessboard(tc.board)

			if result != tc.expected {
				t.Errorf("для доски %dx%d:\nожидалось:\n%s\n\nполучено:\n%s\n",
					tc.board.Size, tc.board.Size, tc.expected, result)
			}

			// Дополнительная проверка: количество символов новой строки
			expectedNewlines := tc.board.Size - 1
			actualNewlines := strings.Count(result, "\n")
			if actualNewlines != expectedNewlines {
				t.Errorf("неправильное количество новых строк: ожидалось %d, получено %d",
					expectedNewlines, actualNewlines)
			}

			// Проверка общего количества символов (без учета \n)
			expectedChars := tc.board.Size * tc.board.Size
			actualChars := len(result) - actualNewlines
			if actualChars != expectedChars {
				t.Errorf("неправильное количество символов: ожидалось %d, получено %d",
					expectedChars, actualChars)
			}
		})
	}
}

func TestGenerateChessboard_EdgeCases(t *testing.T) {
	t.Run("нулевая доска", func(t *testing.T) {
		board := &domain.Board{Size: 0}
		result := GenerateChessboard(board)

		if result != "" {
			t.Errorf("для доски размером 0 ожидалась пустая строка, получено: '%s'", result)
		}
	})

	t.Run("очень большая доска", func(t *testing.T) {
		board := &domain.Board{Size: 100}
		result := GenerateChessboard(board)

		// Проверяем, что строка не пустая и имеет правильную структуру
		if result == "" {
			t.Error("для большой доски ожидалась непустая строка")
		}

		lines := strings.Split(result, "\n")
		if len(lines) != 100 {
			t.Errorf("ожидалось 100 строк, получено %d", len(lines))
		}

		for i, line := range lines {
			if len(line) != 100 {
				t.Errorf("строка %d: ожидалась длина 100, получено %d", i, len(line))
			}
		}
	})
}

func TestGenerateChessboard_Pattern(t *testing.T) {
	// Тестируем правильность шахматного паттерна
	testCases := []struct {
		size     int
		row      int
		col      int
		expected string
	}{
		{size: 8, row: 0, col: 0, expected: whiteSquare}, // Левый верхний угол - белый
		{size: 8, row: 0, col: 1, expected: blackSquare}, // Следующая клетка - черная
		{size: 8, row: 1, col: 0, expected: blackSquare}, // Вторая строка, первая клетка - черная
		{size: 8, row: 1, col: 1, expected: whiteSquare}, // Вторая строка, вторая клетка - белая
		{size: 2, row: 0, col: 0, expected: whiteSquare},
		{size: 2, row: 0, col: 1, expected: blackSquare},
		{size: 2, row: 1, col: 0, expected: blackSquare},
		{size: 2, row: 1, col: 1, expected: whiteSquare},
	}

	for _, tc := range testCases {
		t.Run("проверка паттерна", func(t *testing.T) {
			board := &domain.Board{Size: tc.size}
			result := GenerateChessboard(board)
			lines := strings.Split(result, "\n")

			if tc.row < len(lines) && tc.col < len(lines[tc.row]) {
				actual := string(lines[tc.row][tc.col])
				if actual != tc.expected {
					t.Errorf("для позиции (%d,%d) ожидался '%s', получен '%s'",
						tc.row, tc.col, tc.expected, actual)
				}
			}
		})
	}
}

func TestNewBoardRepository(t *testing.T) {
	repo := NewBoardRepository()

	if repo == nil {
		t.Error("NewBoardRepository вернул nil")
	}

	// Проверяем, что репозиторий работает
	board := repo.GenerateBoard(8)
	if board == nil {
		t.Error("GenerateBoard вернул nil")
		return
	}
	if board.Size != 8 {
		t.Errorf("ожидался размер доски 8, получен %d", board.Size)
	}
}

func TestBoardUsecase_GeneratePattern(t *testing.T) {
	repo := &MockBoardRepository{}
	usecase := NewBoardUsecase(repo)

	pattern := usecase.GeneratePattern()

	if pattern != "" {
		t.Errorf("ожидалась пустая строка, получено: '%s'", pattern)
	}
}

// Бенчмарк тесты
func BenchmarkGenerateChessboard(b *testing.B) {
	boards := []*domain.Board{
		{Size: 8},
		{Size: 16},
		{Size: 64},
		{Size: 100},
	}

	for _, board := range boards {
		b.Run(benchmarkName(board.Size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				GenerateChessboard(board)
			}
		})
	}
}

func BenchmarkValidateSize(b *testing.B) {
	repo := &MockBoardRepository{}
	usecase := NewBoardUsecase(repo)

	// Предварительная проверка что ошибок не будет
	if err := usecase.ValidateSize(8); err != nil {
		b.Skipf("skipping benchmark due to setup error: %v", err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = usecase.ValidateSize(8)
	}
}

// Вспомогательная функция для именования бенчмарков
func benchmarkName(size int) string {
	return string(rune(size)) + "x" + string(rune(size))
}
