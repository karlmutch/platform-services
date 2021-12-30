package platform

// This file contains the implementation of code that checks to ensure
// that the local machine only has one entity accessing a named resource.
// This allows callers of this code to create and test for exclusive
// access to resources, or to check that only one instance of a
// process is running.

import (
	"fmt"
	"net"

	"github.com/go-stack/stack"
	"github.com/karlmutch/kv"
)

type Exclusive struct {
	Name     string
	ReleaseC chan struct{}
	listen   net.Listener
}

func NewExclusive(name string, quitC chan struct{}) (excl *Exclusive, err kv.Error) {

	excl = &Exclusive{
		Name:     name,
		ReleaseC: quitC,
	}

	// Construct an abstract name socket that allows the name to be recycled between process
	// restarts without needing to unlink etc. For more information please see
	// https://gavv.github.io/blog/unix-socket-reuse/, and
	// http://man7.org/linux/man-pages/man7/unix.7.html
	sockName := "@/tmp/"
	sockName += name

	errGo := fmt.Errorf("")
	excl.listen, errGo = net.Listen("unix", sockName)
	if errGo != nil {
		return nil, kv.Wrap(errGo).With("stack", stack.Trace().TrimRuntime())
	}
	go func() {
		go excl.listen.Accept()
		<-excl.ReleaseC
	}()
	return excl, nil
}
