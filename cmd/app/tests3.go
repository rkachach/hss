package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)


func main() {
  sess := session.Must(session.NewSessionWithOptions(session.Options{
      SharedConfigState: session.SharedConfigEnable,
  }))
  service := s3.New(sess, &aws.Config{Endpoint: aws.String("http://localhost:4566")})
  input := s3.ListBucketsInput {

  }
  bucketOut, err := service.ListBuckets(&input)
  if err != nil {
    error.Error(err)
  }
  for _, bucket := range bucketOut.Buckets {
    println(bucket.String())
  }
}
