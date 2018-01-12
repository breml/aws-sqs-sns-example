package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// newHandler creates a new handler which gets notified by SNS and retrieves the respective messages from SQS.
func newHandler(sqsSvc *sqs.SQS, qURL *string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msgs, err := sqsSvc.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl: qURL,
		})
		if err != nil {
			fmt.Println("error receiving messages from SQS: ", err.Error())
			code := http.StatusInternalServerError
			http.Error(w, http.StatusText(code), code)
			return
		}

		sqsSvc.DeleteMessage(&sqs.DeleteMessageInput{})

		fmt.Println("received messages: ", msgs)
	}
}

func main() {
	var (
		port     = flag.String("port", "8080", "the port the app should listen on")
		sqsQueue = flag.String("sqs-queue", "", "the name of the SQS queue to publish on")
	)
	flag.Parse()

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")},
	)
	if err != nil {
		fmt.Println("error creating AWS session: ", err)
		return
	}

	// Set up SQS service
	sqsSvc := sqs.New(sess)

	// Get SQS queue
	queue, err := sqsSvc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(*sqsQueue),
	})
	if err != nil {
		fmt.Println("error getting SQS queue: ", err)
		return
	}

	http.HandleFunc("/", newHandler(sqsSvc, queue.QueueUrl))
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
