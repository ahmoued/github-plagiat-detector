package clone

import (
    "fmt"
    "os/exec"
    "sync"
    "github.com/ahmoued/github-plagiarism-backend/searchgithub"
)

type DownloadResult struct {
    Name     string
    LocalDir string
    Err      error
}



func CloneRepos(repos []searchgithub.RepoInfo) []DownloadResult {
    var wg sync.WaitGroup
    results := make([]DownloadResult, len(repos))

    for i, repo := range repos {
        wg.Add(1)
        go func(i int, repo searchgithub.RepoInfo) {
            defer wg.Done()
            localDir := fmt.Sprintf("./tmp/%s", repo.Name)
            cmd := exec.Command("git", "clone", "--depth", "1", repo.CloneURL, localDir)
            err := cmd.Run()
            results[i] = DownloadResult{
                Name:     repo.Name,
                LocalDir: localDir,
                Err:      err,
            }
        }(i, repo)
    }

    wg.Wait()
    return results
}
