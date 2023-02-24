package leveldbutils

/**
 * @Author: lee
 * @Description:
 * @File: leveldb
 * @Date: 2023-01-12 7:21 下午
 */
import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type LevelDB struct {
	db *leveldb.DB
}

func NewLevelDB(file string, cache int, handles int, readonly bool) (Database, error) {
	return NewCustom(file, func(options *opt.Options) {
		// Ensure we have some minimal caching and file guarantees
		if cache < minCache {
			cache = minCache
		}
		if handles < minHandles {
			handles = minHandles
		}
		// Set default options
		options.OpenFilesCacheCapacity = handles
		options.BlockCacheCapacity = cache / 2 * opt.MiB
		options.WriteBuffer = cache / 4 * opt.MiB // Two of these are used internally
		if readonly {
			options.ReadOnly = true
		}
	})
}

func NewCustom(file string, customize func(options *opt.Options)) (Database, error) {
	options := configureOptions(customize)
	// Open the db and recover any potential corruptions
	db, err := leveldb.OpenFile(file, options)
	if _, corrupted := err.(*errors.ErrCorrupted); corrupted {
		db, err = leveldb.RecoverFile(file, nil)
	}
	if err != nil {
		return nil, err
	}
	// Assemble the wrapper with all the registered metrics
	ldb := &LevelDB{
		db: db,
	}
	return ldb, nil
}

// configureOptions sets some default options, then runs the provided setter.
func configureOptions(customizeFn func(*opt.Options)) *opt.Options {
	// Set default options
	options := &opt.Options{
		Filter:                 filter.NewBloomFilter(10),
		DisableSeeksCompaction: true,
	}
	// Allow caller to make custom modifications to the options
	if customizeFn != nil {
		customizeFn(options)
	}
	return options
}

// Has retrieves if a key is present in the key-value store.
func (db *LevelDB) Has(key []byte) (bool, error) {
	return db.db.Has(key, nil)
}

// Get retrieves the given key if it's present in the key-value store.
func (db *LevelDB) Get(key []byte) ([]byte, error) {
	dat, err := db.db.Get(key, nil)
	if err != nil {
		return nil, err
	}
	return dat, nil
}

// Put inserts the given value into the key-value store.
func (db *LevelDB) Put(key []byte, value []byte) error {
	return db.db.Put(key, value, nil)
}

// Delete removes the key from the key-value store.
func (db *LevelDB) Delete(key []byte) error {
	return db.db.Delete(key, nil)
}

// NewIterator creates a binary-alphabetical iterator over a subset
// of database content with a particular key prefix, starting at a particular
// initial key (or after, if it does not exist).
func (db *LevelDB) NewIterator(prefix []byte, start []byte) iterator.Iterator {
	return db.db.NewIterator(bytesPrefixRange(prefix, start), nil)
}

// NewBatch creates a write-only key-value store that buffers changes to its host
// database until a final write is called.
func (db *LevelDB) NewBatch() Batch {
	return &batch{
		db:    db.db,
		batch: new(leveldb.Batch),
	}
}

// Stat returns a particular internal stat of the database.
func (db *LevelDB) Stat(property string) (string, error) {
	return db.db.GetProperty(property)
}

// Compact flattens the underlying data store for the given key range. In essence,
// deleted and overwritten versions are discarded, and the data is rearranged to
// reduce the cost of operations needed to access them.
//
// A nil start is treated as a key before all keys in the data store; a nil limit
// is treated as a key after all keys in the data store. If both is nil then it
// will compact entire data store.
func (db *LevelDB) Compact(start []byte, limit []byte) error {
	return db.db.CompactRange(util.Range{Start: start, Limit: limit})
}

// Close stops the metrics collection, flushes any pending data to disk and closes
// all io accesses to the underlying key-value store.
func (db *LevelDB) Close() error {
	//db.quitLock.Lock()
	//defer db.quitLock.Unlock()
	//
	//if db.quitChan != nil {
	//	errc := make(chan error)
	//	db.quitChan <- errc
	//	if err := <-errc; err != nil {
	//		db.log.Error("Metrics collection failed", "err", err)
	//	}
	//	db.quitChan = nil
	//}
	return db.db.Close()
}
