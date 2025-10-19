# Go Caching Masterclass 🚀

A hands-on learning journey through caching patterns in Go, using Redis and PostgreSQL.

## 📚 Projects

### Project 1: Cache-Aside Pattern ✅
**Status:** Completed
**Tech Stack:** Go, Redis
**Key Learnings:**
- Cache-aside (lazy loading) pattern
- Cache hits vs misses
- TTL (Time To Live)
- Performance improvement: **~7,400x faster** with Redis cache!

**Results:**
- Without cache: ~2000ms per request
- With cache (hit): ~270µs per request
- Cache miss: ~2001ms (first request only)

**Pattern Implementation:**
```
1. Check Redis cache
2. If HIT → return immediately
3. If MISS → query database
4. Store result in Redis (TTL: 5 min)
5. Return result
```

## 🛠️ Tech Stack
- **Go** 1.25.1
- **Redis** 7.x (in-memory cache)
- **go-redis/v9** (Redis client)

## 📖 How to Run
```bash
# Ensure Redis is running
redis-cli ping  # Should return PONG

# Run Project 1
cd project1-basic-cache
go run main.go
```

## 📊 Performance Comparison

| Scenario | Without Cache | With Cache (Hit) | Improvement |
|----------|--------------|------------------|-------------|
| Request 1 | 2000ms | 2001ms (miss) | - |
| Request 2 | 2000ms | 244µs | 8,197x |
| Request 3 | 2000ms | 266µs | 7,519x |
| Average | 2000ms | ~270µs | ~7,400x |

## 📝 Learning Progress
- [x] Project 1: Basic Cache-Aside Pattern ✅
- [ ] Project 2: Write-Through Cache
- [ ] Project 3: Cache Invalidation Strategies
- [ ] Project 4: Distributed Caching Patterns
- [ ] Project 5: Multi-Layer Cache Architecture

## 💡 Key Takeaways (Project 1)

1. **Cache-Aside Pattern**: Application manages cache explicitly
2. **Lazy Loading**: Data loaded into cache only when requested
3. **TTL Strategy**: 5-minute expiration prevents stale data
4. **Massive Speedup**: Redis in-memory cache is ~7,400x faster than simulated DB
5. **Trade-off**: First request pays database cost, subsequent requests are instant

---

**Learning Date:** $(date +%Y-%m-%d)  
**Repository:** Part of 30-day Go challenge  
**Next:** Project 2 - Write-Through Cache Pattern
