package library

import (
	"encoding/json"
	"io"
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
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		return []byte{}, err
	}
	return content, nil
}
