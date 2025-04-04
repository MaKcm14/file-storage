package main

import "github.com/MaKcm14/file-storage/internal/app"

func main() {
	storage := app.NewService()
	storage.Run()
}
