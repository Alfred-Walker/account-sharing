package log

import (
	"bytes"
	"encoding/csv"
	"errors"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

// getListObjectsV2Output receives S3 service client, bucket name, and key.
// It returns a ListObjectsV2Output object when no error is detected.
func getListObjectsV2Output(svc *s3.S3, bucket string, key string) (*s3.ListObjectsV2Output, error) {
	if svc == nil {
		return nil, errors.New("S3: service client cannot be nil.")
	}

	listObjects, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(key),
	})

	if err != nil {
		log.Printf("S3: error %v at ListObjectsV2", err.Error())
		return nil, err
	}

	return listObjects, nil
}

// GetCsv receives S3 service client, bucket name, and key.
// It returns a buffer of file corresponds to key.
func GetFileBuffer(svc *s3.S3, bucket string, key string) (*bytes.Buffer, error) {
	if svc == nil {
		return nil, errors.New("S3: service client cannot be nil.")
	}

	listObjects, err := getListObjectsV2Output(svc, bucket, key)

	if err != nil {
		log.Printf("S3: error %v at GetFile", err.Error())
		return nil, err
	}

	buffer := new(bytes.Buffer)

	if *listObjects.KeyCount == 0 {
		// no file to return
		return nil, nil
	} else {
		// Get object
		obj, err := svc.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})

		defer obj.Body.Close()

		if err != nil {
			log.Printf("S3: error %v at GetFile", err.Error())
			return nil, err
		}

		buffer.ReadFrom(obj.Body)
		return buffer, nil
	}
}


// UpdateCsv receives S3 service client, bucket name, key, record to add, and header strings of record.
// It adds or updates a log file of csv format.
// It returns an error when it is not nil.
func UpdateCsv(svc *s3.S3, bucket string, key string, record []string, header []string) error {
	if svc == nil {
		return errors.New("S3: service client cannot be nil.")
	}

	listObjects, err := getListObjectsV2Output(svc, bucket, key)

	if err != nil {
		log.Printf("S3: error %v at UpdateCsv", err.Error())
		return err
	}

	buffer := new(bytes.Buffer)

	if *listObjects.KeyCount == 0 {
		bufferWriter := csv.NewWriter(buffer)
		bufferWriter.Write(header) // generates csv geader
		bufferWriter.Write(record) // converts array of string to comma seperated values for 1 row.
		bufferWriter.Flush()       // writes the csv data to the buffered data (buffer.Bytes())
	} else {
		// Get object
		obj, err := svc.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})

		defer obj.Body.Close()

		if err != nil {
			log.Printf("S3: error %v at UpdateCsv", err.Error())
			return err
		}

		buffer.ReadFrom(obj.Body)
		bufferWriter := csv.NewWriter(buffer)
		bufferWriter.Write(record) // converts array of string to comma seperated values for 1 row.
		bufferWriter.Flush()       // writes the csv data to the buffered data (buffer.Bytes())
	}

	_, err = svc.PutObject(&s3.PutObjectInput{
		Body:                 bytes.NewReader(buffer.Bytes()),
		Bucket:               aws.String(bucket),
		Key:                  aws.String(key),
		ACL:                  aws.String("private"),
		ServerSideEncryption: aws.String("AES256"),
	})

	return err
}
