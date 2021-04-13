package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/Alfred-Walker/AccountSharing/constants"
	s3Adapter "github.com/Alfred-Walker/AccountSharing/s3"
	sesAdapter "github.com/Alfred-Walker/AccountSharing/ses"
)

// init initializes rds connection and routing.
func init() {
	log.Printf("mailSender: Lambda cold start...")
}

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest(ctx context.Context) (error) {
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
		if buffer != nil && buffer.String() == "" {
			// send the list of access user
			sesAdapter.SendEmail(mailSvc, buffer.String(), subject)
		} else {
			// no access user on the last day
			sesAdapter.SendEmail(mailSvc, "아무도 접속하지 않았습니다...", subject)
		}
	}
	return nil
}