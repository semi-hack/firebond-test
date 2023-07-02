package routes

import (
	"firebond-test/controllers"

	"github.com/gin-gonic/gin"
)


func Initialize() {
	r := gin.Default()

	r.GET("/rates", controllers.FetchRate)
	r.GET("/rates/history/:cryptocurrency/:fiat", controllers.FetchHistory)
	r.GET("/balance/:address", controllers.FetchAddress)

	r.Run()
}