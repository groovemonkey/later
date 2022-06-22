# Later: Nondeterministically clear your TODO list (and free your mind)

This is just a fun little experiment to learn Go (and play with my favorite datastore, Redis).

Get reminded of a certain thought (or other string) at a random time in the next $DURATION.


## TODO
- make it work
  - delete is broken
  - worker: resolve time value from timestamp and sleep until then

- fmt. to log.
- how does that redis response value work? (*redis.IntCmd and https://redis.uptrace.dev/guide/go-redis.html#redis-nil)
- add error handling to the redis calls
- read in schedulingTimerange from environment variables
- read in redis host/user/pass from environment variables


## Dev/Testing

### Get yourself a redis
```
docker run -d -p 6379:6379 redislabs/redismod
```

### Run this thing
```go run main.go```


## Design

This is essentially a Redis-backed cron app. Here are the Redis primitives we'll be using:

### Creating a task for a user
```
# `taskdetails`: Hash datatype to hold task info, namespaced under a random hash
HMSET $HASH username foo tasktype email msg hello

# `tasks`: Sorted set to hold task/timing info, key=timestamp
ZADD tasks $TIMESTAMP $HASH

# `username_tasks` list, namespaced under the username
LPUSH ${USERNAME}_tasks $HASH
```

### Main Worker-Pool Manager Thread
- Grab the next X tasks, based on timestamp (default=50): (`ZRANGE/ZREVRANGE`)
- check if max number of worker goroutines are already running
  - bad case: max workers already running, and next task in pool is already in the past: LOG A WARNING
- fire off worker goroutines to deal with the tasks in the "next x tasks" list
- if next event time (in redis `tasks` sortedset) is in the future, sleep until $WAKEUP_BEFORE_SECONDS before that task timestamp (configurable): (`time.Sleep(Time.Until($EVENT_TIME - $WAKEUP_BEFORE_SECONDS)))`


### Goroutine worker
- do the event action
- single Redis transaction?
  - delete (ZREM from `tasks` sorted set)
  - delete from `taskdetails` hashmap (HDEL)
  - delete from username_tasks list
- if/when there's a webapp around this, update the task status in the webapp's RDBMS (postgres)
