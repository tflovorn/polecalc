package polecalc

type VectorCache map[float64]map[float64]interface{}

func NewVectorCache() *VectorCache {
	c := make(VectorCache)
	return &c
}

func (c *VectorCache) Contains(k Vector2) bool {
	cacheX, ok := (*c)[k.X]
	if !ok {
		return false
	}
	_, ok = cacheX[k.Y]
	return ok
}

func (c *VectorCache) Get(k Vector2) (interface{}, bool) {
	cacheX, ok := (*c)[k.X]
	if !ok {
		// considering returning cacheX here instead of nil
		// Q: what would cacheX be in this case?
		return nil, false
	}
	cachedObj, ok := cacheX[k.Y]
	return cachedObj, ok
}

func (c *VectorCache) Set(k Vector2, obj interface{}) {
	xCache, ok := (*c)[k.X]
	if !ok {
		// need to create cache for this x
		newCache := make(map[float64]interface{})
		newCache[k.Y] = obj
		(*c)[k.X] = newCache
		return
	}
	xCache[k.Y] = obj
}
