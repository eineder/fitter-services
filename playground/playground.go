package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func main() {
	// Initialize a session
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create a DynamoDB client
	dynamoClient := dynamodb.New(sess)

	// Define the table name
	tableName := "AwsFileToDynamodbStack-swearwordstableAF08B832-10RHLXAX1BECU"

	// Define the list of words to search for
	words := []string{"hello", "world", "example", "Depp"}

	// Construct the PartiQL query
	query := "SELECT * FROM \"" + tableName + "\" WHERE \"word\" IN ("
	parameters := make([]*dynamodb.AttributeValue, 0, len(words))

	for i, word := range words {
		parameters = append(parameters, &dynamodb.AttributeValue{S: aws.String(word)})
		query += "?"
		if i < len(words)-1 {
			query += ", "
		}
	}
	query += ")"

	// Execute the PartiQL query
	input := &dynamodb.ExecuteStatementInput{
		Statement:  aws.String(query),
		Parameters: parameters,
	}
	output, err := dynamoClient.ExecuteStatement(input)
	if err != nil {
		fmt.Println("Error executing statement:", err)
		return
	}

	// Print the retrieved items
	fmt.Println("Retrieved items:")
	for _, item := range output.Items {
		fmt.Println(item)
	}
}
