package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	mu       sync.Mutex // Защищает queue и items
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key Key
	val interface{}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	if c.capacity == 0 {
		// Записать ошибку в log, вернуть ошибку, не буду портить сигнатуру функции и импорт лишний делать
		return false
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	cItem, exists := c.items[key]
	if exists {
		cItemValue, ok := cItem.Value.(*cacheItem)
		if ok {
			cItemValue.val = value
			c.queue.MoveToFront(cItem)
			return true
		}
		// Записать ошибку в log, вернуть ошибку, не буду портить сигнатуру функции и импорт лишний делать
		return true
	}

	if c.queue.Len() == c.capacity {
		lastCItem := c.queue.Back()
		if lastCItem != nil {
			lastCItemValue, ok := lastCItem.Value.(*cacheItem)
			if ok {
				delete(c.items, lastCItemValue.key)
				c.queue.Remove(lastCItem)
			} else {
				// Записать ошибку в log, вернуть ошибку, не буду портить сигнатуру функции и импорт лишний делать
				return false
			}
		}
	}

	newCItem := c.queue.PushFront(&cacheItem{key, value})
	c.items[key] = newCItem

	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	cItem, exists := c.items[key]
	if !exists {
		return nil, false
	}

	c.queue.MoveToFront(cItem)

	cItemValue, ok := cItem.Value.(*cacheItem)
	if ok {
		return cItemValue.val, true
	}
	// Записать ошибку в log, вернуть ошибку, не буду портить сигнатуру функции и импорт лишний делать
	return nil, true
}

func (c *lruCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[Key]*ListItem)

	c.queue = nil
	c.queue = NewList()
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
