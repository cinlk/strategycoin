package app

import (
	"BitCoinProfitStrategy/app/handlers"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine{

	r := gin.New()

	group := r.Group("/api/v1")


	// handlers
	m := group.Group("/profit")
	{
		m.GET("/demo", handlers.FgProfit)
	}


	return r

}



