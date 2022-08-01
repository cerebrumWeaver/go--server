package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
)

type SearchTrsParams struct {
	system.SysAnaTrs
	request.PageInfo
	OrderKey string `json:"orderKey"`
	Desc     bool   `json:"desc"`
}
