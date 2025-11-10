package handler

//type CacheHandler struct {
//	cache *cache.ScheduleCache
//}
//
//func NewCacheHandler(c *cache.ScheduleCache) *CacheHandler {
//	return &CacheHandler{cache: c}
//}
//
//func (h *CacheHandler) GetCache(c *gin.Context) {
//	schedules := h.cache.Get()
//	c.JSON(http.StatusOK, gin.H{
//		"count":     len(schedules),
//		"schedules": schedules,
//	})
//}
