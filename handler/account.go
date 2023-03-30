package handler

import (
	"TrackMaster/model"
	"TrackMaster/pkg"
	"TrackMaster/service"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type AccountHandler struct {
	service service.AccountService
}

func NewAccountHandler(s service.AccountService) AccountHandler {
	return AccountHandler{
		service: s,
	}
}

// Create
// @Tags account
// @Summery create account
// @Produce json
// @Param project body model.Account false "account"
// @Success 200 {array} model.Account "成功"
// @Failure 400 {object} pkg.Error "请求错误"
// @Failure 500 {object} pkg.Error "内部错误"
// @Router /api/v1/accounts [post]
func (h AccountHandler) Create(c *gin.Context) {
	res := pkg.NewResponse(c)
	a := model.Account{}
	err := c.ShouldBindBodyWith(&a, binding.JSON)
	if err != nil {
		res.ToErrorResponse(pkg.NewError(pkg.BadRequest, err.Error()))
		return
	}

	err = h.service.CreateAccount(&a)
	if err != nil {
		res.ToErrorResponse(pkg.NewError(pkg.ServerError, err.Error()))
		return
	}

	res.ToResponse(a)
}

// List
// @Tags account
// @Summery list account
// @Produce json
// @Param page query string false "page"
// @Param pageSize query string false "page size"
// @Param description query string false "description"
// @Success 200 {array} model.Accounts "成功"
// @Failure 400 {object} pkg.Error "请求错误"
// @Failure 500 {object} pkg.Error "内部错误"
// @Router /api/v1/accounts [get]
func (h AccountHandler) List(c *gin.Context) {
	res := pkg.NewResponse(c)
	pager := pkg.Pager{
		Page:     pkg.GetPage(c),
		PageSize: pkg.GetPageSize(c),
	}

	// 如果是一个带query的查询
	desc := c.Query("description")
	a := model.Account{
		Description: desc,
	}

	accounts, totalRow, err := h.service.ListAccount(&a, pager)
	if err != nil {
		res.ToErrorResponse(pkg.NewError(pkg.ServerError, err.Error()))
		return
	}
	res.ToResponseList(accounts, totalRow)
}

// Delete
// @Tags account
// @Summery delete account
// @Produce json
// @Param id path string true "account id"
// @Success 200 {object} object "成功"
// @Failure 400 {object} pkg.Error "请求错误"
// @Failure 500 {object} pkg.Error "内部错误"
// @Router /api/v1/accounts/{id} [delete]
func (h AccountHandler) Delete(c *gin.Context) {
	res := pkg.NewResponse(c)
	id := c.Param("id")
	account := model.Account{
		ID: id,
	}
	err := h.service.DeleteAccount(&account)
	if err != nil {
		res.ToErrorResponse(pkg.NewError(pkg.ServerError, err.Error()))
		return
	}

	res.ToResponse(nil)
}
