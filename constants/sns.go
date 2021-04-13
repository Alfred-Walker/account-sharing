package constants

// IMPORTANT: DO NOT UPLOAD ID OR PASSWORD TO GITHUB OR OTHER REPOSITORIES.
// Use your environment variables or external file to handle private values.
// e.g. BUCKET := os.Getenv("BUCKET")

// You can choose region that supports SNS message transfer at below.
// https://docs.aws.amazon.com/na_en/sns/latest/dg/sns-supported-regions-countries.html

// Here we use tokyo region to send message
const SNS_SESSION_REGION = "ap-northeast-1"

// Alarm message to send via sns
// Be careful of using UCS-2 character length limit  (70)
// https://docs.aws.amazon.com/sns/latest/dg/sms_publish-to-phone.html
const SNS_ALARM_MESSAGE = "YOUR_ALARM_MESSAGE"