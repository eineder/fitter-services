package main

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type ContainsSwearwordsEvent struct {
	Text string `json:"text"`
}

type ContainsSwearwordsResponse struct {
	Swearwords []string `json:"swearwords"`
}

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest(ctx context.Context, event *ContainsSwearwordsEvent) (*ContainsSwearwordsResponse, error) {
	if event == nil {
		return nil, fmt.Errorf("received nil event")
	}

	if event.Text == "" {
		return nil, fmt.Errorf("received empty text")
	}

	swearwordsTableName := os.Getenv("SWEARWORDS_TABLE_NAME")
	if swearwordsTableName == "" {
		msg := "missing SWEARWORDS_TABLE_NAME environment variable"
		fmt.Println(msg)
		return nil, fmt.Errorf(msg)
	}

	distinctWords := distinct(removePunctuations(strings.Split(event.Text, " ")))
	words := remove(&distinctWords, "")
	parameters := getParameters(words)
	sess := session.Must(session.NewSession())
	dynamoClient := dynamodb.New(sess)

	swearwords, err := querySwearwords(parameters, dynamoClient, swearwordsTableName)
	if err != nil {
		fmt.Println("Error querying swearwords ", err)
		return nil, err
	}

	return &ContainsSwearwordsResponse{
		Swearwords: swearwords,
	}, nil
}

func querySwearwords(parameters []*dynamodb.AttributeValue, dynamoClient *dynamodb.DynamoDB, tableName string) ([]string, error) {
	const batchSize = 100
	swearwords := []string{}
	batchParameters := []*dynamodb.AttributeValue{}
	statement := fmt.Sprintf(`SELECT * FROM "%s" WHERE word IN (`, tableName)

	for i, parameter := range parameters {
		batchParameters = append(batchParameters, parameter)
		statement += "?, "
		if len(parameters) == batchSize || i == len(parameters)-1 {
			// remove the last comma and space
			statement = statement[:len(statement)-2]
			statement += ")"
			sw, err := queryBatchSwearwords(dynamoClient, tableName, batchParameters, swearwords, statement)
			if err != nil {
				fmt.Println("Error querying swearwords ", err)
				return nil, err
			}
			swearwords = append(swearwords, sw...)
			batchParameters = make([]*dynamodb.AttributeValue, batchSize)
		}
	}
	return swearwords, nil
}

func queryBatchSwearwords(svc *dynamodb.DynamoDB, tableName string, batchParameters []*dynamodb.AttributeValue, swearwords []string, statement string) ([]string, error) {

	out, err := svc.ExecuteStatement(&dynamodb.ExecuteStatementInput{
		Statement:  &statement,
		Parameters: batchParameters,
	})
	if err != nil {
		fmt.Println("Error executing statement ", err)
		return nil, err
	}

	sw := getSwearwords(out.Items)
	return sw, nil
}

func getSwearwords(items []map[string]*dynamodb.AttributeValue) []string {
	sw := []string{}
	for _, item := range items {
		sw = append(sw, *item["word"].S)
	}
	return sw
}

func getParameters(words *[]string) []*dynamodb.AttributeValue {
	parameters := make([]*dynamodb.AttributeValue, len(*words))
	for i, word := range *words {
		parameters[i] = &dynamodb.AttributeValue{
			S: aws.String(word),
		}
	}

	return parameters
}

func distinct(words []string) []string {
	m := make(map[string]bool)
	a := []string{}
	for _, word := range words {
		if m[word] {
			continue
		}
		m[word] = true
		a = append(a, word)
	}
	return a
}

func removePunctuations(words []string) []string {
	for i, word := range words {
		words[i] = removePunctuation(word)
	}
	return words
}

func removePunctuation(word string) string {
	// Compile the regex
	re := regexp.MustCompile(`^[^a-zA-Z]+|[^a-zA-Z]+$`)
	// Replace non-letter characters at the start and end
	return re.ReplaceAllString(word, "")
}

func remove(arr *[]string, word string) *[]string {
	newArr := []string{}
	for _, w := range *arr {
		if w != word {
			newArr = append(newArr, w)
			break
		}
	}
	return &newArr
}
