# kubelet-bench

An example of Go based e2e benchmark for various Kubelet operations without spinning up whole K8s cluster.

### Motivation

Sometimes you want to test or benchmark simple operation on one of the components from bigger system. In my case I needed to test and benchmark `kubelet` [metric endpoint optimizations](https://github.com/google/cadvisor/pull/2974) we made.

Usually to run newest `kubelet` from source involves running full Kubernetes cluster, ideally on some external cloud provider virtual machines. In order to do so I would probably hit thousands of unrelated issues and bottlenecks only to test very simple operation on `kubelet` which in calling `/cadvisor/metric` endpoint.

Instead of spinning full cluster I decided to use my dev machine, build `kubelet` from source, put in container and use our amazing [e2e framework](https://github.com/efficientgo/e2e) (shameless plug) that supports among other things, programmatic, interactive test/benchmarking with built-in Prometheus monitoring integration. The challenge here is run `kubelet` without any other component (e.g `kube-api`) other than `docker` engine.

It was not trivial for me, since I am not strictly a Kubernetes developer, so I had to reverse engineer many things. It might be not simpler for others too, so I decided to show how to do it in this repo. Enjoy!

### Prerequisite

* Linux machine
* `docker` installed
* `Go` 1.17+
* Cloned https://github.com/kubernetes/kubernetes somewhere locally.

### Usage

1. First of all run `bash build-kubelet.sh <local path to Kubernetes repo>`. This should build kubelet from the source you have and put that in local `docker:latest` image.

2. Check [./kubelet_bench_test.go](./kubelet_bench_test.go) which should be self-explanatory. It contains single Go test that you can run as usual using `go -v test ./...`. This tests starts interactive docker environment with `kubelet` and `Prometheus` (and some docker CRI-O shim) and spams it with couple of `/cadvisor/metrics` calls.

3. After that Prometheus UI should show up in your browser with relevant view on `kubelet` performance.

4. You can kill test by finding test output line `Waiting for user HTTP request on  http://<some local address> ...` and going to this address.
