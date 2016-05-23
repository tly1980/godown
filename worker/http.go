import (
  godown/writer
)

type Stat {

}

type PartWork {
  start int
  length int
}

type HttpWorker {
  url string
  cookie map[string] string
  work chan *PartWork
  wpipe chan *WriteTask
}

func (self *HttpWorker) run(){
  client =  &http.Client{}

  //Initialize the REQUEST
  
  for w := range self.work {
  }
}

func doWork(client *http.Client, pwork *PartWork){
  req, err := http.NewRequest("GET", self.url, nil)
  for k, v := range self.cookie {
    req.Header.Set(k, v)
  }

  
}
