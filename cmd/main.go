package main

import (
	"os"

	"github.com/ganzeroar/test_service/internal/db"
	"github.com/ganzeroar/test_service/internal/handlers"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	file, err := os.OpenFile("./logs/logs.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	viper.SetConfigFile("./envs/.env")
	viper.ReadInConfig()

	port := viper.Get("PORT").(string)
	dbUrl := viper.Get("DB_URL").(string)

	engine := gin.Default()
	db.Init(dbUrl)

	handlers.RegisterRoutes(engine)

	engine.Run(port)
}
