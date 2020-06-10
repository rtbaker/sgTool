package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rtbaker/sgTool/pkg/sgjson"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// main Entry point.
//
// This is a simple tool to view and edit AWS Security Groups
func main() {
	// 1 required argument, one optional
	args := os.Args[1:]

	if len(args) != 1 && len(args) != 2 {
		exitErrorf("Usage: %s <Security Group ID> [Security group definition file]\n", filepath.Base(os.Args[0]))
	}

	sgID := args[0]
	var fileName string

	if len(args) == 2 {
		fileName = args[1]

		// Check file exists
		if !fileExists(fileName) {
			exitErrorf("No such file: %s\n", fileName)
		}
	}

	// Create the AW Session (use AWS_PROFILE env var to specify alternate form the cred's file)
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})

	if err != nil {
		exitErrorf("Error getting AWS session: %s\n", err.Error())
	}

	// Get the Security Group Info
	// Create an EC2 service client.
	svc := ec2.New(sess)

	result, err := svc.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{
		GroupIds: []*string{&sgID},
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "InvalidGroupId.Malformed":
				fallthrough
			case "InvalidGroup.NotFound":
				exitErrorf("%s.", aerr.Message())
			}
		}

		exitErrorf("Unable to get descriptions for security groups, %v", err)
	}

	group, err := sgjson.GroupFromAWS(*result.SecurityGroups[0])

	jsonOut, _ := json.MarshalIndent(group, "", "    ")
	fmt.Println(string(jsonOut))
}

// fileExists Simple check if it's there and a file
func fileExists(filename string) bool {
	var fInfo os.FileInfo
	var err error

	if fInfo, err = os.Stat(filename); os.IsNotExist(err) {
		return false
	}

	return !fInfo.IsDir()
}

// Nicked from the AWS example
func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
