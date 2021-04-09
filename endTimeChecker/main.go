package main

import (
	"context"
	"log"
	"os"

	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"

	"github.com/Alfred-Walker/AccountSharing/constants"
	"github.com/Alfred-Walker/AccountSharing/models"
	"github.com/Alfred-Walker/AccountSharing/rds"
	sqsAdapter "github.com/Alfred-Walker/AccountSharing/sqs"
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
	account := models.GetFirstAccount(rds.RdsGorm).CheckEndTime(rds.RdsGorm)

	if account != nil && account.EndTime == "" {
		// send alarm to users who wait their turn
		err := sendAlarm()

		if err != nil {
			log.Printf("endTimeChecker: Failed to get all messages from SQS")
			log.Printf(err.Error())
		}
	}

	// If no name is provided in the HTTP request body, throw an error
	return echoLambda.ProxyWithContext(ctx, req)
}

func sendAlarm() error {
	// Create a session that gets credential values from ~/.aws/credentials
	// and the default region from ~/.aws/config
	sqsSess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create a Session with a custom region for sns
	snsSess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Region: aws.String(constants.SNS_SESSION_REGION)},
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create an SQS/SNS service client
	sqsClient := sqs.New(sqsSess)
	snsClient := sns.New(snsSess)

	// loop until get no messages
	for {
		// Get list of *sqs.Message (max 10 messages at once)
		messageList, err := sqsAdapter.GetMessageList(sqsClient, aws.Int64(30))

		if err != nil {
			log.Println("SQS: Failed to get message list")
			return err
		}

		if err != nil {
			log.Println("SQS: Failed to get queue url")
			return err
		}

		if len(messageList) == 0 {
			log.Println("SQS: No more messages in the queue")
			break
		}

		// send alarm & delete messages
		for _, message := range messageList {

			// Check your region if you get below error:
			// Invalid parameter: PhoneNumber Reason: +XXXXXXX is not valid to publish to make.
			input := &sns.PublishInput{
				Message:     aws.String(constants.SNS_ALARM_MESSAGE),
				PhoneNumber: aws.String(*message.Body),
			}

			_, err := snsClient.Publish(input)
			if err != nil {
				// delete message & go to next loop
				log.Println("SNS: Publish error:", err)

				sqsAdapter.DeleteMsg(sqsClient, message.ReceiptHandle)
				log.Println("SNS: Message deleted:", err)
				continue
			}

			log.Printf("alarm send - %s", *message.Body)
			sqsAdapter.DeleteMsg(sqsClient, message.ReceiptHandle)
		}
	}

	log.Println("all notification has been sent!")

	return nil
}
