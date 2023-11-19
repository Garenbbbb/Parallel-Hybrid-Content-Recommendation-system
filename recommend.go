package main

import (
	"Work-Stealing-Based-Parallel-Hybrid-Content-Recommendation-system/data"
	"Work-Stealing-Based-Parallel-Hybrid-Content-Recommendation-system/deque"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"reflect"
	"sort"
	"strconv"
	"sync"
	"time"
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

func roundToDecimal(value float64, decimalPlaces int) float64 {
	precision := math.Pow(10, float64(decimalPlaces))
	return math.Round(value*precision) / precision
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
	return similarItems[:int(math.Min(float64(topN), float64(len(similarItems))))]
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
		similarItemsRating = append(similarItemsRating, RecommendationItem{ID: item, SimilarityScore: roundToDecimal(numerator/denomertor, 2)})
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

func resultGatherProcess(finalResult chan Recommendation, taskCount int, resultItem chan Recommendation, resultCart chan Recommendation) {
	resultMap := make(map[string][]RecommendationItem)
	cnt := 0
	for {
		if cnt == taskCount*2 {
			break
		}
		// Use select to wait for messages from either channel
		select {
		case msg1 := <-resultItem:
			cnt += 1
			val, ok := resultMap[msg1.UserID]
			if ok {
				finalResult <- Recommendation{UserID: msg1.UserID, RecommendList: append(msg1.RecommendList, val...)}
			} else {
				resultMap[msg1.UserID] = msg1.RecommendList
			}
		case msg2 := <-resultCart:
			cnt += 1
			val, ok := resultMap[msg2.UserID]
			if ok {
				finalResult <- Recommendation{UserID: msg2.UserID, RecommendList: append(msg2.RecommendList, val...)}
			} else {
				resultMap[msg2.UserID] = msg2.RecommendList
			}
		}
	}
}

const usage = "Usage for large random test: go run recommend.go random (Number of Worker for Content based) (Number of Worker for Collaborative Filter)(REMINDER! set up values in config.json)"
const sampple = "Usage for small sample test: go run recommend.go sample (Number of Worker for Content based) (Number of Worker for Collaborative Filter)"

// Config represents the configuration structure.
type Config struct {
	Tasks        int `json:"Tasks"`
	UserPool     int `json:"UserPool"`
	RatePool     int `json:"RatePool"`
	ItemPool     int `json:"ItemPool"`
	CartNumUpper int `json:"CartNumUpper"`
	CartNumLower int `json:"CartNumLower"`
	FeatureNum   int `json:"FeatureNum"`
	RecommandNum int `json:"RecommandNum"`
}

func loadConfig(filename string) Config {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening config file:", err)
		// Handle the error appropriately, e.g., by using default values or exiting the program
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	config := Config{}
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("Error decoding config JSON:", err)
		// Handle the error appropriately, e.g., by using default values or exiting the program
	}

	return config
}

func main() {

	arg := os.Args
	if len(arg) != 4 {
		fmt.Println(usage)
		fmt.Println(sampple)
		return
	}
	mode := arg[1]

	//Num of workers per Pool
	workersContent, _ := strconv.Atoi(arg[2])
	workersItem, _ := strconv.Atoi(arg[3])
	//Num of recommanded product per task
	var recommandCount int
	//Num of tasks
	var taskCount int
	var userpoolCnt int
	var UserPool []data.Content
	var taskItemPool []data.Content
	var contentCart map[string]data.ItemData
	var taskCartPool []data.ShopingCart

	var config Config

	if mode == "random" {
		config = loadConfig("config.json")
		taskCount = config.Tasks
		recommandCount = config.RecommandNum
		userpoolCnt = config.UserPool
		UserPool = data.CreateRandomUserRatePool(userpoolCnt, config.RatePool, 0.1)
		taskItemPool = data.CreateRandomItemTask(taskCount, config.RatePool, 0.9)
		contentCart = data.CreateRandomShop(config.ItemPool, config.FeatureNum, 0.5)
		taskCartPool = data.CreateRandomTasks(taskCount, config.CartNumLower, config.CartNumUpper, config.FeatureNum, 0.5)
	} else {
		workersContent = 1
		workersItem = 1
		taskCount = 6
		recommandCount = 2
		UserPool = data.USER_POOL
		taskItemPool = data.USER
		contentCart = data.ITEM_POOL
		taskCartPool = data.CART_TASK
		for _, item := range taskCartPool {
			fmt.Println(item.ID + " likes " + item.Items[0].ID)
		}
	}

	//pre-computed data itemI -> ItemJ cosine similarity score based on who rated both of them if haveing n item, there will be (n)*(n-1)/2 by user pool
	a := time.Now()
	similarity_martix := data.ComputeSimilarity_martix(UserPool)
	b := time.Now()
	fmt.Println(b.Sub(a))

	startTime := time.Now()
	fmt.Println("----------------------Similarity Matrix computed Start Processing----------------------")
	// Item-Based Collaborative Filtering
	resultItem := make(chan Recommendation)
	workerPoolItem := newWorkStealingScheduler[deque.TaskItem](workersItem)
	tasksItem := make([]deque.TaskItem, taskCount)
	for i := 0; i < taskCount; i++ {
		tasksItem[i] = deque.TaskItem{Task: deque.Task{ID: i + 1, Count: recommandCount}, Info: taskItemPool[i], Data: &similarity_martix}
	}

	//------------------------------------------------------------------------------------------------
	//content-based
	resultCart := make(chan Recommendation)
	workerPool := newWorkStealingScheduler[deque.TaskCart](workersContent)
	tasks := make([]deque.TaskCart, taskCount)
	for i := 0; i < taskCount; i++ {
		tasks[i] = deque.TaskCart{Task: deque.Task{ID: i + 1, Count: recommandCount}, Info: taskCartPool[i], Data: &contentCart}
	}
	go workerPoolItem.Run(resultItem, &tasksItem)
	go workerPool.Run(resultCart, &tasks)

	finalResult := make(chan Recommendation)

	//combine result from content-based and Item-Based Collaborative Filtering
	go resultGatherProcess(finalResult, taskCount, resultItem, resultCart)

	for i := 0; i < taskCount; i++ {
		fmt.Println(<-finalResult)
	}
	endTime := time.Now()
	fmt.Println(endTime.Sub(startTime).Seconds())
	close(resultItem)
	close(resultCart)
	close(finalResult)

}
