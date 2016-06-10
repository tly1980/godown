package worker

import (
  "testing"

  "github.com/stretchr/testify/assert"
)


func Test_chunk_generator1(t *testing.T) {
  chunk_gen := new_chunk_generator(5, 3)

  assert.Equal(t, true, chunk_gen.has_next(), "should end")
  ch1 := chunk_gen.next()
  assert.Equal(t, true, chunk_gen.has_next(), "should end")
  assert.Equal(t, int64(0), ch1.start, "start should be 0")
  assert.Equal(t, int64(3), ch1.length, "start should be 0")

  assert.Equal(t, true, chunk_gen.has_next(), "should end")
  ch2 := chunk_gen.next()
  assert.Equal(t, false, chunk_gen.has_next(), "should end")
  assert.Equal(t, int64(3), ch2.start, "start should be 3")
  assert.Equal(t, int64(2), ch2.length, "start should be 2")

  ch3 := chunk_gen.next()
  assert.Equal(t, false, chunk_gen.has_next(), "should end")
  var nil_chunk *Chunk
  assert.Equal(t, nil_chunk, ch3, "ch3 should be nil")
}
