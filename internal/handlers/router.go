package handlers

import (
	"github.com/ganzeroar/test_service/internal/handlers/create"
	"github.com/ganzeroar/test_service/internal/handlers/delete"
	"github.com/ganzeroar/test_service/internal/handlers/read"
	"github.com/ganzeroar/test_service/internal/handlers/update"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(engine *gin.Engine) {
	engine.POST("/person", create.CreatePersonHandler)
	engine.GET("/person", read.ReadPersonHandler)
	engine.DELETE("/person/:id", delete.DeletePersonHandler)
	engine.PATCH("/person/:id", update.UpdatePersonHandler)
}
