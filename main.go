package main

import (
	"context"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	service, err := sheets.NewService(context.TODO(), option.WithCredentialsFile("./creds.json"))
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

	e.Logger.Fatal(e.Start(":3000"))
}
