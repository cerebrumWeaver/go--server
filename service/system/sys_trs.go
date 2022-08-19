package system

import (
	"errors"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"gorm.io/gorm"
)

type TrsService struct {
}

func (trsService *TrsService) CheckTrsBeforeAddAnaTrs(anaTrs system.SysAnaTrs) bool {
	return !errors.Is(global.GVA_DB.Where("project_no = ?", anaTrs.ProjectNo).First(&system.SysAnaTrs{}).Error, gorm.ErrRecordNotFound)
}

func (trsService *TrsService) AddAnaTrs(anaTrs system.SysAnaTrs) error {
	//func (trsService *TrsService) AddAnaTrs(cc <-chan system.SysAnaTrs, f chan<- int, r chan<- error) {
	/*if !errors.Is(global.GVA_DB.Where("project_no = ?", anaTrs.ProjectNo).First(&system.SysAnaTrs{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在重复项目号，请修改projectNo")
	}*/
	//anaTrs, _ := <-cc
	fmt.Println("开始添加：")
	//fmt.Println(anaTrs)
	//err := global.GVA_DB.Create(&anaTrs).Error
	//f <- 1
	//r <- err
	return global.GVA_DB.Create(&anaTrs).Error
}

func (s *TrsService) UpdateAnaTrs(trs system.SysAnaTrs) (err error) {
	var oldTrs system.SysAnaTrs
	updateMap := make(map[string]interface{})
	//updateMap["species_name"] = trs.SpeciesName
	//updateMap["sample_name"] = trs.SampleName
	//updateMap["project_type"] = trs.ProjectType
	updateMap["project_no"] = trs.ProjectNo

	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		db := tx.Where("id = ?", trs.ID).Find(&oldTrs)
		/*if oldTrs.ProjectNo != trs.ProjectNo { // 根据唯一ProjectNo项目号来做判别, 由于是自身更新，这里要用"!="符号表示, 否则自身等于自身时就无法更新了
			if !errors.Is(tx.Where("id <> ? AND name = ?", trs.ID, trs.ProjectNo).First(&system.SysAnaTrs{}).Error, gorm.ErrRecordNotFound) {
				global.GVA_LOG.Debug("存在相同ProjectNo项目号修改失败")
				return errors.New("存在相同ProjectNo项目号修改失败")
			}
		}*/
		/*err := utils.ExeSh(trs) // 先注释测试，测试完取消注释
		if err != nil {
			global.GVA_LOG.Error("系统执行脚本出错!!!", zap.Error(err))
			response.FailWithMessage("系统执行脚本出错,请联系管理员", c)
			//log.Fatalf("cmd.Run() failed with %s\n", err)
			return err// 如果脚本执行失败，则后续不执行
		}*/

		txErr := db.Updates(trs).Error
		if txErr != nil {
			global.GVA_LOG.Debug(txErr.Error())
			return txErr
		}
		return nil
	})
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetFileRecordInfoList
//@description: 分页获取数据
//@param: info request.PageInfo
//@return: list interface{}, total int64, err error
func (s *TrsService) GetTrsRecordInfoList(trs system.SysAnaTrs, info request.PageInfo, order string, desc bool) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	//keyword := info.Keyword
	db := global.GVA_DB.Model(&system.SysAnaTrs{})
	var trsLists []system.SysAnaTrs
	// 有新字段，请按如下格式添加
	if trs.ProjectNo != "" {
		db = db.Where("project_no LIKE ?", "%"+trs.ProjectNo+"%")
	}
	if trs.AnalysisType != "" {
		db = db.Where("analysis_type LIKE ?", "%"+trs.AnalysisType+"%")
	}
	if trs.SpeciesNames != "" {
		db = db.Where("species_names LIKE ?", "%"+trs.SpeciesNames+"%")
	}
	if trs.PrjType != "" {
		db = db.Where("prj_type LIKE ?", "%"+trs.PrjType+"%")
	}
	if trs.Category != "" {
		db = db.Where("category LIKE ?", "%"+trs.Category+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return trsLists, total, err
	} else {
		db = db.Limit(limit).Offset(offset)
		if order != "" {
			var OrderStr string
			orderMap := make(map[string]bool, 6) // 注意，5个true，有几个字段就分几个容量
			orderMap["id"] = true
			orderMap["project_no"] = true
			orderMap["analysis_type"] = true
			orderMap["species_names"] = true
			orderMap["prj_type"] = true
			orderMap["category"] = true
			if orderMap[order] {
				if desc {
					OrderStr = order + " desc"
				} else {
					OrderStr = order
				}
			} else {
				err = fmt.Errorf("非法的排序字段: %v", order)
				return trsLists, total, err
			}
			err = db.Order(OrderStr).Find(&trsLists).Error
		} else {
			err = db.Order("updated_at desc").Find(&trsLists).Error
		}
	}
	/*if len(keyword) > 0 {
		db = db.Where("project_no LIKE ?", "%"+keyword+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}*/
	//err = db.Limit(limit).Offset(offset).Order("updated_at desc").Find(&trsLists).Error
	return trsLists, total, err
}

// GetTrsRecordInfoById 通过ID号查找记录
func (s *TrsService) GetTrsRecordInfoById(anaTrs *system.SysAnaTrs) (anaTrsProjectNo *system.SysAnaTrs, err error) {
	var trs system.SysAnaTrs
	if err = global.GVA_DB.Where("id = ?", anaTrs.ID).First(&trs).Error; err != nil {
		return nil, err
	}
	return &trs, err
}
