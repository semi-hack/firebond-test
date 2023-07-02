package routes

import (
	"firebond-test/controllers"

	"github.com/gin-gonic/gin"
)


func Initialize() {
	r := gin.Default()

	r.GET("/rates", controllers.FetchAllRates)
	r.GET("/rates/:cryptocurrency", controllers.FetchCryptoRates)
	r.GET("/rates/:cryptocurrency/:fiat", controllers.FetchRate)
	r.GET("/rates/history/:cryptocurrency/:fiat", controllers.FetchHistory)
	r.GET("/balance/:address", controllers.FetchAddress)

	r.Run()
}