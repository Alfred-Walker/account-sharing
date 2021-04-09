package main

import (
	"log"
	"net/http"

	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/Alfred-Walker/AccountSharing/constants"
	"github.com/Alfred-Walker/AccountSharing/models"
	"github.com/Alfred-Walker/AccountSharing/rds"
)

var e *echo.Echo

// initRoute creates EchoLambda object to setup handlers for GET/POST routes
// It returns the initialized instance of the EchoLambda object.
func initRoute() *echoadapter.EchoLambda {
	log.Printf("endTimeChecker: initializes Echo for route...")
	e = echo.New()

	// CORS default
	e.Use(middleware.CORS())

	// registers GET/POST routes and handlers
	// Q. Why API_GATEWAY_RESOURCE_NAME is required?
	// A. to avoid 403 error, Missing Authentication Token
	// https://aws.amazon.com/ko/premiumsupport/knowledge-center/api-gateway-authentication-token-errors/
	e.GET(constants.API_GATEWAY_RESOURCE_NAME, handleRoot)
	e.GET(constants.API_GATEWAY_RESOURCE_NAME+"/", handleRoot)

	return echoadapter.New(e)
}

// handleRoot is a handler for root path("/").
// Receives an echo HTTP request context.
// It returns account JSON response with statud code.
func handleRoot(c echo.Context) error {
	// checks whether current account endtime is passed or not
	account := models.GetFirstAccount(rds.RdsGorm).CheckEndTime(rds.RdsGorm)

	if account.EndTime == "" {
		// send alarm to users who wait their turn
		err := sendAlarm()

		if err != nil {
			log.Printf("endTimeChecker: Failed to get all messages from SQS")
			log.Printf(err.Error())
		}
	}

	return c.JSON(http.StatusOK, account)
}
