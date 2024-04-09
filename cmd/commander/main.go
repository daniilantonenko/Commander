package main

import (
	"log"

	"commander/config"
	"commander/internal/server"
	database "commander/pkg/db"
)

// @title Commander API
// @version 1.0
// @description This application provides a REST API for managing and executing commands in the form of bash scripts.
func main() {

	cfg := config.ReadConfig()

	db := database.NewDB(&cfg)

	s := server.NewServer(&cfg, db)
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}

	// programmatically set swagger info
	/*docs.SwaggerInfo.Host = host + ":" + port

	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://"+host+":"+port+"/swagger/doc.json"), //The url pointing to API definition
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)*/

}
