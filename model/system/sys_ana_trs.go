package system

import "github.com/flipped-aurora/gin-vue-admin/server/global"

type SysAnaTrs struct {
	global.GVA_MODEL
	SpeciesName string `json:"speciesName" gorm:"comment:物种名字"`
	SampleName  string `json:"sampleName" gorm:"comment:样本名称"`
	ProjectType string `json:"projectType" gorm:"comment:项目类型"`
	ProjectNo   string `json:"projectNo" gorm:"项目号"`
}
