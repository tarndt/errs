// (C)Tylor Arndt 2014, Mozilla Public License (MPL) Version 2.0
// See LICENSE file for details.
// For other licensing options please contact the author.

package errs

import (
	"errors"
	"strings"
	"testing"
)

//These tests are inherently brittle due to being line number sensitive :/

func TestNewErrorMissingMessage(t *testing.T) {
	const expected = "errors_test.go:17 errs.TestNewErrorMissingMessage(): An unknown error occured"
	if actual := New("").Error(); actual != expected {
		t.Errorf("Actual result: %q, not not match expected result: %q", actual, expected)
	}
}

func TestNewErrorString(t *testing.T) {
	const expected = "errors_test.go:24 errs.TestNewErrorString(): TestString"
	if actual := New("TestString").Error(); actual != expected {
		t.Errorf("Actual result: %q, not not match expected result: %q", actual, expected)
	}
}

func TestNewErrorArgs(t *testing.T) {
	const expected = "errors_test.go:31 errs.TestNewErrorArgs(): Testing 1,2,3..."
	if actual := New("Testing %d,2,%d...", 1, 3).Error(); actual != expected {
		t.Errorf("Actual result: %q, not not match expected result: %q", actual, expected)
	}
}

func TestAppendMissingParent(t *testing.T) {
	const expected = "errors_test.go:38 errs.TestAppendMissingParent(): An unknown error occured"
	if actual := Append(nil, "").Error(); actual != expected {
		t.Errorf("Actual result: %q, not not match expected result: %q", actual, expected)
	}
}

func TestAppendMissingMessage(t *testing.T) {
	const expected = "errors_test.go:45 errs.TestAppendMissingMessage(): An unknown error occured; Details:\n\terrors_test.go:45 errs.TestAppendMissingMessage(): Parent Error"
	if actual := Append(New("Parent Error"), "").Error(); actual != expected {
		t.Errorf("Actual result: %q, not not match expected result: %q", actual, expected)
	}
}

func TestAppend(t *testing.T) {
	const expected = "errors_test.go:52 errs.TestAppend(): Child Error; Details:\n\terrors_test.go:52 errs.TestAppend(): Parent Error"
	if actual := Append(New("Parent Error"), "Child Error").Error(); actual != expected {
		t.Errorf("Actual result: %q, not not match expected result: %q", actual, expected)
	}
}

func TestNewAppendChain(t *testing.T) {
	const expected = "errors_test.go:59 errs.TestNewAppendChain(): Child Error; Details:\n\terrors_test.go:65 errs.subA(): subA; Details:\n\terrors_test.go:69 errs.subB(): subB"
	if actual := Append(subA(), "Child Error").Error(); actual != expected {
		t.Errorf("Actual result: %q, not not match expected result: %q", actual, expected)
	}
}

func subA() error {
	return Append(subB(), "sub%s", "A")
}

func subB() error {
	return New("subB")
}

func TestGetRootError(t *testing.T) {
	const expected = "errors_test.go:69 errs.subB(): subB"
	if actual := GetRootErr(subA()).Error(); actual != expected {
		t.Errorf("Actual result: %q, not not match expected result: %q", actual, expected)
	}
}

func TestGetErrorNonPkgErr(t *testing.T) {
	var expected = ErrorLocation{
		funcName: "Unknown Function",
		filename: "Unknown File",
		lineNum:  -1,
	}
	if actual := GetErrorLoc(errors.New("Non pkg error")); actual != expected {
		t.Errorf("Actual result: %+v, not not match expected result: %+v", actual, expected)
	}
}

func TestGetErrorLoc(t *testing.T) {
	var expected = ErrorLocation{
		funcName: "errs.subB",
		filename: "errors_test.go",
		lineNum:  69,
	}
	if actual := GetErrorLoc(GetRootErr(subA())); actual != expected {
		t.Errorf("Actual result: %+v, not not match expected result: %+v", actual, expected)
	}
}

func TestPanicToErr(t *testing.T) {
	defer func() {
		text := PanicToErr(recover()).Error()
		lines := strings.Split(text, "\n")
		//t.Logf(text)
		if lines[0] != "errors_test.go:140 errs.panicyFunc(): A panic occured, but was recovered; Details: Test Panic;" {
			t.Error("Error body did not match")
		} else if strings.TrimSpace(lines[1]) != "*** Stack Trace ***" {
			t.Error("End trace start header missing")
		} else if !strings.Contains(lines[2], "errs/errors_test.go:140") {
			t.Error("Line/Location info for panic location missing")
		} else if strings.TrimSpace(lines[3]) != `panicyFunc: panic("Test Panic")` {
			t.Error("Details from panic location missing")
		} else if !strings.Contains(lines[4], "errs/errors_test.go:135") {
			t.Error("Line/Location info for wrapping closure missing")
		} else if strings.TrimSpace(lines[5]) != "func.002: panicyFunc()" {
			t.Error("Details from wrapping closure missing")
		} else if !strings.Contains(lines[6], "errs/errors_test.go:136") {
			t.Error("Line/Location info for wrapping closure invocation missing")
		} else if strings.TrimSpace(lines[7]) != "TestPanicToErr: }()" {
			t.Error("Details from wrapping closure invocation missing")
		} else if !strings.Contains(lines[8], "testing/testing.go") {
			t.Error("Line/Location info for test runner missing")
		} else if strings.TrimSpace(lines[9]) != "tRunner: test.F(t)" {
			t.Error("Details from test runner")
		} else if !strings.Contains(lines[10], "runtime/proc.c") {
			t.Error("Line/Location info for goexit missing")
		} else if strings.TrimSpace(lines[11]) != "goexit: runtime·goexit(void)" {
			t.Error("Details from goexit")
		} else if strings.TrimSpace(lines[12]) != "*** End Stack Trace ***" {
			t.Error("End trace header missing")
		}
	}()
	func() { //Lets make the trace longer with this extra closure
		panicyFunc()
	}()
}

func panicyFunc() {
	panic("Test Panic")
}
