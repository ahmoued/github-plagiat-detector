package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/ahmoued/github-plagiarism-backend/clone"
	"github.com/ahmoued/github-plagiarism-backend/compare"
	"github.com/ahmoued/github-plagiarism-backend/searchgithub"
	"github.com/ahmoued/github-plagiarism-backend/utils"
	"github.com/ahmoued/github-plagiarism-backend/metrics"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type CompareRequest struct {
    RepoURL string `json:"repo_url"`
}

type CompareResponse struct {
    Results []compare.CompareResultWithMetrics `json:"results"`
}

func compareHandler(w http.ResponseWriter, r *http.Request) {

	
    var req CompareRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }


	parts := strings.Split(strings.Trim(req.RepoURL, "/"), "/")
    if len(parts) < 2 {
        http.Error(w, "Invalid repo URL", http.StatusBadRequest)
        return
    }
    owner := parts[len(parts)-2]
    repoName := parts[len(parts)-1]


	ctx := context.Background()
    token := os.Getenv("GITHUB_TOKEN")
    fmt.Println(token)
    client := searchgithub.NewClient(token)
    repo, readmeContent, err := searchgithub.GetRepoWithReadme(ctx, client, owner, repoName)
    if err != nil {
        http.Error(w, "Failed to fetch repo info", http.StatusInternalServerError)
        return
    }

    inputText := repo.GetName() + " " + repo.GetDescription() + " " + readmeContent


	keywords := utils.ExtractKeywordsFromText(inputText)
    keys := []string{"collaboration", "tool", "doc"}
    if len(keywords) == 0 {
        keywords = []string{repo.GetName()} 
    }


	maxResults := 4
    candidateRepos, err := searchgithub.SearchRepos(client, keys, maxResults)
    if err != nil {
        http.Error(w, "GitHub search failed", http.StatusInternalServerError)
        return
    }


	readmes := searchgithub.FetchReadmes(client, candidateRepos, "")


	minOverlap := 2
    filteredRepoKeys := utils.FilterReposByReadme(readmes, keywords, minOverlap)


	filteredRepos := []searchgithub.RepoInfo{}
    keySet := make(map[string]struct{})
    for _, key := range filteredRepoKeys {
        keySet[key] = struct{}{}
    }
    for _, r := range candidateRepos {
        if _, ok := keySet[r.Owner+"/"+r.Name]; ok {
            filteredRepos = append(filteredRepos, r)
        }
    }


	clonedResults := clone.CloneRepos(candidateRepos)



	inputClone := clone.DownloadResult{
        Name:     repo.GetName(),
        LocalDir: "./tmp/input_repo", 
    }

    inputClone = clone.CloneInputRepo(searchgithub.RepoInfo{Owner: owner, Name: repoName, CloneURL: req.RepoURL})


    results, clonedMetrics, inputMetrics, err := compare.CompareReposCode(inputClone.LocalDir, clonedResults)
    if err!= nil{
        fmt.Println(err.Error())
        http.Error(w, "Error comparing repos", http.StatusInternalServerError)
        return
    }
    fmt.Println(clonedMetrics)
    fmt.Println(inputMetrics)
   

    var wg sync.WaitGroup
    metricResultsChan := make(chan compare.CompareResult, len(clonedMetrics))
    for _, c := range clonedMetrics{
        wg.Add(1)
        go func(repoName string, candidate metrics.Metrics) {
            defer wg.Done()
            sim := metrics.ComputeMetricsSimilarity(inputMetrics, candidate)
            metricResultsChan <- compare.CompareResult{
                Repo:       repoName,
                Similarity: sim,
            }
        }(c.RepoName, c.Metrics)
    }

    wg.Wait()
    close(metricResultsChan)

    var metricResults []compare.CompareResult
    for r := range metricResultsChan {
        metricResults = append(metricResults, r)
    }

    fmt.Println("Metric-based similarities:")
    for _, r := range metricResults {
        fmt.Printf("%s → %.2f%%\n", r.Repo, r.Similarity)
    }


    combinedResults := []compare.CompareResultWithMetrics{}
    for _, r := range results {

        var metricSim float64
        for _, m := range metricResults {
            if m.Repo == r.Repo {
                metricSim = m.Similarity
                break
            }
        }
        combinedResults = append(combinedResults, compare.CompareResultWithMetrics{
            Repo:       r.Repo,
            TokenSimilarity:  r.Similarity,
            MetricsSimilarity: metricSim,
        })
    }

    fmt.Println("combined Results")
    fmt.Println(combinedResults)
    

	resp := CompareResponse{
        Results: combinedResults,
    }
    fmt.Println(resp)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/compare", compareHandler).Methods("POST")

    handler := cors.AllowAll().Handler(r)

    log.Println("Server running on http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", handler))
}
