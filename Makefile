.PHONY: zip
zip: 
	export GOOS=linux GOARCH=amd64 CGO_ENABLED=0
	go build -o ./accountTaker/main ./accountTaker/*.go
	go build -o ./endTimeChecker/main ./endTimeChecker/*.go
	build-lambda-zip.exe -output ./accountTaker/main.zip ./accountTaker/main
	build-lambda-zip.exe -output ./endTimeChecker/main.zip ./endTimeChecker/main

.PHONY: createLambdaAccountTaker
createLambdaAccountTaker: 
	aws lambda create-function --function-name accountTaker --runtime go1.x --zip-file fileb://accountTaker/main.zip --handler main --role "YOUR_ROLE"

.PHONY: createLambdaEndTimeChecker
createLambdaEndTimeChecker: 
	aws lambda create-function --function-name endTimeChecker --runtime go1.x --zip-file fileb://endTimeChecker/main.zip --handler main --role "YOUR_ROLE"

.PHONY: deploy
deploy: zip
	aws lambda update-function-code --region ap-northeast-2 --function-name accountTaker --zip-file fileb://accountTaker/main.zip
	aws lambda update-function-code --region ap-northeast-2 --function-name endTimeChecker --zip-file fileb://endTimeChecker/main.zip