package library

import (
	"encoding/json"
	"io"
	"log/slog"
	"os"
)

func loadDescriptionsFromFile(f string) ([]Descriptions, error) {
	var def []Descriptions

	content, err := loadFile(f)
	if err != nil {
		return def, err
	}

	err = json.Unmarshal(content, &def)
	if err != nil {
		return def, err
	}
	return def, nil
}

func loadFile(file string) ([]byte, error) {
	f, err := os.Open(file)
	if err != nil {
		return []byte{}, err
	}
	defer func() {
		err := f.Close()
		if err != nil {
			slog.Error("error closing file", "error", err)
		}
	}()

	content, err := io.ReadAll(f)
	if err != nil {
		return []byte{}, err
	}
	return content, nil
}
