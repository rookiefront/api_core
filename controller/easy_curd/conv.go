package easy_curd

import (
	"github.com/front-ck996/csy"
	"github.com/rookiefront/api-core/define"
)

func FbFieldConvFront(c *define.BasicContext) {
	data := c.GetReqData()
	c.SendJsonOk(csy.StrFirstToUpper(data["input"].(string)))
}
