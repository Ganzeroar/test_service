package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type agifyResponse struct {
	Count int    `json:"count"`
	Name  string `json:"name"`
	Age   int8   `json:"age"`
}

func GetAge(name string, ch chan int8, chErr chan error) {
	requestURL := fmt.Sprintf("https://api.agify.io/?name=%s", name)

	res, err := http.Get(requestURL)
	if err != nil {
		log.WithFields(log.Fields{"err": err, "res": res}).Warn("GetAge - http.Get")
		chErr <- errors.New("error with agify")
		ch <- 0
		return
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.WithFields(log.Fields{"err": err, "resBody": res}).Warn("GetAge - io.ReadAll")
		chErr <- errors.New("error with agify")
		ch <- 0
		return
	}

	var result agifyResponse
	err = json.Unmarshal(resBody, &result)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Fatal("GetAge - json.Unmarshal")
		chErr <- errors.New("error with agify")
		ch <- 0
		return
	}

	chErr <- nil
	ch <- result.Age
}
