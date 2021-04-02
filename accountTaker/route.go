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

type NewOccupier struct {
	Name string `json:"occupier"`
	EndTime string `json:"endtime"`
}

// initRoute creates EchoLambda object to setup handlers for GET/POST routes
// It returns the initialized instance of the EchoLambda object.
func initRoute() *echoadapter.EchoLambda {
	log.Printf("accountTaker: initializes Echo for route...")

	// creates *echo.Echo object to initialize EchoLambda object to return.
	e = echo.New()

	// CORS default
	e.Use(middleware.CORS())

	// registers GET/POST routes and handlers
	// Q. Why API_GATEWAY_RESOURCE_NAME is required?
	// A. to avoid 403 error, Missing Authentication Token
	// https://aws.amazon.com/ko/premiumsupport/knowledge-center/api-gateway-authentication-token-errors/
	e.GET(constants.API_GATEWAY_RESOURCE_NAME, handleRoot)
	e.GET(constants.API_GATEWAY_RESOURCE_NAME+"/", handleRoot)

	e.POST(constants.API_GATEWAY_RESOURCE_NAME, handleOccupyAccount)
	e.POST(constants.API_GATEWAY_RESOURCE_NAME+"/", handleOccupyAccount)

	e.POST(constants.API_GATEWAY_RESOURCE_NAME+"/release", handleReleaseAccount)
	e.POST(constants.API_GATEWAY_RESOURCE_NAME+"/release/", handleReleaseAccount)

	// creates and returns a new instance of the EchoLambda object
	return echoadapter.New(e)
}

// handleRoot is a handler for root path("/").
// Receives an echo HTTP request context.
// It returns current account usage info JSON response with statud code.
func handleRoot(c echo.Context) error {
	account := models.GetFirstAccount(rds.RdsGorm)
	return c.JSON(http.StatusOK, account)
}

// handleOccupyAccount registers new occupier info to the rds.
// Receives an echo HTTP request context.
// It returns new occupier account JSON response with statud code or returns an error.
func handleOccupyAccount(c echo.Context) error {
	data := NewOccupier{}
	if err := c.Bind(&data); err != nil {
		log.Printf("accountTaker: error at binding NewOccupier")
		return err
	}

	account, err := models.OccupyAccount(rds.RdsGorm, data.Name, data.EndTime)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, account)
}

// handleReleaseAccount reset current occupier info in the rds.
// Receives an echo HTTP request context.
// It returns account JSON response with statud code after release.
func handleReleaseAccount(c echo.Context) error {
	log.Printf("accountTaker: ReleaseAccount Handler")

	account := models.ReleaseAccount(rds.RdsGorm)

	return c.JSON(http.StatusOK, account)
}
