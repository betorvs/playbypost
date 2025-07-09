package utils

import (
	"log/slog"
	"os"
)

func Save(value, file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer func() {
		err := f.Close()
		if err != nil {
			slog.Error("error closing file", "error", err)
		}
	}()
	_, err = f.WriteString(value)
	if err != nil {
		return err
	}
	return nil
}

func Read(file string) (string, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
