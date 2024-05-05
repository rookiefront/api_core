package common

import (
	"errors"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/zh"
)

type _validate struct {
	Uni      *ut.UniversalTranslator
	Trans    ut.Translator
	Validate *validator.Validate
}

var Validate = _validate{}

func init() {
	en := en.New()
	Validate.Uni = ut.New(en, en)

	Validate.Trans, _ = Validate.Uni.GetTranslator("zh")

	Validate.Validate = validator.New()
	en_translations.RegisterDefaultTranslations(Validate.Validate, Validate.Trans)
}

func (that *_validate) Struct(value interface{}) (err error) {
	return that.Translate(that.Validate.Struct(value))
}

func (that *_validate) Translate(err error, tip ...string) error {
	tipMessage := ""
	if len(tip) > 0 {
		tipMessage = tip[0]
	}
	if err != nil {
		var errs validator.ValidationErrors
		errors.As(err, &errs)

		for _, e := range errs {
			return errors.New(tipMessage + " " + e.Translate(that.Trans))
		}
	}
	return nil
}

func (that *_validate) Var(currentValue interface{}, currentVerify string, tip string) error {
	return that.Translate(that.Validate.Var(currentValue, currentVerify), tip)
}
