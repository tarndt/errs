errs: Debuggable errors for Go
====
### Author & Version
Author: [Tylor Arndt]

1.0 - First public release (Bug reports and PRs are welcome)

### Motivation
While the "go way" of explcitly handling errors + using defer has worked very well for me. I have found in large Go applications that errors which "trickle up" from multiple levels of abstraction can be frustrating to track down a source of; this package seeks to make that task easier.

This package also supports a few other useful tasks such as converting panics to errors (`errs.PanicToErr`) and a few functions for manipulating the tracable error it returns such as `errs.GetRootErr` and `errs.GetErrLoc`. 

*All externally exposed error and standard Go errors to this should be 100% interoperable with your existing code base*

### Basic Operations
 * `New(errMsg string, args ...interface{}) error`: Creates a new error in a similar manner to `fmt.Errorf`.
 * `Append(suppliedErr error, errMsg string, args ...interface{}) error`: Append creates an error that contains context information regarding the preceding error that lead to its occurrence. This makes error traces possible.
 * `GetRootErr(err error) error`: GetRootErr returns the original error of an error trace, if the error passed to GetRootErr is not an error trace (errs.ErrsErr) then *it is the root error* and is returned.
* `GetErrLoc(err error) ErrorLocation`: GetErrorLoc returns the location at which the error took place, if the error passed to GetRootErr is not an error trace (errs.ErrsErr) the location is denoted as unknown.
* `PanicToErr(recoverReturnVal interface{}) error`: PanicToErr is passed the result of `recover()` and converts any non-nil value to an error object.

### Configuration Options
* `ErrorFormatter` is a  function used to format errors and can be replaced at runtime.
* `MultiLineErrorJoin` is the string used to join single errors/lines and can also be replaced at runtime.
See [errors.go] for details.

### Example: Tracable Error (using `errs.Append`)
Take for example "permission denied" from a file access in a very large application that does many file accesses. Lets assume a we get a "permission denied" from `os.Open`.

Using this package rather than having `err.Error()` result in `permission denied`, we will get an error trace to work with containing content simmilar to:
```
highLevelJob.go:90 main.runHighLevelJob(): High Level Job: 3785, failed; Details:
    processJob.go:35 mySubPkg.ProcessManyFiles(): File 2 of [fileA,fileB,fileC] could not be processed; Details:
    processJob.go:69 mySubPkg.processFile(): Could not open file; Details:
    permission denied
```
Pseduo Code:
```go
func runHighLevelJob(jobId int) error {
 ...
  if err := processManyFiles(foobar); err != nil {
   return errs.Append(err, "High Level Job: %d, failed!", jobId) 
  }
 ...
}

func ProcessManyFiles(somefilenames []string) error {
  ...
  for i,filename := range somefilenames
   if err := processFile(filename); err != nil {
    return errs.Append(err, "File %d of %+v could not be processed", i)
   }
   ...
  }
}

func processFile(filename string) error {
  ...
  //we are assuming this returns os.ErrPermission for this example
  if err := os.Open(filename); err != nil { 
     return errs.Append(err, "Could not open file")
  }
  ...
}
```
### Example: Converting a panic to an error
```go
func FuncThatShouldNeverPanic() (err error) {
   defer func() {
       if panicErr := errs.PanicToErr(recover()); panicErr != nil {
           err = panicErr
       }
   }
   ...
}
```

[Tylor Arndt]:https://plus.google.com/u/0/+TylorArndt/posts
[errors.go]:https://github.com/tarndt/sema/blob/master/errors.go]
