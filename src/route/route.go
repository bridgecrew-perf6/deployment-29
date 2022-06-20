package route

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"oneclick/controller"
)

func Route() {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard
	r := gin.Default()
	categoryController := controller.NewCategoryController()

	r.POST("/HotUpdateWg/post", categoryController.UpdateWgConfig)
	r.POST("/Configurations/post", categoryController.ShowConfig)
	r.POST("/ParsingState/post", categoryController.ShowStatus)
	r.POST("/IntervalManager/post", categoryController.ExecPing)
	r.POST("/PingIntervalBlock/post", categoryController.UpdateTP)
	r.POST("/PersianUPDATE/post", categoryController.Update)
	r.POST("/LoggingManager/post", categoryController.ShowLog)

	//r.POST("/HotUpdateWg/post", controller.HotUpdateWg)
	//r.POST("/Configurations/post", controller.Configurations)
	//r.POST("/ParsingState/post", controller.ParsingState)
	//r.POST("/IntervalManager/post", controller.IntervalManager)
	//r.POST("/PingIntervalBlock/post", controller.PingIntervalBlock)
	//r.POST("/PersianUPDATE/post", controller.PersianUPDATE)
	//r.POST("/LoggingManager/post", controller.LoggingManager)
	r.Run(":8095")
}
