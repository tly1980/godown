package worker

import (
  "bufio"
  "log"
  "os"
  "testing"
  "path"
  "path/filepath"

  "github.com/stretchr/testify/assert"
)

const TEST_DATA_FOLDER = "../test"
const TEST_BASE_URL = "http://localhost:8080/"


func fread(path string, start int64, length int64) []byte{
  file, err := os.Open(path)
  if (err != nil ) {
    log.Fatal(err)
  }

  buf := make([]byte, length)
  defer file.Close()
  file.Seek(start, os.SEEK_SET)

  reader := bufio.NewReader(file)

  nread, err := reader.Read(buf)

  if int64(nread) != length {
    log.Fatal("read byte count is not equal to requested length!")
  }

  return buf
}

func _test_worker_download(
    t *testing.T, fname string, start int64, length int64){
  url := TEST_BASE_URL + fname
  test_folder, err := filepath.Abs(path.Join("..", "test"))
  if err != nil {
    log.Fatal(err)
  }

  fpath := path.Join(test_folder, fname)

  ch_in := make(chan PartWork)
  ch_out := make(chan PartWork)
  defer close(ch_in)
  defer close(ch_out)

  client := new_http_client()
  http_worker := new_worker(ch_in, ch_out, client)

  pw := PartWork{
    start: uint64(start), length: uint64(length), url: url, 
  }
  go http_worker.run()

  ch_in <- pw
  pw2 := <-ch_out

  buf := fread(fpath, start, length)

  assert.Equal(t, 1, pw2.try_count, "should only try one time")
  assert.Equal(t, buf, pw2.buf[:pw2.length],  "they should be equal")

}

func TestWorkerDo(t *testing.T) {
  _test_worker_download(t, "test1", 0, 100)
  _test_worker_download(t, "test1", 10, 20)
  _test_worker_download(t, "test1", 1234, 2345)
}

// func TestHttpdownGetsize(t *testing.T) {
//   url := "http://localhost:8080/test1"
//   hd := HttpDownloader {client: &http.Client{}}

//   n_size, err := hd.GetSize(url)
//   if n_size != 4096 {
//     t.Error("n_size is not equal to what requested")
//   }
//   if err != nil {
//     t.Error("err is not nil")
//   }
// }