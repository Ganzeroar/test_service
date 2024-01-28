package update

import (
	"net/http"

	"encoding/json"
	"errors"
	"fmt"

	"github.com/ganzeroar/test_service/internal/db"
	"github.com/ganzeroar/test_service/internal/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type UpdatePersonRequest struct {
	Name        string `json:"name" validate:"required,alpha"`
	Surname     string `json:"surname" validate:"required,alpha"`
	Patronymic  string `json:"patronymic" validate:"omitempty,alpha"`
	Age         int8   `json:"age" validate:"omitempty,number"`
	Gender      string `json:"gender" validate:"omitempty,alpha"`
	Nationality string `json:"nationality" validate:"omitempty,alpha"`
}

func UpdatePersonHandler(c *gin.Context) {
	var person models.Person

	personId := c.Param("id")
	if err := db.DB.Where("id = ?", personId).First(&person).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "person not found"})
		return
	}

	requestData := UpdatePersonRequest{}
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
			log.WithFields(log.Fields{"err": err}).Warn("UpdatePersonHandler - else")
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
	}

	updatedPerson := models.Person{
		Name:        requestData.Name,
		Surname:     requestData.Surname,
		Patronymic:  requestData.Patronymic,
		Age:         requestData.Age,
		Gender:      requestData.Gender,
		Nationality: requestData.Nationality,
	}

	db.DB.Model(&person).Updates(&updatedPerson)
	c.JSON(http.StatusOK, gin.H{"data": person})
}
