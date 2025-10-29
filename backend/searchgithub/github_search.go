package searchgithub

import (
	"context"
	"fmt"
    "github.com/google/go-github/v55/github"
	"golang.org/x/oauth2"

	"sync"
)

type RepoInfo struct {
    Owner    string
    Name     string
    CloneURL string
}
 



func SearchRepos(client *github.Client, keywords []string, maxResults int, size int, lang string) ([]RepoInfo, error) {
    ctx := context.Background()
    

    queryLogic := fmt.Sprintf("(%s OR %s) (%s OR %s) (%s OR %s) OR %s OR %s", keywords[0], keywords[1], keywords[2], keywords[3], keywords[4], keywords[5], keywords[6], keywords[7])
    query := fmt.Sprintf("%s language:%s in:name,description,readme", queryLogic, lang)

    if size > 0 {
        minSize := int(size * 8 / 10)
        maxSize := int(size * 12 / 10)
        query += fmt.Sprintf(" size:%d..%d", minSize, maxSize)
    }
    fmt.Println(query)

    opts := &github.SearchOptions{Order: "desc", ListOptions: github.ListOptions{PerPage: maxResults}}

    result, _, err := client.Search.Repositories(ctx, query, opts)
    if err != nil {
        return nil, err
    }
    fmt.Println("good so far")
    fmt.Println(*result.Total)
    repos := []RepoInfo{}
    for i, r := range result.Repositories {
        if i >= maxResults {
            break
        }
        owner := r.GetOwner().GetLogin()
        fmt.Println(owner)
        name := r.GetName()
        repos = append(repos, RepoInfo{
            Owner:    owner,
            Name:     name,
            CloneURL: r.GetCloneURL(),
        })
        fmt.Println(owner, name, "waa")
    }
    return repos, nil
}
func FetchReadmes(client *github.Client, repos []RepoInfo, token string) map[string]string {
    ctx := context.Background()
    /*var client *github.Client
    if token != "" {
        ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
        tc := oauth2.NewClient(ctx, ts)
        client = github.NewClient(tc)
    } else {
        client = github.NewClient(nil)
    }*/

    results := make(map[string]string)
    var mu sync.Mutex
    var wg sync.WaitGroup

    for _, repo := range repos {
        wg.Add(1)
        go func(r RepoInfo) {
            defer wg.Done()
            file, _, err := client.Repositories.GetReadme(ctx, r.Owner, r.Name, nil)
            if err != nil {
                return
            }
            content, err := file.GetContent()
            if err != nil {
                return
            }
            mu.Lock()
            results[r.Owner+"/"+r.Name] = content
            mu.Unlock()
        }(repo)
    }

    wg.Wait()
    return results
}



func GetRepoWithReadme(ctx context.Context, client *github.Client, owner, name string) (*github.Repository, string, int, error) {
    fmt.Println("starting to try to get the repo")
    repo, _, err := client.Repositories.Get(ctx, owner, name)
    fmt.Println("trying to get the repo")
    if err != nil {
        return nil, "", 0, err
    }
    size := repo.GetSize()
    fmt.Println("official github size")
    fmt.Println(size)
    fmt.Println("got the repo")

    file, _, err := client.Repositories.GetReadme(ctx, owner, name, nil)
    if err != nil {
        return repo, "", size, nil
    }
    fmt.Println("got the readme")

    content, err := file.GetContent()
    if err != nil {
        return repo, "", size, err
    }
    fmt.Println("got the content")

    return repo, content, size, nil
}


func NewClient(token string) *github.Client {
    ctx := context.Background()

    if token != "" {
        ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
        tc := oauth2.NewClient(ctx, ts)
        return github.NewClient(tc)
    }

    return github.NewClient(nil)
}
