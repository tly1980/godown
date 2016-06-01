package worker

import "testing"

func TestWorkerDo(t *testing.T) {
  url := "http://localhost:8080/test1"
  http_worker := HttpWorker{
    url: url,
  }

  pw := PartWork{start: 0, length: 100}
  http_worker.init()
  n_read, err := http_worker.do(&pw)
  if n_read != 100 {
    t.Error("n_read is not equal to what requested")
  }
  if err != nil {
    t.Error("err is not nil")
  }
}