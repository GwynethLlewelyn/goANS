package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/op/go-logging"
	"gopkg.in/natefinch/lumberjack.v2"
)

var log = logging.MustGetLogger("goans")	// configuration for the go-logging logger, must be available everywhere.
var logFormat logging.Formatter

/*
			   .__
  _____ _____  |__| ____
 /	   \\__	 \ |  |/	\
|  Y Y	\/ __ \|  |	  |	 \
|__|_|	(____  /__|___|	 /
	  \/	 \/		   \/
*/
// main() starts here.
func main() {
	// Setup the lumberjack rotating logger. This is because we need it for the go-logging logger when writing to files. (copied over from gwyneth 20170813).
	rotatingLogger := &lumberjack.Logger{
		Filename:	"goans.log",
		MaxSize:	1, // megabytes
		MaxBackups:	3,
		MaxAge:		28, //days
	}

	// Set formatting for stderr and file (basically the same).
	logFormat := logging.MustStringFormatter(`%{color}%{time:2006/01/02 15:04:05.0} %{shortfile} - %{shortfunc} ▶ %{level:.4s}%{color:reset} %{message}`) 	// must be initialised or all hell breaks loose

	// Setup the go-logging Logger. Do **not** log to stderr if running as FastCGI!
	// Note: we're just running a standalone server on a plain and simple port (gwyneth 20230522)
	backendFile				:= logging.NewLogBackend(rotatingLogger, "", 0)
	backendFileFormatter	:= logging.NewBackendFormatter(backendFile, logFormat)
	backendFileLeveled 		:= logging.AddModuleLevel(backendFileFormatter)
	backendFileLeveled.SetLevel(logging.DEBUG, "goans")	// we just send debug data to logs if we run as shell; NOTE(gwyneth): testing from where the requests are coming. (20170915)

	backendStderr			:= logging.NewLogBackend(os.Stderr, "", 0)
	backendStderrFormatter	:= logging.NewBackendFormatter(backendStderr, logFormat)
	backendStderrLeveled 	:= logging.AddModuleLevel(backendStderrFormatter)
	backendStderrLeveled.SetLevel(logging.DEBUG, "goans")	// shell is meant to be for debugging mostly
	logging.SetBackend(backendStderrLeveled, backendFileLeveled)

	log.Info("goans started and logging is set up.")
}

// handleHomepage will just show a form to allow entering translations for simple testing.
func handleHomepage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "meow? %q", r.URL.Path[1:])
}
// checkErrPanic logs a fatal error and panics.
func checkErrPanic(err error) {
	if err != nil {
		pc, file, line, ok := runtime.Caller(1)
		log.Panic(filepath.Base(file), ":", line, ":", pc, ok, " - panic:", err)
	}
}
// checkErr checks if there is an error, and if yes, it logs it out and continues.
//
// This is for 'normal' situations when we want to get a log if something goes wrong but do not need to panic
func checkErr(err error) {
	if err != nil {
		pc, file, line, ok := runtime.Caller(1)
		log.Error(filepath.Base(file), ":", line, ":", pc, ok, " - error:", err)
	}
}

// Auxiliary functions for HTTP handling

// checkErrHTTP returns an error via HTTP and also logs the error.
func checkErrHTTP(w http.ResponseWriter, httpStatus int, errorMessage string, err error) {
	if err != nil {
		http.Error(w, fmt.Sprintf(errorMessage, err), httpStatus)
		pc, file, line, ok := runtime.Caller(1)
		log.Error("(", http.StatusText(httpStatus), ") ", filepath.Base(file), ":", line, ":", pc, ok, " - error:", errorMessage, err)
	}
}
// checkErrPanicHTTP returns an error via HTTP and logs the error with a panic.
func checkErrPanicHTTP(w http.ResponseWriter, httpStatus int, errorMessage string, err error) {
	if err != nil {
		http.Error(w, fmt.Sprintf(errorMessage, err), httpStatus)
		pc, file, line, ok := runtime.Caller(1)
		log.Panic("(", http.StatusText(httpStatus), ") ", filepath.Base(file), ":", line, ":", pc, ok, " - panic:", errorMessage, err)
	}
}
// logErrHTTP assumes that the error message was already composed and writes it to HTTP and logs it.
//	this is mostly to avoid code duplication and make sure that all entries are written similarly
func logErrHTTP(w http.ResponseWriter, httpStatus int, errorMessage string) {
	http.Error(w, errorMessage, httpStatus)
	log.Error("(" + http.StatusText(httpStatus) + ") " + errorMessage)
}
// funcName is @Sonia's solution to get the name of the function that Go is currently running.
//	This will be extensively used to deal with figuring out where in the code the errors are!
//	Source: https://stackoverflow.com/a/10743805/1035977 (20170708)
func funcName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}
