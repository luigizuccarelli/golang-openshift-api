package handlers

import (
	"os"
	"testing"

	"github.com/luigizuccarelli/golang-openshift-api/pkg/connectors"
	"github.com/microlib/simple"
)

func TestAPICall(t *testing.T) {

	logger := &simple.Logger{Level: "trace"}
	t.Run("APICall : should pass", func(t *testing.T) {
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		os.Setenv("BASE_URL", "http://test.com")
		con := connectors.NewTestConnector("\"message\":\"hello\"", 200, "none", logger)
		out, e := APICall("api/v1/test", "1234567", con)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		if len(out) == 0 {
			t.Fatalf("Should not be empty")
		}
	})
}

func TestExecJQFilter(t *testing.T) {

	logger := &simple.Logger{Level: "trace"}
	os.Setenv("HOME_DIR", "/home/lzuccarelli")
	conn := connectors.NewTestConnector("\"message\":\"hello\"", 200, "none", logger)
	t.Run("ExecJQFilter : should pass", func(t *testing.T) {
		out, e := ExecJQFilter([]byte("{\"foo\":\"hello\"}"), ".foo", conn)
		conn.Debug("os call %s", out)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		if out == "" {
			t.Fatalf("Expecting oputput - should not be empty")
		}
	})
}

func TestExecRegex(t *testing.T) {

	logger := &simple.Logger{Level: "trace"}
	conn := connectors.NewTestConnector("\"message\":\"hello\"", 200, "none", logger)
	t.Run("ExecRegex : should pass", func(t *testing.T) {
		out := ExecRegex([]byte("{\"foo\":\"hello\"}"), "foo", conn)
		conn.Debug("regex call %s", out)
		if len(out) == 0 {
			t.Fatalf("Expecting oputput - should not be empty")
		}
	})
}

func TestValidateResults(t *testing.T) {

	t.Run("TestValidateResults : check eq should return pass", func(t *testing.T) {
		out := ValidateResults([][]byte{[]byte("test")}, "1", "eq", "test")
		if out != "PASS" {
			t.Fatalf("Expecting to return PASS - got FAILED")
		}
	})
	t.Run("TestValidateResults : check eq should return fail", func(t *testing.T) {
		out := ValidateResults([][]byte{[]byte("test")}, "0", "eq", "test")
		if out != "FAIL" {
			t.Fatalf("Expecting to return FAIL - got PASSED")
		}
	})

	t.Run("TestValidateResults : check lt should return pass", func(t *testing.T) {
		out := ValidateResults([][]byte{[]byte("test")}, "0", "lt", "test")
		if out != "PASS" {
			t.Fatalf("Expecting to return PASS - got FAILED")
		}
	})
	t.Run("TestValidateResults : check lt should return fail", func(t *testing.T) {
		out := ValidateResults([][]byte{[]byte("test")}, "1", "lt", "test")
		if out != "FAIL" {
			t.Fatalf("Expecting to return FAIL - got PASSED")
		}
	})
	t.Run("TestValidateResults : check gt should return pass", func(t *testing.T) {
		out := ValidateResults([][]byte{[]byte("test")}, "2", "gt", "test")
		if out != "PASS" {
			t.Fatalf("Expecting to return PASS - got FAILED")
		}
	})
	t.Run("TestValidateResults : check gt should return fail", func(t *testing.T) {
		out := ValidateResults([][]byte{[]byte("test")}, "1", "gt", "test")
		if out != "FAIL" {
			t.Fatalf("Expecting to return FAIL - got PASSED")
		}
	})
	t.Run("TestValidateResults : check ne should return pass", func(t *testing.T) {
		out := ValidateResults([][]byte{[]byte("test")}, "0", "ne", "test")
		if out != "PASS" {
			t.Fatalf("Expecting to return PASS - got FAILED")
		}
	})
	t.Run("TestValidateResults : check nt should return fail", func(t *testing.T) {
		out := ValidateResults([][]byte{[]byte("test")}, "1", "ne", "test")
		if out != "FAIL" {
			t.Fatalf("Expecting to return FAIL - got PASSED")
		}
	})
}

func TestPipeOutput(t *testing.T) {

	t.Run("TestPipeOutput : should pass", func(t *testing.T) {
		err := PipeOutput("test")
		if err != nil {
			t.Fatalf("Expected to pass")
		}
	})
}
