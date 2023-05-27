package twoqueues

import (
	"container/list"
)

type Cache struct {
	maxBytes      int64
	fifoBytes 	  int64
	lruBytes     int64
	fifoRatio    float32
	lruRatio     float32
	fifoList    *list.List
	fifoCache map[string]*list.Element
	lruList            *list.List
	lruCache         map[string]*list.Element
	OnEvicted     func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

// New is the constructor of lruCache
func New(maxBytes int64, fifoRatio float32, lruRatio float32, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:      maxBytes,
		fifoRatio: 	0.35,
		lruRatio: 		0.65,
		fifoList:    list.New(),
		fifoCache: make(map[string]*list.Element),
		lruList:            list.New(),
		lruCache:         make(map[string]*list.Element),
		OnEvicted:     onEvicted,
	}
}

// Get look up a Key's value
func (c *Cache) Get(key string) (value Value, ok bool) {
	// if it exists in the lruCache, then move it to front
	if ele, ok := c.lruCache[key]; ok {
		c.lruList.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	// if it exists in the fifoCache
	// then move it from fifoCache to lruCache
	// and remove it from fifoCache
	if ele, ok := c.fifoCache[key]; ok {
		// delete it from fifoList (FIFO Queue)
		// and from fifo_lruCache
		c.fifoList.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.fifoCache, kv.key)
		// move to Front of lruList (LRU Queue)
		// and add to lruCache
		c.lruList.MoveToFront(ele)
		c.lruCache[key] = ele
		kv = ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// Add: adds a value to the lruCache
func (c *Cache) Add(key string, value Value) {
	// if it exists in lruCache
	// then move it to front and change it to newValue
	if ele, ok := c.lruCache[key]; ok {
		c.lruList.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.lruBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
		return
	}
	// if it exists in fifo_lruCache
	// then move it from fifoCache to lruCache, and remove it from fifoCache
	// finally change it's value to newValue
	if ele, ok := c.fifoCache[key]; ok {
		// delete it from fifoList (FIFO Queue)
		// and from fifo_lruCache
		c.fifoList.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.fifoCache, kv.key)
		c.fifoBytes -= int64(len(key)) + int64(value.Len())
		// move to Front of lruList (LRU Queue)
		// and add to lruCache
		c.lruList.MoveToFront(ele)
		c.lruCache[key] = ele
		kv = ele.Value.(*entry)
		c.lruBytes += int64(len(key)) + int64(value.Len())
		kv.value = value
		return
	}
	// otherwise
	// add it to fifoList and fifo_lruCache
	ele := c.fifoList.PushBack(&entry{key, value})
	c.fifoCache[key] = ele
	c.fifoBytes -= int64(len(key)) + int64(value.Len())

	// compare lruBytes with maxLRUBytes
	for c.lruBytes > int64(c.lruRatio*float32(c.maxBytes)) {
		c.RemoveLRUOldest()
	}

	// compare fifoBytes with maxFIFOBytes
	for c.fifoBytes > int64(c.fifoRatio*float32(c.maxBytes)) {
		c.RemoveFIFOOldest()
	}
}

// Remove LRU oldest element
func (c *Cache) RemoveLRUOldest() {
	ele := c.lruList.Back()
	if ele != nil {
		c.lruList.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.lruCache, kv.key)
		c.lruBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Remove FIFO oldest element (first in first out)
func (c *Cache) RemoveFIFOOldest() {
	ele := c.fifoList.Front()
	if ele != nil {
		c.fifoList.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.fifoCache, kv.key)
		c.fifoBytes -= int64(len(kv.key)) + int64(kv.value.Len())
	}
}
