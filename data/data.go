package data

import (
	"fmt"
	"math/rand"
	"strconv"
)

// contents := []data.Content{
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

func CreateRandomContent(num int, limit int, prob float64) []Content {
	data := make([]Content, num)
	for i := 0; i < num; i++ {
		data[i] = Content{ID: "Product" + strconv.Itoa(i), Features: CreateRandomFeature(limit, prob)}
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
