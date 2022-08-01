package system

import (
	"bytes"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	systemRes "github.com/flipped-aurora/gin-vue-admin/server/model/system/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

type TrsApi struct {
}

func (a *TrsApi) exeSh(projectNo, speciesName, projectType string) error {
	cmd := exec.Command("sh", "test_RNA_pipeline.sh", projectNo, speciesName, projectType)
	var stdin, stdout, stderr bytes.Buffer
	cmd.Stdin = &stdin
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	fmt.Printf("out:\n%s\nerr:\n%s\n", outStr, errStr)
	return err
}

func (a *TrsApi) AddAnaTrs(c *gin.Context) {
	var anaTrs system.SysAnaTrs
	_ = c.ShouldBindJSON(&anaTrs)
	if err := utils.Verify(anaTrs, utils.TrsVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	projectNo := anaTrs.ProjectNo
	speciesName := anaTrs.SpeciesName
	projectType := anaTrs.ProjectType
	if exits := trsService.CheckTrsBeforeAddAnaTrs(anaTrs); exits {
		global.GVA_LOG.Error("存在重复项目号，请修改projectNo!")
		response.FailWithMessage("存在重复项目号，请修改项目号", c)
		return // 如果projectNo项目号存在，则后续不执行
	}
	err := a.exeSh(projectNo, speciesName, projectType)
	//cmd := exec.Command("sh", "test_RNA_pipeline.sh", "RNAS003", "human", "mRNA-seq")
	//cmd := exec.Command("sh", "test_RNA_pipeline.sh", projectNo, speciesName, projectType)
	//var stdin, stdout, stderr bytes.Buffer
	//cmd.Stdin = &stdin
	//cmd.Stdout = &stdout
	//cmd.Stderr = &stderr
	//err := cmd.Run()
	if err != nil {
		global.GVA_LOG.Error("系统执行脚本出错!!!", zap.Error(err))
		response.FailWithMessage("系统执行脚本出错,请联系管理员", c)
		//log.Fatalf("cmd.Run() failed with %s\n", err)
		return // 如果脚本执行失败，则后续不执行
	}
	//outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	//fmt.Printf("out:\n%s\nerr:\n%s\n", outStr, errStr)
	if err := trsService.AddAnaTrs(anaTrs); err != nil {
		/**
		此处需要添加一个删除文件的linux系统命令

		*/

		global.GVA_LOG.Error("添加失败!", zap.Error(err))
		response.FailWithMessage("添加失败", c)
	} else {
		response.OkWithMessage("添加成功", c)
	}
}

// @Tags Menu
// @Summary 更新转录组信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body system.SysAnaTrs true "物种名字SpeciesName, 样本名称sampleName, 项目类型projectType, 项目号projectNo"
// @Success 200 {object} response.Response{msg=string} "更新转录组信息"
// @Router /trs/UpdateAnaTrs [post]
func (a *TrsApi) UpdateAnaTrs(c *gin.Context) {
	var trs system.SysAnaTrs
	_ = c.ShouldBindJSON(&trs)
	if err := utils.Verify(trs, utils.TrsVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := trsService.UpdateAnaTrs(trs); err != nil {
		global.GVA_LOG.Error("更新失败,疑似项目号已经存在!", zap.Error(err))
		response.FailWithMessage("更新失败,疑似项目号已经存在", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// @Tags ExaFileUploadAndDownload
// @Summary 分页转录组信息列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.PageInfo true "页码, 每页大小"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "分页文件列表,返回包括列表,总数,页码,每页数量"
// @Router /trs/GetFileList [post]
func (a *TrsApi) GetTrsList(c *gin.Context) {
	var pageInfo systemReq.SearchTrsParams // 实体类 下 request文件夹 实体类
	_ = c.ShouldBindJSON(&pageInfo)
	if err := utils.Verify(pageInfo.PageInfo, utils.PageInfoVerify); err != nil {
		global.GVA_LOG.Error("查询失败！", zap.Error(err))
		return
	}
	list, total, err := trsService.GetTrsRecordInfoList(pageInfo.SysAnaTrs, pageInfo.PageInfo, pageInfo.OrderKey, pageInfo.Desc)
	if err != nil {
		global.GVA_LOG.Error("获取转录组信息失败", zap.Error(err))
		response.FailWithMessage("获取转录组信息失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "转录组信息获取成功", c)
	}
}

// GetAnaTrsById 通过唯一id查找记录
func (a *TrsApi) GetAnaTrsById(c *gin.Context) {
	//var anaTrs system.SysAnaTrs
	//_ = c.ShouldBindJSON(&anaTrs)
	//id := c.Query("id")
	//idInt, _ := strconv.Atoi(id)
	var id request.GetById
	_ = c.ShouldBindJSON(&id)
	if err := utils.Verify(id, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	anaTrs := &system.SysAnaTrs{GVA_MODEL: global.GVA_MODEL{ID: uint(id.ID)}}
	if trsInfo, err := trsService.GetTrsRecordInfoById(anaTrs); err != nil {
		global.GVA_LOG.Error("此条记录未找到，无法编辑!", zap.Error(err))
		response.FailWithMessage("此条记录未找到，无法编辑", c)
	} else {
		response.OkWithDetailed(systemRes.SysTrsResponse{Trs: *trsInfo}, "获取成功", c)
	}
}

/*func (a *TrsApi) GetTrsList(c *gin.Context) {
	var pageInfo request.PageInfo
	_ = c.ShouldBindJSON(&pageInfo)
	list, total, err := trsService.GetTrsRecordInfoList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取转录组信息失败", zap.Error(err))
		response.FailWithMessage("获取转录组信息失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "转录组信息获取成功", c)
	}
}*/

// @Summary 下载pdf
// @Router /trs/DownloadPdf
func (a *TrsApi) DownloadPdf(c *gin.Context) {
	id := c.Query("fileName")
	idInt, _ := strconv.Atoi(id)
	trs := &system.SysAnaTrs{GVA_MODEL: global.GVA_MODEL{ID: uint(idInt)}}
	if trsInfo, err := trsService.GetTrsRecordInfoById(trs); err != nil {
		global.GVA_LOG.Error("该项目号未匹配到下载文件!", zap.Error(err))
		response.FailWithMessage("该项目号未匹配到下载文件", c)
	} else {
		fileName := trsInfo.ProjectNo
		cur, _ := os.Getwd()
		filePathEnd := "outs/" + fileName + ".txt"
		filePathInServer := filepath.Join(cur, filePathEnd)
		fmt.Println("我的路径是：\n", fileName)
		_, err := os.Stat(filePathInServer)
		if err != nil {
			global.GVA_LOG.Error("文件不存在!", zap.Error(err))
			response.FailWithMessage("文件不存在", c)
		} else {
			c.Writer.Header().Add("success", "true")
			c.File(filePathInServer)
		}
	}
}

// @Summary 下载待导入的Excel模板
// @Router /trs/DownloadTemplate
func (a *TrsApi) DownloadTemplate(c *gin.Context) {
	excelName := c.Query("excelName")
	cur, _ := os.Getwd()
	filePathEnd := "resource/excel/" + excelName
	excelPath := filepath.Join(cur, filePathEnd)
	if _, err := os.Stat(excelPath); err != nil {
		global.GVA_LOG.Error("模板不存在，联系管理员!", zap.Error(err))
		response.FailWithMessage("模板不存在，联系管理员!", c)
	} else {
		c.Writer.Header().Add("success", "true")
		c.File(excelPath)
	}
}

func (a *TrsApi) ImportExcel(c *gin.Context) {
	fmt.Println("导入excel：api层")
	_, header, err := c.Request.FormFile("file")
	if err != nil {
		global.GVA_LOG.Error("文件传输失败!", zap.Error(err))
		response.FailWithMessage("文件传输失败", c)
		return
	}
	fmt.Println("文件路径：", global.GVA_CONFIG.Excel.Dir+header.Filename)
	//_ = c.SaveUploadedFile(header, global.GVA_CONFIG.Excel.Dir+"ExcelImport.xlsx")
	_ = c.SaveUploadedFile(header, global.GVA_CONFIG.Excel.Dir+header.Filename) // 保存文件 至 指定目录
	if trsList, err := excelService.ImportExcel(global.GVA_CONFIG.Excel.Dir + header.Filename); err != nil {
		global.GVA_LOG.Error("文件内容格式有误!!!", zap.Error(err))
		response.FailWithMessage("文件内容格式有误", c)
	} else {
		for i, trs := range trsList {
			//fmt.Println("Excel内容：\n", trs.SpeciesName, )
			fmt.Printf("Excel第%d行内容：\t%v\t%v\t%v\t%v\n", i, trs.SpeciesName, trs.SampleName, trs.ProjectType, trs.ProjectNo)
			if exits := trsService.CheckTrsBeforeAddAnaTrs(trs); exits {
				global.GVA_LOG.Error("存在重复项目号，请修改projectNo!")
				response.FailWithMessage("存在重复项目号，请修改项目号", c)
				continue // 如果projectNo项目号存在，则后续不执行
			}
			err := a.exeSh(trs.ProjectNo, trs.SampleName, trs.ProjectType)
			if err != nil {
				global.GVA_LOG.Error("系统执行脚本出错!", zap.Error(err))
				response.FailWithMessage("系统执行脚本出错,请联系管理员", c)
				//log.Fatalf("cmd.Run() failed with %s\n", err)
				continue // 如果脚本执行失败，则后续不执行
			}
			if err := trsService.AddAnaTrs(trs); err != nil {
				global.GVA_LOG.Error("添加失败!", zap.Error(err))
				response.FailWithMessage("添加失败", c)
			} else {
				response.OkWithMessage("添加成功", c)
			}
		}
		//response.OkWithMessage("导入成功", c)
	}
}

func (a *TrsApi) SendNotification(c *gin.Context) {
	uuid := utils.GetUserUuid(c)
	if ReqUser, err := userService.GetUserInfo(uuid); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		fmt.Println("获取到的邮箱：\n", ReqUser.Email)
		utils.SendInformation(ReqUser.Email)
		response.OkWithDetailed(gin.H{"userInfo": ReqUser}, "获取成功", c)
	}
}
