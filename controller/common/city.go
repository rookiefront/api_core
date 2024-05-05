package common

import (
	"fmt"
	"github.com/rookiefront/api-core/cmd/model"
	"github.com/rookiefront/api-core/define"
	"github.com/rookiefront/api-core/global"
)

type ApiCity struct {
}
type reqCity struct {
	Code       int    `json:"code"`
	ParentCode int    `json:"parent_code"`
	Name       string `json:"name"`
	Leaf       bool   `json:"leaf" gorm:"-"`
}

func (api ApiCity) Lazy(c *define.BasicContext) {
	var result []reqCity
	parentCode := c.Query("parentCode")
	if parentCode == "" {
		parentCode = "-1"
	}
	global.DB.Model(model.SysCity{}).Find(&result)
	currentList := []reqCity{}
	for _, city := range result {
		if fmt.Sprintf("%d", city.ParentCode) == parentCode {
			currentList = append(currentList, city)
		}
	}

	for i, city := range currentList {
		currentList[i].Leaf = true
		for _, r := range result {
			if r.ParentCode == city.Code {
				currentList[i].Leaf = false
				break
			}
		}
	}
	c.SendJsonOk(currentList)
}
