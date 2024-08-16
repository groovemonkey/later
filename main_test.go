package test

// func test() {
// 	fmt.Println("Some sample future timestamps, in seconds and in Unix datetime format:")
// 	one := generateFutureTimeSeconds(time.Hour * 24 * 14)
// 	fmt.Println(one, time.Unix(one, 0))

// 	two := generateFutureTimeSeconds(time.Hour * 24 * 14)
// 	fmt.Println(two, time.Unix(two, 0))

// 	three := generateFutureTimeSeconds(time.Hour * 24 * 14)
// 	fmt.Println(three, time.Unix(three, 0))

// 	fmt.Println("\nFirst 2 should be the same:")
// 	fmt.Println(createTaskHash("this is a random hash string"))
// 	fmt.Println(createTaskHash("this is a random hash string"))
// 	fmt.Println(createTaskHash("this is a different hash string"))

// }

func main() {
	iterations := 1028
	var message string

	testUser := User{
		name:  "davdav",
		email: "foo@bar.com",
	}

	for i := 0; i < iterations; i++ {
		message = "just a string"
		main.CreateUserTask(client, &testUser, message)
	}
}
