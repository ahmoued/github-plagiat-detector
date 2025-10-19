package main

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "strings"

    "github.com/ahmoued/github-plagiarism-backend/searchgithub"
    "github.com/ahmoued/github-plagiarism-backend/utils"
    "github.com/ahmoued/github-plagiarism-backend/clone"
    "github.com/ahmoued/github-plagiarism-backend/compare"

    "github.com/gorilla/mux"
    "github.com/rs/cors"
)

type CompareRequest struct {
    RepoURL string `json:"repo_url"`
}

type CompareResponse struct {
    Results []compare.CompareResult `json:"results"`
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
    client := searchgithub.NewClient("")
    repo, readmeContent, err := searchgithub.GetRepoWithReadme(ctx, client, owner, repoName)
    if err != nil {
        http.Error(w, "Failed to fetch repo info", http.StatusInternalServerError)
        return
    }

    inputText := repo.GetName() + " " + repo.GetDescription() + " " + readmeContent


	keywords := utils.ExtractKeywordsFromText(inputText)
    if len(keywords) == 0 {
        keywords = []string{repo.GetName()} 
    }


	maxResults := 50
    candidateRepos, err := searchgithub.SearchRepos(client, keywords, maxResults, "")
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


	clonedResults := clone.CloneRepos(filteredRepos)


	inputClone := clone.DownloadResult{
        Name:     repo.GetName(),
        LocalDir: "./tmp/input_repo", 
    }

	inputCode, _ := compare.ReadCodeFiles(inputClone.LocalDir)
    results := []compare.CompareResult{}
    for _, c := range clonedResults {
        code, _ := compare.ReadCodeFiles(c.LocalDir)
        sim := compare.ComputeSimilarity(inputCode, code)
        results = append(results, compare.CompareResult{
            Repo:       c.Name,
            Similarity: sim,
        })
    }


	resp := CompareResponse{
        Results: results,
    }
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
