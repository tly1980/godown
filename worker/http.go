import (
  "fmt"
  "strconv"

  "godown/writer"
)

const SIZE_WRITEBLOCK = 1024 * 100 * 4 

type Stat {

}

type PartWork {
  start int
  length int
  range_header string
}


type HttpWorker {
  url string
  cookie map[string] string
  work chan *PartWork
  wpipe chan *WriteTask
  client *http.Client
  buf []byte
}

func (self *HttpWorker) init(){
  self.buf := make([]byte, SIZE_WRITEBLOCK);
  self.client := &http.Client{}
}

func (self *HttpWorker) run(){
  //Initialize the REQUEST

  for w := range self.work {
    wtask := self.do(w)
    self.wpipe <- wtask
  }
}

func do(pwork *PartWork) *writer.WriteTask {
  req, err := http.NewRequest("GET", self.url, nil)
  for k, v := range self.cookie {
    req.Header.Set(k, v)
  }

  req.Header.Set('Range', pwork.range_header)
  resp, err := self.client.Do(req)

  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    log.Printf("Failed: {}", err)
  }

  wtask = &writer.WriteTask{
      &body,
      pwork.start
  }
  return wtask 
}

const RANGE_BATCH = 1024 * 1024 * 4

type GetOpsManager {
  url string
  cookie map[string] string
  work chan *PartWork
  client *http.Client
}

func run(){

  req, err := http.NewRequest("HEAD", self.url, nil)
  for k, v := range self.cookie {
    req.Header.Set(k, v)
  }

  resp, err := self.client.Do(req)
  cnt_length_str = resp.Header.Get('Content-Length')
  number, _ := strconv.Atoi(cnt_length_str, 10, 0)
}
