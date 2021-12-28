package main

import (
	"flag"
	"os"
	"testing"
	"time"

	"github.com/karlmutch/envflag"
)

func init() {
}

var (
	parsedFlags = false

	TestStopC = make(chan bool)

	TestRunMain string
)

// TestRunMain can be used to run the server in production mode as opposed to
// funit or unit testing mode.  Traditionally gathering coverage data and running
// in production are done seperately.  This unit test allows the runner to do
// both at the same time.  To do this a test binary is generated using the command
//
// cd $(GOROOT)/src/github.com/karlmutch/platform-services/cmd/experimentsrv
// go test -coverpkg="." -c -o bin/experimentsrv-cpu-run-coverage
//
// Then the resulting /bin/experimentsrv-cpu-run-coverage binary is run as through it were a traditional
// server binary for the go experiment server using the command below.  The difference being that the
// binary now has coverage instrumentation.  In order to collect the coverage run any production
// workload and use cases you need then CTRL-C the server.
//
// ./bin/experimentsrv-cpu-run-coverage -test.run "^TestRunMain$" -test.coverprofile=system.out
//
// As an additional feature coverage files have is that they can also be merged using
// commands similar to the following:
//
// $ go get github.com/wadey/gocovmerge
// $ gocovmerge unit.out system.out > all.out
// $ go tool cover -html all.out
//
// Using the coverage merge tool testing done using a fully deployed system with
// real projects, proxies, projects, and workloads along with integration testing can be merged
// together from different test steps in an integration and test pipeline.
//

// TestMain is invoked by the GoLang entry point for the runtime of compiled GoLang
// programs when the compiled and linked image has been run using the 'go test'
// command
//
// This function will invoke the applications entry point to initiate the normal execution flow
// of the server with the tests remaining under the scheduling control of the
// GoLang test runtime. For more information please read https://golang.org/pkg/testing/
//
func TestMain(m *testing.M) {

	// Only perform this Parsed check inside the test framework. Do not be tempted
	// to do this in the main of our production package
	//
	if !flag.Parsed() {
		envflag.Parse()
	}
	parsedFlags = true

	quitC := make(chan struct{})
	doneC := make(chan struct{})

	resultCode := -1
	{
		// Start the server under test
		go func() {
			logger.Info("Starting Server")
			if errs := EntryPoint(quitC, doneC); len(errs) != 0 {
				for _, err := range errs {
					logger.Error(err.Error())
				}
				os.Exit(-1)
			}

			<-quitC

			// When using benchmarking in production mode, that is no tests running the
			// user can park the server on a single unit test that only completes when this
			// channel is close, which happens only when there is a quitC from the application
			// due to a CTRL-C key sequence or kill -n command
			//
			// If the test was not selected for by the tester then this will be essentially a
			// NOP
			//
			close(TestStopC)

			logger.Info("forcing test mode server down")
			func() {
				defer func() {
					recover()
				}()
				close(quitC)
			}()

		}()

		// Wait for the server to signal it is ready for work
		<-doneC

		// Wait for any depent modules to initialize completely
		listenerC := make(chan bool)
		defer close(listenerC)

		modules := &Modules{}
		modules.AddListener(listenerC)

		running := func() (running bool) {
			giveUpAt := time.Now().Add(30 * time.Second)
			for {
				select {
				case <-time.After(5 * time.Second):
					if giveUpAt.Before(time.Now()) {
						return false
					}
				case up := <-listenerC:
					if up {
						return true
					}
				}
			}
		}()

		if !running {
			logger.Fatal("not all modules initialized within the server")
			os.Exit(-1)
		}

		if len(TestRunMain) != 0 {
			<-TestStopC
		} else {
			logger.Info("running tests")
			resultCode = m.Run()
			logger.Info("finished running tests")

			close(quitC)
		}
	}

	logger.Info("waiting for server down to complete")

	// Wait until the main server is shutdown
	<-quitC

	if resultCode != 0 {
		os.Exit(resultCode)
	}
}
