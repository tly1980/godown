package worker

import (
  "bufio"
  "fmt"
  "log"
  "os"
  "path"
  "path/filepath"
  "testing"

  "github.com/stretchr/testify/assert"
)

const TEST_DATA_FOLDER = "../test"
const TEST_BASE_URL = "http://localhost:8080/test/"
const TMP_BASE = "/tmp"


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
    log.Printf("%v\n", string(buf[:nread]))
    log.Fatal("read byte count is not equal to requested length:", nread, length, len(string(buf[:nread])))
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
  fmt.Printf("fpath: %v\n", fpath)

  ch_in := make(chan PartWork)
  ch_out := make(chan PartWork)
  ch_done := make(chan string)
  defer func () {
    close(ch_in)
    close(ch_out)
    // Don't get from ch_done before closing in and out channels
    <-ch_done
    close(ch_done)
  }()

  client := new_http_client()
  http_worker := new_worker("test_worker", ch_in, ch_out, ch_done, client)

  pw := PartWork{
    start: start, length: length, url: url,
  }
  go http_worker.run()

  ch_in <- pw
  pw2 := <-ch_out

  buf := fread(fpath, start, length)

  assert.Equal(t, 1, pw2.try_count, "should only try one time")
  assert.Equal(t, buf, pw2.buf,  "they should be equal")

}

func TestWorker_run(t *testing.T) {
  _test_worker_download(t, "test1", 0, 100)
  _test_worker_download(t, "test1", 10, 20)
  _test_worker_download(t, "test1", 1234, 2345)
  _test_worker_download(t, "test1", 0, 4096)
}

func TestHttpdown_get_size(t *testing.T) {
  url := TEST_BASE_URL + "test1"
  cookie := make(map[string]string)
  hd := NewHttpDownloader(1, url, "a", cookie)

  n_size, err := hd.get_size()
  if n_size != 4096 {
    t.Error("n_size is not equal to what requested")
  }
  if err != nil {
    t.Error("err is not nil")
  }
}

func _test_httpdown_doanlod(t *testing.T, fname string) {
  url := TEST_BASE_URL + fname
  dst := path.Join(TMP_BASE, fname)
  cookie := make(map[string]string)
  hd := NewHttpDownloader(2, url, dst, cookie)
  fmt.Printf("dst: %s\n", dst)
  hd.Fetch()
  hd.WaitTillDone()
}

func TestHttpdown_download(t *testing.T) {
  _test_httpdown_doanlod(t, "test1")
}
