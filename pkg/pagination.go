package pkg

import (
	"github.com/gin-gonic/gin"
	"os"
)

func GetPage(c *gin.Context) int {
	page := StrTo(c.Query("page")).MustInt()
	if page <= 0 {
		return 1
	}
	return page
}

func GetPageSize(c *gin.Context) int {
	pageSize := StrTo(c.Query("pageSize")).MustInt()
	if pageSize <= 0 {
		return StrTo(os.Getenv("DEFAULT_PAGE_SIZE")).MustInt()
	}
	maxPageSize := StrTo(os.Getenv("MAX_PAGE_SIZE")).MustInt()
	if pageSize > maxPageSize {
		return maxPageSize
	}
	return pageSize
}

func GetPageOffset(page, pageSize int) int {
	result := 0
	if page > 0 {
		result = (page - 1) * pageSize
	}
	return result
}
