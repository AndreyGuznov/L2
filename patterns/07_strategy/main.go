package main

import "fmt"

/*
Стратегия - это поведенческий паттерн проектирования, который определяет семейство схожих алгоритмов и помещает каждый из них в
собственный класс,после чего алгоритмы можно взаимозаменять прямо во время исполнения программы.

Плюсы:
	- Возможность замены алгоритмов в рантайме
	- Отделение алгоритмов от остальной логики, сокрытие самих алгоритмов

Минусы:
	- Усложнение кода, засчет введения дополнительных объектов
	- Клиент должен знать в чем состоит отличие алгоритмов, чтобы выбрать нужный
	Для примера реалзиации паттерна рассмотрим разработку In-Memory Cache.
	Так как кеш находится внутри памяти - размер его ограничен.
	Когда память заполнится (либо чуть раньше) потребуется удалить
	какие-то записи. Сделать это можно несколькими способами - алгоритмами:
	- Least Recently Used (LRU) - убрать запись, использовавшуюсь наиболее
	давно
	- FIFO - убрать запись, созданную раньше остальных
	- Least Frequently Used (LFU) - убрать запись наименее часто
	использовавшуюся
	При этом требуется отделить алгоритмы от кеша, для возможности
	замены алгоритма в рантайме и для возможности их модификации без
	изменения кода кеша.
*/

// Интерфейс стратегии
type evictionAlgo interface {
	evict(c *cache)
}

type fifo struct {
}

func (l *fifo) evict(c *cache) {
	fmt.Println("Evicting by fifo strategy")
}

type lru struct {
}

func (l *lru) evict(c *cache) {
	fmt.Println("Evicting by lru strategy")
}

type lfu struct {
}

func (l *lfu) evict(c *cache) {
	fmt.Println("Evicting by lfu strategy")
}

//  Контекст

type cache struct {
	storage      map[string]string
	evictionAlgo evictionAlgo
	capacity     int
	maxCapacity  int
}

func initCache(e evictionAlgo) *cache {
	storage := make(map[string]string)
	return &cache{
		storage:      storage,
		evictionAlgo: e,
		capacity:     0,
		maxCapacity:  2,
	}
}

func (c *cache) setEvictionAlgo(e evictionAlgo) {
	c.evictionAlgo = e
}

func (c *cache) add(key, value string) {
	if c.capacity == c.maxCapacity {
		c.evict()
	}
	c.capacity++
	c.storage[key] = value
}

func (c *cache) del(key string) {
	delete(c.storage, key)
}

func (c *cache) evict() {
	c.evictionAlgo.evict(c)
	c.capacity--
}

func main() {
	lfu := &lfu{}
	cache := initCache(lfu)

	cache.add("a", "1")
	cache.add("b", "2")

	cache.add("c", "3")

	lru := &lru{}
	cache.setEvictionAlgo(lru)

	cache.add("d", "4")

	fifo := &fifo{}
	cache.setEvictionAlgo(fifo)

	cache.add("e", "5")

}
