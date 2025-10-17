package main

import (
	"chessboard/internal/delivery/console"
	"chessboard/internal/usecase"
)

func main() {
	repo := usecase.NewBoardRepository()
	service := usecase.NewBoardUsecase(repo)
	handler := console.NewBoardHandler(service)

	// Обработка пользовательского ввода и отображение доски
	handler.HandleUserInput()
}
