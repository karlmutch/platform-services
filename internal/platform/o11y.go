package platform

import (
	"context"

	"github.com/go-stack/stack"
	"github.com/honeycombio/opencensus-exporter/honeycomb"
	"github.com/karlmutch/kv"

	"go.opencensus.io/trace"
)

func StartOpenCensus(ctx context.Context, apiKey string, dataset string) (err kv.Error) {

	if len(apiKey) != 0 {
		exporter := honeycomb.NewExporter(apiKey, dataset)
		defer func() {
			go func() {
				select {
				case <-ctx.Done():
					exporter.Close()
				}
			}()
		}()

		trace.RegisterExporter(exporter)

		trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

		return nil
	}

	return kv.NewError("apiKey is missing, please supply a value for this parameter").With("stack", stack.Trace().TrimRuntime())
}
