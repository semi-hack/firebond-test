package controllers

import (
	"context"
	"firebond-test/configs"
	"firebond-test/services"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/gin-gonic/gin"
)

func FetchRate(c *gin.Context) {
	crypto := c.Param("cryptocurrency")
	fiat := c.Param("fiat")

	rate, err := services.GetRate(fiat, crypto)

	if err != nil {
		c.JSON(400, gin.H{"message": "error finding rates"})
		return
	}

	c.JSON(200, rate)

}

func FetchCryptoRates(c *gin.Context) {
	crypto := c.Param("cryptocurrency")

	rate, err := services.GetRates(crypto)

	if err != nil {
		c.JSON(400, gin.H{"message": "error finding rates"})
		return
	}

	c.JSON(200, rate)
}

func FetchAllRates(c *gin.Context) {
	rate, err := services.GetAllRate()

	if err != nil {
		c.JSON(400, gin.H{"message": "error finding rates"})
		return
	}

	c.JSON(200, gin.H{"status": true, "data": rate})
}

func FetchAddress(c *gin.Context) {
	address := c.Param("address")

	client, err := ethclient.Dial(configs.EnvAlchemyURI())
	if err != nil {
		c.JSON(400, gin.H{"message": "conncetion error"})
	}

	ethAddress := common.HexToAddress(address)

	// Get the balance of the address
	balance, err := client.BalanceAt(context.Background(), ethAddress, nil)
	if err != nil {
		c.JSON(400, gin.H{"message": "failed to get balance"})
	}

	// Return the balance as a response
	c.JSON(200, gin.H{
		"address": address,
		"balance": balance.String(),
	})

}

func FetchHistory(c *gin.Context) {
	crypto := c.Param("cryptocurrency")
	fiat := c.Param("fiat")

	rates, err := services.GetHistory(fiat, crypto)

	if err != nil {
		c.JSON(400, gin.H{"message": "error finding rate history"})
		return
	}

	c.JSON(200, gin.H{"status": true, "data": rates})
}
