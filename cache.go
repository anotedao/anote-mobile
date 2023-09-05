package main

import (
	"time"
)

type Cache struct {
	StatsCache *Stats
}

func (c *Cache) loadStatsCache() {
	c.StatsCache = nil
	c.StatsCache = getStats()
}

func (c *Cache) start() {
	for {
		c.loadStatsCache()

		time.Sleep(time.Second * 60)
	}
}

func initCache() *Cache {
	c := &Cache{}
	c.StatsCache = &Stats{}
	go c.start()

	return c
}
