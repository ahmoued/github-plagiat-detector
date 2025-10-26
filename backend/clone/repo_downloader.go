package clone

import (
    "fmt"
    "io/fs"
    "os/exec"
    "sync"
    "github.com/ahmoued/github-plagiarism-backend/searchgithub"
    "path/filepath"
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

func DirSize(path string) (int64, error) {
    var size int64
    err := filepath.WalkDir(path, func(_ string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }
        if !d.IsDir() {
            info, err := d.Info()
            if err != nil {
                return err
            }
            size += info.Size()
        }
        return nil
    })
    return size, err
}


func CloneInputRepo(repo searchgithub.RepoInfo) DownloadResult {
        
        
            
            localDir := "./tmp/input_repo"
            cmd := exec.Command("git", "clone", "--depth", "1", repo.CloneURL, localDir)
            err := cmd.Run()
            return DownloadResult{
                Name:     repo.Name,
                LocalDir: localDir,
                Err:      err,
            }
        
    }

