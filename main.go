package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"hugoutil/content"
	"hugoutil/watson"

	"github.com/spf13/pflag"
)

var yaml = pflag.BoolP("yaml", "y", false, "use YAML instead of TOML for writing frontmatter")
var verbose = pflag.BoolP("verbose", "v", false, "output more info about what's going on")
var wat = pflag.BoolP("watson", "w", false, "add metadata from Watson analysis")
var interact = pflag.BoolP("interactive", "i", false, "prompt to choose keywords and category from Watson")
var addtags = pflag.String("tag", "", "comma-separated list of tags to add")
var deltags = pflag.String("untag", "", "comma-separated list of tags to remove")
var addcats = pflag.String("categorize", "", "comma-separated list of categories to add")
var delcats = pflag.String("uncategorize", "", "comma-separated list of categories to remove")
var help = pflag.BoolP("help", "h", false, "get help")

// cslParse parses a comma-separated list from the command line, returning a slice of whitespace-trimmed strings
func cslParse(x string) []string {
	vals := strings.Split(x, ",")
	result := make([]string, 0, len(vals))
	for _, v := range vals {
		result = append(result, strings.TrimSpace(v))
	}
	return result
}

func main() {
	pflag.Parse()

	if *help || len(pflag.Args()) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: hugoutil [OPTION]... [FILE]...\n\n")
		pflag.PrintDefaults()
		return
	}

	if *interact {
		*wat = true
	}

	config, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Config error: %v\n", err)
		return
	}

	var watsonService *watson.Watson
	if *wat {
		w, err := watson.NewWatson(config.APIKey, config.APIURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Problem using Watson: %v\n", err)
			return
		}
		watsonService = &w
		runtime.Breakpoint()
	}

	for _, file := range pflag.Args() {
		if *verbose && !*wat {
			fmt.Printf("Processing %s\n", file)
		}
		art, err := content.Read(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", file, err)
			return
		}

		// Manipulate here
		if watsonService != nil {
			res, err := watsonService.Analyze(art.HTML())
			if err != nil {
				fmt.Fprintf(os.Stderr, "Problem using Watson: %v\n", err)
				return
			}
			if *interact {
				fmt.Printf("Updating %s (%s)\n", file, art.Header["title"])
				res = watson.Interact(res)
			}
			art.Tags.AddAll(res.Keywords)
			art.Categories.AddAll(res.Categories)
		}

		if *delcats != "" {
			art.Categories.RemoveAll(cslParse(*delcats))
		}
		if *deltags != "" {
			art.Tags.RemoveAll(cslParse(*deltags))
		}
		if *addcats != "" {
			art.Categories.AddAll(cslParse(*addcats))
		}
		if *addtags != "" {
			art.Tags.AddAll(cslParse(*addtags))
		}

		if *yaml {
			err = art.Write(file, content.YAML)
		} else {
			err = art.Write(file, content.TOML)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", file, err)
			return
		}
	}
}
