package model

import (
	"database/sql/driver"
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

type PrimarykeyType uint

type Model struct {
	ID           PrimarykeyType `json:"id" gorm:"primarykey;int"`
	CreatedAt    time.Time      `json:"-"`
	UpdatedAt    time.Time      `json:"-"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	CreateUserID uint           `json:"-" gorm:"column:c_id;index"`
}

type DataJSON map[string]any
type DataJSONArray []map[string]any

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