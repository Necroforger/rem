package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

type Name struct {
	Date        int64
	Description string
}

var (
	newFlag   = flag.Bool("n", false, "create new reminder")
	clearFlag = flag.Bool("c", false, "clear reminders")
	listFlag  = flag.Bool("l", false, "list all reminders")
	helpFlag  = flag.Bool("h", false, "help")
	pathFlag  = flag.String("p", "", "path of reminder json file")
)

var (
	path = os.Getenv("GOPATH")
)

func newReminder() {
	var text string

	// Attempt to open the reminders file in write only mode. If it does not exist, create it.
	file, err := os.OpenFile(path+"/reminders.json", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fmt.Print("Description: ")

	scanner := bufio.NewScanner(os.Stdin)
	if ok := scanner.Scan(); ok {
		text = scanner.Text()
	}

	err = json.NewEncoder(file).Encode(Name{time.Now().Unix(), text})
	if err != nil {
		log.Fatal(err)
	}
}

func removeReminders() {
	err := os.Remove(path + "/reminders.json")
	if err != nil {
		log.Fatal(err)
	}
}

func listReminders() error {
	file, err := os.Open(path + "/reminders.json")
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(file)

	for {
		var content Name
		switch decoder.Decode(&content) {
		case nil:
			t := time.Unix(content.Date, 0)
			fmt.Printf("%s: %s\n", t.Format("2006-01-02 15:04:05"), content.Description)
		case io.EOF:
			return nil
		default:
			return err
		}
	}
}

func main() {
	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		flag.PrintDefaults()
	}

	if *pathFlag != "" {
		path = *pathFlag
	}

	if *newFlag {
		newReminder()
	}

	if *clearFlag {
		removeReminders()
	}

	if *helpFlag {
		flag.PrintDefaults()
	}

	if *listFlag {
		err := listReminders()
		if err != nil {
			log.Fatal(err)
		}
	}
}
