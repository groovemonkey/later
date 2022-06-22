package main

import (
	// "context"
	"context"
	"crypto/sha1"
	"fmt"
	"math/rand"
	"time"

	"github.com/go-redis/redis/v9"
)

const (
	// Two weeks
	schedulingTimerange = time.Hour * 24 * 14
	workerLoopBatchSize = 50
)

var ctx = context.Background()

type user struct {
	name  string
	email string
}

type task struct {
	hash          string
	scheduledTime string
	username      string
	email         string
	message       string
}

func main() {
	fmt.Println("Testing Golang Redis")

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	dave := user{
		name:  "dave",
		email: "dave@example.org",
	}

	fmt.Println("\nWriting to redis with createUserTask()")
	taskHash, err := createUserTask(client, &dave, "This is a task message! Woohoo! Test all kinds of symbols and stuff in here.")
	if err != nil {
		panic("Error while trying to createUserTask()")
	}
	fmt.Println("Taskhash is ", taskHash)

	fmt.Println("\nRetrieving task...")
	task, err := getTaskDetails(client, taskHash)
	if err != nil {
		panic("Error while retrieving task with getTaskDetails()")
	}
	fmt.Println("Task is:", task)

	fmt.Println("\nSending task email:")
	sendTaskEmailTEST(&task)

	fmt.Println("Running worker loop!")
	err = runWorkerLoop(client)
	if err != nil {
		fmt.Println("Workerloop exited with an error.")
	}
}

// Return some random time (in seconds) in the future, limited by $timeRange
// Based on https://stackoverflow.com/questions/43495745/how-to-generate-random-date-in-go-lang
func generateFutureTimeSeconds(timeRange time.Duration) int64 {
	// add a bit of time to the minimum, to prevent immediate notifications
	min := time.Now().Add(time.Hour * 24)
	max := min.Add(timeRange)
	delta := max.Unix() - min.Unix()

	// We return time in seconds
	seconds := rand.Int63n(delta) + min.Unix()
	return seconds
}

// Create a user task, and store it in the appropriate places in Redis. Return the task's hash and an optional error.
func createUserTask(rdb *redis.Client, user *user, message string) (string, error) {
	// Pick a time
	scheduledTimeSecs := generateFutureTimeSeconds(schedulingTimerange)

	// create a unique hash for this task, based on scheduled time, username, and task message
	stringToHash := fmt.Sprintf("%d-%s-%s", scheduledTimeSecs, user.name, message)
	taskHash := hashString(stringToHash)

	// `taskdetails`: Hash datatype to hold task info, namespaced under the taskHash
	// Answers: What is the task with hash X, and what would I need to know to run it?
	// HMSET $HASH username foo tasktype email msg hello
	rdb.HSet(ctx, taskHash,
		"hash", taskHash,
		"scheduledTimeSecs", scheduledTimeSecs,
		"username", user.name,
		"email", user.email,
		"message", message,
	)

	// `tasks`: Sorted set to hold task/timing info, key=timestamp
	// Answers: What are the next scheduled tasks we need to run, and when??
	// ZADD tasks $TIMESTAMP $HASH
	rdb.ZAdd(ctx, "tasks",
		redis.Z{
			Score:  float64(scheduledTimeSecs),
			Member: taskHash,
		})

	// `username_tasks` list, namespaced under the username
	// Answers: What are User X's pending tasks?
	// LPUSH ${USERNAME}_tasks $HASH
	rdb.LPush(ctx, fmt.Sprintf("%s_tasks", user.name), taskHash)

	return taskHash, nil
}

// Create a SHA1 hash of the string passed in
func hashString(hashString string) string {
	hasher := sha1.New()
	hasher.Write([]byte(hashString))
	bs := hasher.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

// Use a taskHash to look up a task in Redis, and create/return a task struct based on that data.
func getTaskDetails(rdb *redis.Client, taskHash string) (task, error) {
	result, err := rdb.HGetAll(ctx, taskHash).Result()
	if err != nil {
		return task{}, err
	}

	return task{
		hash:          result["key"],
		email:         result["email"],
		scheduledTime: result["scheduledTimeSecs"],
		username:      result["username"],
		message:       result["message"],
	}, nil
}

func sendTaskEmailTEST(task *task) {
	taskMessage := fmt.Sprintf("Processing task %s for user %s: sending email to %s with message: %s", task.hash, task.username, task.email, task.message)
	fmt.Println(taskMessage)
}

// A simple "worker" loop that grabs a task and fires off a worker goroutine
func runWorkerLoop(rdb *redis.Client) error {
	// Eternal Loop
	for {
		// Get the next N task hashes
		for _, taskHash := range workerGrabTaskHashBatch(rdb) {
			fmt.Println("Workerloop got taskHash: ", taskHash)

			// Fire off a goroutine to handle it
			go handleTask(rdb, taskHash)
			// sendTaskEmailTEST(&task)
		}

		// TODO sleep until next task hash is due.
		// Careful, a race condition could develop here if handleTask goroutines don't clean up in time
		// Solution, if this actually becomes a problem: make an in-progress queue
		fmt.Println("runWorkerLoop is sleeping...")
		time.Sleep(5 * time.Second)
	}
}

// Returns the next batch of taskHashes
// Say the name of this function 3 times in a dark room to invoke a portal to the off-by-oneth dimension.
func workerGrabTaskHashBatch(rdb *redis.Client) []string {
	var taskHashes []string

	rangeByOpts := redis.ZRangeBy{
		Min:    "-inf",
		Max:    "+inf",
		Offset: 0,
		Count:  workerLoopBatchSize,
	}
	result, err := rdb.ZRangeByScoreWithScores(ctx, "tasks", &rangeByOpts).Result()
	if err != nil {
		// TODO more fine-grained error handling
		return taskHashes
	}

	for _, zItem := range result {
		// zItems have a Score and a Member (taskHash)
		taskHashes = append(taskHashes, zItem.Member.(string))
		fmt.Println("Processing Zitem: ", zItem)
	}
	return taskHashes
}

// handleTask is designed to be fired off as a goroutine to handle one task
// It currently sends a simple (fmt.Println) email to prove it works
func handleTask(rdb *redis.Client, taskHash string) {
	// Make a task from the hash
	task, err := getTaskDetails(rdb, taskHash)
	if err != nil {
		fmt.Println("Error in handleTask: ", err)
	}
	sendTaskEmailTEST(&task)

	// Clean up task
}

func deleteTask(task *task) error {

}
