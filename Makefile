.PHONY: zip
zip: 
	export GOOS=linux GOARCH=amd64 CGO_ENABLED=0
	go build -o ./accountTaker/main ./accountTaker/*.go
	go build -o ./endTimeChecker/main ./endTimeChecker/*.go
	build-lambda-zip.exe -output ./accountTaker/main.zip ./accountTaker/main
	build-lambda-zip.exe -output ./endTimeChecker/main.zip ./endTimeChecker/main

.PHONY: createLambdaAccountTaker
createLambdaAccountTaker: 
	aws lambda create-function --function-name accountTaker --runtime go1.x --zip-file fileb://accountTaker/main.zip --handler main --role YOUR_ROLE

.PHONY: createLambdaEndTimeChecker
createLambdaEndTimeChecker: 
	aws lambda create-function --function-name endTimeChecker --runtime go1.x --zip-file fileb://endTimeChecker/main.zip --handler main --role YOUR_ROLE

.PHONY: createAlarmQueue
createAlarmQueue: 
	aws sqs create-queue --queue-name YOUR_QUEUE_NAME --region ap-northeast-2

# A FIFO queue name must end with the .fifo suffix.
.PHONY: createFifoAlarmQueue
createFifoAlarmQueue: 
	aws sqs create-queue --queue-name YOUR_QUEUE_NAME.fifo --region ap-northeast-2 --attributes FifoQueue=true

.PHONY: purgeAlarmQueue
purgeAlarmQueue: 
	aws sqs purge-queue --queue-url YOUR_QUEUE_URL

.PHONY: createS3Bucket
createS3Bucket: 
	aws s3api create-bucket --bucket YOUR_BUCKET_NAME --region ap-northeast-2 --create-bucket-configuration LocationConstraint=ap-northeast-2

# Be sure that your region support to send SNS message
.PHONY: testSendingMsgToPhone
testSendingMsgToPhone:
	aws sns publish --phone-number +PHONE_NUMBER --message "published" --region ap-northeast-1

.PHONY: deploy
deploy: zip
	aws lambda update-function-code --region ap-northeast-2 --function-name accountTaker --zip-file fileb://accountTaker/main.zip
	aws lambda update-function-code --region ap-northeast-2 --function-name endTimeChecker --zip-file fileb://endTimeChecker/main.zip