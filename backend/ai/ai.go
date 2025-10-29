package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type GeminiKeywords struct {
	Keywords []string `json:"keywords"`
}

func GetKeywordsFromGemini(repoName, description, readme, apiKey string) ([]string, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("empty API key")
	}

	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent?key=" + apiKey

	fullPrompt := fmt.Sprintf(`
Extract 8 functionality keywords from this repository:
Name: %s
Description: %s
README: %s
The words need to describe the functionality of the repo so that github search repos API can fetch relevant repos
for the top 3 most important keywords give also deviations of it as another separate keyword in example if the word editor is important also give edit or if document is important also give doc so that github search repos api can get relevant repos
Dont unclude keywords that are not very descriptive of the project like db or auth unless the murpose of the repo include something like that. in example a project whose goal is to faciliate the auth process or to introduce a new login method would need those auth and login words. Another project that simply handles login but thats not the goal shouldnt have them
the alterations only happen on the top 3 words that are very descriptive of the functionality, the words you give should be top1 then alteration of top1 then top2 then alteration of top2 then top 3 then alteration of top3 then top 7 then top 8
dont give any words related to he stack used as they dont tell about the functionality
this is to be used for github search repos api so the words you give need to cover the full scope of the goal of the project. Dont neglect a word that will make the github api miss some context
give top priority to the name of the repo , if the name is foot-manager then foot and football and manager definitely need to be in the keywords list.
read the whole data you get. phrase what does this project do exactly then give keywords so that someone who saw the keywords only could guess the purpose of the project. always keep in mind that those keywords will be used bu git search repos api to narrow the search for possible plagiat copies 

Return ONLY valid JSON in this format:
{"keywords": ["keyword1", "keyword2", ..., "keyword10"]}
Do NOT include any text before or after the JSON. Dont write the word json, this needs to be ready to use in code
give only one word keywords.
`, repoName, description, readme)

	reqBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": fullPrompt},
				},
			},
		},
	}

	jsonData, _ := json.Marshal(reqBody)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to call Gemini API: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("gemini returned status %d: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse Gemini envelope JSON: %v\nBody: %s", err, string(body))
	}

	// Extract text from possible Gemini response paths
	var text string
	if cands, ok := result["candidates"].([]interface{}); ok && len(cands) > 0 {
		if cand0, ok := cands[0].(map[string]interface{}); ok {
			if content, ok := cand0["content"].(map[string]interface{}); ok {
				if parts, ok := content["parts"].([]interface{}); ok && len(parts) > 0 {
					if p0, ok := parts[0].(map[string]interface{}); ok {
						if t, ok := p0["text"].(string); ok {
							text = t
						}
					}
				}
			}
		}
	}
	if text == "" {
		if out, ok := result["output"].(string); ok {
			text = out
		}
	}

	if strings.TrimSpace(text) == "" {
		// fallback: try to extract JSON substring directly from raw body
		s := strings.TrimSpace(string(body))
		if i := strings.Index(s, "{"); i >= 0 {
			j := strings.LastIndex(s, "}")
			if j > i {
				text = s[i : j+1]
			}
		}
	}

	// Clean Markdown or backtick fences
	clean := strings.TrimSpace(text)
	clean = strings.TrimPrefix(clean, "```")
	clean = strings.TrimPrefix(clean, "json")
	clean = strings.TrimSuffix(clean, "```")
	clean = strings.Trim(clean, "` \t\n\r")

	// Parse JSON
	var geminiResp GeminiKeywords
	if err := json.Unmarshal([]byte(clean), &geminiResp); err != nil {
		return nil, fmt.Errorf("failed to parse Gemini output JSON: %v\nRaw text: %s", err, clean)
	}

	return geminiResp.Keywords, nil
}
