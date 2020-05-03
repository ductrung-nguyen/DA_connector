package base

import "da_connector/context"

type SimpleCallback func()

type BaseObject struct {
	name string
	// TODO: UserContext here
	context context.UserContext

	cache map[string]interface{}
}

// New creates a new baseObject
func New(name string, context context.UserContext) BaseObject {
	return BaseObject{name: name, context: context}
}

// ClearCache Clears the object's internal cache.
func (obj BaseObject) ClearCache() {
	obj.cache = make(map[string]interface{})
}

// GetCache Retrieves an item from the internal cache.
// Key to retrieve
// Either a callback or an explicit default value
func (obj BaseObject) GetCache(key string, defaultValue interface{}) (result interface{}) {
	if v, ok := obj.cache[key]; ok {
		switch x := v.(type) {
		case SimpleCallback:
			x()
			result = nil
			break
		default:
			result = v
		}
		return
	}
	return defaultValue
}

// GetCacheItem Retrieves a keyed item from inside a cache item.
func (obj BaseObject) GetCacheItem(key string, keyInCacheItem string, defaultCache interface{}, defaultItem interface{}) (result interface{}) {
	if cache := obj.GetCache(key, defaultCache); cache == nil {
		return defaultItem
	} else {
		switch v := cache.(type) {
		case map[string]interface{}:
			result = v[keyInCacheItem]
			break
		default:
			result = defaultItem
		}
		return
	}
}

// SetCache Sets a specific cache item, for when a cacheable value was a by-product.
func (obj BaseObject) SetCache(key string, value interface{}) {
	obj.cache[key] = value
}

// GetContext returns the user context
func (obj BaseObject) GetContext() context.UserContext {
	return obj.context
}

// GetName Protected as a derived class may want to offer the name under a different name.
func (obj BaseObject) GetName() string {
	return obj.name
}
