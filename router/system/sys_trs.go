package system

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type TrsRouter struct {
}

func (s *TrsRouter) InitTrsRouter(Router *gin.RouterGroup) {
	trsRouter := Router.Group("trs").Use(middleware.OperationRecord())
	trsRouterWithoutRecord := Router.Group("trs")
	trsApi := v1.ApiGroupApp.SystemApiGroup.TrsApi
	{
		trsRouter.POST("addAnaTrs", trsApi.AddAnaTrs)
		trsRouter.POST("importExcel", trsApi.ImportExcel)
		trsRouter.POST("updateAnaTrs", trsApi.UpdateAnaTrs)
	}
	{
		trsRouterWithoutRecord.POST("getTrsList", trsApi.GetTrsList)            //查询查找
		trsRouterWithoutRecord.GET("downloadPdf", trsApi.DownloadPdf)           // 下载pdf
		trsRouterWithoutRecord.GET("downloadTemplate", trsApi.DownloadTemplate) // 下载待导入excel模板
		trsRouterWithoutRecord.POST("getAnaTrsById", trsApi.GetAnaTrsById)      // 根据id号查找
		trsRouterWithoutRecord.GET("sendNotification", trsApi.SendNotification) // 给当前登录账号发送邮件
	}
}
