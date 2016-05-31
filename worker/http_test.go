package worker

import "testing"

func TestWorkerDo(t *testing.T) {
  url := "http://localhost:8080/test1"
  http_worker := HttpWorker{
    url,
    make(map[string]string),
    nil,
    nil,
    nil,
    nil,
  }

  pw := PartWork{start: 0, length: 100}
  http_worker.init()
  http_worker.do(&pw)
}