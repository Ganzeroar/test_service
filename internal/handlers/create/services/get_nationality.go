package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type nationalizeResponse struct {
	Count   int    `json:"count"`
	Name    string `json:"name"`
	Country []struct {
		CountryId   string  `json:"country_id"`
		Probability float32 `json:"probability"`
	} `json:"country"`
}

func GetNationality(surname string, ch chan string, chErr chan error) {
	requestURL := fmt.Sprintf("https://api.nationalize.io/?name=%s", surname)

	res, err := http.Get(requestURL)
	if err != nil {
		log.WithFields(log.Fields{"err": err, "res": res}).Warn("GetNationality - http.Get")
		chErr <- errors.New("error with nationality")
		ch <- ""
		return
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.WithFields(log.Fields{"err": err, "resBody": res}).Warn("GetNationality - io.ReadAll")
		chErr <- errors.New("error with nationality")
		ch <- ""
		return
	}

	var result nationalizeResponse
	err = json.Unmarshal(resBody, &result)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Fatal("GetNationality - json.Unmarshal")
		chErr <- errors.New("error with nationality")
		ch <- ""
		return
	}

	indexOfMaxProbabilityCountry := getindexOfMaxProbabilityCountry(result)

	chErr <- nil
	ch <- result.Country[indexOfMaxProbabilityCountry].CountryId
}

func getindexOfMaxProbabilityCountry(result nationalizeResponse) int {
	var indexOfMaxProbabilityCountry int = 0
	var maxProbability float32 = result.Country[0].Probability
	for index, element := range result.Country {
		if element.Probability > maxProbability {
			indexOfMaxProbabilityCountry = index
		}
	}
	return indexOfMaxProbabilityCountry

}
