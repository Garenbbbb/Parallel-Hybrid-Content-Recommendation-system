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

type Content struct {
	ID       string
	Features map[string]float64
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

func CreateRandomContent(num int, limit int, prob float64) map[string]Content {

	data := make(map[string]Content)
	for i := 0; i < num; i++ {
		data["Product"+strconv.Itoa(i)] = Content{ID: "Product" + strconv.Itoa(i), Features: CreateRandomFeature(limit, prob)}
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

func CreateRandomUserRatePool(user int, item int) []Content {
	return []Content{
		{ID: "User1", Features: map[string]float64{"Item1": 5, "Item2": 0, "Item3": 2, "Item4": 5, "Item5": 5, "Item6": 5, "Item7": 2, "Item8": 5}},
		{ID: "User2", Features: map[string]float64{"Item1": 3, "Item2": 2, "Item3": 5, "Item4": 1, "Item5": 5, "Item6": 3, "Item7": 1, "Item8": 2}},
		{ID: "User3", Features: map[string]float64{"Item1": 2, "Item2": 5, "Item3": 3, "Item4": 4, "Item5": 2, "Item6": 5, "Item7": 6, "Item8": 1}},
		{ID: "User4", Features: map[string]float64{"Item1": 6, "Item2": 2, "Item3": 6, "Item4": 5, "Item5": 2, "Item6": 2, "Item7": 1, "Item8": 5}},
		{ID: "User5", Features: map[string]float64{"Item1": 1, "Item2": 5, "Item3": 1, "Item4": 2, "Item5": 5, "Item6": 5, "Item7": 2, "Item8": 4}},
		{ID: "User6", Features: map[string]float64{"Item1": 2, "Item2": 6, "Item3": 2, "Item4": 5, "Item5": 4, "Item6": 6, "Item7": 2, "Item8": 5}},
		{ID: "User7", Features: map[string]float64{"Item1": 0, "Item2": 8, "Item3": 5, "Item4": 0, "Item5": 5, "Item6": 0, "Item7": 5, "Item8": 0}},
		{ID: "User8", Features: map[string]float64{"Item1": 3, "Item2": 6, "Item3": 1, "Item4": 5, "Item5": 6, "Item6": 0, "Item7": 2, "Item8": 5}},
		{ID: "User9", Features: map[string]float64{"Item1": 5, "Item2": 3, "Item3": 2, "Item4": 5, "Item5": 1, "Item6": 5, "Item7": 4, "Item8": 5}},
		{ID: "User10", Features: map[string]float64{"Item1": 5, "Item2": 0, "Item3": 2, "Item4": 5, "Item5": 5, "Item6": 5, "Item7": 2, "Item8": 5}},
		{ID: "User11", Features: map[string]float64{"Item1": 5, "Item2": 0, "Item3": 2, "Item4": 5, "Item5": 5, "Item6": 5, "Item7": 2, "Item8": 5}},
		{ID: "User12", Features: map[string]float64{"Item1": 3, "Item2": 2, "Item3": 5, "Item4": 1, "Item5": 5, "Item6": 3, "Item7": 1, "Item8": 2}},
		{ID: "User13", Features: map[string]float64{"Item1": 2, "Item2": 5, "Item3": 3, "Item4": 4, "Item5": 2, "Item6": 5, "Item7": 6, "Item8": 1}},
		{ID: "User14", Features: map[string]float64{"Item1": 6, "Item2": 2, "Item3": 6, "Item4": 5, "Item5": 2, "Item6": 2, "Item7": 1, "Item8": 5}},
		{ID: "User15", Features: map[string]float64{"Item1": 1, "Item2": 5, "Item3": 1, "Item4": 2, "Item5": 5, "Item6": 5, "Item7": 2, "Item8": 4}},
		{ID: "User16", Features: map[string]float64{"Item1": 2, "Item2": 6, "Item3": 2, "Item4": 5, "Item5": 4, "Item6": 6, "Item7": 2, "Item8": 5}},
		{ID: "User17", Features: map[string]float64{"Item1": 0, "Item2": 8, "Item3": 5, "Item4": 0, "Item5": 5, "Item6": 0, "Item7": 5, "Item8": 0}},
		{ID: "User18", Features: map[string]float64{"Item1": 3, "Item2": 6, "Item3": 1, "Item4": 5, "Item5": 6, "Item6": 0, "Item7": 2, "Item8": 5}},
		{ID: "User19", Features: map[string]float64{"Item1": 5, "Item2": 3, "Item3": 2, "Item4": 5, "Item5": 1, "Item6": 5, "Item7": 4, "Item8": 5}},
	}
}

func ComputeSimilarity_martix(userRating []Content) map[string]Content {
	matrix := make(map[string]Content)
	matrix["similarity"] = Content{Features: make(map[string]float64)}
	itemCount := len(userRating[0].Features)
	for i := 1; i < itemCount; i++ {
		for j := i + 1; j <= itemCount; j++ {
			// Calculate cosine similarity
			item1 := "Item" + strconv.Itoa(i)
			item2 := "Item" + strconv.Itoa(j)
			similarity := cosineSimilarity(item1, item2, userRating)
			// Store the result in the map
			key := fmt.Sprintf("%s_to_%s", item1, item2)
			matrix["similarity"].Features[key] = similarity
		}
	}
	return matrix
}

func CreateRandomItemTask(num int, item int) []Content {
	return []Content{
		{ID: "User1", Features: map[string]float64{"Item1": 5, "Item2": 0, "Item3": 2, "Item4": 5, "Item5": 5, "Item6": 5, "Item7": 2, "Item8": 5}},
		{ID: "User2", Features: map[string]float64{"Item1": 3, "Item2": 2, "Item3": 5, "Item4": 1, "Item5": 5, "Item6": 3, "Item7": 1, "Item8": 2}},
		{ID: "User3", Features: map[string]float64{"Item1": 2, "Item2": 5, "Item3": 3, "Item4": 4, "Item5": 2, "Item6": 5, "Item7": 6, "Item8": 1}},
		{ID: "User4", Features: map[string]float64{"Item1": 6, "Item2": 2, "Item3": 6, "Item4": 5, "Item5": 2, "Item6": 2, "Item7": 1, "Item8": 5}},
		{ID: "User5", Features: map[string]float64{"Item1": 1, "Item2": 5, "Item3": 1, "Item4": 2, "Item5": 5, "Item6": 5, "Item7": 2, "Item8": 4}},
		{ID: "User6", Features: map[string]float64{"Item1": 2, "Item2": 6, "Item3": 2, "Item4": 5, "Item5": 4, "Item6": 6, "Item7": 2, "Item8": 5}},
		{ID: "User7", Features: map[string]float64{"Item1": 0, "Item2": 8, "Item3": 5, "Item4": 0, "Item5": 5, "Item6": 0, "Item7": 5, "Item8": 0}},
		{ID: "User8", Features: map[string]float64{"Item1": 3, "Item2": 6, "Item3": 1, "Item4": 5, "Item5": 6, "Item6": 0, "Item7": 2, "Item8": 5}},
		{ID: "User9", Features: map[string]float64{"Item1": 5, "Item2": 3, "Item3": 2, "Item4": 5, "Item5": 1, "Item6": 5, "Item7": 4, "Item8": 5}},
	}
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
