package main

import (
	"os"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func main() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{SharedConfigState: session.SharedConfigEnable}))
	ec2Service := ec2.New(sess)

	result, err := ec2Service.DescribeInstances(nil)
	if err != nil {
		fmt.Println("Error: ", err)
	} else {
		instanceName := os.Args[1]
		for _, v := range result.Reservations {
			for _, instance := range v.Instances {
				for _, value := range instance.Tags {
					if *value.Value == instanceName {
						publicDns := *instance.PublicDnsName
						sshCmd := "ssh ubuntu@" + publicDns
						fmt.Println(sshCmd)
					}
				}
			}
		}
	}
}
