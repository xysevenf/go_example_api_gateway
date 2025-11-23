package rate_limiter

import (
	"encoding/binary"

	"github.com/coocood/freecache"
)

const limitsCacheSize = 4 * 1024 * 1024

type freecacheLimitsStorage struct {
	cache *freecache.Cache
}

func (f *freecacheLimitsStorage) GetLimitCounter(key string) (int, bool) {
	v, error := f.cache.Get([]byte(key))
	if error == nil {
		return int(BytesToInt(v)), true
	}
	return 0, false
}

func (f *freecacheLimitsStorage) SetLimitCounter(key string, val int, ttl int) bool {
	intVal := IntToBytes(int64(val))
	error := f.cache.Set([]byte(key), intVal, ttl)
	return error == nil
}

func NewFreecacheLimitsStorage() *freecacheLimitsStorage {
	return &freecacheLimitsStorage{cache: freecache.NewCache(limitsCacheSize)}
}

func IntToBytes(val int64) []byte {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, uint64(val))
	return bytes
}

func BytesToInt(bytes []byte) int64 {
	return int64(binary.LittleEndian.Uint64(bytes))
}
