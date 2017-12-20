// dynamodDB credential data model

package credentialstore

import (
    "log"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/awserr"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
    "github.com/satori/go.uuid"
)

type CredentialRecord struct {
    credentialName string
    credentialType string
    credentialData string
}

// using dynamodbiface so mocking can be used for unit testing
type MyDynamo struct {
    Db dynamodbiface.DynamoDBAPI
}

var Dyna *MyDynamo
var describeTableOutput *dynamodb.DescribeTableOutput = nil

// configure dynamodb connection
func ConfigureDynamoDB() error {
    Dyna = new(MyDynamo)
    awsSession, err := session.NewSession(&aws.Config{Region: aws.String(GetConfig().region)})
    if err != nil {
        return err
    }
    var svc *dynamodb.DynamoDB = dynamodb.New(awsSession)
    Dyna.Db = dynamodbiface.DynamoDBAPI(svc)
    return nil
}

// initialize dynamoDB with schema
func InitModel() error {
    model := &dynamodb.CreateTableInput{
        AttributeDefinitions: []*dynamodb.AttributeDefinition {
            {
                AttributeName: aws.String("credentialId"),
                AttributeType: aws.String("S"),
            },
        },
        KeySchema: []*dynamodb.KeySchemaElement {
            {
                AttributeName: aws.String("credentialId"),
                KeyType:       aws.String("HASH"),
            },
        },
        ProvisionedThroughput: &dynamodb.ProvisionedThroughput {
            ReadCapacityUnits: aws.Int64(5),
            WriteCapacityUnits: aws.Int64(5),
        },
        TableName: aws.String("cred"),
    }
    
    // see if table exists
    if nil == describeTableOutput {
        describeTableOutput, err := Dyna.Db.DescribeTable(&dynamodb.DescribeTableInput{TableName: model.TableName})

        if err != nil {
            if aerr, ok := err.(awserr.Error); ok {
                switch aerr.Code() {
                case dynamodb.ErrCodeResourceNotFoundException:
                    // table does not exist, create
                    result, err := Dyna.Db.CreateTable(model)

                    if err != nil {
                        log.Print(err.Error())
                        return err
                    }

                    log.Print(result)
                default:
                    log.Print(err.Error())
                }
            }
        } else {
            log.Print(describeTableOutput)
        }
    }
    return nil
}

// create a new credential record in the DB
// returns the new record ID
func CreateCred(cred CredentialRecord) (string, error) {
    err := InitModel()
    var newId string
    if nil != err {
        return "", err
    }
    newId = uuid.NewV4().String()
    newCred := &dynamodb.PutItemInput{
        Item: map[string]*dynamodb.AttributeValue{
            "credentialID": {
                S: aws.String(newId),
            },
            "credentialData": {
                S: aws.String(cred.credentialData),
            },
            "name": {
                S: aws.String(cred.credentialName),
            },
            "credType": {
                S: aws.String(cred.credentialType),
            },
        },
        ReturnConsumedCapacity: aws.String("TOTAL"),
        TableName: aws.String("cred"),
    }

    result, err := Dyna.Db.PutItem(newCred)
    
    if err != nil {
        log.Print(err.Error())
    } else {
        log.Print(result)
    }
    
    return newId, err
}