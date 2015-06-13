package queue 
import (
	"errors"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"time"
	"math/rand"
	"strconv"
)
type Queue struct {
	backend *leveldb.DB
	random *rand.Rand
	Prefix []byte
}

var (
	DBCache map[string]*leveldb.DB = make(map[string]*leveldb.DB)
        ValueNotFound error = errors.New("No values are found.")
)

func New(path string, prefix string) (q *Queue, err error) {
	backend := DBCache[path]
	if backend == nil {
		backend, err = leveldb.OpenFile(path, nil)
		if err != nil {
			return
		}
		DBCache[path] = backend
	}
	random := rand.New(rand.NewSource(time.Now().Unix()))
	q = &Queue{backend: backend, random: random, Prefix: []byte(prefix)}
	return
}

func CloseAll() {
	for _, db := range DBCache {
		db.Close()
	}
}

func (q *Queue) Push(val []byte) (err error) {
	key := q.generateKey()
	return q.backend.Put(key, val, nil)
}

func (q *Queue) generateKey() (key []byte) {
	key = q.Prefix
	nano := time.Now().UnixNano()
	key = strconv.AppendInt(key, nano, 16)
	key = strconv.AppendInt(key, q.random.Int63(), 16)
	return
}

func (q *Queue) Pop() (val []byte, err error) {
	backend := q.backend
	iter := backend.NewIterator(util.BytesPrefix(q.Prefix), nil)
	if iter.Next() {
		val := iter.Value()
		err = backend.Delete(iter.Key(), nil)
		if err != nil {
			return val, err
		}
		return val, err
	}

	return val, ValueNotFound
}

func (q *Queue) List()(list [][]byte) {
	iter := q.backend.NewIterator(util.BytesPrefix(q.Prefix), nil)
	for iter.Next() {
		list = append(list, iter.Value())
	}
	return
}
