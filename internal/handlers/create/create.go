package create

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/ganzeroar/test_service/internal/db"
	"github.com/ganzeroar/test_service/internal/handlers/create/services"
	"github.com/gin-gonic/gin"

	"github.com/ganzeroar/test_service/internal/models"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
)

type createPersonRequest struct {
	Name       string `json:"name" validate:"required,alpha"`
	Surname    string `json:"surname" validate:"required,alpha"`
	Patronymic string `json:"patronymic" validate:"omitempty,alpha"`
}

func CreatePersonHandler(c *gin.Context) {
	requestData := createPersonRequest{}

	if err := c.BindJSON(&requestData); err != nil {
		var jsonErr *json.UnmarshalTypeError
		var jsonSyntaxErr *json.SyntaxError

		if errors.As(err, &jsonErr) {
			errorMsg := fmt.Sprintf("Field %s must be %s", jsonErr.Field, jsonErr.Type)
			c.AbortWithStatusJSON(http.StatusBadRequest, errorMsg)
			return
		} else if errors.As(err, &jsonSyntaxErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, "Incorrect JSON")
			return
		} else {
			log.WithFields(log.Fields{"err": err}).Warn("CreatePersonHandler - else")
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
	}

	if err := validateData(requestData); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	person, err := createPersonObj(requestData)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		log.WithFields(log.Fields{"err": err}).Info("createPersonObj - err")
		return
	}

	savePerson(person)

	c.JSON(http.StatusCreated, &person)
}

func validateData(requestData createPersonRequest) error {
	validate := *validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(requestData)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			if err.Tag() == "required" {
				text := fmt.Sprintf("Field %s is %s", err.Field(), err.Tag())
				return errors.New(text)
			} else if err.Tag() == "alpha" {
				text := fmt.Sprintf("Field %s must be %s", err.Field(), err.Tag())
				return errors.New(text)
			}
		}
	}
	return nil

}

func createPersonObj(requestData createPersonRequest) (models.Person, error) {
	var person models.Person

	person.Name = requestData.Name
	person.Surname = requestData.Surname
	person.Patronymic = requestData.Patronymic

	ageCh := make(chan int8)
	ageErrCh := make(chan error)
	genderCh := make(chan string)
	genderErrCh := make(chan error)
	nationalityCh := make(chan string)
	nationalityErrCh := make(chan error)

	go services.GetAge(requestData.Name, ageCh, ageErrCh)
	go services.GetGender(requestData.Name, genderCh, genderErrCh)
	go services.GetNationality(requestData.Surname, nationalityCh, nationalityErrCh)

	if err := <-ageErrCh; err != nil {
		return person, err
	}
	person.Age = <-ageCh

	if err := <-genderErrCh; err != nil {
		return person, err
	}
	person.Gender = <-genderCh

	if err := <-nationalityErrCh; err != nil {
		return person, err
	}
	person.Nationality = <-nationalityCh

	return person, nil
}

func savePerson(person models.Person) {
	db.DB.Create(&person)
}
