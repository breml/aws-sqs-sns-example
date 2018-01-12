package main

import (
	"flag"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func main() {
	var (
		sqsQueue = flag.String("sqs-queue", "", "the name of the SQS queue to publish on")
		snsTopic = flag.String("sns-topic", "", "the ARN of the SNS topic to publish on")
	)
	flag.Parse()

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")},
	)
	if err != nil {
		fmt.Println("error creating AWS session: ", err)
		return
	}

	// Create a service clients
	sqsSvc := sqs.New(sess)
	snsSvc := sns.New(sess)

	// Get SQS queue
	queue, err := sqsSvc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(*sqsQueue),
	})
	if err != nil {
		fmt.Println("error getting SQS queue: ", err)
		return
	}

	// Send a message to SQS
	out, err := sqsSvc.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String("test message"),
		QueueUrl:    queue.QueueUrl,
	})
	if err != nil {
		fmt.Println("error sending message: ", err.Error())
		return
	}
	fmt.Println("message sent to queue: ", out)

	// Announce new message on SNS
	resp, err := snsSvc.Publish(&sns.PublishInput{
		Message:  aws.String(fmt.Sprintf("new message in '%s': %s", queue.QueueUrl, out)),
		TopicArn: aws.String(*snsTopic),
	})
	if err != nil {
		fmt.Println("error publishing to SNS topic:", err.Error())
		return
	}
	fmt.Println("message published to SNS: ", resp.String())
}
