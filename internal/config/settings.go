package config

import (
	"fmt"
	"log/slog"
	"os"
)

type ConfigOpt func(*Settings, *slog.Logger) error

// configEnv defines the logic of ENV vars' configurations.
func configEnv(key string) (string, error) {
	const op = "conf.config-env"

	val := os.Getenv(key)

	if len(val) == 0 {
		return "", fmt.Errorf("error of the %s: %w: the %s var is empty or didn't set: check it and try again",
			op, ErrEnvVar, key)
	}

	return val, nil
}

// ConfigSocket defines the logic of storage's socket configuration.
func ConfigSocket(st *Settings, log *slog.Logger) error {
	socket, err := configEnv("SOCKET")

	if err != nil {
		log.Error(err.Error())
		return err
	}
	st.Socket = socket

	return nil
}
