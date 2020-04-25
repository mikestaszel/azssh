package main

import (
	"bufio"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"os"
	"os/exec"
	"strings"
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

	return nil, fmt.Errorf("could not find instance")
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
	instanceState, err := getInstanceState(*instance)
	if err != nil {
		return err
	}
	if instanceState == "shutting-down" || instanceState == "terminated" || instanceState == "stopping" || instanceState == "stopped" {
		fmt.Errorf("instance stopped or terminated")
	}

	instanceId := *instance.InstanceId
	input := &ec2.StartInstancesInput{
		InstanceIds: []*string {
			aws.String(instanceId),
		},
	}

	_, err = ec2Service.StartInstances(input)
	if err != nil {
		return err
	}
	fmt.Println("starting...")
	return nil
}

// Stop an instance.
func stopInstance(ec2Service *ec2.EC2, instance *ec2.Instance) error {
	instanceState, err := getInstanceState(*instance)
	if err != nil {
		return err
	}

	if instanceState == "shutting-down" || instanceState == "terminated" || instanceState == "stopping" || instanceState == "stopped" {
		fmt.Errorf("instance stopped or terminated")
	}

	instanceId := *instance.InstanceId
	input := &ec2.StopInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceId),
		},
	}

	_, err = ec2Service.StopInstances(input)
	if err != nil {
		return err
	}
	fmt.Println("stopping...")
	return nil
}

// Reboot an instance.
func rebootInstance(ec2Service *ec2.EC2, instance *ec2.Instance) error {
	instanceState, err := getInstanceState(*instance)
	if err != nil {
		return err
	}
	if instanceState == "shutting-down" || instanceState == "terminated" || instanceState == "stopping" || instanceState == "stopped" {
		fmt.Errorf("instance stopped or terminated")
	}

	instanceId := *instance.InstanceId
	input := &ec2.RebootInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceId),
		},
	}

	_, err = ec2Service.RebootInstances(input)
	if err != nil {
		return err
	}
	fmt.Println("restarting...")
	return nil
}

// Terminate an instance.
func terminateInstance(ec2Service *ec2.EC2, instance *ec2.Instance) error {
	instanceState, err := getInstanceState(*instance)
	if err != nil {
		return err
	}
	if instanceState == "terminated" {
		fmt.Errorf("instance already terminated")
	}

	instanceId := *instance.InstanceId
	input := &ec2.TerminateInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceId),
		},
	}

	_, err = ec2Service.TerminateInstances(input)
	if err != nil {
		return err
	}
	fmt.Println("terminating...")
	return nil
}

func printUsage() {
	usageStr := "azssh: utility to manage EC2 instances\n\n" +
		"Usage:\n" +
		"azssh ssh instance_name        SSH into instance\n" +
		"azssh ls                       list instances\n" +
		"azssh dns                      print public DNS of instance\n" +
		"azssh start instance_name      start instance\n" +
		"azssh stop instance_name       stop instance\n" +
		"azssh restart instance_name    restart instance\n" +
		"azssh terminate instance_name  terminate instance\n" +
		"azssh help                     print this help text"
	fmt.Println(usageStr)
}

// Main function.
func main() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{SharedConfigState: session.SharedConfigEnable}))
	ec2Service := ec2.New(sess)

	if len(os.Args) == 3 && os.Args[1] == "ssh" {
		instanceName := os.Args[2]
		instance, _ := getInstanceByName(ec2Service, instanceName)

		publicDns, _ := getInstancePublicDns(*instance)

		if publicDns == "" {
			fmt.Println("Instance does not have a public DNS address or is stopped.")
			os.Exit(1)
		}

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Username? [ubuntu]: ")
		username, _ := reader.ReadString('\n')
		username = strings.TrimSuffix(username, "\n")
		if username == "" {
			username = "ubuntu"
		}

		sshCmd := fmt.Sprintf("ssh %s@%s", username, publicDns)
		fmt.Println("running command:", sshCmd)

		cmd := exec.Command("ssh", fmt.Sprintf("%s@%s", username, publicDns))
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		_ = cmd.Run()
	} else if len(os.Args) == 2 && os.Args[1] == "ls" {
		instances, err := getInstances(ec2Service)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if len(instances) == 0 {
			fmt.Println("no instances detected.")
		} else {
			fmt.Println("instances:")
			for _, instance := range instances {
				instanceName, _ := getInstanceName(instance)

				if instanceName == "" {
					instanceName, _ = getInstancePublicDns(instance)
				}

				instanceState, _ := getInstanceState(instance)
				fmt.Printf("%s - %s\n", instanceName, instanceState)
			}
		}
	} else if len(os.Args) == 3 && os.Args[1] == "dns" {
		instanceName := os.Args[2]
		instance, err := getInstanceByName(ec2Service, instanceName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		publicDns, err := getInstancePublicDns(*instance)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(publicDns)
	} else if len(os.Args) == 3 && (os.Args[1] == "start") {
		instanceName := os.Args[2]
		instance, err := getInstanceByName(ec2Service, instanceName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = startInstance(ec2Service, instance)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else if len(os.Args) == 3 && (os.Args[1] == "stop") {
		instanceName := os.Args[2]
		instance, err := getInstanceByName(ec2Service, instanceName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = stopInstance(ec2Service, instance)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else if len(os.Args) == 3 && (os.Args[1] == "terminate") {
		instanceName := os.Args[2]
		instance, err := getInstanceByName(ec2Service, instanceName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = terminateInstance(ec2Service, instance)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else if len(os.Args) == 3 && (os.Args[1] == "reboot" || os.Args[1] == "restart") {
		instanceName := os.Args[2]
		instance, err := getInstanceByName(ec2Service, instanceName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = rebootInstance(ec2Service, instance)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else if len(os.Args) == 2 && (os.Args[1] == "help") {
		printUsage()
	} else {
		fmt.Println("Invalid parameters!")
		printUsage()
		os.Exit(1)
	}
}
