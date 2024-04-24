package consistency

import (
	"errors"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

const MAX_VIRTUAL_NODE = 5

type HashFunc func(data []byte) uint32

// 哈希环 用于负载均衡
type HashRing struct {
	hashFunc HashFunc
	keys     []uint32 // sorted
	hashMap  map[uint32]uint32
	sync.Mutex
}

func NewRing(hf HashFunc) *HashRing {
	h := &HashRing{
		hashFunc: hf,
		keys:     make([]uint32, 0),
		hashMap:  make(map[uint32]uint32),
	}

	if hf == nil {
		h.hashFunc = crc32.ChecksumIEEE
	}

	return h
}

func (r *HashRing) Add(node string) {
	r.Lock()
	defer r.Unlock()
	for i := 0; i < MAX_VIRTUAL_NODE; i++ {
		str := node + strconv.Itoa(i)
		hash := r.hashFunc([]byte(str))
		r.hashMap[hash] = r.hashFunc([]byte(node))
		r.keys = append(r.keys, hash)
	}

	sort.Slice(r.keys, func(i, j int) bool {
		return r.keys[i] < r.keys[j]
	})
}

func (r *HashRing) Remove(node string) {
	r.Lock()
	defer r.Unlock()
	length := len(r.keys)
	if length == 0 {
		return
	}

	for i := 0; i < MAX_VIRTUAL_NODE; i++ {
		str := node + strconv.Itoa(i)
		hash := r.hashFunc([]byte(str))
		delete(r.hashMap, hash)

		idx, ret := r.serach(hash)
		if ret {
			r.keys = append(r.keys[0:idx], r.keys[idx+1:]...)
		}
	}
}

func (r *HashRing) serach(hash uint32) (int, bool) {
	length := len(r.keys)
	idx := sort.Search(length, func(i int) bool {
		return r.keys[i] == hash
	})
	if idx >= length {
		return 0, false
	}

	return idx, true
}

func (r *HashRing) GetNode(node string) (uint32, error) {
	r.Lock()
	defer r.Unlock()
	length := len(r.keys)
	if length == 0 {
		return 0, errors.New("Not Found")
	}

	hash := r.hashFunc([]byte(node))
	idx := sort.Search(length, func(i int) bool {
		return r.keys[i] >= hash
	})

	if idx == length {
		idx = idx % length
	}

	return r.hashMap[r.keys[idx]], nil
}
