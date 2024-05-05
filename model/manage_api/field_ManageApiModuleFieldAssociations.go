package manage_api

import (
	"database/sql/driver"
	"encoding/json"
)

type ManageApiModuleFieldAssociations struct {
	//Belongs To | Has One | Has Many | Many To Many
	Type         string `json:"type"`
	Table        string `json:"table"`
	Field        string `json:"field"`
	CurrentField string `json:"current_field"`
	Module       string `json:"module"`

	Form ManageApiModuleFieldAssociationsFrom `json:"form"`
}

type ManageApiModuleFieldAssociationsFrom struct {
	Width int `json:"width"`
}

func (j *ManageApiModuleFieldAssociations) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	if len(bytes) == 0 {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

func (j ManageApiModuleFieldAssociations) Value() (driver.Value, error) {
	marshal, err := json.Marshal(j)
	if err != nil {
		return nil, err
	}
	return marshal, nil
}
