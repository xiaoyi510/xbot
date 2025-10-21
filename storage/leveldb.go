package storage

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// LevelDB LevelDB 存储实现
type LevelDB struct {
	db   *leveldb.DB
	path string
}

// NewLevelDB 创建 LevelDB 存储
func NewLevelDB(path string) (*LevelDB, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}

	return &LevelDB{
		db:   db,
		path: path,
	}, nil
}

// Get 获取值
func (l *LevelDB) Get(key string) ([]byte, error) {
	data, err := l.db.Get([]byte(key), nil)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return data, nil
}

// Set 设置值
func (l *LevelDB) Set(key string, value []byte) error {
	return l.db.Put([]byte(key), value, nil)
}

// Delete 删除值
func (l *LevelDB) Delete(key string) error {
	return l.db.Delete([]byte(key), nil)
}

// Has 判断是否存在
func (l *LevelDB) Has(key string) bool {
	has, err := l.db.Has([]byte(key), nil)
	if err != nil {
		return false
	}
	return has
}

// Keys 获取所有以 prefix 开头的 key
func (l *LevelDB) Keys(prefix string) ([]string, error) {
	var keys []string

	iter := l.db.NewIterator(util.BytesPrefix([]byte(prefix)), nil)
	defer iter.Release()

	for iter.Next() {
		keys = append(keys, string(iter.Key()))
	}

	return keys, iter.Error()
}

// Close 关闭数据库
func (l *LevelDB) Close() error {
	return l.db.Close()
}
