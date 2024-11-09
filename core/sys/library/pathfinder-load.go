package library

import "encoding/json"

func loadAncestriesPF(f string) ([]PFAncestry, error) {
	def := []PFAncestry{}

	// load from file
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

func loadBackgroundsPF(f string) ([]PFBackground, error) {
	def := []PFBackground{}

	// load from file
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

func loadClassesPF(f string) ([]PFClass, error) {
	def := []PFClass{}

	// load from file
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
