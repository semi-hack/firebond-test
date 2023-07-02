package main

import (
	"firebond-test/configs"
	"firebond-test/routes"
	"firebond-test/services"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	configs.ConnectDB()

	c := cron.New()

	go Cron()

	_, err := c.AddFunc("*/5 * * * *", Cron)
	if err != nil {
		fmt.Println("Error scheduling cron job:", err)
		return
	}

	c.Start()

	routes.Initialize()

	r.Run()

	select {}
}

func Cron() {
	cryptoCurrencies := []string{"BTC", "SOL", "ETH", "XRP", "LTC"}
	fiatCurrencies := []string{"NGN", "USD", "EUR", "JPY", "GBP"}

	var crypto string
	var fiat string

	combinations := generateCombinations(cryptoCurrencies, fiatCurrencies)

	for _, combination := range combinations {

		crypto = combination[1]
		fiat = combination[0]

		err := services.GetMarketRates(crypto, fiat)

		if err != nil {
			fmt.Println("ERROR")
		}

		time.Sleep(5 * time.Second)
	}

}

func generateCombinations(cryptoCurrencies, fiatCurrencies []string) [][]string {
	combinations := [][]string{}
	for _, crypto := range cryptoCurrencies {
		for _, fiat := range fiatCurrencies {
			combination := []string{crypto, fiat}
			combinations = append(combinations, combination)
		}
	}
	return combinations
}
