package app

import (
	"github.com/ArtemRotov/computer-club-administrator/internal/controller"
	"github.com/ArtemRotov/computer-club-administrator/internal/service"
	"os"
)

func Run(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// services
	manager := service.NewManagerService()

	// controllers
	fileController := controller.NewTextFileHandler(manager)

	// handle
	return fileController.Handle(file)
}
