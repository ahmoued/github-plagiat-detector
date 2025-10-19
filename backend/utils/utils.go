package utils

import (
    "strings"
)

var stopwords = map[string]bool{
    "the": true, "and": true, "for": true, "with": true, "from": true, "this": true,
    "that": true, "a": true, "an": true, "in": true, "on": true, "by": true,
}

func ExtractKeywordsFromText(text string) []string {
    words := strings.Fields(strings.ToLower(text))
    keywords := []string{}
    seen := map[string]bool{}
    for _, w := range words {
        w = strings.Trim(w, ".,:;()[]{}\"'")
        if len(w) > 3 && !stopwords[w] && !seen[w] {
            keywords = append(keywords, w)
            seen[w] = true
        }
    }
    return keywords
}

func FilterReposByReadme(candidateReadmes map[string]string, keywords []string, minOverlap int) []string {
    filtered := []string{}
    for repo, readme := range candidateReadmes {
        count := 0
        lowerReadme := strings.ToLower(readme)
        for _, kw := range keywords {
            if strings.Contains(lowerReadme, strings.ToLower(kw)) {
                count++
            }
        }
        if count >= minOverlap {
            filtered = append(filtered, repo)
        }
    }
    return filtered
}