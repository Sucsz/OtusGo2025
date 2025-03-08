package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type WordCount struct {
	Word  string
	Count int
}

func Top10(str string) []string {
	top10Words := make([]string, 0, 10)
	mapWordCount := FreqAnalysis(str)
	sortedSlice := SortFreqAnalysisMap(mapWordCount)

	if len(sortedSlice) > 10 {
		sortedSlice = sortedSlice[:10]
	}
	for i := 0; i < len(sortedSlice); i++ {
		top10Words = append(top10Words, sortedSlice[i].Word)
	}

	return top10Words
}

func FreqAnalysis(str string) map[string]int {
	mapWordCount := make(map[string]int)
	splitStr := strings.Fields(str)

	for _, word := range splitStr {
		mapWordCount[word]++
	}

	return mapWordCount
}

func SortFreqAnalysisMap(m map[string]int) []WordCount {
	sortedSlice := make([]WordCount, 0, len(m))

	for k, v := range m {
		sortedSlice = append(sortedSlice, WordCount{Word: k, Count: v})
	}

	sort.Slice(sortedSlice, func(i, j int) bool {
		// Sort.Slice(...) требует Less function, поэтому используя > сортировка по возрастанию, < - по убыванию
		if sortedSlice[i].Count == sortedSlice[j].Count {
			return sortedSlice[i].Word < sortedSlice[j].Word
		}
		return sortedSlice[i].Count > sortedSlice[j].Count
	})

	return sortedSlice
}
