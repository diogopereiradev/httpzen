package benchmark_module

import (
	"sync"
	"time"

	request_module "github.com/diogopereiradev/httpzen/internal/request"
)

type Metrics struct {
	ExecutedThreads       int32
	TotalRequests         int
	TotalErrors           int
	TotalSuccess          int
	TotalBytesSent        int
	TotalBytesReceived    int
	TotalDuration         int
	Duration              int
	RequestsMinLatency    float64
	RequestsMaxLatency    float64
	RequestsPerSecond     int
}

type BenchmarkOptions struct {
	Request       request_module.RequestOptions
	ThreadsAmount int
	Duration      int
}

type BenchmarkResult struct {
	Metrics Metrics
	mutex   sync.Mutex
}

func initialResultModel() *BenchmarkResult {
	return &BenchmarkResult{
		Metrics: Metrics{
			ExecutedThreads:       0,
			TotalRequests:         0,
			TotalErrors:           0,
			TotalSuccess:          0,
			TotalBytesSent:        0,
			TotalBytesReceived:    0,
			TotalDuration:         0,
			Duration:              0,
			RequestsMinLatency:    0,
			RequestsMaxLatency:    0,
			RequestsPerSecond:     0,
		},
	}
}

func RunBenchmark(options BenchmarkOptions, metrics *Metrics) {
	m := initialResultModel()
	m.Metrics.TotalDuration = options.Duration

	options.runThreads(m, metrics)
}

func (o *BenchmarkOptions) runThreads(model *BenchmarkResult, metrics *Metrics) {
	stop := make(chan struct{})
	done := make(chan struct{})

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for i := 0; i < o.Duration; i++ {
			<-ticker.C
			model.Metrics.Duration++
			*metrics = model.Metrics
		}
		close(done)
	}()

	for i := 0; i < o.ThreadsAmount; i++ {
		go func() {
			model.Metrics.ExecutedThreads++
			for {
				select {
				case <-stop:
					return
				default:
					res := o.doRequest(model)
					model.mutex.Lock()

					model.Metrics.TotalRequests++
					model.Metrics.TotalBytesReceived += len(res.Result)
					model.Metrics.TotalBytesSent += len(o.Request.Body) + len(o.Request.Headers) + len(o.Request.Url) + len(o.Request.Method)
	        model.Metrics.RequestsPerSecond = model.Metrics.TotalRequests / o.Duration

					if res.ExecutionTime < model.Metrics.RequestsMinLatency || model.Metrics.RequestsMinLatency == 0 {
						model.Metrics.RequestsMinLatency = res.ExecutionTime
					}

					if res.ExecutionTime > model.Metrics.RequestsMaxLatency {
						model.Metrics.RequestsMaxLatency = res.ExecutionTime
					}

					*metrics = model.Metrics
					model.mutex.Unlock()
				}
			}
		}()
	}
	<-done
	close(stop)
}

func (o *BenchmarkOptions) doRequest(model *BenchmarkResult) *request_module.RequestResponse {
	resp := request_module.RunRequest(request_module.RequestOptions{
		Method:      o.Request.Method,
		Url:         o.Request.Url,
		Headers:     o.Request.Headers,
		Body:        o.Request.Body,
		Timeout:     time.Duration(1 * time.Minute),
		BypassError: true,
	})

	model.mutex.Lock()
	defer model.mutex.Unlock()

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		model.Metrics.TotalSuccess++
	} else {
		model.Metrics.TotalErrors++
	}
	return &resp
}
