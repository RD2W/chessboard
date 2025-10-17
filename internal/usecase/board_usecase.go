package usecase

import (
	"fmt"
	"strings"

	"chessboard/internal/domain"
)

const (
	whiteSquare = " "
	blackSquare = "#"
)

type boardUsecase struct {
	repo domain.BoardRepository
}

func NewBoardUsecase(repo domain.BoardRepository) domain.BoardService {
	return &boardUsecase{repo: repo}
}

func (uc *boardUsecase) CreateBoard(size int) *domain.Board {
	if err := uc.ValidateSize(size); err != nil {
		// Возвращаем доску размером по умолчанию при ошибке валидации
		return uc.repo.GenerateBoard(domain.DefaultBoardSize)
	}
	return uc.repo.GenerateBoard(size)
}

func (uc *boardUsecase) ValidateSize(size int) error {
	if size < domain.MinBoardSize {
		return fmt.Errorf("размер доски не может быть меньше %d", domain.MinBoardSize)
	}
	if size > domain.MaxBoardSize {
		return fmt.Errorf("размер доски не может превышать %d", domain.MaxBoardSize)
	}
	return nil
}

func (uc *boardUsecase) GeneratePattern() string {
	return ""
}

// BoardRepository реализация
type boardRepository struct{}

func NewBoardRepository() domain.BoardRepository {
	return &boardRepository{}
}

func (r *boardRepository) GenerateBoard(size int) *domain.Board {
	return &domain.Board{Size: size}
}

// GenerateChessboard генерирует строку с шахматной доской
func GenerateChessboard(board *domain.Board) string {
	var result strings.Builder

	for i := 0; i < board.Size; i++ {
		for j := 0; j < board.Size; j++ {
			// Логика чередования: сумма координат четная - пробел, нечетная - #
			if (i+j)%2 == 0 {
				result.WriteString(whiteSquare)
			} else {
				result.WriteString(blackSquare)
			}
		}
		// Добавляем символ новой строки после каждой строки, кроме последней
		if i < board.Size-1 {
			result.WriteString("\n")
		}
	}

	return result.String()
}
