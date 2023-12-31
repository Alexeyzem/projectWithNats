package cashe

import (
	"context"
	"errors"
	"sync"
	"time"
	"wb_project/pkg/fullrepo"
	"wb_project/pkg/logging"
	"wb_project/pkg/user"
)

type CacheData struct {
	Data       user.User
	Expiration int64
	Created    time.Time
}

type UserCache struct {
	sync.RWMutex
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	Data              map[string]CacheData
}

//defaultExpiration - container lifetime,
//cleanupInterval - delete interval
func NewCache(defaultExpiration, cleanupInterval time.Duration) *UserCache {
	lg := logging.GetLogger()
	lg.Info("create Cache")
	data := make(map[string]CacheData)
	cache := UserCache{
		Data:              data,
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
	}

	if cleanupInterval > 0 {
		cache.startGC()
	}
	return &cache
}

func (c *UserCache) Set(key string, value user.User, duration time.Duration) {
	if duration == 0 {
		duration = c.defaultExpiration
	}
	if duration < 0 {
		duration = -1 * duration
	}
	expiration := time.Now().Add(duration).UnixNano()
	c.Lock()
	defer c.Unlock()
	c.Data[key] = CacheData{
		Data:       value,
		Expiration: expiration,
		Created:    time.Now(),
	}
}

func (c *UserCache) Get(key string) (user.User, bool) {
	c.RLock()
	defer c.RUnlock()
	data, ok := c.Data[key]
	if !ok {
		return data.Data, false
	} else {
		if time.Now().UnixNano() > data.Expiration { //просроченный кэш
			return data.Data, false
		}
	}
	return data.Data, ok
}

func (c *UserCache) startGC() {
	go c.gC()
}
func (c *UserCache) gC() {
	for {
		<-time.After(c.cleanupInterval)
		if c.Data == nil {
			return
		}
		if keys := c.expiredKeys(); len(keys) != 0 {
			c.clearData(keys)
		}
	}
}
func (c *UserCache) expiredKeys() (keys []string) {
	c.RLock()

	defer c.RUnlock()

	for k, i := range c.Data {
		if time.Now().UnixNano() > i.Expiration && i.Expiration > 0 {
			keys = append(keys, k)
		}
	}

	return nil
}
func (c *UserCache) clearData(keys []string) {
	c.Lock()
	defer c.Unlock()
	for _, k := range keys {
		delete(c.Data, k)
	}
}

func (c *UserCache) Delete(key string) error {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.Data[key]; !ok {
		return errors.New("key not founded")
	}
	delete(c.Data, key)
	return nil
}

func (c *UserCache) LoadFile(repo *fullrepo.FullRepo, duration time.Duration) error {
	users, err := repo.RepoU.FindAll(context.Background())
	if err != nil {
		return err
	}
	for i, u := range users {
		if i >= 1000 {
			break
		}
		d, err := repo.RepoD.FindOne(context.Background(), u.ID)
		if err != nil {
			return err
		}
		it, err := repo.RepoI.FindAllOfOneUser(context.Background(), u.TrackNumber)
		if err != nil {
			return err
		}
		p, err := repo.RepoP.FindOne(context.Background(), u.ID)
		if err != nil {
			return err
		}
		u.Deliv = d
		u.Items = it
		u.Payment = p
		c.Set(u.OrderUid, u, duration)
	}
	return nil
}
