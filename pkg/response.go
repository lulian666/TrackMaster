package pkg

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Ctx *gin.Context
}

type Pager struct {
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
	TotalRow int64 `json:"totalRow"`
}

func NewResponse(c *gin.Context) *Response {
	return &Response{
		Ctx: c,
	}
}

func (r *Response) ToResponse(data interface{}) {
	if data == nil {
		data = gin.H{}
	}
	r.Ctx.JSON(http.StatusOK, data)
}

func (r *Response) ToResponseList(list interface{}, totalRow int64) {
	r.Ctx.JSON(http.StatusOK, gin.H{
		"data": list,
		"pager": Pager{
			Page:     GetPage(r.Ctx),
			PageSize: GetPageSize(r.Ctx),
			TotalRow: totalRow,
		},
	})
}

func (r *Response) ToErrorResponse(err *Error) {
	response := gin.H{"code": err.Code, "msg": err.Msg}
	details := err.Details
	if len(details) > 0 {
		response["details"] = details
	}
	r.Ctx.JSON(err.StatusCode(), response)
}
