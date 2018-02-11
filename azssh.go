package main

import (
	"os"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/aws"
	"os/exec"
)

// Get a list of instances.
func getInstances(ec2Service *ec2.EC2) ([]ec2.Instance, error) {
	result, err := ec2Service.DescribeInstances(nil)
	output := make([]ec2.Instance, 0)

	if err == nil {
		for _, v := range result.Reservations {
			for _, instance := range v.Instances {
				output = append(output, *instance)
			}
		}
	}

	return output, err
}

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

// Get the instance's name.
func getInstanceName(instance ec2.Instance) (string, error) {
	for _, value := range instance.Tags {
		if *value.Key == "Name" {
			return *value.Value, nil
		}
	}
	return "", fmt.Errorf("could not find instance")
}

// Get the instance's public DNS address.
func getInstancePublicDns(instance ec2.Instance) (string, error) {
	return *instance.PublicDnsName, nil
}

// Get the instance state. One of: pending | running | shutting-down | terminated | stopping | stopped
func getInstanceState(instance ec2.Instance) (string, error) {
	return *instance.State.Name, nil
}

// Start a stopped instance.
func startInstance(ec2Service *ec2.EC2, instance *ec2.Instance) error {
	instanceState, _ := getInstanceState(*instance)
	if instanceState == "shutting-down" || instanceState == "terminated" || instanceState == "stopping" || instanceState == "stopped" {
		fmt.Errorf("instance stopped or terminated")
	}

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
	fmt.Println("starting...")
	return nil
}

// Stop an instance.
func stopInstance(ec2Service *ec2.EC2, instance *ec2.Instance) error {
	instanceState, _ := getInstanceState(*instance)
	if instanceState == "shutting-down" || instanceState == "terminated" || instanceState == "stopping" || instanceState == "stopped" {
		fmt.Errorf("instance stopped or terminated")
	}

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
	fmt.Println("stopping...")
	return nil
}

// Reboot an instance.
func rebootInstance(ec2Service *ec2.EC2, instance *ec2.Instance) error {
	instanceState, _ := getInstanceState(*instance)
	if instanceState == "shutting-down" || instanceState == "terminated" || instanceState == "stopping" || instanceState == "stopped" {
		fmt.Errorf("instance stopped or terminated")
	}

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
	fmt.Println("restarting...")
	return nil
}

func printUsage() {
	fmt.Println("azssh: utility to manage EC2 instances")
	fmt.Println("usage: azssh ssh instanceName        SSH into instance")
	fmt.Println("   or: azssh ls                      list instances")
	fmt.Println("   or: azssh dns                     print public DNS of instance")
	fmt.Println("   or: azssh up instanceName         start this instance public DNS of instance")
	fmt.Println("   or: azssh down instanceName       print public DNS of instance")
	fmt.Println("   or: azssh reboot instanceName     print public DNS of instance")
	fmt.Println("   or: azssh help                    print this help text")
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

		publicDns, _ := getInstancePublicDns(*instance)
		sshCmd := "ssh ubuntu@" + publicDns
		fmt.Println("running command: ", sshCmd)

		cmd := exec.Command("ssh", "ubuntu@" + publicDns)
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		cmd.Run()
	} else if len(os.Args) == 2 && os.Args[1] == "ls" {
		instances, _ := getInstances(ec2Service)
		fmt.Println("instances:")
		for _, instance := range instances {
			instanceName, _ := getInstanceName(instance)
			instanceState, _ := getInstanceState(instance)
			fmt.Printf("%s - %s\n", instanceName, instanceState)
		}
	} else if len(os.Args) == 3 && os.Args[1] == "dns" {
		instanceName := os.Args[2]
		instance, _ := getInstanceByName(ec2Service, instanceName)
		publicDns, _ := getInstancePublicDns(*instance)
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
	} else if len(os.Args) == 2 && (os.Args[1] == "help")  {
		printUsage()
	} else {
		fmt.Println("Invalid parameters!")
		printUsage()
		os.Exit(1)
	}
}
