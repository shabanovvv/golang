package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	item, ok := c.items[key]
	if ok {
		item.Value = value
		c.queue.MoveToFront(item)
	} else {
		item = c.queue.PushFront(value)
		c.items[key] = item
	}

	if c.capacity < len(c.items) {
		backItem := c.queue.Back()
		if backItem != nil {
			c.queue.Remove(backItem)
			delete(c.items, key)
		}
	}

	return ok
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	if item, ok := c.items[key]; ok {
		c.queue.MoveToFront(item)

		return item.Value, true
	}

	return nil, false
}

func (c *lruCache) Clear() {
	c.queue = NewList()
	c.items = make(map[Key]*ListItem)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
