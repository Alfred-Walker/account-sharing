.PHONY: zip
zip: 
	export GOOS=linux GOARCH=amd64 CGO_ENABLED=0
	go build -o main *.go
	build-lambda-zip.exe -output main.zip main

.PHONY: createLambdaFunction
createLambdaFunction: 
	aws lambda create-function --function-name shareApplication --runtime go1.x --zip-file fileb://main.zip --handler main --role arn:aws:iam::494131715847:role/aws_study_admin

.PHONY: deploy
deploy: zip
	aws lambda update-function-code --region ap-northeast-2 --function-name noswagg --zip-file fileb://main.zip