package cost

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// InstanceTypeInventory type for housing a map of instance types divided by whether they are latest generation or not
type InstanceTypeInventory struct {
	LatestGeneration map[string][]InstanceInfo `json:"latest_generation"`
	OlderGeneration  map[string][]InstanceInfo `json:"older_generation"`
}

// InstanceInfo type for holding some metadata for an instance
type InstanceInfo struct {
	InstanceID   string `json:"instance_id"`
	InstanceName string `json:"instance_name,omitempty"`
}

var (
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger

	policy string
)

func init() {

	// A habbit of mine.  I like to have a nice log line when needed.
	infoLogger = log.New(os.Stdout,
		"[INFO]: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	warningLogger = log.New(os.Stderr,
		"[WARNING]: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	errorLogger = log.New(os.Stderr,
		"[ERROR]: ",
		log.Ldate|log.Ltime|log.Lshortfile)

}

// Ec2Modernization outputs all ec2 instances bucketed by latest generation criteria
func Ec2Modernization() {

	svc := ec2.New(session.Must(session.NewSession()))

	// Look up all supported instance types
	diti := &ec2.DescribeInstanceTypesInput{}
	dito, err := svc.DescribeInstanceTypes(diti)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				errorLogger.Fatalln(aerr.Error())
			}
		} else {
			errorLogger.Fatalln(err.Error())
		}
		return
	}

	// build a map of instance type to latest genration boolean and cache in memory
	instanceTypesMetadata := make(map[string]bool)
	for _, instanceType := range dito.InstanceTypes {
		instanceTypesMetadata[*instanceType.InstanceType] = *instanceType.CurrentGeneration
	}

	instanceInfos := make(map[string][]InstanceInfo)
	err = svc.DescribeInstancesPages(&ec2.DescribeInstancesInput{}, func(page *ec2.DescribeInstancesOutput, lastPage bool) bool {
		for _, reservation := range page.Reservations {
			for _, instance := range reservation.Instances {
				// It is useful to have a tag of key Name but not required
				var name = ""
				for _, tag := range instance.Tags {
					if *tag.Key == "Name" {
						name = *tag.Value
					}
				}

				instanceInfo := InstanceInfo{
					InstanceID:   *instance.InstanceId,
					InstanceName: name,
				}

				// if the instance type exists in the map append the instance info to the list
				// else initialize the instance type into the map
				if value, exists := instanceInfos[*instance.InstanceType]; exists {
					instanceInfos[*instance.InstanceType] = append(value, instanceInfo)
				} else {
					instanceInfos[*instance.InstanceType] = []InstanceInfo{instanceInfo}
				}
			}
		}
		return !lastPage
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				errorLogger.Fatalln(aerr.Error())
			}
		} else {
			errorLogger.Fatalln(err.Error())
		}
		return
	}

	// loop over all the instance types we have gathered and divide them into latest generation or not
	latest := make(map[string][]InstanceInfo)
	notLatest := make(map[string][]InstanceInfo)
	for key, value := range instanceInfos {
		if instanceTypesMetadata[key] {
			latest[key] = value
		} else {
			notLatest[key] = value
		}
	}

	instanceTypeInventory := InstanceTypeInventory{
		LatestGeneration: latest,
		OlderGeneration:  notLatest,
	}

	// print the raw json.  Do not use a logger from above so we can get the raw json string
	rawJSON, _ := json.Marshal(instanceTypeInventory)
	fmt.Println(string(rawJSON))
}
