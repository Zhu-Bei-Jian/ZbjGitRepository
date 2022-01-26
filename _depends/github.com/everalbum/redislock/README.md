redislock
=========

Pessimistic locking for Go using Redis.

### Installation

    go get -u github.com/everalbum/redislock

### Documentation

http://godoc.org/github.com/everalbum/redislock

### Example

```go
lock, ok, err := redislock.TryLock(conn, "user:123")
if err != nil {
	log.Fatal("Error while attempting lock")
}
if !ok {
	// User is in use - return to avoid duplicate work, race conditions, etc.
	return
}
defer lock.Unlock()

// Do something with the user.
```
