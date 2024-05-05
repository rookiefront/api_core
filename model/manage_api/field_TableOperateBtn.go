package manage_api

import (
	"database/sql/driver"
	"encoding/json"
)

type ManageApiModuleFieldTableOperateBtn struct {
	Name       string `json:"name"`
	Permission string `json:"permission"`
	Location   string `json:"location"`
	Component  string `json:"component"`
	Api        string `json:"api"`
}

func (j *ManageApiModuleFieldTableOperateBtn) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	if len(bytes) == 0 {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

func (j ManageApiModuleFieldTableOperateBtn) Value() (driver.Value, error) {
	marshal, err := json.Marshal(j)
	if err != nil {
		return nil, err
	}
	return marshal, nil
}
