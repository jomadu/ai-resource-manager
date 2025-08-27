// Package cache provides caching interfaces and implementations for ARM.
// 
// The cache system is composed of:
// - Cache: Basic key-value storage with TTL (metadata, version lists)
// - RepositoryCache: Git repository management (clone, fetch, extract)
// - RulesetCache: Registry-agnostic content storage (extracted files)
// - CacheKeyGenerator: SHA256-based cache key generation
package cache
