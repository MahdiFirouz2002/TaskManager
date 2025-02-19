package concurrentlogger

import (
	"log"
	"nikandishan/structs/task"
	"sync"
)

var (
	UpdateTaskChan = make(chan task.Task, 100)
)

func StartLogger(wg *sync.WaitGroup) {
	const workerCount = 3

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go taskWorker(wg, i)
	}
}

func taskWorker(wg *sync.WaitGroup, workerID int) {
	defer wg.Done()
	for updateItem := range UpdateTaskChan {
		log.Printf("[Worker %d] Task updated: ID=%s, Title=%s", workerID, updateItem.ID, updateItem.Title)
	}
}
