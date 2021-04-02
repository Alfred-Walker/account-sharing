# account-sharing


<!-- ABOUT THE PROJECT -->
## About

[**account-sharing**](https://github.com/Alfred-Walker/account-sharing/) is a serverless project example for account sharing based on cooperative schedulilng.

It contains:
* AWS Lambda Function to update account sharing status based on AWS API Gateway's REST API supports
  - AWS Lambda proxy routing on Echo framework context
  - Simple MySQL RDS model for sharing account (single row)
  - Makefile supports (build & deploy & lambda function creation)
  

<br/>

## Main End-points
* accountTaker <br/>
  * just replace all 'study' in path with 'question' to access end-points for `Question` instead `Study`.

|  HTTP |  Path |  Method |  Purpose |
| --- | --- | --- | --- |
|**GET** |/{API_GATEWAY_RESOURCE_NAME}|Read|Retrieve a current account occupier and end time|
|**POST** |/{API_GATEWAY_RESOURCE_NAME}|Update|Update sharing account occupation info (occupier's name, end time)|
|**POST** |/{API_GATEWAY_RESOURCE_NAME}/release|Update|Initialize occupier and endtime info|

* endTimeChecker (to be described later. not necessary at all in current stage.) <br/>

|  HTTP |  Path |  Method |  Purpose |
| --- | --- | --- | --- |
<br/>


## Getting Started
<!-- GETTING STARTED -->

### AWS RDS MySQL
* Configure AWS RDS MySQL DB & Create table for sharing accounts (Required)
  * Be sure that sync DB authentification info with settings in a file `constants\rds.go` that contains user name, password, DB name, and host info.
  * You can create table using sample query written in `models\accounts.go`.
  * Be careful not to upload the rds.go file to public repositories.
```sh
const DB_USERNAME = "YOUR_DB_USERNAME"
const DB_PASSWORD = "YOUR_DB_PASSWORD"
const DB_NAME = "accountSharing"
const DB_HOST = "my-practice-db.c7hdfqvga72s.ap-northeast-2.rds.amazonaws.com"
const DB_PORT = "3306"

### AWS API Gateway & Lambda
```
* Create Lambda function using Makefile or AWS-CLI
  * Be sure to use your own AWS role defined.
```sh
	aws lambda create-function --function-name accountTaker --runtime go1.x --zip-file fileb://accountTaker/main.zip --handler main --role "YOUR_ROLE"
	aws lambda create-function --function-name endTimeChecker --runtime go1.x --zip-file fileb://endTimeChecker/main.zip --handler main --role "YOUR_ROLE"
```
* Configure AWS API Gateway API with methods (accountTaker)
  * Create API resources with methods connected to lambda function, accountTaker.
  * Activate API Gateway CORS when creating API resources.
  * Connects lambda function to API methods. (and checks Lambda proxy, too.)
  * Don't forget to deploy your API before testing.
```sh
/
  /accountsharing
    ANY
    OPTIONS
      /{proxy+}
        ANY
        OPTIONS
```

* Register AWS CloudWatch event rule for lambda function, endTimeChecker
  * Add EventBridge(CloudWatch Events) trigger to the endTimeChecker lambda function.
  * For cron job, generate rules written in regex
    * e.g. rate(1 minute)

### Golang Project
```sh
make deploy
```
(If you don't want to use Makefile, you can build, generate zip file, and deploy manually by typing each commands in Makefile.)
* In case of build-lambda-zip.exe, please see below AWS official guidelines.
* (https://docs.aws.amazon.com/ko_kr/lambda/latest/dg/golang-package.html)



## Example
* [**Account Sharing Example**](https://alfred-walker.github.io/account-sharing-app/) (Korean, Asia/Seoul) <br/>


<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to be learn, inspire, and create. Any contributions you make are **greatly appreciated**.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request


<!-- CONTACT -->
## Contact

studio.alfred.walker@gmail.com

<!-- ACKNOWLEDGEMENTS -->
## Acknowledgements
[othneildrew/Best-README-Template](https://github.com/othneildrew/Best-README-Template)
[Img Shields](https://shields.io)


<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->

