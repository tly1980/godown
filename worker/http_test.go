package worker

import (
  "bufio"
  "fmt"
  "log"
  "os"
  "testing"
  "path"
  "path/filepath"

  "github.com/stretchr/testify/assert"
)

const TEST_DATA_FOLDER = "../test"

func read(path string, start int64, length int64) []byte{
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
    log.Fatal("asd")
  }

  return buf
}

func TestWorkerDo(t *testing.T) {
  ch_in := make(chan PartWork)
  ch_out := make(chan PartWork)

  base_url := "http://localhost:8080/"
  base_path, err := filepath.Abs(path.Join("..", "test"))
  if err != nil {
    log.Fatal(err)
  }
  fname := "test1"

  url := base_url + fname
  fpath := path.Join(base_path, fname)

  fmt.Println("url: ", url)
  client := new_http_client()
  http_worker := new_worker(ch_in, ch_out, client)

  pw := PartWork{
    start: 0, length: 100, url: url, 
  }
  go http_worker.run()

  ch_in <- pw
  pw2 := <-ch_out

  defer close(ch_in)
  defer close(ch_out)

  buf := read(fpath, 0, 100)

  //fmt.Printf("%v, %v\n", pw2.try_count, pw2.buf[:pw2.length])

  assert.Equal(t, pw2.try_count, 1,  "they should be equal")
  assert.Equal(t, buf, pw2.buf[:pw2.length],  "they should be equal")
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