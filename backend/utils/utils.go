package utils


import (
    "bufio"
    "math"
    "os"
    "path/filepath"
    "regexp"
    "sort"
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

func DetectLanguage(repoPath string) string {
    exts := map[string]int{}

    filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
        if err != nil || info.IsDir() {
            return nil
        }
        ext := filepath.Ext(path)
        if ext != "" {
            exts[ext]++
        }
        return nil
    })

    // Pick the most frequent extension
    var maxCount int
    var mainExt string
    for ext, count := range exts {
        if count > maxCount {
            maxCount = count
            mainExt = ext
        }
    }

    switch mainExt {
    case ".go":
        return "go"
    case ".py":
        return "python"
    case ".js", ".jsx":
        return "javascript"
    case ".ts", ".tsx":
        return "typescript"
    case ".java":
        return "java"
    default:
        return "unknown"
    }
}




// KeywordScore holds a word and its TF-IDF score
type KeywordScore struct {
    Word  string
    Score float64
}

// ExtractFunctionalKeywords extracts functional keywords from repo metadata
// metadataText: concatenated repo name, description, README
// srcDir: path to source code files
// topN: number of keywords to return
func ExtractFunctionalKeywords(metadataText, srcDir string, topN int) ([]string, error) {
    // Preprocess metadata
    metadataText = strings.ToLower(metadataText)
    words := tokenize(metadataText, ignoreWords)

    // Compute term frequency (TF) in metadata
    tf := make(map[string]float64)
    totalWords := float64(len(words))
    for _, w := range words {
        tf[w]++
    }
    for k := range tf {
        tf[k] = tf[k] / totalWords
    }

    // Compute document frequency (DF) across source code files
    df := make(map[string]float64)
    numDocs := 0.0
    err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
        if err != nil || info.IsDir() || (!strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, ".py") &&
            !strings.HasSuffix(path, ".js") && !strings.HasSuffix(path, ".ts")) {
            return nil
        }
        numDocs++
        seen := make(map[string]struct{})
        file, err := os.Open(path)
        if err != nil {
            return nil
        }
        defer file.Close()

        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
            lineWords := tokenize(scanner.Text(), ignoreWords)
            for _, w := range lineWords {
                seen[w] = struct{}{}
            }
        }
        for w := range seen {
            df[w]++
        }
        return nil
    })
    if err != nil {
        return nil, err
    }

    // Compute TF-IDF
    tfidf := make(map[string]float64)
    for w, freq := range tf {
        idf := 1.0
        if df[w] > 0 {
            idf = math.Log((numDocs + 1) / (df[w] + 1))
        }
        tfidf[w] = freq * idf
    }

    // Sort by TF-IDF score
    sorted := make([]KeywordScore, 0, len(tfidf))
    for w, score := range tfidf {
        sorted = append(sorted, KeywordScore{w, score})
    }
    sort.Slice(sorted, func(i, j int) bool {
        return sorted[i].Score > sorted[j].Score
    })

    // Return top N words
    result := []string{}
    for i := 0; i < len(sorted) && i < topN; i++ {
        result = append(result, sorted[i].Word)
    }
    return result, nil
}

// tokenize splits text into lowercase words, removing punctuation
/*func tokenize(text string) []string {
    re := regexp.MustCompile(`[a-zA-Z]+`)
    matches := re.FindAllString(text, -1)
    for i := range matches {
        matches[i] = strings.ToLower(matches[i])
    }
    return matches
}*/

func tokenize(text string, ignore map[string]bool) []string {
    re := regexp.MustCompile(`[a-zA-Z]+`)
    matches := re.FindAllString(text, -1)
    var tokens []string
    for _, word := range matches {
        word = strings.ToLower(word)
        if !ignore[word] {
            tokens = append(tokens, word)
        }
    }
    return tokens
}


var ignoreWords = map[string]bool{
    "docker": true, "database": true, "sql": true, "nosql": true,
    "relational": true, "github": true, "multi-model": true,
    "deployment": true, "backend": true, "js": true, "env": true,
    "and": true, "build": true, "kubectl": true,
}
