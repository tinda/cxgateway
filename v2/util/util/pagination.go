package util

import (
	"github.com/codingXiang/cxgateway/v2/middleware/pagination"
	"github.com/jinzhu/gorm"
	"strconv"
)

func Pagination(in *gorm.DB, data map[string]interface{}) (*gorm.DB, map[string]interface{}) {
	if data == nil {
		return in, data
	} else if data[pagination.PAGE_SIZE] == nil && data[pagination.PAGE] == nil {
		return in, data
	} else {
		var (
			pageSize = 10
			page     = 1
		)

		if in := data[pagination.PAGE_SIZE]; in != nil {
			pageSize, _ = strconv.Atoi(in.(string))
			delete(data, pagination.PAGE_SIZE)
		}
		if in := data[pagination.PAGE]; in != nil {
			page, _ = strconv.Atoi(in.(string))
			delete(data, pagination.PAGE)
		}
		return in.Limit(pageSize).Offset((page - 1) * pageSize), data
	}
}
