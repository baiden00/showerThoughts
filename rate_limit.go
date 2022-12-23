package main

// To be used to rate limit in upcoming pub/sub project. Just rough ideas, not meant to be excellent code.
//Idea here is to store an object and remove it when the time has passed.
import (
	"fmt"
	"sync"
	"time"
)

// Item to be inserted into the Cache
type Item struct {
	value interface{}
	name  string
	ttl   time.Time
}

// Cache Storage for items
type Cache struct {
	storage       map[string]*Item
	lock          sync.Mutex
	done          chan bool
	deletedStream []int
}

//Continuously scan through map to evict expired entries. If there's nothing in the map send a signal to channel to trigger the main thread to exit.
//In the actual code, the main thread won't exit since a server will be running. Or will it ? IDK.
func (c *Cache) evict() {
	for {
		if len(c.storage) == 0 {
			c.done <- true
			//close(c.done)
			//break
		}
		for key, val := range c.storage {
			secs := time.Since(val.ttl).Seconds()
			if secs >= 5 {
				c.lock.Lock()
				fmt.Printf("Deleting:  %v \n", key)
				c.deletedStream = append(c.deletedStream, time.Now().Second())
				delete(c.storage, key)
				c.lock.Unlock()
			}
		}
	}
}

func add(a int, b int) int {
	return a + b
}

func main() {
	fmt.Println("Hello World")

	var items []*Item

	start := time.Now()
	items = append(items, &Item{value: 3, name: "three", ttl: time.Now().Add(4 * time.Second)})
	items = append(items, &Item{value: 4, name: "four", ttl: time.Now().Add(3 * time.Second)})
	items = append(items, &Item{value: 5, name: "five", ttl: time.Now().Add(3 * time.Second)})
	items = append(items, &Item{value: 6, name: "six", ttl: time.Now().Add(2 * time.Second)})

	testCache := &Cache{
		storage:       make(map[string]*Item),
		lock:          sync.Mutex{},
		done:          make(chan bool, 1),
		deletedStream: []int{},
	}

	for _, item := range items {
		testCache.storage[item.name] = item
	}

	fmt.Println(len(testCache.storage))
	fmt.Println(testCache.storage)

	go testCache.evict()

	// Waiting for signal to end
	<-testCache.done

	//timer
	fmt.Println(time.Since(start))

	for _, val := range testCache.deletedStream {
		fmt.Println(val)
	}
}
