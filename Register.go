package api_core

import (
	"github.com/rookiefront/api-core/initialize"
)

func Register(callBack func()) {
	initialize.Init(callBack)
}
