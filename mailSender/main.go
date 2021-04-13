package main

import (
	"bufio"
	"context"
	"encoding/csv"
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

	// Be careful of file extension
	buffer, err := s3Adapter.GetFileBuffer(s3Svc, constants.OCCUPATION_LOG_BUCKET, lastDay + ".csv")
	var textBody string

	if err != nil {
		log.Printf("mailSender: Failed to get a buffer from S3")
	} else {
		// Create anSES service client
		mailSvc := ses.New(sess)

		// Set subject for last day's access log
		subject := fmt.Sprintf("[Account-Sharing] 접속 로그 (%s)", lastDay)

		// Send buffer contents to mail
		if buffer != nil && buffer.String() != "" {
			// csv reader setup
			rdr := csv.NewReader(bufio.NewReader(buffer))
 
			// read all csv contents
			rows, _ := rdr.ReadAll()

			// simple html tags
			textBody = "<table>"

			for i, row := range rows {
				if i != 0 {
					textBody += "<tr>"
				}
				
				for j := range row {
					if i == 0 {
						textBody += "<th>"
					} else {
						textBody += "<td>"
					}

					textBody += rows[i][j]
					// fmt.Printf("%s ", rows[i][j])

					if i == 0 {
						textBody += "</th>"
					} else {
						textBody += "</td>"
					}
				}

				if i != 0 {
					textBody += "</tr>"
				}
			}

			textBody += "</table>"

		} else {
			// no access user on the last day
			textBody = "아무도 접속하지 않았습니다..."
		}

		// send email
		sesAdapter.SendEmail(mailSvc, textBody, subject)
	}
	return nil
}