package auth

import (
	"sort"
	"strings"

	"github.com/spaceuptech/space-cloud/config"
)

func sortFileRule(rules []*config.FileRule) {

	sort.Slice(rules, func(i, j int) bool {
		return rules[i].Prefix < rules[j].Prefix
	})
	var splitKey int
	for key, val := range rules {
		if strings.Index(val.Prefix, "{") != -1 {
			splitKey = key
			break
		}
	}
	ar1 := rules[:splitKey]
	ar2 := rules[splitKey:]
	rules = append(bubbleSortFileRule(ar1), bubbleSortFileRule(ar2)...)
}

func bubbleSortFileRule(arr []*config.FileRule) []*config.FileRule {
	var lenArr []int
	for _, value := range arr {
		lenArr = append(lenArr, strings.Count(value.Prefix, "/"))
	}

	for i := 0; i < len(lenArr)-1; i++ {
		for j := 0; j < len(lenArr)-i-1; j++ {
			if lenArr[j] < lenArr[j+1] {
				temp := arr[j]
				arr[j] = arr[j+1]
				arr[j+1] = temp
				num := lenArr[j]
				lenArr[j] = lenArr[j+1]
				lenArr[j+1] = num
			}
		}
	}
	return arr
}

func sortPubsubRule(rules []*config.PubsubRule) {

	sort.Slice(rules, func(i, j int) bool {
		return rules[i].Subject < rules[j].Subject
	})
	var splitKey int
	for key, val := range rules {
		if strings.Index(val.Subject, "{") != -1 {
			splitKey = key
			break
		}
	}
	ar1 := rules[:splitKey]
	ar2 := rules[splitKey:]
	rules = append(bubbleSortPubsubRule(ar1), bubbleSortPubsubRule(ar2)...)
}

func bubbleSortPubsubRule(arr []*config.PubsubRule) []*config.PubsubRule {
	var lenArr []int
	for _, value := range arr {
		lenArr = append(lenArr, strings.Count(value.Subject, "/"))
	}

	for i := 0; i < len(lenArr)-1; i++ {
		for j := 0; j < len(lenArr)-i-1; j++ {
			if lenArr[j] < lenArr[j+1] {
				temp := arr[j]
				arr[j] = arr[j+1]
				arr[j+1] = temp
				num := lenArr[j]
				lenArr[j] = lenArr[j+1]
				lenArr[j+1] = num
			}
		}
	}
	return arr
}
