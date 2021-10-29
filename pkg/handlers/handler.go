package handlers

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"

	"github.com/luigizuccarelli/golang-openshift-api/pkg/connectors"
)

const (
	CONTENTTYPE     string = "Content-Type"
	APPLICATIONJSON string = "application/json"
)

func APICall(api string, token string, conn connectors.Client) ([]byte, error) {
	req, _ := http.NewRequest("GET", os.Getenv("BASE_URL")+api, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := conn.Do(req)
	if err != nil {
		conn.Error(fmt.Sprintf("Function APICall http request %v", err))
		return []byte(""), err
	}

	defer resp.Body.Close()
	body, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		conn.Error(fmt.Sprintf("Function APICall %v", e))
		return []byte(""), err
	}
	conn.Debug(fmt.Sprintf("Function APICall response %s", string(body)))
	return body, nil
}

func ExecJQFilter(input []byte, filter string, conn connectors.Client) (string, error) {
	// clean up first
	err := ioutil.WriteFile("/tmp/test.json", input, 0755)
	if err != nil {
		conn.Error("cleanup up tmp file %v", err)
	}
	result, err := ExecOS(os.Getenv("HOME_DIR"), "jq", []string{filter, "/tmp/test.json"}, false)
	conn.Debug("Filter response %s", result)
	if err != nil {
		conn.Error("filter error %v", err)
		return "", err
	}
	return result, nil
}

func ExecRegex(input []byte, regex string, conn connectors.Client) [][]byte {
	reg := regexp.MustCompile(regex)
	res := reg.FindAll(input, -1)
	for x, _ := range res {
		conn.Info("Result %v", string(res[x]))
	}
	return res
}

func ValidateResults(input [][]byte, occurence string, operand string, compare string) string {
	var result string = ""
	switch operand {
	case "eq":
		if occurence == "all" {
			for x, _ := range input {
				if string(input[x]) != compare {
					return "FAIL"
				}
			}
			return "PASS"
		} else {
			occ, _ := strconv.Atoi(occurence)
			if occ == len(input) {
				result = "PASS"
			} else {
				result = "FAIL"
			}
		}
	case "lt":
		occ, _ := strconv.Atoi(occurence)
		if occ < len(input) {
			result = "PASS"
		} else {
			result = "FAIL"
		}
	case "gt":
		occ, _ := strconv.Atoi(occurence)
		if occ > len(input) {
			result = "PASS"
		} else {
			result = "FAIL"
		}
	case "ne":
		occ, _ := strconv.Atoi(occurence)
		if occ != len(input) {
			result = "PASS"
		} else {
			result = "FAIL"
		}
	}
	return result
}

// ExecOS - call to shell
func ExecOS(path string, c string, params []string, trim bool) (string, error) {
	var stdout, stderr bytes.Buffer
	var out string = ""
	cmd := exec.Command(c, params...)
	cmd.Dir = path
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	if err != nil {
		return errStr, err
	}
	if trim {
		out = outStr[:len(outStr)-1]
	} else {
		out = outStr
	}
	return out, nil
}

func PipeOutput(input string) error {
	r, w := io.Pipe()
	go func() {
		fmt.Fprint(w, input)
		w.Close()
	}()

	_, err := io.Copy(os.Stdout, r)
	if err != nil {
		return err
	}
	return nil
}
