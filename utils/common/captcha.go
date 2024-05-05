package common

import "github.com/mojocn/base64Captcha"

type captcha struct {
	Driver base64Captcha.Driver
	Store  base64Captcha.Store
}

var Captcha = &captcha{}
var store = base64Captcha.DefaultMemStore

func init() {
	//driver := base64Captcha.DefaultDriverDigit
	//80, 240, 5, OptionShowSineLine | OptionShowSlimeLine | OptionShowHollowLine, nil, []string{"3Dumb.ttf"}
	driver := &base64Captcha.DriverMath{
		Height:     80,
		Width:      240,
		NoiseCount: 2,
		//ShowLineOptions: base64Captcha.OptionShowSineLine | base64Captcha.OptionShowSlimeLine | base64Captcha.OptionShowHollowLine,
		//BgColor:         tt.fields.BgColor,
	}
	driver.ConvertFonts()

	Captcha = NewCaptcha(driver, store)
}

func NewCaptcha(driver base64Captcha.Driver, store base64Captcha.Store) *captcha {
	return &captcha{Driver: driver, Store: store}
}

func (c *captcha) Generate() (id, b64s string, err error) {
	id, content, answer := c.Driver.GenerateIdQuestionAnswer()
	item, err := c.Driver.DrawCaptcha(content)
	if err != nil {
		return "", "", err
	}
	c.Store.Set(id, answer)
	b64s = item.EncodeB64string()
	return
}

func (c *captcha) Verify(id, answer string, clear bool) (match bool) {
	match = c.Store.Get(id, clear) == answer
	return
}
