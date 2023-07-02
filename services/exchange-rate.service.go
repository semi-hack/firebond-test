package services

import (
	"context"
	"encoding/json"
	"errors"
	"firebond-test/configs"
	"firebond-test/models"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrorCouldNotFindPair               = "failed to find crypto fiat pair"
	ErrorCouldNotFindRates              = "failed to find exchange rates"
	ErrorFailedToDeserialize            = "failed to deserialize"
	ErrorCouldNotCreateExchangeRatePair = "failed to create exchange rate pair"
)

var ExchangeRateCollection *mongo.Collection = configs.GetCollection(configs.DB, "exchange-rates")
var HistoryCollection *mongo.Collection = configs.GetCollection(configs.DB, "exchange-rate-histories")

var ctx = context.TODO()

func Create(fiat, crypto string, rate float64) (primitive.ObjectID, error) {

	var existingRate *models.ExchangeRate

	exchangeRate := models.ExchangeRate{
		ID:        primitive.NewObjectID(),
		Crypto:    crypto,
		Fiat:      fiat,
		Rate:      rate,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	query := bson.M{"fiat": fiat, "crypto": crypto}

	err := ExchangeRateCollection.FindOne(ctx, query).Decode(&existingRate)

	if existingRate != nil {
		update := bson.M{"$set": bson.M{"rate": rate, "updatedat": time.Now()}}

		err := ExchangeRateCollection.FindOneAndUpdate(ctx, query, update).Decode(&existingRate)

		if err != nil {
			return primitive.NilObjectID, errors.New(ErrorCouldNotFindPair)
		}

		_, err = HistoryCollection.InsertOne(ctx, exchangeRate)
		if err != nil {
			return primitive.NilObjectID, errors.New(ErrorCouldNotCreateExchangeRatePair)
		}

		oid := existingRate.ID
		return oid, nil
	}

	result, err := ExchangeRateCollection.InsertOne(ctx, exchangeRate)
	if err != nil {
		return primitive.NilObjectID, errors.New(ErrorCouldNotCreateExchangeRatePair)
	}
	_, err = HistoryCollection.InsertOne(ctx, exchangeRate)
	if err != nil {
		return primitive.NilObjectID, errors.New(ErrorCouldNotCreateExchangeRatePair)
	}
	oid := result.InsertedID.(primitive.ObjectID)
	return oid, nil
}

func GetMarketRates(fiat, crypto string) error {
	baseURL := "https://min-api.cryptocompare.com/data/price"
	params := url.Values{}
	params.Set("fsym", crypto)
	params.Set("tsyms", fiat)

	requestURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	response, err := http.Get(requestURL)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	// Parse the JSON response
	var data map[string]float64
	json.Unmarshal(body, &data)

	var lastCurrency string
	var lastValue float64

	// Accessing the dynamic key and value
	for currency, value := range data {
		lastCurrency = currency
		lastValue = value
	}

	_, err = Create(lastCurrency, crypto, lastValue)
	if err != nil {
		return errors.New(ErrorCouldNotFindPair)
	}

	return err
}

func GetRate(fiat, crypto string) (*models.ExchangeRate, error) {
	var rate *models.ExchangeRate

	query := bson.M{"fiat": fiat, "crypto": crypto}

	err := ExchangeRateCollection.FindOne(ctx, query).Decode(&rate)

	if err != nil {
		return nil, errors.New(ErrorCouldNotFindPair)
	}

	return rate, nil
}

func GetRates(crypto string) ([]*models.ExchangeRate, error) {
	var rates []*models.ExchangeRate

	result, err := ExchangeRateCollection.Find(ctx, bson.M{ "crypto": crypto })

	if err != nil {
		return nil, errors.New(ErrorCouldNotFindPair)
	}

	err = result.All(ctx, &rates)

	if err != nil {
		return nil, err
	}

	return rates, nil
}

func GetAllRate() ([]*models.ExchangeRate, error) {
	var rates []*models.ExchangeRate
	result, err := ExchangeRateCollection.Find(ctx, bson.M{})

	if err != nil {
		return nil, errors.New(ErrorCouldNotFindPair)
	}

	err = result.All(ctx, &rates)
	if err != nil {
		return nil, errors.New(ErrorFailedToDeserialize)
	}
	return rates, nil
}

func GetHistory(fiat, crypto string) ([]*models.ExchangeRate, error) {
	var rates []*models.ExchangeRate

	startTime := time.Now().Add(-24 * time.Hour)
    endTime := time.Now()

	query := bson.M{"fiat": fiat, "crypto": crypto, "createdat": bson.M{
		"$gte": startTime,
		"$lte": endTime,
	},}

	result, err := HistoryCollection.Find(ctx, query)

	fmt.Println(result)

	if err != nil {
		return nil, errors.New(ErrorCouldNotFindPair)
	}

	err = result.All(ctx, &rates)

	if err != nil {
		return nil, errors.New(ErrorFailedToDeserialize)
	}

	return rates, nil
}
