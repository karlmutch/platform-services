package main

// This file contains tests for the DB layer.  In order for the test to run
// the DB flags must be set appropriately

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/go-test/deep"

	"github.com/golang/protobuf/proto"

	model "github.com/karlmutch/platform-services/internal/experiment"
	grpc "github.com/karlmutch/platform-services/internal/gen/experimentsrv"
	"github.com/karlmutch/platform-services/internal/platform"

	"github.com/go-stack/stack"
	"github.com/karlmutch/errors"
)

func TestDBA(t *testing.T) {

	timeout := time.Duration(time.Minute)
	giveUp := time.Now().Add(timeout)
	msg := fmt.Sprintf("first attempt failed, retrying the DB connection for %v", timeout)

	err := errors.New("")
	for {
		select {
		case <-time.After(time.Second):
			if err = model.GetDBStatus(); err != nil {
				if time.Now().After(giveUp) {
					t.Error(errors.Wrap(err).With("stack", stack.Trace().TrimRuntime()))
					return
				}
				if len(msg) != 0 {
					t.Log(msg)
					msg = ""
				}
				continue
			}
			return
		}
	}
}

// diffExp is provided so that grpc clone operations which do not respect the
// representation of zero length arrays versus nil's can be worked around during
// testing
//
func diffExp(l *grpc.Experiment, r *grpc.Experiment) (diffs []string) {
	rc := proto.Clone(r).(*grpc.Experiment)
	if len(rc.InputLayers) == 0 {
		if l.InputLayers != nil {
			rc.InputLayers = map[uint32]*grpc.InputLayer{}
		}
	}
	if len(rc.OutputLayers) == 0 {
		if l.OutputLayers != nil {
			rc.OutputLayers = map[uint32]*grpc.OutputLayer{}
		}
	}

	return deep.Equal(l, rc)
}

func newTestExperiment() (out *grpc.Experiment) {
	return &grpc.Experiment{
		Uid:          "test-only-" + platform.GetPseudoUUID(),
		Name:         "test-only-" + platform.GetPseudoUUID(),
		Description:  "test-only-" + platform.GetPseudoUUID(),
		InputLayers:  map[uint32]*grpc.InputLayer{},
		OutputLayers: map[uint32]*grpc.OutputLayer{},
	}
}

func TestDBExperimentSimple(t *testing.T) {

	if err := model.GetDBStatus(); err != nil {
		t.Error(errors.Wrap(err).With("stack", stack.Trace().TrimRuntime()).Error())
		return
	}

	in := newTestExperiment()

	exp, err := model.InsertExperiment(context.Context(context.Background()), in)
	if err != nil {
		t.Error(err.Error())
		return
	}

	// To check equivalence between the supplied data and the apparently written
	// data copy the two fields we know changed into the original data and then
	// do the deep comparison
	in.Created = exp.Created
	if diff := diffExp(in, exp); len(diff) != 0 {
		t.Error(errors.New(strings.Join(diff, ", ")).With("stack", stack.Trace().TrimRuntime()))
		return
	}

	selected, err := model.SelectExperiment(context.Context(context.Background()), 0, in.Uid)
	if err != nil {
		t.Error(err.With("uid", in.Uid).Error())
		return
	}
	if selected == nil {
		t.Error(errors.New("SelectExperimentWide returned no data unexpectedly").With("uid", in.Uid).Error())
		return
	}
	if diff := diffExp(in, selected); len(diff) != 0 {
		t.Error(errors.New(strings.Join(diff, ", ")).With("stack", stack.Trace().TrimRuntime()))
		return
	}

	// Now do a wide select even though we have no layers to test the simple case
	wide, err := model.SelectExperimentWide(context.Context(context.Background()), in.Uid)
	if err != nil {
		t.Error(err.With("uid", in.Uid).Error())
		return
	}
	if wide == nil {
		t.Error(errors.New("SelectExperimentWide returned no data unexpectedly").With("uid", in.Uid).Error())
		return
	}
	if diff := diffExp(in, wide); len(diff) != 0 {
		t.Error(errors.New(strings.Join(diff, ", ")).With("stack", stack.Trace().TrimRuntime()))
		return
	}

	if err = model.DeactivateExperiment(context.Context(context.Background()), in.Uid); err != nil {
		t.Error(err.With("uid", in.Uid).Error())
		return
	}

	// Try reinserting and make sure it fails
	if _, err = model.InsertExperiment(context.Context(context.Background()), in); err == nil {
		t.Error("failed tests due to reinsertion of a duplicate experiment working")
		return
	}
}

func TestDBExperimentWide(t *testing.T) {

	if err := model.GetDBStatus(); err != nil {
		t.Error(errors.Wrap(err).With("stack", stack.Trace().TrimRuntime()).Error())
		return
	}

	in := newTestExperiment()

	in.InputLayers[1] = &grpc.InputLayer{
		Name:   platform.GetPseudoUUID(),
		Type:   grpc.InputLayer_Enumeration,
		Values: []string{},
	}
	in.OutputLayers[0] = &grpc.OutputLayer{
		Name:   platform.GetPseudoUUID(),
		Type:   grpc.OutputLayer_Enumeration,
		Values: []string{},
	}
	in.OutputLayers[1] = &grpc.OutputLayer{
		Name:   platform.GetPseudoUUID(),
		Type:   grpc.OutputLayer_Probability,
		Values: []string{platform.GetPseudoUUID()},
	}
	in.OutputLayers[3] = &grpc.OutputLayer{
		Name:   platform.GetPseudoUUID(),
		Type:   grpc.OutputLayer_Raw,
		Values: []string{platform.GetPseudoUUID(), platform.GetPseudoUUID()},
	}

	exp, err := model.InsertExperiment(context.Context(context.Background()), in)
	if err != nil {
		t.Error(err.Error())
		return
	}
	// To check equivalence between the supplied data and the apparently written
	// data copy the two fields we know changed into the original data and then
	// do the deep comparison
	in.Created = exp.Created
	if diff := diffExp(in, exp); len(diff) != 0 {
		t.Error(errors.New(strings.Join(diff, ", ")).With("stack", stack.Trace().TrimRuntime()))
		return
	}

	// Now do a wide select to include our layers
	wide, err := model.SelectExperimentWide(context.Context(context.Background()), in.Uid)
	if err != nil {
		t.Error(err.With("uid", in.Uid).Error())
		return
	}
	if wide == nil {
		t.Error(errors.New("SelectExperimentWide returned no data unexpectedly").With("uid", in.Uid).Error())
		return
	}
	if diff := diffExp(in, wide); len(diff) != 0 {
		t.Error(errors.New(strings.Join(diff, ", ")).With("stack", stack.Trace().TrimRuntime()))
		return
	}

	if err = model.DeactivateExperiment(context.Context(context.Background()), in.Uid); err != nil {
		t.Error(err.With("uid", in.Uid).Error())
		return
	}
}
