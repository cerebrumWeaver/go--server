package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

type SysAnaTrs struct {
	global.GVA_MODEL
	ProjectNo           string `json:"projectNo" gorm:"项目编号"`
	RawCleanDir         string `json:"rawCleanDir" gorm:"数据目录"`
	AnalysisType        string `json:"analysisType" gorm:"分析类型"`
	SampleInfo          string `json:"sampleInfo" gorm:"样本信息"`
	CmpGroup            string `json:"cmpGroup" gorm:"比较组"`
	SpeciesNames        string `json:"speciesNames" gorm:"物种名字"`
	DifferenceThreshold string `json:"differenceThreshold" gorm:"差异阈值"`
	SubProjectNo        string `json:"subProjectNo" gorm:"子项目编号"`
	PrjType             string `json:"prjType" gorm:"项目类型"`
	ParaNum             int    `json:"paraNum" gorm:"线程数"`
	SeqType             string `json:"seqType" gorm:"测序类型"`
	Strand              bool   `json:"strand" gorm:"链特异性"`
	GeneId              string `json:"geneId" gorm:"基因或转录本"`
	FeatureType         string `json:"featureType" gorm:"结构类型"`
	Category            string `json:"category" gorm:"动植物"`
	PpiSpecies          string `json:"ppiSpecies" gorm:"近缘物种"`
	State               int    `json:"state" gorm:"分析状态"`
}
type TrsApi struct {
}
