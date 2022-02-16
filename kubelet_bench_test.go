package kubelet_bench

import (
	"testing"
	"time"

	"github.com/efficientgo/e2e"
	e2einteractive "github.com/efficientgo/e2e/interactive"
	e2emonitoring "github.com/efficientgo/e2e/monitoring"
	"github.com/efficientgo/tools/core/pkg/testutil"
)

// Run `bash build-kubelet.sh <path to kubernetes>` to build kubelet:latest image from source before running this test.
// This allows development and usage of different vesrions of kubelet.
func TestKubeletMetrics_E2e(t *testing.T) {
	e, err := e2e.NewDockerEnvironment("kubeletbench")
	testutil.Ok(t, err)
	t.Cleanup(e.Close)

	// Increase to 1 second to scrape kubelet every 1 second. This will be the source of kubelet traffic and test bench.
	mon, err := e2emonitoring.Start(e, e2emonitoring.WithScrapeInterval(1*time.Second))
	testutil.Ok(t, err)

	// CRI docker shim, because why make things too simple?
	// Kubelet dropped Docker support, so we need to start extra shim to make it use docker.
	// Kudos to https://github.com/Mirantis/cri-dockerd.

	testutil.Ok(t, e2e.StartAndWaitReady(e.Runnable("cri-dockerd").
		Init(e2e.StartOptions{
			Image:      "quay.io/bwplotka/cri-dockerd:v0.2.0",
			Command:    e2e.NewCommand("", "--docker-endpoint=unix:///var/run/docker.sock"),
			Privileged: true,
			Volumes:    []string{"/var/run:/var/run:rw", "/var/lib/docker/:/var/lib/docker:rw"},
		})))
	// Obviously this setup is crippled. It probably cannot create any pod etc. But what we care in this test is to
	// make sure kubelet have cadvisor running and gathering cgroups in my local system.
	// On my local machine this modest 932996 (~1MB) bytes of /cadvisor/metrics response.
	kubelet := e2e.NewInstrumentedRunnable(e, "kubelet").
		WithPorts(map[string]int{"http": 10250}, "http").
		WithMetricPath("/metrics/cadvisor").
		WithMetricScheme("https").
		Init(e2e.StartOptions{
			Image: "kubelet:latest",
			Command: e2e.NewCommand("",
				"--fail-swap-on=false",
				"--container-runtime-endpoint=unix:///var/run/cri-dockerd.sock",
			),
			Privileged: true,
			Volumes: []string{
				"/sys:/sys:ro",
				"/var/run:/var/run:ro",
				// I don't know why beside CRI-O kubelet still touches /var/lib/docker images... maybe tech debt?
				// 20:54:50 kubelet: E0215 19:54:50.427959  1 cri_stats_provider.go:455] "Failed to get the info of the filesystem with mountpoint"
				// err="failed to get device for dir \"/var/lib/docker\": stat failed on /var/lib/docker with error: no such file or directory" mountpoint="/var/lib/docker"
				"/var/lib/docker/:/var/lib/docker:ro",
			},
		})
	testutil.Ok(t, e2e.StartAndWaitReady(kubelet))

	// Open monitoring page with kubelet performance metrics.
	testutil.Ok(t, mon.OpenUserInterfaceInBrowser("/graph?g0.expr=rate(container_cpu_usage_seconds_total%7Bname%3D\"kubeletbench-kubelet\"%7D%5B1m%5D)&g0.tab=0&g0.stacked=0&g0.range_input=1h&g1.expr=container_memory_working_set_bytes%7Bname%3D\"kubeletbench-kubelet\"%7D&g1.tab=0&g1.stacked=0&g1.range_input=1h"))
	testutil.Ok(t, e2einteractive.RunUntilEndpointHit())
}
