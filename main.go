package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const (
	// SpreadsheetID for the Animal Crossing data spreadsheet - https://docs.google.com/spreadsheets/d/13d_LAJPlxMa_DubPTuirkIV4DERBMXbrWQsmSh8ReK4
	SpreadsheetID = "13d_LAJPlxMa_DubPTuirkIV4DERBMXbrWQsmSh8ReK4"
)

var (
	// Cache stores Google Sheets API responses in memory
	Cache = make(map[string]interface{})
)

func main() {
	if err := initViper(); err != nil {
		log.Fatal(err)
	}

	service, err := sheets.NewService(context.TODO(), option.WithCredentialsJSON(getCredsJSON()))
	if err != nil {
		log.Fatalf("failed to setup sheets service: %v", err)
	}

	e := echo.New()
	e.Debug = true

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{DisableStackAll: true}))
	e.Use(middleware.Secure())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "REQUEST: method=${method}, status=${status}, uri=${uri}, latency=${latency_human}\n",
	}))

	e.GET("/search/:category", SearchHandler(service))
	e.GET("/seasonal/:hemisphere/:category", SeasonalHandler(service))

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", viper.GetInt("port"))))
}

func initViper() error {
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// Prepare for Environment variables
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Defaults
	viper.SetDefault("port", 3000)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Config file was found but another error was produced
			return errors.Wrap(err, "viper error")
		}
	}

	return nil
}

func getCredsJSON() []byte {
	privateKey := strings.ReplaceAll(viper.GetString("google.private_key"), "\n", "\\n")
	raw := fmt.Sprintf(`{
		"type": "service_account",
		"project_id": "cohesive-armor-129523",
		"private_key_id": "%s",
		"private_key": "%s",
		"client_email": "%s",
		"client_id": "%s",
		"auth_uri": "https://accounts.google.com/o/oauth2/auth",
		"token_uri": "https://oauth2.googleapis.com/token",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/animal-crossing-sheets%%%%40cohesive-armor-129523.iam.gserviceaccount.com"
	  }`,
		viper.GetString("google.private_key_id"),
		privateKey,
		viper.GetString("google.client_email"),
		viper.GetString("google.client_id"),
	)

	return []byte(raw)
}
