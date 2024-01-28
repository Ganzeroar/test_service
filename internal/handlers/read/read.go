package read

import (
	"errors"
	"net/http"

	"strconv"

	"github.com/ganzeroar/test_service/internal/db"
	"github.com/ganzeroar/test_service/internal/models"
	"github.com/gin-gonic/gin"
)

var persons []models.Person

func ReadPersonHandler(c *gin.Context) {
	requestURLQuery := c.Request.URL.Query()

	if err := getPersons(requestURLQuery); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
	}

	c.JSON(http.StatusOK, &persons)
}

func getPersons(requestURLQuery map[string][]string) error {
	if !urlHasLimitAndOffset(requestURLQuery) {
		return errors.New("url must has limit and offset")
	}
	limit, err := strconv.Atoi(requestURLQuery["limit"][0])
	if err != nil {
		return errors.New("limit must be integer")
	}
	offset, err := strconv.Atoi(requestURLQuery["offset"][0])
	if err != nil {
		return errors.New("offset must be integer")
	}

	if hasFilteringParams(requestURLQuery) {
		query := createFilterQuery(requestURLQuery)
		db.DB.Limit(limit).Offset(offset).Where(query).Find(&persons)
	} else {
		db.DB.Limit(limit).Offset(offset).Find(&persons)
	}
	return nil
}

func urlHasLimitAndOffset(query map[string][]string) bool {
	_, hasLimit := query["limit"]
	_, hasOffset := query["offset"]
	return hasLimit && hasOffset
}

func hasFilteringParams(query map[string][]string) bool {
	possibleKeys := [6]string{
		"name",
		"surname",
		"patronymic",
		"age",
		"gender",
		"nationality",
	}
	for _, v := range possibleKeys {
		if _, ok := query[v]; ok {
			return true
		}
	}
	return false
}

func createFilterQuery(requestURLQuery map[string][]string) map[string]interface{} {
	query := make(map[string]interface{})
	for k, v := range requestURLQuery {
		if k == "limit" || k == "offset" {
			continue
		}
		query[k] = v[0]
	}
	return query
}
