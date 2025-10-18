package console

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"chessboard/internal/domain"
	"chessboard/internal/usecase"
)

type BoardHandler struct {
	boardService domain.BoardService
}

func NewBoardHandler(service domain.BoardService) *BoardHandler {
	return &BoardHandler{boardService: service}
}

func (h *BoardHandler) CreateAndDisplayBoard(size int) {
	board := h.boardService.CreateBoard(size)
	chessboard := usecase.GenerateChessboard(board)

	fmt.Printf("Шахматная доска %dx%d:\n", board.Size, board.Size)
	fmt.Println(chessboard)
}

func (h *BoardHandler) HandleUserInput() {
	if len(os.Args) > 1 {
		input := os.Args[1]
		size, err := h.parseBoardSizeStrict(input)
		if err != nil {
			fmt.Printf("Ошибка: %s. Используется размер по умолчанию %dx%d.\n",
				err.Error(), domain.DefaultBoardSize, domain.DefaultBoardSize)
			size = domain.DefaultBoardSize
		}
		h.CreateAndDisplayBoard(size)
	} else {
		fmt.Printf("Используется размер по умолчанию %dx%d\n",
			domain.DefaultBoardSize, domain.DefaultBoardSize)
		h.CreateAndDisplayBoard(domain.DefaultBoardSize)
	}
}

// parseBoardSizeStrict парсит и валидирует размер доски со строгой проверкой
func (h *BoardHandler) parseBoardSizeStrict(input string) (int, error) {
	input = strings.TrimSpace(input)

	if strings.HasPrefix(input, "-") {
		return 0, fmt.Errorf("отрицательные числа не поддерживаются: '%s'", input)
	}

	if strings.Contains(input, ".") || strings.Contains(input, ",") {
		return 0, fmt.Errorf("дробные числа не поддерживаются: '%s'", input)
	}

	size, err := strconv.Atoi(input)
	if err != nil {
		return 0, fmt.Errorf("неверный формат числа: '%s'", input)
	}

	// Валидируем размер через сервис
	if err := h.boardService.ValidateSize(size); err != nil {
		return 0, err
	}

	return size, nil
}
