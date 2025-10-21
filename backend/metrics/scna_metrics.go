package main

import (
    "bufio"
    "fmt"
    "io/fs"
    "os"
    "path/filepath"
    "regexp"
    "strings"
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

func extractMetrics(code string) Metrics {
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

    // Regex patterns
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

func main() {
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
            metrics := extractMetrics(string(data))
            fmt.Println(path, metrics)
        }
        return nil
    })
}
