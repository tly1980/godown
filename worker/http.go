package worker

import (
  "bufio"
  "godown/writer"
  "net/http"
  "log"
)

const SIZE_WRITEBLOCK = 1024 * 100 * 4 

type PartWork struct{
  start int
  length int
  range_header string
}


type HttpWorker struct {
  url string
  cookie map[string] string
  work chan *PartWork
  wpipe chan *writer.WriteTask
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

  req.Header.Set("Range", pwork.range_header)
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

// const RANGE_BATCH = 1024 * 1024 * 4

// type GetOpsManager {
//   url string
//   cookie map[string] string
//   work chan *PartWork
//   client *http.Client
// }

// func run(){

//   req, err := http.NewRequest("HEAD", self.url, nil)
//   for k, v := range self.cookie {
//     req.Header.Set(k, v)
//   }

//   resp, err := self.client.Do(req)
//   cnt_length_str = resp.Header.Get('Content-Length')
//   number, _ := strconv.Atoi(cnt_length_str, 10, 0)
// }
