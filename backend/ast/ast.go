// astcompare/astcompare.go
package astcompare

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/ahmoued/github-plagiarism-backend/clone"
)

type Result struct {
	RepoName   string  `json:"repo"`
	Similarity float64 `json:"similarity"`
}

var supportedExts = map[string]bool{
	".js":  true,
	".ts":  true,
	".tsx": true,
	".jsx": true,
	".go":  true,
}

// CompareReposFull compares ALL supported files (all-vs-all) between input and each candidate
func CompareReposFull(inputRepoDir string, candidates []clone.DownloadResult) ([]Result, error) {
	tmpDir := "./tmp/asts"
	os.MkdirAll(tmpDir, 0755)

	inputASTDir := filepath.Join(tmpDir, "input")
	os.RemoveAll(inputASTDir)
	os.MkdirAll(inputASTDir, 0755)
	exportRepoASTs(inputRepoDir, inputASTDir)

	// Get list of all input files (for all-vs-all)
	var inputFiles []string
	filepath.Walk(inputRepoDir, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && isSupportedFile(path) {
			inputFiles = append(inputFiles, path)
		}
		return nil
	})

	var results []Result
	var wg sync.WaitGroup
	resultsChan := make(chan Result, len(candidates))

	for _, cand := range candidates {
		if cand.Err != nil {
			continue
		}
		wg.Add(1)
		go func(c clone.DownloadResult) {
			defer wg.Done()

			candASTDir := filepath.Join(tmpDir, c.Name)
			os.RemoveAll(candASTDir)
			os.MkdirAll(candASTDir, 0755)
			exportRepoASTs(c.LocalDir, candASTDir)

			// Get all candidate files
			var candFiles []string
			filepath.Walk(c.LocalDir, func(path string, info os.FileInfo, err error) error {
				if err == nil && !info.IsDir() && isSupportedFile(path) {
					candFiles = append(candFiles, path)
				}
				return nil
			})

			var scores []float64

			// All-vs-all: every input file vs every candidate file
			for _, inPath := range inputFiles {
				inRel, _ := filepath.Rel(inputRepoDir, inPath)
				inJSON := filepath.Join(inputASTDir, inRel+".json")
				if _, err := os.Stat(inJSON); os.IsNotExist(err) {
					continue
				}

				for _, candPath := range candFiles {
					candRel, _ := filepath.Rel(c.LocalDir, candPath)
					candJSON := filepath.Join(candASTDir, candRel+".json")
					if _, err := os.Stat(candJSON); os.IsNotExist(err) {
						continue
					}

					if score, err := runAPTED(inJSON, candJSON); err == nil {
						scores = append(scores, score)
					}
				}
			}

			// Take average (or max, if you prefer)
			avgScore := 0.0
			if len(scores) > 0 {
				for _, s := range scores {
					avgScore += s
				}
				avgScore /= float64(len(scores))
			}

			resultsChan <- Result{
				RepoName:   c.Name,
				Similarity: avgScore,
			}
		}(cand)
	}

	wg.Wait()
	close(resultsChan)

	for r := range resultsChan {
		results = append(results, r)
	}

	return results, nil
}

// --- Helpers (same as before) ---

func exportRepoASTs(repoDir, outDir string) {
	_ = filepath.Walk(repoDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !isSupportedFile(path) {
			return nil
		}
		relPath, _ := filepath.Rel(repoDir, path)
		jsonPath := filepath.Join(outDir, relPath+".json")
		_ = os.MkdirAll(filepath.Dir(jsonPath), 0755)
		_ = runASTExport(path, jsonPath)
		return nil
	})
}

func isSupportedFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return supportedExts[ext]
}

func runASTExport(srcFile, jsonFile string) error {
	cmd := exec.Command("python3", "ast_export.py", srcFile, jsonFile)
	return cmd.Run()
}

func runAPTED(json1, json2 string) (float64, error) {
	cmd := exec.Command("python3", "apted_runner.py", json1, json2)
	out, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("apted failed: %w", err)
	}
	var score float64
	_, err = fmt.Sscanf(string(out), "%f", &score)
	return score, err
}