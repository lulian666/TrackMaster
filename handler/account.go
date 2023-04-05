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
// @Param projectID body string true "project ID"
// @Param id body string true "page size"
// @Param description body string true "description"
// @Success 200 {array} model.Account "成功"
// @Failure 400 {object} pkg.Error "请求错误"
// @Failure 500 {object} pkg.Error "内部错误"
// @Router /api/v2/accounts [post]
func (h AccountHandler) Create(c *gin.Context) {
	res := pkg.NewResponse(c)
	a := model.Account{}
	err := c.ShouldBindBodyWith(&a, binding.JSON)
	if err != nil {
		res.ToErrorResponse(pkg.NewError(pkg.BadRequest, err.Error()))
		return
	}

	err1 := h.service.CreateAccount(&a)
	if err1 != nil {
		res.ToErrorResponse(err1)
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
// @Param project query string true "project id"
// @Success 200 {object} model.Accounts "成功"
// @Failure 400 {object} pkg.Error "请求错误"
// @Failure 500 {object} pkg.Error "内部错误"
// @Router /api/v2/accounts [get]
func (h AccountHandler) List(c *gin.Context) {
	res := pkg.NewResponse(c)
	pager := pkg.Pager{
		Page:     pkg.GetPage(c),
		PageSize: pkg.GetPageSize(c),
	}

	projectID := c.Query("project")
	if projectID == "" {
		res.ToErrorResponse(pkg.NewError(pkg.BadRequest, "project required in query"))
		return
	}

	p := model.Project{
		ID: projectID,
	}

	// 如果是一个带query的查询
	desc := c.Query("description")
	a := model.Account{
		Description: desc,
	}

	accounts, totalRow, err := h.service.ListAccount(&p, &a, pager)
	if err != nil {
		res.ToErrorResponse(err)
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
// @Router /api/v2/accounts/{id} [delete]
func (h AccountHandler) Delete(c *gin.Context) {
	res := pkg.NewResponse(c)
	id := c.Param("id")
	account := model.Account{
		ID: id,
	}
	err := h.service.DeleteAccount(&account)
	if err != nil {
		res.ToErrorResponse(err)
		return
	}

	res.ToResponse(nil)
}
