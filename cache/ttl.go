package cache

import (
	"time"

	"github.com/jellydator/ttlcache/v3"
	"github.com/nawaltni/tracker/domain"
)

var _ domain.UserCache = (*UserCache)(nil)

// UserCache is a cache for users
type UserCache struct {
	c   *ttlcache.Cache[int, domain.User]
	ttl time.Duration
}

// NewUserCache creates a new UserCache
func NewUserCache() *UserCache {
	cache := ttlcache.New[int, domain.User]()
	return &UserCache{c: cache, ttl: 60 * time.Minute}
}

// Set sets a user in the cache
func (c *UserCache) Set(id int, user domain.User) {
	c.c.Set(id, user, c.ttl)
}

// Get gets a user from the cache
func (c *UserCache) Get(id int) *domain.User {
	item := c.c.Get(id)
	if item == nil {
		return nil
	}

	user := item.Value()

	return &user
}
