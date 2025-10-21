package searchgithub

import (
	"context"
	"fmt"
	//"time"

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

/*func ProgressiveSearch(client *github.Client, keywords []string, maxResults int) ([]RepoInfo, error) {
	ctx := context.Background()
	found := make(map[string]RepoInfo)
	var mu sync.Mutex
	var wg sync.WaitGroup
	resultsCh := make(chan RepoInfo, maxResults*len(keywords)) // buffered channel

	// Helper to generate all subsets of size n
	subsets := func(words []string, n int) [][]string {
		if n > len(words) {
			return nil
		}
		var res [][]string
		var helper func(start int, path []string)
		helper = func(start int, path []string) {
			if len(path) == n {
				tmp := make([]string, n)
				copy(tmp, path)
				res = append(res, tmp)
				return
			}
			for i := start; i < len(words); i++ {
				helper(i+1, append(path, words[i]))
			}
		}
		helper(0, []string{})
		return res
	}

	// Worker function for each keyword subset
	sem := make(chan struct{}, 5) // max 5 concurrent requests

    searchWorker := func(subset []string) {
        defer wg.Done()
        sem <- struct{}{}           // acquire
        defer func() { <-sem }()   // release after finishing
        time.Sleep(2 * time.Second)

        query := strings.Join(subset, " ") + " in:name,description,readme"
        opts := &github.SearchOptions{Sort: "stars", Order: "desc", ListOptions: github.ListOptions{PerPage: 50}}
        result, _, err := client.Search.Repositories(ctx, query, opts)
        if err != nil {
            fmt.Println("GitHub search error:", err)
            return
        }

        for _, r := range result.Repositories {
            resultsCh <- RepoInfo{
                Owner:    r.GetOwner().GetLogin(),
                Name:     r.GetName(),
                CloneURL: r.GetCloneURL(),
            }
        }
    }


	// Launch goroutines for all subsets
	for k := len(keywords); k >= 3; k-- {
		for _, subset := range subsets(keywords, k) {
			wg.Add(1)
			go searchWorker(subset)
		}
	}

	// Close channel once all goroutines are done
	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	// Collect results and deduplicate
	for r := range resultsCh {
		key := r.Owner + "/" + r.Name
		mu.Lock()
		if _, exists := found[key]; !exists && len(found) < maxResults {
			found[key] = r
		}
		mu.Unlock()
	}

	// Convert map to slice
	repos := make([]RepoInfo, 0, len(found))
	for _, r := range found {
		repos = append(repos, r)
        fmt.Println(r.Name)
	}
        fmt.Println(repos)
        fmt.Println(len(repos))

	return repos, nil
}
*/


func SearchRepos(client *github.Client, keywords []string, maxResults int) ([]RepoInfo, error) {
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


    query := strings.Join(keywords, " ") + " in:name,description,readme"
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



func GetRepoWithReadme(ctx context.Context, client *github.Client, owner, name string) (*github.Repository, string, error) {
    fmt.Println("starting to try to get the repo")
    repo, _, err := client.Repositories.Get(ctx, owner, name)
    fmt.Println("trying to get the repo")
    if err != nil {
        return nil, "", err
    }
    fmt.Println("got the repo")

    file, _, err := client.Repositories.GetReadme(ctx, owner, name, nil)
    if err != nil {
        return repo, "", nil
    }
    fmt.Println("got the readme")

    content, err := file.GetContent()
    if err != nil {
        return repo, "", err
    }
    fmt.Println("got the content")

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
