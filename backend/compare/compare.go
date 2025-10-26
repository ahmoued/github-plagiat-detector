package compare

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/ahmoued/github-plagiarism-backend/clone"
	"github.com/ahmoued/github-plagiarism-backend/metrics"
)




type CompareResult struct {
    Repo       string  `json:"repo"`
    Similarity float64 `json:"similarity"` 
}

type CompareResultWithMetrics struct {
    Repo       string  `json:"repo"`
    TokenSimilarity float64 `json:"similarity"`
    MetricsSimilarity float64 `json:"metrics_similarity"`
    ASTSimilarity float64 `json:"ast_similarity"`

}


/*type Metrics struct {
    NumLines          int
    AvgLineLength     float64
    NumFunctions      int
    NumClasses        int
    NumLoops          int
    NumConditionals   int
    NumImports        int
    NumComments       int
    CyclomaticEstimate int
}*/





func ReadCodeFiles(path string) (string, metrics.Metrics, error) {
    var codeStrings []string
    var repoMetrics metrics.Metrics
    var fileCount int
    err := filepath.Walk(path, func(p string, info fs.FileInfo, err error) error {
        if err != nil || info.IsDir() {
            return err
        }
        ext := filepath.Ext(p)
        if ext == ".go" || ext == ".js" || ext == ".ts" || ext == ".py" || ext == ".vue"{
            data, err := os.ReadFile(p)
            if err != nil {
                return err
            }
            cleaned := normalizeCode(string(data))
            codeStrings = append(codeStrings, cleaned)

            

            m := metrics.ExtractMetrics(string(data))
            repoMetrics.NumLines += m.NumLines
            repoMetrics.AvgLineLength += m.AvgLineLength
            repoMetrics.NumFunctions += m.NumFunctions
            repoMetrics.NumClasses += m.NumClasses
            repoMetrics.NumLoops += m.NumLoops
            repoMetrics.NumConditionals += m.NumConditionals
            repoMetrics.NumImports += m.NumImports
            repoMetrics.NumComments += m.NumComments
            repoMetrics.CyclomaticEstimate += m.CyclomaticEstimate

            fileCount++
        }
        return nil
    })
    if err != nil {
        return "", repoMetrics, err
    }
    if fileCount > 0 {
        repoMetrics.AvgLineLength /= float64(fileCount)
    }
    fmt.Println("metrics of this repo")
    fmt.Println()
    return strings.Join(codeStrings, "\n"), repoMetrics, nil

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

func CompareReposCode(inputPath string, cloned []clone.DownloadResult) ([]CompareResult, []struct{
    RepoName string
    Metrics metrics.Metrics
}, metrics.Metrics, error) {
    inputCode, inputMetrics, err := ReadCodeFiles(inputPath)
    fmt.Println("input metrics")
    fmt.Println(inputMetrics)
    if err != nil {
        fmt.Println(err.Error())
    }
    results := []CompareResult{}
    allMetrics :=[]struct{
        RepoName string
        Metrics metrics.Metrics
    }{}
    

    for _, c := range cloned {
        code, clonedmetrics, err := ReadCodeFiles(c.LocalDir)
        if err != nil{
            fmt.Println(err.Error())
            return nil, nil, inputMetrics, err
        }
        fmt.Println(c.Name)
        fmt.Println(clonedmetrics)
        sim := ComputeSimilarity(inputCode, code)
        results = append(results, CompareResult{
            Repo:       c.Name,
            Similarity: sim,
        })
        allMetrics = append(allMetrics, struct{RepoName string; Metrics metrics.Metrics}{RepoName: c.Name, Metrics: clonedmetrics})
    }

    return results, allMetrics, inputMetrics, nil
}
