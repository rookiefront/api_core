package define

import (
	"encoding/json"
)

type ResultJSON struct {
	Err  error       `json:"err"`
	Data interface{} `json:"data"`
}

func (j ResultJSON) MarshalJSON() ([]byte, error) {
	h := map[string]interface{}{
		"err":  j.Err,
		"data": j.Data,
	}

	if j.Err != nil {
		h["err"] = j.Err.Error()
	}

	return json.Marshal(h)
}

func NewResultJSON() ResultJSON {
	return ResultJSON{
		Data: nil,
	}
}

func (c ResultJSON) ToJSON() string {
	marshal, err := json.Marshal(c)
	if err != nil {
		return "{}"
	}
	return string(marshal)
}

func (c ResultJSON) ToTagJSON() string {
	marshal, err := json.Marshal(c)
	if err != nil {
		return "result=> {} <=result"
	}
	return "result=>" + string(marshal) + "<=result"
}
