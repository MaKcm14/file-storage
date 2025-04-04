package app

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/MaKcm14/file-storage/internal/config"
)

// Service defines the main storage's server services.
type Service struct {
	logFile *os.File
	log     *slog.Logger
}

func NewService() Service {
	date := strings.Split(time.Now().String()[:19], " ")

	mainLogFile, err := os.Create(fmt.Sprintf("../../logs/price-service-main-logs_%s___%s.txt",
		date[0], strings.Join(strings.Split(date[1], ":"), "-")))

	if err != nil {
		panic(fmt.Sprintf("error of creating the main-log-file: %v", err))
	}

	log := slog.New(slog.NewTextHandler(mainLogFile, &slog.HandlerOptions{Level: slog.LevelInfo}))

	log.Info("main application's configuring begun")

	st, err := config.NewSettings(log)

	if err != nil {
		panic(err)
	}
	_ = st

	return Service{
		logFile: mainLogFile,
	}
}

// Run starts the storage's server.
func (s *Service) Run() {
	defer s.close()
	defer s.log.Info("the app was FULLY STOPPED")
}

func (s *Service) close() {
	s.logFile.Close()
}
