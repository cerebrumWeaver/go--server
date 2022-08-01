package system

import (
	"errors"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/xuri/excelize/v2"
)

type ExcelService struct {
}

// @Tags excel
// @Summary 导入Excel文件
// @Security ApiKeyAuth
// @accept multipart/form-data
// @Produce  application/json
// @Param file formData file true "导入Excel文件"
// @Success 200 {object} response.Response{msg=string} "导入Excel文件"
// @Router /excel/importExcel [post]
func (exa *ExcelService) ImportExcel(fileName string) ([]system.SysAnaTrs, error) {
	/*_, header, err := c.Request.FormFile("file")
	if err != nil {
		global.GVA_LOG.Error("导入文件失败!", zap.Error(err))
		response.FailWithMessage("导入文件失败", c)
		return
	}*/
	//fmt.Println("文件名：", header.Filename)
	//_ = c.SaveUploadedFile(header, global.GVA_CONFIG.Excel.Dir+header.Filename)
	if menus, err := exa.ParseExcel2InfoList(fileName); err != nil {
		/*global.GVA_LOG.Error("导入文件失败!", zap.Error(err))
		response.FailWithMessage("导入文件失败", err.Error())*/
		return nil, err
	} else {
		//response.OkWithMessage("导入成功", c)
		return menus, err
	}
}

// ParseExcel2InfoList 解析 Excel内容
func (exa *ExcelService) ParseExcel2InfoList(fileName string) ([]system.SysAnaTrs, error) {
	skipHeader := true
	fixedHeader := []string{"物种名字", "样本名称", "项目类型", "项目号"}
	file, err := excelize.OpenFile(fileName)
	if err != nil {
		return nil, err
	}
	//menus := make([]system.SysBaseMenu, 0)
	menus := make([]system.SysAnaTrs, 0)
	rows, err := file.Rows("Sheet1") // 打开一个sheet表单
	if err != nil {
		return nil, err
	}
	for rows.Next() { // 如果找到下一个行元素，Next将返回true
		row, err := rows.Columns() // 返回当前行的列值
		if err != nil {
			return nil, err
		}
		if skipHeader { // 判断模板第一行是否为指定字段
			if exa.compareStrSlice(row, fixedHeader) {
				skipHeader = false
				continue
			} else {
				return nil, errors.New("Excel格式错误")
			}
		}
		if len(row) != len(fixedHeader) { // 某一行数据 有 遗漏缺失或者多余，跳过此行数据处理
			continue
		}
		//id, _ := strconv.Atoi(row[0])
		//hidden, _ := strconv.ParseBool(row[3])
		//sort, _ := strconv.Atoi(row[5])
		/*menu := system.SysBaseMenu{
			GVA_MODEL: global.GVA_MODEL{
				ID: uint(id),
			},
			Name:      row[1],
			Path:      row[2],
			Hidden:    hidden,
			ParentId:  row[4],
			Sort:      sort,
			Component: row[6],
		}*/
		menu := system.SysAnaTrs{
			/*GVA_MODEL: global.GVA_MODEL{
				ID: uint(id),
			},*/
			SpeciesName: row[0],
			SampleName:  row[1],
			ProjectType: row[2],
			ProjectNo:   row[3],
		}
		menus = append(menus, menu)
	}
	return menus, nil
}

func (exa *ExcelService) compareStrSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	if (b == nil) != (a == nil) {
		return false
	}
	for key, value := range a {
		if value != b[key] {
			return false
		}
	}
	return true
}
