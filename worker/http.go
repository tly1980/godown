package worker

import (
  "bufio"
  "fmt"
  "log"
  "net/http"
  "io"
  "strconv"
  "net"
  "time"
  //"godown/writer"
)

const SIZE_WRITEBLOCK = 1024 * 100 * 4
const RNAGE_BYTES_FMT = "bytes=%d-%d"
const HTTP_TIMEOUT_SECS = 30
const WORKER_MAX_RETIREX = 3

type PartWork struct{
  url string
  cookie map[string] string
  start uint64
  length uint64
  n_read uint64
  buf []byte
  try_count int
}

type FetchWork struct {
  id string
  src_url string
  dst_url string
  cookie map[string] string
  try_count int
  status string
  size uint64
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
  ch_in chan PartWork
  ch_out chan PartWork
  ch_done chan bool
  client *http.Client
  buf []byte
}

func new_worker(
  ch_in chan PartWork, ch_out chan PartWork, ch_done chan bool,
  client *http.Client) *HttpWorker {

  return &HttpWorker{
    ch_in: ch_in,
    ch_out: ch_out,
    ch_done: ch_done,
    client: client,
    buf: make([]byte, SIZE_WRITEBLOCK),
  }
}


func (self *HttpWorker) run(){
  //Initialize the REQUEST
  for w := range self.ch_in {
    for i := 0; i < WORKER_MAX_RETIREX; i++ {
      w.try_count++
      n_read, err := self.do(&w)
      if (err == nil ){
        w.buf = self.buf
        w.n_read = uint64(n_read)
        break
      }else{
        log.Printf("err: %s", err)
      }
    }
    self.ch_out <- w
  }

  self.ch_done <- true
}


type ChunkGenerator struct{
  start uint64
  total_size uint64
  block_size uint64
}

type Chunk struct{
  start uint64
  length uint64
}

func new_chunk_generator(total_size uint64, block_size uint64) *ChunkGenerator {
  return &ChunkGenerator {
    start: uint64(0),
    total_size: total_size,
    block_size: block_size,
  }
}

func (self *ChunkGenerator) next() *Chunk {
  if !self.has_next() {
    return nil
  }

  ret := Chunk {}

  if ( self.total_size - self.start  > self.block_size ) {
    ret.start = self.start
    ret.length = self.block_size

  } else {
    ret.start = self.start
    ret.length = self.total_size - self.start
  }

  self.start += ret.length

  return &ret
}

func (self *ChunkGenerator) has_next() bool {
  if self.start >= self.total_size {
    return false
  }else{
    return true
  }
}

func (self *HttpWorker) do(pwork *PartWork) (int, error) {
  req, err := http.NewRequest("GET", pwork.url, nil)

  if err != nil {
    log.Printf("Failed to create http request to:%v, %v", pwork.url, err)
    return -1, err
  }

  for k, v := range pwork.cookie {
    req.Header.Set(k, v)
  }

  range_val := fmt.Sprintf(
      RNAGE_BYTES_FMT, pwork.start, pwork.start + pwork.length - 1)

  req.Header.Set("Range", range_val)
  resp, err := self.client.Do(req)

  if err != nil {
    return -1, err
  }

  defer resp.Body.Close()
  reader := bufio.NewReader(resp.Body)

  n_read, err := reader.Read(self.buf)
  if err != io.EOF && err != nil {
    log.Printf("Failed in reading from resp: %v", err)
    return -1, err
  }

  if uint64(n_read) != pwork.length {
    return -1, fmt.Errorf(
        "n_read: %v != PartWork.length: %v (requested)",
        n_read, pwork.length)
  }

  return n_read, nil
}

type HttpDownloader struct {
  worker_count int
  workers [] *HttpWorker
  src_url string
  dst_url string
  cookie map[string] string
  ch_pw_in chan PartWork
  ch_pw_out chan PartWork
  ch_done chan bool
  ch_done_worker chan bool
  client *http.Client
  try_count int
  status string
  res_size uint64
}

func NewHttpDownloader(
    worker_count int,
    src_url string, dst_url string,
    cookie map[string] string) *HttpDownloader {
  ret := & HttpDownloader{
    src_url: src_url,
    dst_url: dst_url,
    worker_count: worker_count,
    cookie: cookie,
    ch_done: make(chan bool),
    ch_pw_in: make(chan PartWork),
    ch_pw_out: make(chan PartWork),
    client: new_http_client(),
  }

  return ret
}


func (self *HttpDownloader) init(){
  for i :=0; i < self.worker_count; i++ {
    w := new_worker(self.ch_pw_in, self.ch_pw_out, self.ch_done, self.client)
    go w.run()
    self.workers = append(self.workers, w)
  }

  self.status = "init"
}



func (self *HttpDownloader) Fetch() {
  res_size, err := self.get_size()
  if ( err != nil ){
    log.Fatal(err)
  }
  self.res_size = res_size
}



func (self *HttpDownloader) get_size() (uint64, error) {
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
    size, err := strconv.ParseUint(val[0], 10, 64)

    if (err == nil ) {
      return size, nil
    }
  }

  return 0, err
}

