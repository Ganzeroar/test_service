package delete

import (
	"net/http"

	"github.com/ganzeroar/test_service/internal/db"
	"github.com/ganzeroar/test_service/internal/models"
	"github.com/gin-gonic/gin"
)

func DeletePersonHandler(c *gin.Context) {
	var person models.Person

	personId := c.Param("id")
	if err := db.DB.Where("id = ?", personId).First(&person).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "person not found"})
		return
	}
	db.DB.Delete(&person)

	c.Status(http.StatusNoContent)
}
