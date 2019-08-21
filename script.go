package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"time"
)

func train(filename string) map[string]int {
	WORDS := make(map[string]int)
	symbols := regexp.MustCompile("[àáéíóúâôêçñãõa-z]+")
	if content, err := ioutil.ReadFile(filename); err == nil {
		for _, w := range symbols.FindAllString(strings.ToLower(string(content)), -1) {
			WORDS[w]++
		}
	} else {
		panic("Failed loading data from training file")
	}
	return WORDS
}

func edits1(word string, ch chan string) {
	const alphabet = "abcdefghijklmnopqrstuvwxyzáéíóúâôêçñãõ"
	type Pair struct{ a, b string }

	var splits []Pair
	for i := 0; i < len(word)+1; i++ {
		splits = append(splits, Pair{word[:i], word[i:]})
	}

	for _, split := range splits {
		if len(split.b) > 0 {
			ch <- split.a + split.b[1:]
		}

		if len(split.b) > 1 {
			ch <- split.a + string(split.b[1]) + string(split.b[0]) + split.b[2:]
		}

		for _, c := range alphabet {
			if len(split.b) > 0 {
				ch <- split.a + string(c) + split.b[1:]
			}
		}

		for _, c := range alphabet {
			ch <- split.a + string(c) + split.b
		}
	}
}

func edits2(word string, ch chan string) {
	ch1 := make(chan string, 1024*1024)
	go func() { edits1(word, ch1); ch1 <- "" }()
	for e1 := range ch1 {
		if e1 == "" {
			break
		}
		edits1(e1, ch)
	}
}

func best(word string, edits func(string, chan string), model map[string]int) string {
	ch := make(chan string, 1024*1024)
	go func() { edits(word, ch); ch <- "" }()
	maxFreq := 0
	correction := ""
	for word := range ch {
		if word == "" {
			break
		}
		if freq, present := model[word]; present && freq > maxFreq {
			maxFreq, correction = freq, word
		}
	}
	return correction
}

func correct(word string, model map[string]int) string {
	if _, present := model[word]; present {
		return word
	}

	if correction := best(word, edits1, model); correction != "" {
		return correction
	}

	if correction := best(word, edits2, model); correction != "" {
		return correction
	}

	return word
}

func main() {
	suffix := "en"
	if input, err := ioutil.ReadFile(fmt.Sprintf("misspelling_%s.txt", suffix)); err == nil {
		var totalTime float64
		pattern := regexp.MustCompile("(?P<Misspell>[àáéíóúâôêçñãõa-z]+)->(?P<Correct>[àáéíóúâôêçñãõa-z,]+)")
		model := train(fmt.Sprintf("dictionary_%s.txt", suffix))
		totalFixed := 0
		totalLines := 0
		for _, lines := range pattern.FindAllStringSubmatch(strings.ToLower(string(input)), -1) {
			word := lines[1]
			corrections := strings.Split(lines[2], ",")
			startTime := time.Now()
			fixed := correct(word, model)
			totalTime += time.Now().Sub(startTime).Seconds()
			for _, correction := range corrections {
				if fixed == correction {
					totalFixed++
					break
				} else {
					fmt.Printf("\nExpect %s but got %s\n", lines[0], fixed)
				}
			}
			totalLines++
		}
		fmt.Printf("\nTime : %v seconds\n", totalTime)
		fmt.Printf("Precision : %v out of %v (%v%%)\n", totalFixed, totalLines, float32(totalFixed)/float32(totalLines)*100.0)
		fmt.Println("Finished")
	} else {
		panic("Failed loading data from misspellings file")
	}
}
