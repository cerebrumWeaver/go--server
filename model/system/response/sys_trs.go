package response

import "github.com/flipped-aurora/gin-vue-admin/server/model/system"

type SysTrsResponse struct {
	Trs system.SysAnaTrs `json:"trs"`
}
