package log

import (
	"errors"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"

	"github.com/Alfred-Walker/AccountSharing/constants"
)

// GetCsv receives S3 service client, bucket name, and key.
// It returns a buffer of file corresponds to key.
func SendEmail(svc *ses.SES, textBody string, subject string) error {
	// The subject line for the email.
	Subject := subject

	//The email body for recipients with non-HTML email clients.
	TextBody := textBody

	// The HTML body for the email.
	HtmlBody := 
		"<div>" + 
		textBody +
		"</div>" +
		"<p>This email was sent with " +
		"<a href='https://aws.amazon.com/ses/'>Amazon SES</a> using the " +
		"<a href='https://aws.amazon.com/sdk-for-go/'>AWS SDK for Go</a>.</p>"

	// The character encoding for the email.
	CharSet := "UTF-8"

	if svc == nil {
		return errors.New("S3: service client cannot be nil.")
	}

	// Assemble the email.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(constants.EMAIL_RECIPIENT),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(HtmlBody),
				},
				Text: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(TextBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(Subject),
			},
		},
		Source: aws.String(constants.EMAIL_SENDER),
	}

	// Attempt to send the email.
	_, err := svc.SendEmail(input)

	// Display error messages if they occur.
	if err != nil {
		log.Printf("SES: Error at %v", err.Error())
		return err
	}

	log.Println("Email Sent to address: " + constants.EMAIL_RECIPIENT)

	return nil
}
