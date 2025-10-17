package domain

const (
	DefaultBoardSize = 8
	MinBoardSize     = 4
	MaxBoardSize     = 100
)

// Board представляет шахматную доску
type Board struct {
	Size int
}

// BoardRepository определяет контракт для работы с досками
type BoardRepository interface {
	GenerateBoard(size int) *Board
}

// BoardService определяет бизнес-логику для работы с досками
type BoardService interface {
	CreateBoard(size int) *Board
	ValidateSize(size int) error
	GeneratePattern() string
}
