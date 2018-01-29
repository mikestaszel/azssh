package main

import (
	"os"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/aws"
	"os/exec"
)

// Get an instance by name.
func getInstanceByName(ec2Service *ec2.EC2, instanceName string) (*ec2.Instance, error) {
	result, err := ec2Service.DescribeInstances(nil)

	if err == nil {
		for _, v := range result.Reservations {
			for _, instance := range v.Instances {
				for _, value := range instance.Tags {
					if *value.Value == instanceName {
						return instance, nil
					}
				}
			}
		}
	}

	return nil, err
}

// Get the instance's public DNS address.
func getPublicDns(instance *ec2.Instance) (string, error) {
	return *instance.PublicDnsName, nil
}

// Start a stopped instance.
func startInstance(ec2Service *ec2.EC2, instance *ec2.Instance) error {
	// TODO: first check the state of the machine - if already running or starting, raise an error
	instanceId := *instance.InstanceId
	input := &ec2.StartInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceId),
		},
	}

	_, err := ec2Service.StartInstances(input)
	if err != nil {
		return err
	}
	return nil
}

// Stop an instance.
func stopInstance(ec2Service *ec2.EC2, instance *ec2.Instance) error {
	// TODO: first check the state of the machine - if already stopped or stopping, raise an error
	instanceId := *instance.InstanceId
	input := &ec2.StopInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceId),
		},
	}

	_, err := ec2Service.StopInstances(input)
	if err != nil {
		return err
	}
	return nil
}

// Reboot an instance.
func rebootInstance(ec2Service *ec2.EC2, instance *ec2.Instance) error {
	// TODO: first check the state of the machine - if not running, raise an error
	instanceId := *instance.InstanceId
	input := &ec2.RebootInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceId),
		},
	}

	_, err := ec2Service.RebootInstances(input)
	if err != nil {
		return err
	}
	return nil
}

// Main function.
func main() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{SharedConfigState: session.SharedConfigEnable}))
	ec2Service := ec2.New(sess)

	if len(os.Args) == 3 && os.Args[1] == "ssh" {
		instanceName := os.Args[2]
		instance, _ := getInstanceByName(ec2Service, instanceName)

		// TODO: if the instance is off, prompt user - should it be turned on? if yes, wait 60 seconds then retry
		// TODO: do something if publicDns is an empty string!

		publicDns, _ := getPublicDns(instance)
		sshCmd := "ssh ubuntu@" + publicDns
		fmt.Println("running command: ", sshCmd)

		cmd := exec.Command("ssh", "ubuntu@" + publicDns)
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		cmd.Run()
	} else if len(os.Args) == 3 && os.Args[1] == "dns" {
		instanceName := os.Args[2]
		instance, _ := getInstanceByName(ec2Service, instanceName)
		publicDns, _ := getPublicDns(instance)
		fmt.Println(publicDns)
	} else if len(os.Args) == 3 && (os.Args[1] == "up" || os.Args[1] == "start") {
		instanceName := os.Args[2]
		instance, _ := getInstanceByName(ec2Service, instanceName)
		startInstance(ec2Service, instance)
	} else if len(os.Args) == 3 && (os.Args[1] == "down" || os.Args[1] == "shutdown" || os.Args[1] == "stop") {
		instanceName := os.Args[2]
		instance, _ := getInstanceByName(ec2Service, instanceName)
		stopInstance(ec2Service, instance)
	} else if len(os.Args) == 3 && (os.Args[1] == "reboot" || os.Args[1] == "restart") {
		instanceName := os.Args[2]
		instance, _ := getInstanceByName(ec2Service, instanceName)
		rebootInstance(ec2Service, instance)
	} else {
		// TODO: Write some help text!
		fmt.Println("Invalid parameters!")
		os.Exit(1)
	}
}
