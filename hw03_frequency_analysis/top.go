package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

func Top10(text string) []string {
	sliceStrings := strings.Fields(text)
	mapText := make(map[string]int)
	var normalizedWord string

	for _, v := range sliceStrings {
		normalizedWord = normalizeWord(v)
		if len(normalizedWord) > 0 {
			mapText[normalizedWord]++
		}
	}

	keys := sortedWords(mapText)

	return top10Words(keys)
}

// приведение слов к нормальному виду.
func normalizeWord(word string) string {
	// Приводим к нижнему регистру
	word = strings.ToLower(word)

	// исключение для много повторений тире ----
	if len(strings.Trim(word, "-")) == 0 && len(word) > 1 {
		return strings.TrimSpace(word)
	}

	// Удаляем знаки препинания с начала слова
	re := regexp.MustCompile(`^[^\wа-яА-Я]+`)
	word = re.ReplaceAllString(word, "")

	// Удаляем знаки препинания с конца слова
	re = regexp.MustCompile(`[^\wа-яА-Я]+$`)
	word = re.ReplaceAllString(word, "")

	return strings.TrimSpace(word)
}

// сортировка ключей в мапе.
func sortedWords(mapText map[string]int) []string {
	// Создаем срез ключей
	keys := make([]string, 0, len(mapText))
	for k := range mapText {
		keys = append(keys, k)
	}

	// Сортируем ключи по убыванию их значений
	sort.Slice(keys, func(i, j int) bool {
		if mapText[keys[i]] == mapText[keys[j]] {
			return keys[i] < keys[j] // лексикографически
		}
		return mapText[keys[i]] > mapText[keys[j]]
	})

	return keys
}

// Получаем топ-10 слов.
func top10Words(keys []string) []string {
	var topWords []string
	for i := 0; i < 10 && i < len(keys); i++ {
		topWords = append(topWords, keys[i])
	}

	return topWords
}
