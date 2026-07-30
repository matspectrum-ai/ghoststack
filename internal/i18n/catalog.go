package i18n

import "sync"

type Catalog struct {
	mu       sync.RWMutex
	messages map[string]map[string]string
}

func NewCatalog() *Catalog {
	return &Catalog{messages: make(map[string]map[string]string)}
}

func (c *Catalog) AddLocale(locale string, messages map[string]string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.messages[locale] = messages
}

func (c *Catalog) Translate(locale, key string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	msgs, ok := c.messages[locale]
	if !ok {
		return key
	}
	if msg, ok := msgs[key]; ok {
		return msg
	}
	return key
}
