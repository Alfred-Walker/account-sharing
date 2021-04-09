package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"

	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/nyaruka/phonenumbers"

	"github.com/Alfred-Walker/AccountSharing/constants"
	"github.com/Alfred-Walker/AccountSharing/models"
	"github.com/Alfred-Walker/AccountSharing/rds"
	sqsAdapter "github.com/Alfred-Walker/AccountSharing/sqs"
)

var e *echo.Echo

type NewOccupier struct {
	Name    string `json:"occupier"`
	EndTime string `json:"endtime"`
}

type PhoneNumber struct {
	Number  string `json:"number"`
	Country string `json:"country"`
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

	e.POST(constants.API_GATEWAY_RESOURCE_NAME+"/alarm", handleSetAlarm)
	e.POST(constants.API_GATEWAY_RESOURCE_NAME+"/alarm/", handleSetAlarm)

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

// handleSetAlarm is a handler for SNS alarm registration.
// Receives an echo HTTP request context.
// It returns the result of alarm registration with statud code.
func handleSetAlarm(c echo.Context) error {
	generatePhoneNum := func(number string, countryCode string) (string, error) {
		parsed, err := phonenumbers.Parse(number, strings.ToUpper(countryCode))

		if err != nil {
			return "", errors.New("phoneNumGenerator: number is invalid.")
		}

		return strconv.Itoa(int(parsed.GetCountryCode())) + strconv.FormatUint(parsed.GetNationalNumber(), 10), nil
	}

	// data binding
	data := PhoneNumber{}
	if err := c.Bind(&data); err != nil {
		log.Printf("accountTaker: error at binding request data")
		return err
	}

	// generate phone num from the request data
	generated, err := generatePhoneNum(data.Number, data.Country)

	if err != nil {
		log.Printf(err.Error())
		return err
	}

	// Create a session that gets credential values from ~/.aws/credentials
	// and the default region from ~/.aws/config
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create an SQS service client
	svc := sqs.New(sess)
	queueName := constants.SQS_QUEUE
	queueUrl, err := sqsAdapter.GetQueueUrl(svc, &queueName)

	if err != nil {
		log.Println("SQS: Failed to get queue url")
		return err
	}

	// Send generated number to the queue
	err = sqsAdapter.SendMsg(svc, queueUrl, &generated)

	if err != nil {
		log.Println("SQS: Failed to send the message.")
		return err
	}

	return c.JSON(http.StatusOK, "Alarm Registered!")
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
