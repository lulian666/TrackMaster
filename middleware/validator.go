package middleware

import (
	"TrackMaster/pkg"
	"github.com/gin-gonic/gin"
)

func validate[T any](c *gin.Context, v T) {
	isValid, err := pkg.BindAndValid(c, v)
	if !isValid {
		res := pkg.NewResponse(c)
		res.ToErrorResponse(
			pkg.NewError(pkg.InvalidParams, "invalid parameters").WithDetails(
				err.Errors()...,
			))
		res.Ctx.Abort()
	}
}

func Validator(v interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		validate(c, v)
		c.Next()
	}
}
