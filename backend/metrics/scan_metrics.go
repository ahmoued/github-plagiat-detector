package metrics

import (
    "fmt"
    "io/fs"
    "os"
    "path/filepath"
    "regexp"
    "strings"
    "math"
)

type Metrics struct {
    NumLines          int
    AvgLineLength     float64
    NumFunctions      int
    NumClasses        int
    NumLoops          int
    NumConditionals   int
    NumImports        int
    NumComments       int
    CyclomaticEstimate int
}

func ExtractMetrics(code string) Metrics {
    var m Metrics
    lines := strings.Split(code, "\n")
    totalLen := 0
    for _, line := range lines {
        trimmed := strings.TrimSpace(line)
        if trimmed == "" {
            continue
        }
        m.NumLines++
        totalLen += len(trimmed)
        if strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "/*") {
            m.NumComments++
        }
    }
    if m.NumLines > 0 {
        m.AvgLineLength = float64(totalLen) / float64(m.NumLines)
    }


    funcRegex := regexp.MustCompile(`\b(func|def|function)\b`)
    classRegex := regexp.MustCompile(`\b(class|type|struct)\b`)
    loopRegex := regexp.MustCompile(`\b(for|while|foreach|range)\b`)
    condRegex := regexp.MustCompile(`\b(if|else if|switch|case)\b`)
    importRegex := regexp.MustCompile(`\b(import|require|from|use)\b`)
    branchRegex := regexp.MustCompile(`\b(if|for|while|case|&&|\|\|)\b`)

    m.NumFunctions = len(funcRegex.FindAllString(code, -1))
    m.NumClasses = len(classRegex.FindAllString(code, -1))
    m.NumLoops = len(loopRegex.FindAllString(code, -1))
    m.NumConditionals = len(condRegex.FindAllString(code, -1))
    m.NumImports = len(importRegex.FindAllString(code, -1))
    m.CyclomaticEstimate = len(branchRegex.FindAllString(code, -1))

    return m
}

func ScanRepoMetrics() {
    repoPath := "./myrepo"
    filepath.WalkDir(repoPath, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }
        if !d.IsDir() {
            data, err := os.ReadFile(path)
            if err != nil {
                return err
            }
            metrics := ExtractMetrics(string(data))
            fmt.Println(path, metrics)
        }
        return nil
    })
}





func ComputeMetricsSimilarity(input, candidate Metrics) float64 {
    
    single := func(inp, cand float64) float64 {
       
        if cand == 0 {
            if inp == 0 {
                return 1.0
            }
            return 0.0 
        }
        d := math.Abs(inp - cand) / cand 
        if d >= 1 {
            return 0.0 
        }
        return 1.0 - d 
    }

    wNumLines := 0.10
    wAvgLineLen := 0.05
    wNumFuncs := 0.15
    wNumClasses := 0.15
    wNumLoops := 0.10
    wNumCond := 0.10
    wNumImports := 0.05
    wNumComments := 0.05
    wCyclomatic := 0.25


    score := 0.0
    score += wNumLines * single(float64(input.NumLines), float64(candidate.NumLines))
    score += wAvgLineLen * single(input.AvgLineLength, candidate.AvgLineLength)
    score += wNumFuncs * single(float64(input.NumFunctions), float64(candidate.NumFunctions))
    score += wNumClasses * single(float64(input.NumClasses), float64(candidate.NumClasses))
    score += wNumLoops * single(float64(input.NumLoops), float64(candidate.NumLoops))
    score += wNumCond * single(float64(input.NumConditionals), float64(candidate.NumConditionals))
    score += wNumImports * single(float64(input.NumImports), float64(candidate.NumImports))
    score += wNumComments * single(float64(input.NumComments), float64(candidate.NumComments))
    score += wCyclomatic * single(float64(input.CyclomaticEstimate), float64(candidate.CyclomaticEstimate))

    return score * 100.0
}
