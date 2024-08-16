# Later: Nondeterministically clear your TODO list (and free your mind)

This is just a fun little experiment to learn Go (and play with my favorite datastore, Redis).

Get reminded of a certain thought (or other string) at a random time in the next $DURATION.


## TODO
- write test function to fill up DB
  - see how things work with the waitgroup -- does the function wait for workers (and limit them to approx 50?)

### Main Worker-Pool Manager Thread
- check if max number of worker goroutines are already running
  - bad case: max workers already running, and next task in pool is already in the past: LOG A WARNING
- fire off worker goroutines to deal with the tasks in the "next x tasks" list
- if next event time (in redis `tasks` sortedset) is in the future, sleep until $WAKEUP_BEFORE_SECONDS before that task timestamp (configurable): (`time.Sleep(Time.Until($EVENT_TIME - $WAKEUP_BEFORE_SECONDS)))`


### Various small stuff
- turn constants into env vars
  - read in schedulingTimerange from environment variables
  - read in redis host/user/pass from environment variables

- move testing code from main() into a separate test script
- remove verbose println logging. Change real (WARN/ERR) logs from fmt. to log.
- how does that redis response value work? (*redis.IntCmd and https://redis.uptrace.dev/guide/go-redis.html#redis-nil)
- add error handling to the redis calls

### Features
- email sending
- SMS sending
- web/API frontend, user registration, CRUD on tasks, see past tasks, etc.
- iPhone app which uses the frontend API?

### Infra
- run on AWS
- Frontend and backend go binaries designed to run separately, frontend hits backend via HTTP API
- Nomad running binaries packaged into Alpine containers (or even running as raw_binary jobtype in Nomad!)
  - one or more of each container/binary (parallel processing against redis needs to be safe for this to work)
  - redis container with persistent volume
  - managed postgres instance or maybe just a container w/ persistent volume
- Safe crash + restart behavior for data persistence (temp queues for in-flight work that get carefully flushed on startup?)


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
