package leveldbutils

/**
 * @Author: lee
 * @Description:
 * @File: batch
 * @Date: 2023-01-13 4:13 下午
 */
import (
	"github.com/syndtr/goleveldb/leveldb"
)

type batch struct {
	db    *leveldb.DB
	batch *leveldb.Batch
	size  int
}

// Put inserts the given value into the batch for later committing.
func (b *batch) Put(key, value []byte) error {
	b.batch.Put(key, value)
	b.size += len(key) + len(value)
	return nil
}

// Delete inserts the a key removal into the batch for later committing.
func (b *batch) Delete(key []byte) error {
	b.batch.Delete(key)
	b.size += len(key)
	return nil
}

// ValueSize retrieves the amount of data queued up for writing.
func (b *batch) ValueSize() int {
	return b.size
}

// Write flushes any accumulated data to disk.
func (b *batch) Write() error {
	return b.db.Write(b.batch, nil)
}

// Reset resets the batch for reuse.
func (b *batch) Reset() {
	b.batch.Reset()
	b.size = 0
}

// Replay replays the batch contents.
func (b *batch) Replay(w KeyValueWriter) error {
	return b.batch.Replay(&replayer{writer: w})
}
