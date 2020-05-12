package main

import (
	"flag"

	"github.com/jmoney8080/aws-cost-maintenance/cost"
)

var (
	policy string
)

func init() {
	flag.StringVar(&policy, "policy", "", "The policy to run: ")
}

func main() {
	flag.Parse()

	switch policy {
	case "Ec2Modernization":
		cost.Ec2Modernization()
	default:
		flag.PrintDefaults()
	}
}
