package main

import (
	"chessboard/internal/delivery/console"
	"chessboard/internal/usecase"
	"fmt"
	"os"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Printf("chessboard version %s, commit %s, built %s\n", version, commit, date)
		return
	}

	repo := usecase.NewBoardRepository()
	service := usecase.NewBoardUsecase(repo)
	handler := console.NewBoardHandler(service)

	// Обработка пользовательского ввода и отображение доски
	handler.HandleUserInput()
}
