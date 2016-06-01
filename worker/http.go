package worker

import (
  "bufio"
  "fmt"
  "log"
  "net/http"
  "io"
  //"godown/writer"
)

const SIZE_WRITEBLOCK = 1024 * 100 * 4 
const RNAGE_BYTES_FMT = "bytes=%d-%d"

type PartWork struct{
  start int
  length int
  buf []byte
}


type HttpWorker struct {
  url string
  cookie map[string] string
  ch_partwork_get chan PartWork
  ch_partwork_store chan PartWork
  client *http.Client
  buf []byte
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
  //log.Printf("self.buf", self.buf)
  n_read, err := reader.Read(self.buf)
  if err != io.EOF && err != nil {
    log.Printf("Failed in reading from resp: %v", err)
    return -1, err
  }

  if n_read != pwork.length {
    return -1, fmt.Errorf("n_read: %v != PartWork.length: %v (requested)",
        n_read, pwork.length)
  }

  log.Printf("%v", string(self.buf[0:n_read]))

  return n_read, nil
}

type HttpManager struct {
  workers [] *HttpWorker

  ch_partwork_get chan PartWork
  ch_partwork_store chan PartWork
}

