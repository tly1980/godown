package worker

import (
  "bufio"
  "fmt"
  "log"
  "net/http"
  "io"
  "strconv"
  //"godown/writer"
)

const SIZE_WRITEBLOCK = 1024 * 100 * 4 
const RNAGE_BYTES_FMT = "bytes=%d-%d"
const HTTP_TIMEOUT_SECS = 30

type PartWork struct{
  url string
  cookie map[string] string
  start uint64
  length uint64
  buf []byte
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
  ch_partwork_get chan PartWork
  ch_partwork_store chan PartWork
  client *http.Client
  buf []byte
}

func new_worker(url string, client *http.Client){
  return &HttpWorker{
    url: url,
    client: client
  }
}

func (self *HttpWorker) init(){
  self.client = &http.Client{}
  self.buf = make([]byte, SIZE_WRITEBLOCK);
}

func (self *HttpWorker) run(){
  //Initialize the REQUEST

  for w := range self.ch_partwork_get {
    self.do(&w)
    w.buf = self.buf
    self.ch_partwork_store <- w
  }
}

func (self *HttpWorker) do(pwork *PartWork) (int, error) {
  req, err := http.NewRequest("GET", self.url, nil)

  if err != nil {
    log.Printf("Failed to create http request to:%v, %v", self.url, err)
    return -1, err
  }

  for k, v := range self.cookie {
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

  log.Printf("%v", string(self.buf[0:n_read]))

  return n_read, nil
}

type HttpDownloader struct {
  worker_count int
  workers [] *HttpWorker
  cookie map[string] string
  ch_work chan string
  ch_partwork_get chan PartWork
  ch_partwork_store chan PartWork
  ch_done chan bool
  ch_done_worker chan bool
  client *http.Client
}

func NewHttpDownloader(worker_count int, ch_work chan, ch_done chan) * {
  hdownloaer = & HttpDownloader{
    worker_count: worker_count,
    ch_work: ch_work,
    ch_done: ch_done,
    client: new_http_client()
  }
}

func (self *NewHttpDownloader) Init() {
  for (var i = 0; i < self.worker_count; i ++ ){

  }
}


func (self *HttpDownloader) get_size(url string) (uint64, error) {
  req, err := http.NewRequest("HEAD", url, nil)

  if err != nil {
    log.Printf("Failed to create http request to:%v, %v", url, err)
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


func (self *HttpDownloader) GetSize(url string) (uint64, error) {
  req, err := http.NewRequest("HEAD", url, nil)

  if err != nil {
    log.Printf("Failed to create http request to:%v, %v", url, err)
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

func (self *HttpDownloader) Fetch(url string, path string) (uint64, error) {
}
