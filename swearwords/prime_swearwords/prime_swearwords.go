package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
)

func HandleRequest(ctx context.Context, event *any) (any, error) {
	if event == nil {
		return nil, fmt.Errorf("received nil event")
	}

	fmt.Printf("%+v\n", event)
	fmt.Println("Prime_Swearwords called")

	session := session.Must(session.NewSession(&aws.Config{
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second, // Increase the timeout to 30 seconds
		},
	}))

	db := dynamodb.New(session)

	errPrime := EnsureTablePrimed(db, session)
	if errPrime != nil {
		fmt.Printf("Got error priming table")
		return nil, errPrime
	}

	return nil, nil
}

func EnsureTablePrimed(db *dynamodb.DynamoDB, session *session.Session) error {
	isPrimed, err := isTablePrimed(db)
	if err != nil {
		fmt.Printf("Got error checking if table is primed")
		return err
	}

	if isPrimed {
		return nil
	}

	bucket := aws.String(os.Getenv("BUCKET_NAME"))
	key := aws.String(os.Getenv("BUCKET_KEY"))
	fmt.Printf("getting file content from bucket %s and key %s\n", *bucket, *key)
	content, readErr := getFileContent(bucket, key, session)
	if readErr != nil {
		fmt.Printf("Got no file content")
		return readErr
	}
	fmt.Println("Got file content")

	// Split content into lines
	lines := strings.Split(*content, "\n")
	fmt.Printf("Adding  %d words to table\n", len(lines))

	var items []*dynamodb.WriteRequest
	for _, line := range lines {
		item := &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: map[string]*dynamodb.AttributeValue{
					"word": {
						S: aws.String(line),
					},
				},
			},
		}
		items = append(items, item)
	}

	err = BatchWriteItems(db, items)
	if err != nil {
		return err
	}
	fmt.Printf("Added  %d words to table\n", len(lines))

	_, err = db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("SWEARWORDS_TABLE_NAME")),
		Item: map[string]*dynamodb.AttributeValue{
			"word": {
				S: aws.String("###table_primed###"),
			},
		},
	})
	if err != nil {
		fmt.Printf("Got error calling PutItem: %s", err)
		return err
	}

	return nil
}
func BatchWriteItems(db *dynamodb.DynamoDB, items []*dynamodb.WriteRequest) error {
	SWEARWORDS_TABLE_NAME := os.Getenv("SWEARWORDS_TABLE_NAME")
	var batchItems []*dynamodb.WriteRequest
	for i, item := range items {
		batchItems = append(batchItems, item)
		if len(batchItems) == 25 || i == len(items)-1 {
			batchInput := &dynamodb.BatchWriteItemInput{
				RequestItems: map[string][]*dynamodb.WriteRequest{
					SWEARWORDS_TABLE_NAME: batchItems,
				},
			}
			_, err := db.BatchWriteItem(batchInput)
			if err != nil {
				fmt.Printf("Got error calling BatchWriteItem: %s", err)
				return err
			}
			batchItems = nil
		}
	}
	return nil
}

func isTablePrimed(db *dynamodb.DynamoDB) (bool, error) {
	table_name := aws.String(os.Getenv("SWEARWORDS_TABLE_NAME"))
	result, err := db.GetItem(&dynamodb.GetItemInput{
		TableName: table_name,
		Key: map[string]*dynamodb.AttributeValue{
			"word": {
				S: aws.String("###table_primed###"),
			},
		},
	})
	if err != nil {
		fmt.Printf("Got error calling GetItem: %s", err)
		return true, err
	}
	fmt.Println("GetItem succeeded:")
	fmt.Printf("%+v\n", result)
	return result.Item != nil, nil
}

func getFileContent(bucket *string, key *string, session *session.Session) (*string, error) {

	s3svc := s3.New(session)
	fmt.Println("Created S3 service")

	fmt.Printf("Getting object %s from bucket %s\n", *key, *bucket)
	result, err := s3svc.GetObject(&s3.GetObjectInput{
		Bucket: bucket,
		Key:    key,
	})
	fmt.Println("Got object")
	if err != nil {
		fmt.Printf("Got error calling GetObject: %s", err)
		return nil, err
	}

	defer result.Body.Close()

	// Read the file contents into a byte slice
	fmt.Println("Reading file")
	fileBytes, err := ioutil.ReadAll(result.Body)
	if err != nil {
		fmt.Printf("Got error reading file: %s", err)
		return nil, err
	}
	fmt.Println("Read file")

	// Convert the byte slice to a string
	fileContents := string(fileBytes)

	return &fileContents, nil
}

func main() {
	lambda.Start(HandleRequest)
}
