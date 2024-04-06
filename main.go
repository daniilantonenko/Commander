package main

import (
	"log"
	"net/http"
	"os/exec"
	"sync"

	"commander/docs"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

var mu sync.Mutex
var cmdMap map[int]*exec.Cmd

// @title Commander API
// @version 1.0
// @description This application provides a REST API for managing and executing commands in the form of bash scripts.
func main() {

	config()

	database()

	host := viper.GetString("SERVER_HOST")
	port := viper.GetString("SERVER_PORT")

	cmdMap = make(map[int]*exec.Cmd)

	router := mux.NewRouter()
	router.HandleFunc("/commands", createCommand).Methods("POST")
	router.HandleFunc("/commands", getCommands).Methods("GET")
	router.HandleFunc("/commands/{id}", getCommand).Methods("GET")
	router.HandleFunc("/commands/{id}", stopCommand).Methods("DELETE")

	// programmatically set swagger info
	docs.SwaggerInfo.Host = host + ":" + port

	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://"+host+":"+port+"/swagger/doc.json"), //The url pointing to API definition
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(":"+port, router))

	log.Println("App is runing!")
}

var db *gorm.DB

func database() {
	var err error
	dbuser := viper.GetString("DB_USER")
	dbpassword := viper.GetString("DB_PASS")
	dbname := viper.GetString("DB_NAME")
	dbhost := viper.GetString("DB_HOST")

	dsn := "host=" + dbhost + " user=" + dbuser + " password=" + dbpassword + " dbname=" + dbname + " sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&Command{})

	log.Println("Connected to database")
}

// Reading configuration from file or environment variables.
func config() {
	viper.SetConfigFile(".env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	log.Println("Environment variables loaded")
}
