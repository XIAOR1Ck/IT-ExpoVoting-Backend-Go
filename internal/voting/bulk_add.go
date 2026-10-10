package voting

import (
	"os"
	"encoding/json"
)

func LoadProjects(path string) ([]Project, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var projects []Project
	if err := json.Unmarshal(data, &projects); err != nil {
		return nil, err
	}

	return projects, nil
}
