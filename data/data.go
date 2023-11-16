package data

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
)

// contents := map[string] data.Content{
// 	{ID: "Item1", Features: map[string]float64{"Feature1": 0.2, "Feature2": 0.8}},
// 	{ID: "Item2", Features: map[string]float64{"Feature1": 0.1, "Feature2": 0.3,"Feature3": 0.4}},
// 	{ID: "Item3", Features: map[string]float64{"Feature1": 0.2, "Feature2": 0.1}},
// 	{ID: "Item4", Features: map[string]float64{"Feature1": 0.3, "Feature2": 0.2, "Feature3": 0.4, "Feature4": 0.3}},
// }

// taskPool := []data.ShopingCart{
// 	{ID: "User1", Items: []data.Content{{ID: "Item1", Features: map[string]float64{"Feature1": 0.2, "Feature2": 0.8}}}},
// 	{ID: "User2", Items: []data.Content{{ID: "Item2", Features: map[string]float64{"Feature1": 0.5}}}},
// 	{ID: "User3", Items: []data.Content{{ID: "Item3", Features: map[string]float64{"Feature1": 0.6, "Feature2": 0.3}}}},
// 	{ID: "User5", Items: []data.Content{{ID: "Item5", Features: map[string]float64{"Feature1": 0.4, "Feature2": 0.3, "Feature3": 0.4, "Feature4": 0.3}}}},
// }

// rating := []Content{
// 	{ID: "User1", Features: map[string]float64{"Item1": 5, "Item2": 0, "Item3": 2, "Item4": 5, "Item5": 5, "Item6": 5, "Item7": 2, "Item8": 5}},
// 	{ID: "User2", Features: map[string]float64{"Item1": 3, "Item2": 2, "Item3": 5, "Item4": 1, "Item5": 5, "Item6": 3, "Item7": 1, "Item8": 2}},
// 	{ID: "User3", Features: map[string]float64{"Item1": 2, "Item2": 5, "Item3": 3, "Item4": 4, "Item5": 2, "Item6": 5, "Item7": 6, "Item8": 1}},
// 	{ID: "User4", Features: map[string]float64{"Item1": 6, "Item2": 2, "Item3": 6, "Item4": 5, "Item5": 2, "Item6": 2, "Item7": 1, "Item8": 5}},
// 	{ID: "User5", Features: map[string]float64{"Item1": 1, "Item2": 5, "Item3": 1, "Item4": 2, "Item5": 5, "Item6": 5, "Item7": 2, "Item8": 4}},
// }

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
	itemCount := len(userRating[0].Features)
	for i := 1; i < itemCount; i++ {
		for j := i + 1; j <= itemCount; j++ {
			// Calculate cosine similarity
			item1 := "Item" + strconv.Itoa(i)
			item2 := "Item" + strconv.Itoa(j)
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
