package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"time"
)

var ErrFileWrite = "file write error: "

var actions = []string{
	"logged in",
	"logged out",
	"created record",
	"deleted record",
	"updated account",
}

type logItem struct {
	action    string
	timestamp time.Time
}

type User struct {
	id    int
	email string
	logs  []logItem
}

func (u User) GetActivityInfo() string {
	output := fmt.Sprintf("UID: %d; Email: %s;\nActivity log:\n", u.id, u.email)
	for index, item := range u.logs {
		output += fmt.Sprintf("%d. [%s] at %s\n", index, item.action, item.timestamp.Format(time.RFC3339))
	}

	return output
}

func main() {
	rand.Seed(time.Now().Unix())

	startTime := time.Now()

	wg := &sync.WaitGroup{}

	users := make(chan User)

	for i := 0; i < runtime.NumCPU()-1; i++ {
		wg.Add(1)
		go func(num int, wg *sync.WaitGroup) {
			defer wg.Done()
			fmt.Printf("Starting worker #%d\n", num)
			defer fmt.Printf("Stoping worker #%d\n", num)
			for u := range users {
				if err := saveUserInfo(u); err != nil {
					fmt.Printf("error saving user info %d\n", u.id)
				}
			}
		}(i, wg)
	}

	for i := 0; i < 1000; i++ {
		users <- generateUser(i)
	}

	close(users)

	wg.Wait()

	fmt.Printf("DONE: Time Elapsed: %.2f seconds\n", time.Since(startTime).Seconds())
}

func saveUserInfo(user User) error {
	fmt.Printf("WRITTINF FILE FOR UID %d\n", user.id)

	fileName := fmt.Sprintf("users/uid%d.txt\n", user.id)
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	if _, err := file.WriteString(user.GetActivityInfo()); err != nil {
		log.Fatal(ErrFileWrite, err)
	}
	time.Sleep(time.Second)

	return nil
}

func generateUser(id int) User {
	return User{
		id:    id + 1,
		email: fmt.Sprintf("user%dcompany.com", id+1),
		logs:  generateLogs(100),
	}
}

func generateLogs(count int) []logItem {
	logs := make([]logItem, count)

	for i := 0; i < count; i++ {
		logs[i] = logItem{
			action:    actions[rand.Intn(len(actions)-1)],
			timestamp: time.Now(),
		}
	}

	return logs
}
