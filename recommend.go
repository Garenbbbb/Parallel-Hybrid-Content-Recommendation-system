package main

import (
	"Work-Stealing-Based-Parallel-Hybrid-Content-Recommendation-system/data"
	"Work-Stealing-Based-Parallel-Hybrid-Content-Recommendation-system/deque"
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"sort"
	"sync"
)

type Recommendation struct {
	UserID        string
	RecommendList []RecommendationItem
}

type Worker[T any] struct {
	ID    int
	Deque deque.Deque[T]
}

// WorkStealingScheduler represents a work-stealing scheduler with multiple workers
type WorkStealingScheduler[T any] struct {
	Workers []Worker[T]
}

type RecommendationItem struct {
	ID              string
	SimilarityScore float64
}

func newWorkStealingScheduler[T any](workers int) *WorkStealingScheduler[T] {
	workerPool := WorkStealingScheduler[T]{}
	for i := 0; i < workers; i++ {
		workerPool.Workers = append(workerPool.Workers, Worker[T]{ID: i})
	}
	return &workerPool
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

func FindTopSimilarItems(user data.ShopingCart, content *map[string]data.ItemData, topN int) []RecommendationItem {
	var similarItems []RecommendationItem

	user_feature := CartFeatures(&user)
	for item := range *content {
		// Check if the item is not in the user's cart
		if !contains(&user.Items, item) {
			// Compute similarity
			similarity := Similarity(user_feature, (*content)[item].Features)
			similarItems = append(similarItems, RecommendationItem{ID: item, SimilarityScore: similarity})
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

func FindTopSimilarItemsUnRated(task data.Content, rateList *data.RateData, topN int) []RecommendationItem {
	var unRated []string
	var rated []string
	var similarItemsRating []RecommendationItem
	for item := range task.Features {
		if task.Features[item] == 0 {
			unRated = append(unRated, item)
		} else {
			rated = append(rated, item)
		}
	}
	for _, item := range unRated {
		numerator := 0.0
		denomertor := 0.0
		for _, rated := range rated {
			key := fmt.Sprintf("%s_to_%s", item, rated)
			value, exists := (*rateList).Rating[key]
			if exists {
				numerator += task.Features[rated] * value
				denomertor += value
			} else {
				continue
			}
		}
		similarItemsRating = append(similarItemsRating, RecommendationItem{ID: item, SimilarityScore: numerator / denomertor})
	}
	sort.Slice(similarItemsRating, func(i, j int) bool {
		return similarItemsRating[i].SimilarityScore > similarItemsRating[j].SimilarityScore
	})
	// Return the top N items
	return similarItemsRating[:int(math.Min(float64(topN), float64(len(similarItemsRating))))]
}

// Content-based filtering
func processTask[T any](task T, workerId int) Recommendation {

	t := reflect.TypeOf(task)
	switch t {
	case reflect.TypeOf(deque.TaskCart{}):
		task, _ := any(task).(deque.TaskCart)
		topSimilarItems := FindTopSimilarItems(task.Info, task.Data, task.Count)
		recommendationResult := Recommendation{UserID: task.Info.ID, RecommendList: topSimilarItems}
		// jsonData, _ := json.Marshal(topSimilarItems)
		// fmt.Println("Worker " + strconv.Itoa(workerId) + " Processed Task" + strconv.Itoa(task.Task.ID)) //+ " Result: " + string(jsonData) + "recommanded to " + task.Info.ID)
		return recommendationResult
	case reflect.TypeOf(deque.TaskItem{}):
		task, _ := any(task).(deque.TaskItem)
		topSimilarItems := FindTopSimilarItemsUnRated(task.Info, task.Data, task.Count)
		recommendationResult := Recommendation{UserID: task.Info.ID, RecommendList: topSimilarItems}
		return recommendationResult
	default:
		// Handle other types or provide a default case
		fmt.Println("Unsupported task type")
	}
	return Recommendation{}
}

// workerProcessTasks simulates a worker processing tasks
func workerProcessTasks[T any](worker *Worker[T], workerPool *WorkStealingScheduler[T], result chan Recommendation) {
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
			result <- processTask(task, worker.ID)
		}
	}
}

func (ws *WorkStealingScheduler[T]) Run(result chan Recommendation, tasks *[]T) {
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
			workerProcessTasks(worker, ws, result)
		}(&ws.Workers[i])
	}
	// Wait for all workers to finish
	wg.Wait()
}

func main() {
	//Num of workers per Pool
	workersContent := 4
	workersItem := 4
	//Num of recommanded product per task
	recommandCount := 3
	//Num of tasks
	taskCount := 10

	UserPool := data.CreateRandomUserRatePool(100, 100, 0.1)

	//pre-computed data itemI -> ItemJ cosine similarity score based on who rated both of them if haveing n item, there will be (n)*(n-1)/2 by user pool
	similarity_martix := data.ComputeSimilarity_martix(UserPool)

	// Item-Based Collaborative Filtering
	taskItemPool := data.CreateRandomItemTask(taskCount, 100, 0.9)

	resultItem := make(chan Recommendation)
	workerPoolItem := newWorkStealingScheduler[deque.TaskItem](workersItem)
	tasksItem := make([]deque.TaskItem, taskCount)
	for i := 0; i < taskCount; i++ {
		tasksItem[i] = deque.TaskItem{Task: deque.Task{ID: i + 1, Count: recommandCount}, Info: taskItemPool[i], Data: &similarity_martix}
	}
	go workerPoolItem.Run(resultItem, &tasksItem)

	for i := 0; i < taskCount; i++ {
		fmt.Println(<-resultItem)
	}

	close(resultItem)
	//------------------------------------------------------------------------------------------------
	//content-based
	contentCart := data.CreateRandomContent(1000, 10, 0.5)
	taskCartPool := data.CreateRandomTasks(taskCount, 2, 5, 10, 0.5)

	resultCart := make(chan Recommendation)
	workerPool := newWorkStealingScheduler[deque.TaskCart](workersContent)
	tasks := make([]deque.TaskCart, taskCount)
	for i := 0; i < taskCount; i++ {
		tasks[i] = deque.TaskCart{Task: deque.Task{ID: i + 1, Count: recommandCount}, Info: taskCartPool[i], Data: &contentCart}
	}
	go workerPool.Run(resultCart, &tasks)

	for i := 0; i < taskCount; i++ {
		fmt.Println(<-resultCart)
	}
	close(resultCart)

	// go func() {
	// 	time.Sleep(2 * time.Second)
	// 	ch1 <- "Message from channel 1"
	// }()

	// // Goroutine to send a message on ch2 after 3 seconds
	// go func() {
	// 	time.Sleep(3 * time.Second)
	// 	ch2 <- "Message from channel 2"
	// }()

	// // Use select to wait for messages from either channel
	// select {
	// case msg1 := <-ch1:
	// 	fmt.Println(msg1)
	// case msg2 := <-ch2:
	// 	fmt.Println(msg2)
	// }
}
