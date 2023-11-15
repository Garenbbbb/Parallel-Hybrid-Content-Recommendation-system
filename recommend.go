package main

import (
	"Work-Stealing-Based-Parallel-Hybrid-Content-Recommendation-system/data"
	"Work-Stealing-Based-Parallel-Hybrid-Content-Recommendation-system/deque"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"sort"
	"strconv"
	"sync"
)

type User struct {
	ID          string
	Preferences map[string]float64
}

// type Recommendation struct {
// 	UserID      string
// 	Recommended []string
// }

type Worker[T any] struct {
	ID    int
	Deque deque.Deque[T]
}

// WorkStealingScheduler represents a work-stealing scheduler with multiple workers
type WorkStealingScheduler[T any] struct {
	Workers []Worker[T]
}

type Recommendation struct {
	ID              string
	SimilarityScore float64
}

func newWorkStealingCartScheduler[T any](workers int) *WorkStealingScheduler[T] {
	workerPool := WorkStealingScheduler[T]{}
	for i := 0; i < workers; i++ {
		workerPool.Workers = append(workerPool.Workers, Worker[T]{ID: i})
	}
	return &workerPool
}

// Create a new worker with a unique ID
func newWorker[T any](id int) Worker[T] {
	return Worker[T]{ID: id}
}

func contains(items *[]data.Content, itemID string) bool {
	for _, item := range *items {
		if item.ID == itemID {
			return true
		}
	}
	return false
}

func CartFeatures(cart *data.ShopingCart) map[string]float64 {
	cartFeatures := make(map[string]float64)
	featureCounts := make(map[string]int)

	for _, item := range cart.Items {
		for key, value := range item.Features {
			cartFeatures[key] += value
			featureCounts[key] += 1
		}
	}
	for key := range cartFeatures {
		cartFeatures[key] /= float64(featureCounts[key])
	}
	return cartFeatures
}

// Calculate and return the cosine similarity
func Similarity(features1 map[string]float64, features2 map[string]float64) float64 {
	dotProduct := 0.0
	magnitude1 := 0.0
	magnitude2 := 0.0

	for key := range features1 {
		dotProduct += features1[key] * features2[key]
		magnitude1 += math.Pow(features1[key], 2)
		magnitude2 += math.Pow(features2[key], 2)
	}
	magnitude1 = math.Sqrt(magnitude1)
	magnitude2 = math.Sqrt(magnitude2)
	if magnitude1 == 0 || magnitude2 == 0 {
		return 0.0
	}
	return roundToDecimal(dotProduct/(magnitude1*magnitude2), 3)
}

func FindTopSimilarItems(user data.ShopingCart, content *[]data.Content, topN int) []Recommendation {
	var similarItems []Recommendation

	user_feature := CartFeatures(&user)
	for _, item := range *content {
		// Check if the item is not in the user's cart
		if !contains(&user.Items, item.ID) {
			// Compute similarity
			similarity := Similarity(user_feature, item.Features)
			similarItems = append(similarItems, Recommendation{ID: item.ID, SimilarityScore: similarity})
		}
	}
	// Sort
	sort.Slice(similarItems, func(i, j int) bool {
		return similarItems[i].SimilarityScore > similarItems[j].SimilarityScore
	})
	// Return the top N items
	return similarItems[:topN]
}

func roundToDecimal(value float64, decimalPlaces int) float64 {
	precision := math.Pow(10, float64(decimalPlaces))
	return math.Round(value*precision) / precision
}

// Content-based filtering
func processTask[T any](task T, workerId int, content *[]data.Content) string {

	t := reflect.TypeOf(task)
	switch t {
	case reflect.TypeOf(deque.TaskCart{}):
		task, _ := any(task).(deque.TaskCart)
		topSimilarItems := FindTopSimilarItems(task.Info, content, 3)
		jsonData, _ := json.Marshal(topSimilarItems)
		return "Worker " + strconv.Itoa(workerId) + " Processed Task" + strconv.Itoa(task.Task.ID) + " Result: " + string(jsonData) + "recommanded to " + task.Info.ID
	default:
		// Handle other types or provide a default case
		fmt.Println("Unsupported task type")
	}
	return "2"
}

func (ws *WorkStealingScheduler[T]) Run(result chan string, tasks *[]T, content *[]data.Content) {
	// Distribute tasks to workers
	for i, task := range *tasks {
		worker := &ws.Workers[i%len(ws.Workers)]
		worker.Deque.PushBack(task)
	}
	var wg sync.WaitGroup
	for i := range ws.Workers {
		wg.Add(1)
		go func(worker *Worker[T]) {
			defer wg.Done()
			workerProcessTasks(worker, ws, result, content)
		}(&ws.Workers[i])
	}
	// Wait for all workers to finish
	wg.Wait()
}

// workerProcessTasks simulates a worker processing tasks
func workerProcessTasks[T any](worker *Worker[T], workerPool *WorkStealingScheduler[T], result chan string, content *[]data.Content) {
	for {
		task, ok := worker.Deque.PopFront()
		if !ok {
			// No tasks in own deque, try stealing from other workers
			otherWorkerIndex := rand.Intn(len(workerPool.Workers))
			otherWorker := &workerPool.Workers[otherWorkerIndex]
			if otherWorker.ID != worker.ID {
				stolenTask, ok := otherWorker.Deque.PopBack()
				if ok {
					worker.Deque.PushFront(stolenTask)
				}
			}
		} else {
			// Process the task
			result <- processTask(task, worker.ID, content)
		}
	}
}

func main() {

	contents := data.CreateRandomContent(1000, 10, 0.5)
	taskPool := data.CreateRandomTasks(10, 2, 5, 10, 0.5)

	result := make(chan string)
	workers := 4
	workerPool := newWorkStealingCartScheduler[deque.TaskCart](workers)
	taskCount := len(taskPool)
	tasks := make([]deque.TaskCart, taskCount)
	for i := 0; i < taskCount; i++ {
		tasks[i] = deque.TaskCart{Task: deque.Task{ID: i + 1}, Info: taskPool[i]}
	}
	go workerPool.Run(result, &tasks, &contents)

	for i := 0; i < taskCount; i++ {
		fmt.Println(<-result)
	}

	close(result)

	// Item-Based Collaborative Filtering
	// users := []User{
	// 	{ID: "User1", Preferences: map[string]float64{"Item1": 5, "Item2": 0, "Item3": 2, "Item4": 5, "Item5": 5, "Item6": 5, "Item7": 2, "Item8": 5}},
	// 	{ID: "User2", Preferences: map[string]float64{"Item1": 3, "Item2": 2, "Item3": 5, "Item4": 1, "Item5": 5, "Item6": 3, "Item7": 1, "Item8": 2}},
	// 	{ID: "User3", Preferences: map[string]float64{"Item1": 2, "Item2": 5, "Item3": 3, "Item4": 4, "Item5": 2, "Item6": 5, "Item7": 6, "Item8": 1}},
	// 	{ID: "User4", Preferences: map[string]float64{"Item1": 6, "Item2": 2, "Item3": 6, "Item4": 5, "Item5": 2, "Item6": 2, "Item7": 1, "Item8": 5}},
	// 	{ID: "User5", Preferences: map[string]float64{"Item1": 1, "Item2": 5, "Item3": 1, "Item4": 2, "Item5": 5, "Item6": 5, "Item7": 2, "Item8": 4}},
	// 	{ID: "User6", Preferences: map[string]float64{"Item1": 2, "Item2": 6, "Item3": 2, "Item4": 5, "Item5": 4, "Item6": 6, "Item7": 2, "Item8": 5}},
	// 	{ID: "User7", Preferences: map[string]float64{"Item1": 0, "Item2": 8, "Item3": 5, "Item4": 0, "Item5": 5, "Item6": 0, "Item7": 5, "Item8": 0}},
	// 	{ID: "User8", Preferences: map[string]float64{"Item1": 3, "Item2": 6, "Item3": 1, "Item4": 5, "Item5": 6, "Item6": 0, "Item7": 2, "Item8": 5}},
	// 	{ID: "User9", Preferences: map[string]float64{"Item1": 5, "Item2": 3, "Item3": 2, "Item4": 5, "Item5": 1, "Item6": 5, "Item7": 4, "Item8": 5}},
	// }

	// itemCount := 8
	// //compute itemI -> ItemJ cosine similarity score based on who rated both of them if haveing n item, there will be (n)*(n-1)/2
	// workerPoolItem := WorkStealingScheduler{}
	// for i := 0; i < workers; i++ {
	// 	workerPoolItem.Workers = append(workerPool.Workers, newWorker(i))
	// }
	// taskCountItem := itemCount * (itemCount - 1) / 2

	// tasksItem := make([]deque.TaskItem, taskCountItem)
	// index := 0
	// for i := 0; i < itemCount-1; i++ {
	// 	for j := i + 1; j < itemCount; j++ {
	// 		tasksItem[index] = deque.TaskItem{ID: i + 1, Task: [2]int{i, j}}
	// 		index += 1
	// 	}
	// }
	// cal_done := sync.WaitGroup{}
	// cal_done.Add(index)

	// similarity_martix := make(chan [3]float32)

	// go workerPoolItem.Run2(similarity_martix, &tasksItem, &users)

}

// func (ws *WorkStealingScheduler) Run2(result chan [3]float32, tasks *[]deque.TaskItem, content *[]User) {
// 	// Distribute tasks to workers
// 	for i, task := range *tasks {
// 		worker := &ws.Workers[i%len(ws.Workers)]
// 		worker.Deque.PushBack(task)
// 	}
// 	var wg sync.WaitGroup
// 	for i := range ws.Workers {
// 		wg.Add(1)
// 		go func(worker *Worker) {
// 			defer wg.Done()
// 			workerProcessTasks(worker, ws, result, content)
// 		}(&ws.Workers[i])
// 	}
// 	// Wait for all workers to finish
// 	wg.Wait()
// }
