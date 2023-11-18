package data

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
)

type Content struct {
	ID       string
	Features map[string]float64
}

type ItemData struct {
	Features map[string]float64
}

type RateData struct {
	Rating map[string]float64
}

type ShopingCart struct {
	ID    string
	Items []Content
}

var USER_POOL = []Content{
	{ID: "UserA", Features: map[string]float64{"Apple": 7, "Orange": 7, "Banana": 5, "Spinach": 5, "Tomato": 6}},
	{ID: "UserB", Features: map[string]float64{"Apple": 8, "Orange": 6, "Banana": 5, "Spinach": 6, "Tomato": 5}},
	{ID: "UserC", Features: map[string]float64{"Apple": 6, "Orange": 5, "Banana": 4, "Spinach": 4, "Tomato": 4}},
	{ID: "UserD", Features: map[string]float64{"Apple": 5, "Orange": 7, "Banana": 6, "Spinach": 5, "Tomato": 5}},
	{ID: "UserE", Features: map[string]float64{"Apple": 8, "Orange": 7, "Banana": 6, "Spinach": 4, "Tomato": 6}},
}

var USER = []Content{
	{ID: "User1", Features: map[string]float64{"Apple": 7, "Orange": 0, "Banana": 7, "Spinach": 0, "Tomato": 0}},
	{ID: "User2", Features: map[string]float64{"Apple": 6, "Orange": 6, "Banana": 0, "Spinach": 6, "Tomato": 0}},
	{ID: "User3", Features: map[string]float64{"Apple": 0, "Orange": 0, "Banana": 0, "Spinach": 6, "Tomato": 8}},
	{ID: "User4", Features: map[string]float64{"Apple": 0, "Orange": 0, "Banana": 8, "Spinach": 0, "Tomato": 0}},
}

var ITEM_POOL = map[string]ItemData{
	"Peach":     {Features: map[string]float64{"Sweet": 0.8, "Healthy": 0.6, "Fruit": 1, "Soft": 0.2}},
	"BlueBerry": {Features: map[string]float64{"Sweet": 0.6, "Healthy": 0.8, "Fruit": 1, "Soft": 0.8}},
	"Kiwi":      {Features: map[string]float64{"Sweet": 0.5, "Healthy": 0.9, "Fruit": 1, "Soft": 0.8}},
	"Avocado":   {Features: map[string]float64{"Sweet": 0.3, "Healthy": 0.7, "Fruit": 1, "Soft": 0.4}},
	"Cucumber":  {Features: map[string]float64{"Sweet": 0.1, "Healthy": 0.8, "Fruit": 1, "Soft": 0.1}},
	"carrot":    {Features: map[string]float64{"Sweet": 0.2, "Healthy": 0.8, "Fruit": 1, "Soft": 0.1}},
	"Burger":    {Features: map[string]float64{"Sweet": 0.1, "Healthy": 0.2, "Fruit": 0.2, "Soft": 0.5}},
	"Fires":     {Features: map[string]float64{"Sweet": 0.3, "Healthy": 0.1, "Fruit": 0.1, "Soft": 0.6}},
	"Pizza":     {Features: map[string]float64{"Sweet": 0.4, "Healthy": 0.1, "Fruit": 0.1, "Soft": 0.7}},
}

var CART_TASK = []ShopingCart{
	{ID: "User1", Items: []Content{{ID: "Apple", Features: map[string]float64{"Sweet": 0.6, "Healthy": 0.8, "Fruit": 1, "Soft": 0.3}}}},
	{ID: "User2", Items: []Content{{ID: "Celery", Features: map[string]float64{"Sweet": 0.2, "Healthy": 0.8, "Fruit": 1, "Soft": 0.2}}}},
	{ID: "User3", Items: []Content{{ID: "strawberry", Features: map[string]float64{"Sweet": 0.5, "Healthy": 0.6, "Fruit": 1, "Soft": 0.7}}}},
	{ID: "User4", Items: []Content{{ID: "corn", Features: map[string]float64{"Sweet": 0.5, "Healthy": 0.6, "Fruit": 1, "Soft": 0.4}}}},
}

func CreateRandomFeature(limit int, prob float64) map[string]float64 {
	features := make(map[string]float64)
	for i := 1; i <= limit; i++ {
		randomValue := rand.Float64()
		featureName := fmt.Sprintf("Feature%d", i)
		if randomValue > prob {
			features[featureName] = 0
			continue
		}
		randomValue = float64(int(rand.Float64()*100)) / 100
		features[featureName] = randomValue
	}
	return features
}

func CreateRandomContent(num int, limit int, prob float64) map[string]ItemData {

	data := make(map[string]ItemData)
	for i := 0; i < num; i++ {
		data["Item"+strconv.Itoa(i)] = ItemData{Features: CreateRandomFeature(limit, prob)}
	}
	return data
}

func CreateRandomCart(num int, limit int, prob float64) []Content {
	charSet := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	lenC := len(charSet)
	data := make([]Content, num)
	for i := 0; i < num; i++ {
		data[i] = Content{ID: "Item-" + string(charSet[rand.Intn(lenC)]) + string(charSet[rand.Intn(lenC)]), Features: CreateRandomFeature(limit, prob)}
	}
	return data
}

func CreateRandomTasks(num int, item_lower int, item_upper int, feature_limit int, prob float64) []ShopingCart {
	data := make([]ShopingCart, num)
	for i := 0; i < num; i++ {
		data[i] = ShopingCart{ID: "User" + strconv.Itoa(i), Items: CreateRandomCart(rand.Intn(item_upper-item_lower+1)+item_lower, feature_limit, prob)}
	}
	return data
}

func CreateRandomRating(item int, prob float64) map[string]float64 {
	rating := make(map[string]float64)
	for i := 1; i <= item; i++ {
		randomValue := rand.Float64()
		featureName := fmt.Sprintf("Item%d", i)
		if randomValue <= prob {
			rating[featureName] = 0
			continue
		}
		rating[featureName] = float64(rand.Intn(10))
	}
	return rating
}

// prob to rate be 0
func CreateRandomUserRatePool(user int, item int, prob float64) []Content {
	data := make([]Content, user)
	for i := 0; i < user; i++ {
		data[i] = Content{ID: "RandomUser" + strconv.Itoa(i), Features: CreateRandomRating(item, prob)}
	}
	return data
}

func ComputeSimilarity_martix(userRating []Content) RateData {
	matrix := RateData{Rating: make(map[string]float64)}
	itemCount := userRating[0].Features
	var keys []string
	for key := range itemCount {
		keys = append(keys, key)
	}
	for i := 0; i < len(keys)-1; i++ {
		for j := i + 1; j < len(keys); j++ {
			// Calculate cosine similarity
			item1 := keys[i]
			item2 := keys[j]
			similarity := cosineSimilarity(item1, item2, userRating)
			// Store the result in the map
			key1 := fmt.Sprintf("%s_to_%s", item1, item2)
			matrix.Rating[key1] = similarity
			key2 := fmt.Sprintf("%s_to_%s", item2, item1)
			matrix.Rating[key2] = similarity
		}
	}
	return matrix
}

// prob of 0 rate
func CreateRandomItemTask(num int, item int, prob float64) []Content {
	data := make([]Content, num)
	for i := 0; i < num; i++ {
		data[i] = Content{ID: "User" + strconv.Itoa(i), Features: CreateRandomRating(item, prob)}
	}
	return data
}

func cosineSimilarity(item1 string, item2 string, ratings []Content) float64 {
	dotProduct := 0.0
	magnitude1 := 0.0
	magnitude2 := 0.0

	for _, user := range ratings {
		// Skip missing ratings represented by "-"
		if user.Features[item1] == 0 || user.Features[item2] == 0 {
			continue
		}

		dotProduct += user.Features[item1] * user.Features[item2]
		magnitude1 += math.Pow(user.Features[item1], 2)
		magnitude2 += math.Pow(user.Features[item2], 2)
	}

	magnitude1 = math.Sqrt(magnitude1)
	magnitude2 = math.Sqrt(magnitude2)

	if magnitude1 == 0 || magnitude2 == 0 {
		return 0.0
	}

	return dotProduct / (magnitude1 * magnitude2)
}
