package main

import (
	"flag"
	"log"
	"os"

	"github.com/jmoney8080/aws-cost-maintenance/cost"
)

var (
	// Info Logger
	Info *log.Logger
	// Warning Logger
	Warning *log.Logger
	// Error Logger
	Error *log.Logger

	policy string
)

func init() {

	// A habbit of mine.  I like to have a nice log line when needed.
	Info = log.New(os.Stdout,
		"[INFO]: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(os.Stderr,
		"[WARNING]: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(os.Stderr,
		"[ERROR]: ",
		log.Ldate|log.Ltime|log.Lshortfile)

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
