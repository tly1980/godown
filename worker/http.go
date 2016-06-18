package worker

import (
  "bufio"
  "fmt"
  "io"
  "log"
  "net"
  "net/http"
  "os"
  "strconv"
  "time"
  "sync"

  "github.com/golang/glog"

  //"godown/writer"
)

const SIZE_WRITEBLOCK = 1024 * 100 * 4
const RNAGE_BYTES_FMT = "bytes=%d-%d"
const HTTP_TIMEOUT_SECS = 30
const WORKER_MAX_RETIREX = 3
const DEFAULT_REMOTE_CHUNK_SIZE = SIZE_WRITEBLOCK * 60

type PartWork struct{
  url string
  cookie map[string] string
  start int64
  length int64
  buf []byte
  try_count int
}


func new_http_client() *http.Client {
  netTransport := &http.Transport{
    Dial: (&net.Dialer{
      Timeout: 5 * time.Second,
    }).Dial,
    TLSHandshakeTimeout: 5 * time.Second,
  }
  return &http.Client{
    Timeout: time.Second * HTTP_TIMEOUT_SECS,
    Transport: netTransport,
  }
}

type HttpWorker struct {
  name string
  ch_in chan PartWork
  ch_out chan PartWork
  ch_done chan string
  client *http.Client
}

func new_worker(
  name string,  ch_in chan PartWork, ch_out chan PartWork, ch_done chan string,
  client *http.Client) *HttpWorker {

  return &HttpWorker{
    name: name,
    ch_in: ch_in,
    ch_out: ch_out,
    ch_done: ch_done,
    client: client,
  }
}


func (self *HttpWorker) run(){
  //Initialize the REQUEST
  for w := range self.ch_in {
    for i := 0; i < WORKER_MAX_RETIREX; i++ {
      w.try_count++
      buf, err := self.do(&w)
      if (err != nil ){
        log.Printf("err: %s", err)
        continue
      }else{
        w.buf = buf
        fmt.Println("done\n")
        break
      }
    }
    self.ch_out <- w
  }

  self.ch_done <- self.name
}


func (self *HttpWorker) do(pwork *PartWork) ([]byte, error) {
  req, err := http.NewRequest("GET", pwork.url, nil)

  glog.Info("9999 worker: %s: %v", self.name, pwork)

  if err != nil {
    log.Printf("Failed to create http request to:%v, %v", pwork.url, err)
    return nil, err
  }

  for k, v := range pwork.cookie {
    req.Header.Set(k, v)
  }

  range_val := fmt.Sprintf(
      RNAGE_BYTES_FMT, pwork.start, pwork.start + pwork.length - 1)

  req.Header.Set("Range", range_val)
  resp, err := self.client.Do(req)
  log.Printf("range_val: %s\n", range_val)

  if err != nil {
    return nil, err
  }

  defer resp.Body.Close()
  n_size := int(pwork.length)
  reader := bufio.NewReaderSize(resp.Body, n_size)

  buf, err := reader.Peek(n_size)
  if err != io.EOF && err != nil {
    log.Printf("Failed in reading from resp: %v", err)
    return nil, err
  }

  return buf, nil
}

type HttpDownloader struct {
  sync.Mutex
  worker_count int
  workers [] *HttpWorker
  src_url string
  dst_url string
  cookie map[string] string
  ch_pw_in chan PartWork
  ch_pw_out chan PartWork
  ch_done chan string
  client *http.Client
  try_count int
  status string
  res_size int64
  file *os.File
  chunk_generator ChunkGenerator
  write_count int64
  expect_write_count int64
  is_produced_end bool
  _thread_count int
}

func  NewHttpDownloader(
    worker_count int,
    src_url string, dst_url string,
    cookie map[string] string) *HttpDownloader {
  ret := & HttpDownloader{
    src_url: src_url,
    dst_url: dst_url,
    worker_count: worker_count,
    cookie: cookie,
    ch_done: make(chan string),
    ch_pw_in: make(chan PartWork),
    ch_pw_out: make(chan PartWork),
    client: new_http_client(),
  }

  return ret
}

func (self *HttpDownloader) _init() error {
  file, err := os.Create(self.dst_url)
  if err != nil {
    glog.Errorf("Failed to create file: %s, %s", self.dst_url, err)
    return err
  }
  self.file = file

  for i :=0; i < self.worker_count; i++ {
    wname := fmt.Sprintf("w:%d", i)
    w := new_worker(
        wname, self.ch_pw_in, self.ch_pw_out, self.ch_done, self.client)
    go w.run()
    self.workers = append(self.workers, w)
  }

  self.status = "init"
  return nil
}

func (self *HttpDownloader) _thread(fn func(), name string) {
  self._thread_count ++
  go func() {
    fn()
    self.ch_done <- name
  }()
}

func (self *HttpDownloader) clean_up() {
  close(self.ch_pw_in)
  close(self.ch_pw_out)
}

func (self *HttpDownloader) Fetch() {
  res_size, err := self.get_size()
  if ( err != nil ){
    log.Fatal(err)
  }
  self.res_size = res_size

  err = self._init()
  if err != nil {
    log.Fatal(err)
  }

  self._thread(self._chunk_producer, "_chunk_producer")
  self._thread(self._storer, "_storer")
  self._thread(self._checker, "_checker")
}

func (self *HttpDownloader) WaitTillDone() {
  for i := 0; i < self.worker_count + self._thread_count ; i++ {
    finisher := <- self.ch_done
    fmt.Printf("done: %s\n", finisher)
  }
}


func (self *HttpDownloader) _chunk_producer() {
  chunk_generator := new_chunk_generator(self.res_size, DEFAULT_REMOTE_CHUNK_SIZE)

  for chunk_generator.has_next() {
    chunk := chunk_generator.next()
    pw := PartWork {
      url: self.src_url,
      cookie: self.cookie,
      start: chunk.start,
      length: chunk.length,
      try_count: 0,
    }

    fmt.Printf("aaa pw: %v\n", pw)

    self.ch_pw_in <- pw
    fmt.Printf("bbb pw: %v\n", chunk_generator.has_next())
    self.Lock()
    self.expect_write_count += 1
    self.Unlock()
  }

  fmt.Printf("ccc\n")

  self.Lock()
  fmt.Printf("ddd\n")
  self.is_produced_end = true
  self.Unlock()
  fmt.Printf("eee\n")

   fmt.Println("gen end")
}

func (self *HttpDownloader) _storer() {
  defer self.file.Close()

  writer := bufio.NewWriter(self.file)
  defer writer.Flush()

  for pw := range self.ch_pw_out { 
    fmt.Printf("got pw out, start:%v, length: %v\n", pw.start, pw.length)
    self.file.Seek(pw.start, 0)
    writer.Write(pw.buf)
    fmt.Println("finished writing")
    self.Lock()
    self.write_count++
    self.Unlock()
    fmt.Printf("write: %v - %v \n", self.is_produced_end, self.write_count)
  }
}

func (self *HttpDownloader) _checker() {
  defer self.clean_up()

  for {
     time.Sleep(time.Millisecond * 500)
     if self.is_produced_end && self.write_count >= self.expect_write_count {
      fmt.Printf("stop\n")
      break
    }
  }
}


func (self *HttpDownloader) get_size() (int64, error) {
  req, err := http.NewRequest("HEAD", self.src_url, nil)

  if err != nil {
    log.Printf("Failed to create http request to:%v, %v", self.src_url, err)
    return 0, err
  }

  for k, v := range self.cookie {
    req.Header.Set(k, v)
  }

  resp, err := self.client.Do(req)

  if err != nil {
    return 0, err
  }

  defer resp.Body.Close()

  if val, ok := resp.Header["Content-Length"]; ok {
    size, err := strconv.ParseInt(val[0], 10, 64)

    if (err == nil ) {
      return size, nil
    }
  }

  return 0, err
}

