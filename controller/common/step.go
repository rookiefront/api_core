package common

import (
	"fmt"
	"github.com/front-ck996/csy"
	"github.com/rookiefront/api-core/cmd/model"
	"github.com/rookiefront/api-core/define"
	"github.com/rookiefront/api-core/global"
	"strings"
	"time"
)

type ApiStep struct {
}

func findAllParentIDs(id model.PrimarykeyType, steps []model.SysStep) []model.PrimarykeyType {
	var parentIDs []model.PrimarykeyType
	for _, step := range steps {
		if step.Model.ID == id {
			parentIDs = append(parentIDs, step.ParentId)
			parentIDs = append(parentIDs, findAllParentIDs(step.ParentId, steps)...)
		}
	}
	return parentIDs
}

func (api ApiStep) RefreshFulPath(c *define.BasicContext) {
	var list []model.SysStep
	global.DB.Find(&list)
	list2 := list
	for _, v := range list {
		if v.CreatedAt.Format("2006") == "0001" {
			v.CreatedAt = time.Now()
		}
		ids := csy.SliceReverse(findAllParentIDs(v.Model.ID, list))
		var ids2 []string
		for _, v := range ids {
			if v != 0 {
				ids2 = append(ids2, fmt.Sprintf("%d", v))
			}
		}
		ids2 = append(ids2, fmt.Sprintf("%d", v.Model.ID))
		ids2 = csy.SliceUnique[string](ids2)
		fullPath := strings.Join(ids2, ",")
		v.Leaf = true
		for _, c1 := range list2 {
			if c1.ParentId == v.ID {
				v.Leaf = false
				if v.Leaf != c1.Leaf {
					global.DB.Save(&v)
				}
				break
			}
		}
		if v.FullPath != fullPath {
			v.FullPath = fullPath
			global.DB.Save(&v)

		}
	}
	c.SendJsonOk("")
}
