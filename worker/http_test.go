package worker

import "testing"

func TestWorkerDo(t *testing.T) {
  url := "http://localhost:8080/nsw_prices.csv"
  http_worker := HttpWorker{
    url,
    make(map[string]string),
    nil,
    nil,
    nil,
    nil,
  }

  pw := PartWork{0, 100, "bytes=0-99"}
  http_worker.init()
  http_worker.do(&pw)
}