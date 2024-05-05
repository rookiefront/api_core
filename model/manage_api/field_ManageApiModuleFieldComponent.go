package manage_api

import (
	"database/sql/driver"
	"encoding/json"
)

type ManageApiModuleFieldComponent struct {
	Name string `json:"name"`
	Bind string `json:"bind"`
}

func (j *ManageApiModuleFieldComponent) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	if len(bytes) == 0 {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

func (j ManageApiModuleFieldComponent) Value() (driver.Value, error) {
	marshal, err := json.Marshal(j)
	if err != nil {
		return nil, err
	}
	return marshal, nil
}
