package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func main() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{SharedConfigState: session.SharedConfigEnable}))
	ec2_service := ec2.New(sess)

	result, err := ec2_service.DescribeInstances(nil)
	if err != nil {
		fmt.Println("Error: ", err)
	} else {
		fmt.Println("Success!")
		fmt.Println(result)
	}
}
