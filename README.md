# AWS SQS-SNS Example

A simple example to demonstrate how two apps can communicate over [AWS SQS](https://aws.amazon.com/sqs/) and notify each other about new messages through [AWS SNS](https://aws.amazon.com/sns/).

It consists of two apps: A listener and a talker. The talker sends a message every time it's run and the listener is a web app which listens for new messages from SNS (SNS needs to be configured to POST to it).

## Build

Run `make` to build both binaries

# Run

1. Run listener: `WS_ACCESS_KEY_ID=<your-access-key-id> AWS_SECRET_ACCESS_KEY=<your-secret-access-key> ./listener -sqs-queue '<your-sqs-queue-name>'`
1. Run listener: `WS_ACCESS_KEY_ID=<your-access-key-id> AWS_SECRET_ACCESS_KEY=<your-secret-access-key> ./talker -sqs-queue '<your-sqs-queue-name>' -sns-topic '<your-sns-topic-arn>'`

Then listener should then log the new message in the queue.
