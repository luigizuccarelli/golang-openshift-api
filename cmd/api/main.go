package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/luigizuccarelli/golang-openshift-api/pkg/connectors"
	"github.com/luigizuccarelli/golang-openshift-api/pkg/handlers"
	"github.com/luigizuccarelli/golang-openshift-api/pkg/schema"
	"github.com/luigizuccarelli/golang-openshift-api/pkg/validator"
	"github.com/microlib/simple"
)

var (
	loglevel    string
	rulesconfig string
)

func init() {
	flag.StringVar(&loglevel, "l", "info", "Set log level [info,debug,trace]")
	flag.StringVar(&rulesconfig, "c", "input-rules.json", "Use a rules config json file")
}

func main() {

	var rules []schema.Rule
	flag.Parse()
	if rulesconfig == "" || loglevel == "" {
		flag.Usage()
		os.Exit(1)
	}

	logger := &simple.Logger{Level: loglevel}
	err := validator.ValidateEnvars(logger)
	if err != nil {
		os.Exit(-1)
	}
	conn := connectors.NewClientConnections(logger)
	data, err := ioutil.ReadFile(rulesconfig)
	if err != nil {
		conn.Error("Failed to read input file %v", err)
		os.Exit(1)
	}

	err = json.Unmarshal(data, &rules)
	if err != nil {
		conn.Error("Failed to marshal json struct %v", err)
		os.Exit(1)
	}
	// first get token
	token, err := handlers.ExecOS(os.Getenv("HOME_DIR"), "oc", []string{"whoami", "-t"}, true)
	conn.Info("Token %s", token)
	if err != nil {
		conn.Info("Please login into openshift cluster 'oc logn -u http://cluster:port'")
		conn.Error("error from call %v", err)
		os.Exit(1)
	}
	fmt.Println("")
	// loop through each rule and test
	for x, _ := range rules {
		// clean up before we start
		e := os.Remove("/tmp/test.json")
		if e != nil {
			conn.Error("could not clean %v", e)
		}

		conn.Info("Executing test : %s ", rules[x].Name)
		response, e := handlers.APICall(rules[x].Api, token, conn)
		if e != nil {
			conn.Error("call error %v", e)
			conn.Info("FAIL")
			os.Exit(1)
		}

		filterRes, err := handlers.ExecJQFilter(response, rules[x].Filter, conn)
		if err != nil {
			conn.Error("jq filter error %v", e)
			conn.Info("FAIL")
			os.Exit(1)
		}

		regexRes := handlers.ExecRegex([]byte(filterRes), rules[x].Regex, conn)
		result := handlers.ValidateResults(regexRes, rules[x].Occurence, rules[x].Operand, rules[x].Regex)
		conn.Info(result)
		fmt.Println("")
	}
}
