package compare 

import(
	"strings"
	"io/fs"
	"os"
	"path/filepath"
	"github.com/ahmoued/github-plagiarism-backend/clone"
)

type CompareResult struct {
    Repo       string  `json:"repo"`
    Similarity float64 `json:"similarity"` 
}

func ReadCodeFiles(path string) (string, error) {
    var codeStrings []string
    err := filepath.Walk(path, func(p string, info fs.FileInfo, err error) error {
        if err != nil || info.IsDir() {
            return err
        }
        ext := filepath.Ext(p)
        if ext == ".go" || ext == ".js" || ext == ".ts" {
            data, err := os.ReadFile(p)
            if err != nil {
                return err
            }
            cleaned := normalizeCode(string(data))
            codeStrings = append(codeStrings, cleaned)
        }
        return nil
    })
    if err != nil {
        return "", err
    }
    return strings.Join(codeStrings, "\n"), nil
}

func normalizeCode(s string) string {
    lines := strings.Split(s, "\n")
    var cleaned []string
    for _, l := range lines {
        l = strings.TrimSpace(l)
        if strings.HasPrefix(l, "//") || l == "" {
            continue
        }
        cleaned = append(cleaned, l)
    }
    return strings.Join(cleaned, "\n")
}



func ComputeSimilarity(a, b string) float64 {
    aTokens := strings.Fields(a)
    bTokens := strings.Fields(b)

    if len(aTokens) == 0 || len(bTokens) == 0 {
        return 0
    }

    matches := 0
    tokenSet := make(map[string]bool)
    for _, t := range aTokens {
        tokenSet[t] = true
    }
    for _, t := range bTokens {
        if tokenSet[t] {
            matches++
        }
    }

    
    total := (len(aTokens) + len(bTokens)) / 2
    return float64(matches) / float64(total) * 100
}

func CompareReposCode(inputPath string, cloned []clone.DownloadResult) []CompareResult {
    inputCode, _ := ReadCodeFiles(inputPath)
    results := []CompareResult{}

    for _, c := range cloned {
        code, _ := ReadCodeFiles(c.LocalDir)
        sim := ComputeSimilarity(inputCode, code)
        results = append(results, CompareResult{
            Repo:       c.Name,
            Similarity: sim,
        })
    }

    return results
}
