package worker

type ChunkGenerator struct{
  start int64
  total_size int64
  block_size int64
}

type Chunk struct{
  start int64
  length int64
}

func new_chunk_generator(total_size int64, block_size int64) *ChunkGenerator {
  return &ChunkGenerator {
    start: 0,
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
