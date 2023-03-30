package pkg

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"strings"
)

type ValidError struct {
	Key     string
	Message string
}

type ValidErrors []*ValidError

func (v *ValidError) Error() string {
	return v.Message
}

func (v ValidErrors) Error() string {
	return strings.Join(v.Errors(), ",")
}

func (v ValidErrors) Errors() []string {
	var errs []string
	for _, err := range v {
		errs = append(errs, err.Error())
	}
	return errs
}

func BindAndValid(c *gin.Context, v interface{}) (bool, ValidErrors) {
	var errs ValidErrors
	if err := c.ShouldBindBodyWith(v, binding.JSON); err != nil {
		verrs, ok := err.(validator.ValidationErrors)
		if !ok {
			return false, errs
		}
		for _, v := range verrs {
			errs = append(errs, &ValidError{
				Key:     v.Tag(),
				Message: v.Error(),
			})
		}
		return false, errs
	}
	return true, nil
}
