package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
)

var (
	config       *Config     = &Config{}
	tasks        []*Task     = make([]*Task, 0)
	logFile      *os.File    = nil
	logFileMutex *sync.Mutex = &sync.Mutex{}
)

func init() {
	if err := config.ReadFile("config.yml"); err != nil {
		log.Fatal(err)
	}

	f, err := os.OpenFile(config.LogFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0777)

	if err != nil {
		log.Fatal(err)
	}

	logFile = f
}

func main() {
	for _, c := range config.Tasks {
		task := NewTask(c)

		if err := task.Initialize(); err != nil {
			log.Fatal(err)
		}

		tasks = append(tasks, task)

		log.Printf("Initialized task '%s'\n", c.Label)
	}

	for _, task := range tasks {
		task.Start()

		log.Printf("Started task '%s'\n", task.Config.Label)
	}

	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)
	<-s

	for _, task := range tasks {
		if err := task.Stop(); err != nil {
			log.Println(err)
		}

		if err := task.Close(); err != nil {
			log.Println(err)
		}

		log.Printf("Stopped task '%s'\n", task.Config.Label)
	}

	if err := logFile.Close(); err != nil {
		log.Println(err)
	}
}
