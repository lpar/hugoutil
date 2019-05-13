package watson

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"sort"
	"strings"

	nlu "github.com/watson-developer-cloud/go-sdk/naturallanguageunderstandingv1"
)

// Watson is a wrapper around the Watson Natural Language Understanding API
type Watson struct {
	Service *nlu.NaturalLanguageUnderstandingV1
	Options *nlu.AnalyzeOptions
}

// Results represents a set of extracted categories and keywords.
type Results struct {
	Entities   []string
	Keywords   []string
	Categories []string
}

// NewWatson returns a Watson object, used to wrap the Watson Natural Language Understanding API and make it easier
// to call. You need to provide the API key and API URL obtained from IBM Cloud Dashboard.
func NewWatson(apikey string, apiurl string) (Watson, error) {
	w := Watson{}
	var err error
	w.Service, err = nlu.NewNaturalLanguageUnderstandingV1(
		&nlu.NaturalLanguageUnderstandingV1Options{
			URL:       apiurl,
			Version:   "2018-11-16",
			IAMApiKey: apikey,
		},
	)
	if err != nil {
		return w, err
	}
	w.Options = w.Service.NewAnalyzeOptions(&nlu.Features{
		Entities:   &nlu.EntitiesOptions{},
		Categories: &nlu.CategoriesOptions{},
		Keywords:   &nlu.KeywordsOptions{},
	})
	return w, nil
}

const maxKeywords = 10

// Analyze runs the HTML of a Hugo post through Watson Natural Language Understanding,
// and returns a simplified set of candidate categories, keywords and subjects (entities).
func (w Watson) Analyze(html string) (Results, error) {
	w.Options.HTML = &html
	resp, err := w.Service.Analyze(w.Options)
	if err != nil {
		return Results{}, err
	}
	watres := w.Service.GetAnalyzeResult(resp)
	// Preprocess and return the results
	results := Results{
		Categories: []string{},
		Keywords:   []string{},
		Entities:   []string{},
	}
	// Sort the categories and keywords by score/relevance
	sort.Slice(watres.Categories, func(i, j int) bool {
		return *watres.Categories[i].Score > *watres.Categories[j].Score
	})
	sort.Slice(watres.Keywords, func(i, j int) bool {
		return *watres.Keywords[i].Relevance > *watres.Keywords[j].Relevance
	})
	// Keep just the most detailed level of each category since Hugo doesn't have hierarchical categories by default
	for _, cat := range watres.Categories {
		results.Categories = append(results.Categories, path.Base(*cat.Label))
	}
	// Keep the top 10 keywords
	for i, k := range watres.Keywords {
		results.Keywords = append(results.Keywords, *k.Text)
		if i == maxKeywords {
			break
		}
	}
	for _, subj := range watres.Entities {
		results.Entities = append(results.Entities, *subj.Text)
	}
	return results, nil
}

func Interact(resin Results) Results {
	cats := append(resin.Entities, resin.Categories...)
	for i, c := range cats {
		fmt.Printf("%d: %s\n", i, c)
		if i == 9 {
			break
		}
	}
	for i, k := range resin.Keywords {
		fmt.Printf("%c: %s\n", i+'a', k)
		if i == 25 {
			break
		}
	}
	fmt.Printf("Select a category by number and any number of keywords by letter\n> ")
	reader := bufio.NewReader(os.Stdin)
	text := ""
	for {
		text, _ = reader.ReadString('\n')
		if strings.TrimSpace(text) != "" {
			break
		}
	}
	selectedCategories := []string{}
	for _, c := range text {
		if c >= '0' && c <= '9' {
			selectedCategories = append(selectedCategories, cats[c-'0'])
		}
	}
	selectedKeywords := []string{}
	for _, c := range text {
		if c >= 'a' && c <= 'z' {
			selectedKeywords = append(selectedKeywords, resin.Keywords[c-'a'])
		}
	}
	return Results{
		Categories: selectedCategories,
		Keywords:   selectedKeywords,
	}
}
