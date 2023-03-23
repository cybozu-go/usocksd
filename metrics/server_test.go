package metrics

import (
	"context"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"

	"github.com/cybozu-go/well"
	"github.com/prometheus/client_golang/prometheus"
)

func TestMetricsServer(t *testing.T) {
	t.Parallel()

	addr := "http://localhost:30080/metrics"
	env := well.NewEnvironment(context.Background())
	s := &Server{
		Env: env,
	}
	ln, err := net.Listen("tcp", ":30080")
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Serve(ln); err != nil {
		t.Fatal(err)
	}
	dummyMetric := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "dummy_counter",
	})
	if err := Registry.Register(dummyMetric); err != nil {
		t.Error(err)
	}
	dummyMetric.Add(1)

	resp, err := http.Get(addr)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	output := string(body)
	if !strings.Contains(output, "dummy_counter 1") {
		t.Errorf("could not get expected metric, got: %s", output)
	}
}
