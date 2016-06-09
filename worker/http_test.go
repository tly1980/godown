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
  fmt.Printf("fpath: %v\n", fpath)

  ch_in := make(chan PartWork)
  ch_out := make(chan PartWork)
  ch_done := make(chan bool)
  defer func () {
    close(ch_in)
    close(ch_out)
    // Don't get from ch_done before closing in and out channels
    <-ch_done
    close(ch_done)
  }()

  client := new_http_client()
  http_worker := new_worker(ch_in, ch_out, ch_done, client)

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

func TestWorker_run(t *testing.T) {
  _test_worker_download(t, "test1", 0, 100)
  _test_worker_download(t, "test1", 10, 20)
  _test_worker_download(t, "test1", 1234, 2345)
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

func Test_chunk_generator1(t *testing.T) {
  chunk_gen := new_chunk_generator(5, 3)
  ch1 := chunk_gen.next()
  assert.Equal(t, uint64(0), ch1.start, "start should be 0")
  assert.Equal(t, uint64(3), ch1.length, "start should be 0")

  ch2 := chunk_gen.next()
  assert.Equal(t, uint64(3), ch2.start, "start should be 3")
  assert.Equal(t, uint64(2), ch2.length, "start should be 2")

  ch3 := chunk_gen.next()
  var exp *Chunk
  assert.Equal(t, exp, ch3, "ch3 should be nil")
}
