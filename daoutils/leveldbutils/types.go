package leveldbutils

/**
 * @Author: lee
 * @Description:
 * @File: types
 * @Date: 2023-01-12 7:21 下午
 */
import (
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
	"io"
)

type LevelDBConfig struct {
	Path    string `mapstructure:"path"        json:"path"       yaml:"path"`
	Cache   int    `mapstructure:"cache"        json:"cache"       yaml:"cache"`
	Handles int    `mapstructure:"handles"        json:"handles"       yaml:"handles"`
}

type Database interface {
	Reader
	Writer
	Batcher
	Stater
	Compacter

	io.Closer
}

// Reader contains the methods required to read data from both key-value as well as
// immutable ancient data.
type Reader interface {
	KeyValueReader
	NewIterator(prefix []byte, start []byte) iterator.Iterator
}

// Writer contains the methods required to write data to both key-value as well as
// immutable ancient data.
type Writer interface {
	KeyValueWriter
}

// KeyValueReader wraps the Has and Get method of a backing data store.
type KeyValueReader interface {
	// Has retrieves if a key is present in the key-value data store.
	Has(key []byte) (bool, error)

	// Get retrieves the given key if it's present in the key-value data store.
	Get(key []byte) ([]byte, error)
}

// KeyValueWriter wraps the Put method of a backing data store.
type KeyValueWriter interface {
	// Put inserts the given value into the key-value data store.
	Put(key []byte, value []byte) error

	// Delete removes the key from the key-value data store.
	Delete(key []byte) error
}

// Stater wraps the Stat method of a backing data store.
type Stater interface {
	// Stat returns a particular internal stat of the database.
	Stat(property string) (string, error)
}

// Compacter wraps the Compact method of a backing data store.
type Compacter interface {
	// Compact flattens the underlying data store for the given key range. In essence,
	// deleted and overwritten versions are discarded, and the data is rearranged to
	// reduce the cost of operations needed to access them.
	//
	// A nil start is treated as a key before all keys in the data store; a nil limit
	// is treated as a key after all keys in the data store. If both is nil then it
	// will compact entire data store.
	Compact(start []byte, limit []byte) error
}

// bytesPrefixRange returns key range that satisfy
// - the given prefix, and
// - the given seek position
func bytesPrefixRange(prefix, start []byte) *util.Range {
	r := util.BytesPrefix(prefix)
	r.Start = append(r.Start, start...)
	return r
}
