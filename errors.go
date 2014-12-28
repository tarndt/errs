// (C)Tylor Arndt 2014, Mozilla Public License (MPL) Version 2.0
// See LICENSE file for details.
// For other licensing options please contact the author.

package errs

import (
	"bytes"
	"fmt"
	"path"
	"runtime"
	"runtime/debug"
)

//ErrsErr is the error interface implementation used by this package
type ErrsErr struct {
	Msg    string
	Loc    ErrorLocation
	Parent error
}

var _ error = &ErrsErr{} //Ensure ErrsErr implements error at compile time

//MultiLineErrorJoin is the string used to join single errors/lines;
//replace as desired at runtime (pref. at program init for thread-saftey).
var MultiLineErrorJoin = "; Details:\n\t"

//Implements error interface, contructs multi-error traces as necessary
func (err *ErrsErr) Error() string {
	if err.Parent == nil {
		return err.Msg
	}
	//Build a multi-error message
	var msgBuf bytes.Buffer
	var ok bool
	for {
		msgBuf.WriteString(err.Msg)
		msgBuf.WriteString(MultiLineErrorJoin)
		err, ok = err.Parent.(*ErrsErr)
		switch {
		case !ok:
			msgBuf.WriteString(err.Parent.Error())
		case err.Parent == nil:
			msgBuf.WriteString(err.Msg)
		default:
			continue
		}
		break
	}
	return msgBuf.String()
}

//ErrorDetails describes where an error occured
type ErrorLocation struct {
	funcName, filename string
	lineNum            int
}

//ErrorFormatter is used to build the textual representation of errors from
//the error message text and the information regarding the location of the error
//in source;
//replace as desired at runtime (pref. at program init for thread-saftey).
var ErrorFormatter func(errMsg string, loc ErrorLocation) string = defFmtErrMsg

//default ErrorFormatter implementation
func defFmtErrMsg(errMsg string, loc ErrorLocation) string {
	return fmt.Sprintf("%s:%d %s(): %s", loc.filename, loc.lineNum, loc.funcName, errMsg)
}

//New creates a new error in a similar manner to fmt.Errorf
func New(errMsg string, args ...interface{}) error {
	loc := getCallersInfo(2)
	return &ErrsErr{
		Msg: ErrorFormatter(fmtErrMsg(errMsg, args), loc),
		Loc: loc,
	}
}

//Append creates an error that contains context information regarding the
//preceding error that lead to its occurrence. This makes error traces possible.
func Append(suppliedErr error, errMsg string, args ...interface{}) error {
	loc := getCallersInfo(2)
	return &ErrsErr{
		Msg:    ErrorFormatter(fmtErrMsg(errMsg, args), loc),
		Loc:    loc,
		Parent: suppliedErr,
	}
}

//GetRootErr returns the original error of an error trace, if the error passed
//to GetRootErr is not an error trace (errs.ErrsErr) it is the root error and is
//returned.
func GetRootErr(err error) error {
	if err == nil {
		return nil
	}
	var errWithParent *ErrsErr
	var ok bool
	for {
		if errWithParent, ok = err.(*ErrsErr); !ok || errWithParent.Parent == nil {
			break
		}
		err = errWithParent.Parent
	}
	return err
}

//GetErrorLoc returns the location at which the error took place, if the error
//passed to GetRootErr is not an error trace (errs.ErrsErr) the location is
//denoted as unknown.
func GetErrLoc(err error) ErrorLocation {
	if err != nil {
		if errWithLoc, ok := err.(*ErrsErr); ok {
			return errWithLoc.Loc
		}
	}
	return ErrorLocation{
		funcName: "Unknown Function",
		filename: "Unknown File",
		lineNum:  -1,
	}
}

//PanicToErr is passed the result of "recover()" and converts any non-nil value
//to an error object.
func PanicToErr(recoverReturnVal interface{}) error {
	const msgFmt = "%s:%d %s(): A panic occured, but was recovered; Details: %+v;\n\t*** Stack Trace ***\n\t%s*** End Stack Trace ***\n"
	if recoverReturnVal != nil {
		loc := getCallersInfo(4)
		prettyStack := bytes.Join(bytes.Split(debug.Stack(), []byte{'\n'})[6:], []byte{'\n', '\t'})
		return fmt.Errorf(msgFmt, loc.filename, loc.lineNum, loc.funcName, recoverReturnVal, prettyStack)
	}
	return nil
}

func fmtErrMsg(msg string, args []interface{}) string {
	if len(msg) == 0 {
		msg = "An unknown error occured"
	} else if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	return msg
}

func getCallersInfo(depth int) ErrorLocation {
	funcName, filename, lineNum := "Unknown Function", "Unknown File", -1
	var programCounter uintptr
	var isKnownFunc bool
	if programCounter, filename, lineNum, isKnownFunc = runtime.Caller(depth); isKnownFunc {
		filename = path.Base(filename)
		funcName = path.Base(runtime.FuncForPC(programCounter).Name())
	}
	return ErrorLocation{
		funcName: funcName,
		filename: filename,
		lineNum:  lineNum,
	}
}
