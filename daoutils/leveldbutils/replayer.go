package leveldbutils

/**
 * @Author: lee
 * @Description:
 * @File: replayer
 * @Date: 2023-01-13 4:18 下午
 */

// replayer is a small wrapper to implement the correct replay methods.
type replayer struct {
	writer  KeyValueWriter
	failure error
}

// Put inserts the given value into the key-value data store.
func (r *replayer) Put(key, value []byte) {
	// If the replay already failed, stop executing ops
	if r.failure != nil {
		return
	}
	r.failure = r.writer.Put(key, value)
}

// Delete removes the key from the key-value data store.
func (r *replayer) Delete(key []byte) {
	// If the replay already failed, stop executing ops
	if r.failure != nil {
		return
	}
	r.failure = r.writer.Delete(key)
}
