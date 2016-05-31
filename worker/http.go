package worker

import (
  "bufio"
  "fmt"
  "log"
  "net/http"

  //"godown/writer"
)

const SIZE_WRITEBLOCK = 1024 * 100 * 4 

type PartWork struct{
  start int
  length int
  buf []byte
}


type HttpWorker struct {
  url string
  cookie map[string] string
  work chan PartWork
  wpipe chan PartWork
  client *http.Client
  buf []byte
}

func (self *HttpWorker) init(){
  self.client = &http.Client{}
  self.buf = make([]byte, SIZE_WRITEBLOCK);
}

func (self *HttpWorker) run(){
  //Initialize the REQUEST

  for w := range self.work {
//  
    log.Printf("{}", w)
  }
}

func (self *HttpWorker) do(pwork *PartWork) (int, error) {
  req, err := http.NewRequest("GET", self.url, nil)

  if err != nil {
    log.Fatal("aaa %v", err)
  }

  for k, v := range self.cookie {
    req.Header.Set(k, v)
  }

  range_val := fmt.Sprintf("bytes=%d-%d", 
      pwork.start, pwork.start + pwork.length - 1)

  log.Printf("range_val: %v", range_val)

  req.Header.Set("Range", range_val)
  resp, err := self.client.Do(req)

  if err != nil {
    log.Fatal("bbb %v", err)
  }

  defer resp.Body.Close()
  reader := bufio.NewReader(resp.Body)

  n_read, err := reader.Read(self.buf)
  if err != nil {
    log.Printf("Failed: %v", err)
  }

  if n_read != pwork.length {
    log.Printf("Failed: n_read != requested length. ", 
      n_read, pwork.length)
  }

  log.Printf("%v", string(self.buf[0:n_read]))

  return n_read, err
}
