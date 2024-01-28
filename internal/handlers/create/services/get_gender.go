package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type genderizeResponse struct {
	Count       int     `json:"count"`
	Name        string  `json:"name"`
	Gender      string  `json:"gender"`
	Probability float32 `json:"probability"`
}

func GetGender(name string, ch chan string, chErr chan error) {
	requestURL := fmt.Sprintf("https://api.genderize.io/?name=%s", name)

	res, err := http.Get(requestURL)
	if err != nil {
		log.WithFields(log.Fields{"err": err, "res": res}).Warn("GetGender - http.Get")
		chErr <- errors.New("error with gender")
		ch <- ""
		return
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.WithFields(log.Fields{"err": err, "resBody": res}).Warn("GetGender - io.ReadAll")
		chErr <- errors.New("error with gender")
		ch <- ""
		return
	}

	var result genderizeResponse
	err = json.Unmarshal(resBody, &result)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Fatal("GetGender - json.Unmarshal")
		chErr <- errors.New("error with gender")
		ch <- ""
		return
	}

	chErr <- nil
	ch <- result.Gender
}
