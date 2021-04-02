package main

import (
	"context"
	"log"
	"os"

	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/Alfred-Walker/AccountSharing/models"
	"github.com/Alfred-Walker/AccountSharing/rds"
)

// EchoLambda proxy instance to path route
var echoLambda *echoadapter.EchoLambda

// init initializes rds connection and routing.
func init() {
	log.Printf("endTimeChecker: Lambda cold start...")

	// initializes timezone
	os.Setenv("TZ", "Asia/Seoul")

	// initializes EchoLambda instance
	echoLambda = initRoute()

	// initializes GORM instance for RDS
	err := rds.InitRDS()

	if err != nil {
		panic(err)
	}
}

func main() {
	lambda.Start(echoLambdaHandler)
}

// echoLambdaHandler receives context and an API Gateway proxy event,
// and ProxyWithContext transforms them into an http.Request object, and sends it to the echo.Echo for routing.
// It returns a proxy response object generated from the ProxyWithContext.
func echoLambdaHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// check and update endTime
	models.GetFirstAccount(rds.RdsGorm).CheckEndTime(rds.RdsGorm)
	
	// If no name is provided in the HTTP request body, throw an error
	return echoLambda.ProxyWithContext(ctx, req)
}
