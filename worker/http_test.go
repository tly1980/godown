package worker

import (
  "testing"
  "net/http"
)

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

func TestHttpdownGetsize(t *testing.T) {
  url := "http://localhost:8080/test1"
  hd := HttpDownloader {client: &http.Client{}}

  n_size, err := hd.GetSize(url)
  if n_size != 4096 {
    t.Error("n_size is not equal to what requested")
  }
  if err != nil {
    t.Error("err is not nil")
  }
}