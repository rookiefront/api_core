package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"time"
)

type PrimarykeyType uint

var hash map[string]any
var dicts map[string][]SysDictItem

type Model struct {
	ID           PrimarykeyType `json:"id" gorm:"primarykey;int"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	CreateUserID int            `json:"-" gorm:"column:c_id;index"`
}
type DataArrayString []string
type DataJSON map[string]any
type DataJSONArray []map[string]any

func init() {
	dicts = map[string][]SysDictItem{}
}

// 解析数据
func (j *DataJSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		//return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
		return nil
	}
	if len(bytes) == 0 {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

// Value方法是将自定义结构体转译成数据库能识别储存的编码
func (j DataJSON) Value() (driver.Value, error) {
	marshal, err := json.Marshal(j)
	if err != nil {
		return nil, err
	}
	return marshal, nil
}

// 解析数据
func (j *DataJSONArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		//return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
		return nil
	}
	if len(bytes) == 0 {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

// Value方法是将自定义结构体转译成数据库能识别储存的编码
func (j DataArrayString) Value() (driver.Value, error) {
	marshal, err := json.Marshal(j)
	if err != nil {
		return nil, err
	}
	return marshal, nil
}

// 解析数据
func (j *DataArrayString) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		//return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
		return nil
	}
	if len(bytes) == 0 {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

// Value方法是将自定义结构体转译成数据库能识别储存的编码
func (j DataJSONArray) Value() (driver.Value, error) {

	marshal, err := json.Marshal(j)
	if err != nil {
		return nil, err
	}
	return marshal, nil
}

func (m Model) IdTure() bool {
	if m.ID != 0 {
		return true
	}
	return false
}
func (m *Model) GetId() PrimarykeyType {
	return m.ID
}

func SetHash(h map[string]any) {
	hash = h
}

func GetModel(tableName string, isArray bool) (interface{}, error) {
	if isArray {
		tableName += "_arr"
	}
	if _, ok := hash[tableName]; ok {
		return hash[tableName], nil
	}
	return "", errors.New("未找到")
}

func RegisterSystemStatus(dictType string, value []SysDictItem) {
	for i, _ := range value {
		value[i].DictType = dictType
	}
	dicts[dictType] = value
}
func GetSystemStatus(dictType string) []SysDictItem {
	if _, ok := dicts[dictType]; ok {
		return dicts[dictType]
	}
	return nil
}
