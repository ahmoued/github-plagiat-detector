package searchgithub

import (
    "context"
    
    "strings"
    "github.com/google/go-github/v55/github"
    "golang.org/x/oauth2"
    
    "sync"
    
)

type RepoInfo struct {
    Owner    string
    Name     string
    CloneURL string
}

func SearchRepos(client *github.Client, keywords []string, maxResults int, token string) ([]RepoInfo, error) {
    ctx := context.Background()
    /*var client *github.Client
    if token != "" {
        ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
        tc := oauth2.NewClient(ctx, ts)
        client = github.NewClient(tc)
    } else {
        client = github.NewClient(nil)
    }
*/ 


    query := strings.Join(keywords, " ") + " in:name,description"
    opts := &github.SearchOptions{Sort: "stars", Order: "desc", ListOptions: github.ListOptions{PerPage: maxResults}}

    result, _, err := client.Search.Repositories(ctx, query, opts)
    if err != nil {
        return nil, err
    }

    repos := []RepoInfo{}
    for i, r := range result.Repositories {
        if i >= maxResults {
            break
        }
        owner := r.GetOwner().GetLogin()
        name := r.GetName()
        repos = append(repos, RepoInfo{
            Owner:    owner,
            Name:     name,
            CloneURL: r.GetCloneURL(),
        })
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



func GetRepoWithReadme(ctx context.Context, client *github.Client, owner, name string) (*github.Repository, string, error) {
    repo, _, err := client.Repositories.Get(ctx, owner, name)
    if err != nil {
        return nil, "", err
    }

    file, _, err := client.Repositories.GetReadme(ctx, owner, name, nil)
    if err != nil {
        return repo, "", nil
    }

    content, err := file.GetContent()
    if err != nil {
        return repo, "", err
    }

    return repo, content, nil
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
