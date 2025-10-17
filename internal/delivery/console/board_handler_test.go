package console

import (
	"errors"
	"testing"

	"chessboard/internal/domain"
)

// MockBoardService для тестирования
type MockBoardService struct {
	validateError error
}

func (m *MockBoardService) CreateBoard(size int) *domain.Board {
	return &domain.Board{Size: size}
}

func (m *MockBoardService) ValidateSize(size int) error {
	return m.validateError
}

func (m *MockBoardService) GeneratePattern() string {
	return ""
}

func TestParseBoardSizeStrict(t *testing.T) {
	// Создаем мок сервиса без ошибок валидации
	mockService := &MockBoardService{}
	handler := NewBoardHandler(mockService)

	testCases := []struct {
		name        string
		input       string
		expected    int
		expectError bool
		errorMsg    string
	}{
		// Валидные случаи
		{
			name:        "валидный размер",
			input:       "8",
			expected:    8,
			expectError: false,
		},
		{
			name:        "минимальный размер",
			input:       "1",
			expected:    1,
			expectError: false,
		},
		{
			name:        "максимальный размер",
			input:       "100",
			expected:    100,
			expectError: false,
		},
		{
			name:        "число с пробелами",
			input:       "  8  ",
			expected:    8,
			expectError: false,
		},
		{
			name:        "двузначное число",
			input:       "25",
			expected:    25,
			expectError: false,
		},

		// Невалидные случаи - парсинг
		{
			name:        "отрицательное число",
			input:       "-5",
			expected:    0,
			expectError: true,
			errorMsg:    "отрицательные числа не поддерживаются",
		},
		{
			name:        "дробное число с точкой",
			input:       "5.5",
			expected:    0,
			expectError: true,
			errorMsg:    "дробные числа не поддерживаются",
		},
		{
			name:        "дробное число с запятой",
			input:       "3,14",
			expected:    0,
			expectError: true,
			errorMsg:    "дробные числа не поддерживаются",
		},
		{
			name:        "не число",
			input:       "abc",
			expected:    0,
			expectError: true,
			errorMsg:    "неверный формат числа",
		},
		{
			name:        "число с буквами",
			input:       "12a",
			expected:    0,
			expectError: true,
			errorMsg:    "неверный формат числа",
		},
		{
			name:        "пустая строка",
			input:       "",
			expected:    0,
			expectError: true,
			errorMsg:    "неверный формат числа",
		},
		{
			name:        "только пробелы",
			input:       "   ",
			expected:    0,
			expectError: true,
			errorMsg:    "неверный формат числа",
		},
		{
			name:        "специальные символы",
			input:       "!@#",
			expected:    0,
			expectError: true,
			errorMsg:    "неверный формат числа",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := handler.parseBoardSizeStrict(tc.input)

			if tc.expectError {
				if err == nil {
					t.Errorf("ожидалась ошибка для ввода '%s', но ошибки нет", tc.input)
				}
				if tc.errorMsg != "" && err != nil {
					if !contains(err.Error(), tc.errorMsg) {
						t.Errorf("ожидалась ошибка с текстом '%s', получено: '%s'", tc.errorMsg, err.Error())
					}
				}
			} else {
				if err != nil {
					t.Errorf("неожиданная ошибка для ввода '%s': %v", tc.input, err)
				}
				if result != tc.expected {
					t.Errorf("для ввода '%s' ожидался размер %d, получен %d",
						tc.input, tc.expected, result)
				}
			}
		})
	}
}

func TestParseBoardSizeStrict_WithValidationError(t *testing.T) {
	// Мок сервиса, который возвращает ошибку валидации
	mockService := &MockBoardService{
		validateError: errors.New("размер доски не может быть меньше 1"),
	}
	handler := NewBoardHandler(mockService)

	testCases := []struct {
		name     string
		input    string
		errorMsg string
	}{
		{
			name:     "валидация через сервис",
			input:    "0",
			errorMsg: "размер доски не может быть меньше 1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := handler.parseBoardSizeStrict(tc.input)

			if err == nil {
				t.Errorf("ожидалась ошибка валидации для ввода '%s', но ошибки нет", tc.input)
			}

			if result != 0 {
				t.Errorf("при ошибке валидации ожидался результат 0, получен %d", result)
			}

			if err != nil && !contains(err.Error(), tc.errorMsg) {
				t.Errorf("ожидалась ошибка с текстом '%s', получено: '%s'", tc.errorMsg, err.Error())
			}
		})
	}
}

func TestCreateAndDisplayBoard(t *testing.T) {
	mockService := &MockBoardService{}
	handler := NewBoardHandler(mockService)

	// Этот тест проверяет, что функция не паникует при разных размерах
	testCases := []struct {
		name string
		size int
	}{
		{"маленькая доска", 2},
		{"стандартная доска", 8},
		{"большая доска", 20},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Просто проверяем, что функция выполняется без паники
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("функция CreateAndDisplayBoard запаниковала с размером %d: %v", tc.size, r)
				}
			}()

			handler.CreateAndDisplayBoard(tc.size)
		})
	}
}

// Вспомогательная функция для проверки содержания подстроки
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && (contains(s[1:], substr) || contains(s[:len(s)-1], substr)))
}

// Тест для проверки создания хэндлера
func TestNewBoardHandler(t *testing.T) {
	mockService := &MockBoardService{}
	handler := NewBoardHandler(mockService)

	if handler == nil {
		t.Error("NewBoardHandler вернул nil")
	}

	if handler.boardService != mockService {
		t.Error("BoardService не был правильно установлен в хэндлере")
	}
}
