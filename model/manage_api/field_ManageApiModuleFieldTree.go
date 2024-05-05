package manage_api

import (
	"database/sql/driver"
	"encoding/json"
)

type ManageApiModuleFieldTree struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	PID   string `json:"pid"`
}

func (j *ManageApiModuleFieldTree) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	if len(bytes) == 0 {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

func (j ManageApiModuleFieldTree) Value() (driver.Value, error) {
	marshal, err := json.Marshal(j)
	if err != nil {
		return nil, err
	}
	return marshal, nil
}
