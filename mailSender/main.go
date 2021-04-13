package main

import (
	"context"
	"fmt"
	"log"
	"time"

	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/Alfred-Walker/AccountSharing/constants"
	s3Adapter "github.com/Alfred-Walker/AccountSharing/s3"
	sesAdapter "github.com/Alfred-Walker/AccountSharing/ses"
)

// EchoLambda proxy instance to path route
var echoLambda *echoadapter.EchoLambda

// init initializes rds connection and routing.
func init() {
	log.Printf("mailSender: Lambda cold start...")
}

func main() {
	lambda.Start(echoLambdaHandler)
}

// echoLambdaHandler receives context and an API Gateway proxy event,
// and ProxyWithContext transforms them into an http.Request object, and sends it to the echo.Echo for routing.
// It returns a proxy response object generated from the ProxyWithContext.
func echoLambdaHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// get last day date
	loc, _ := time.LoadLocation("Asia/Seoul")
	localTime := time.Now().Local().AddDate(0, 0, -1).In(loc)
	lastDay := localTime.Format("2006-01-02") 
	
	// Create a session that gets credential values from ~/.aws/credentials
	// and the default region from ~/.aws/config
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create an S3 service client
	s3Svc := s3.New(sess)

	buffer, err := s3Adapter.GetFileBuffer(s3Svc, constants.OCCUPATION_LOG_BUCKET, lastDay)

	if err != nil {
		log.Printf("mailSender: Failed to get a buffer from S3")
	} else {
		// Create anSES service client
		mailSvc := ses.New(sess)

		// Set subject for last day's access log
		subject := fmt.Sprintf("[Account-Sharing] 접속 로그 (%s)", lastDay)

		// Send buffer contents to mail
		sesAdapter.SendEmail(mailSvc, buffer, subject)
	}

	// If no name is provided in the HTTP request body, throw an error
	return echoLambda.ProxyWithContext(ctx, req)
}
