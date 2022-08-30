package content

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/russross/blackfriday/v2"
	"gopkg.in/yaml.v2"
)

type Article struct {
	Header     map[string]interface{}
	Body       []byte
	Tags       StringSet
	Categories StringSet
}

// HTML renders an article to HTML in a way that can be content-analyzed by Watson.
// This includes taking any existing title and description and adding them to the markdown body
// before converting it to HTML.
func (art *Article) HTML() string {
	title := art.Header["title"].(string)
	var md bytes.Buffer
	if title != "" {
		md.WriteString("# ")
		md.WriteString(title)
		md.WriteString("\n\n")
	}
	desc, ok := art.Header["description"]
	if ok {
		sdesc := desc.(string)
		md.WriteString(sdesc)
		md.WriteString("\n\n")
	}
	md.Write(art.Body)
	body := blackfriday.Run(md.Bytes())
	return string(body)
}

var tomlSeparator = []byte("+++\n")
var yamlSeparator = []byte("---\n")

type Format int

const (
	TOML Format = iota
	YAML
)

func Read(filename string) (*Article, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return &Article{}, fmt.Errorf("can't read file %s: %w", filename, err)
	}
	art, err := Parse(data)
	if err != nil {
		return &Article{}, fmt.Errorf("%s: %w", filename, err)
	}
	return art, nil
}

func Parse(data []byte) (*Article, error) {
	art := Article{}
	var err error
	separator := data[0:4]
	var format Format
	if bytes.Equal(separator, tomlSeparator) {
		format = TOML
	} else if bytes.Equal(separator, yamlSeparator) {
		format = YAML
	} else {
		return &art, fmt.Errorf("no YAML or TOML frontmatter found")
	}
	chunks := bytes.SplitN(data, separator, 3)
	switch format {
	case TOML:
		err = toml.Unmarshal(chunks[1], &art.Header)
	case YAML:
		err = yaml.Unmarshal(chunks[1], &art.Header)
	}
	art.Body = chunks[2]
	if err != nil {
		return &art, fmt.Errorf("parse error: %w", err)
	}
	art.Tags = NewStringSet(ToStringSlice(art.Header["tags"]))
	art.Categories = NewStringSet(ToStringSlice(art.Header["categories"]))
	return &art, nil
}

func exists(name string) (bool, error) {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err != nil, err
}

func (a *Article) writeHeader(writer io.Writer, format Format) error {
	var werr error
	check := func(err error) {
		if err != nil {
			werr = fmt.Errorf("error writing header: %w", err)
		}
	}
	var err error
	switch format {
	case TOML:
		_, err = writer.Write(tomlSeparator)
		check(err)
		err = toml.NewEncoder(writer).Encode(a.Header)
		check(err)
		_, err = writer.Write(tomlSeparator)
		check(err)
	case YAML:
		_, err = writer.Write(yamlSeparator)
		check(err)
		err = yaml.NewEncoder(writer).Encode(a.Header)
		check(err)
		_, err = writer.Write(yamlSeparator)
		check(err)
	}
	if werr != nil {
		return fmt.Errorf("error marshalling frontmatter: %w", err)
	}
	return nil
}

func (a *Article) Write(filename string, format Format) error {
	var perms os.FileMode = 0600
	var bakfile string
	exists, err := exists(filename)
	a.Header["tags"] = a.Tags.Slice()
	a.Header["categories"] = a.Categories.Slice()
	if err != nil {
		return fmt.Errorf("can't check existence of %s: %w", filename, err)
	}
	if exists {
		finfo, err := os.Stat(filename)
		if err != nil {
			return fmt.Errorf("can't check permissions of %s: %w", bakfile, err)
		}
		perms = finfo.Mode().Perm()
		bakfile = filename + ".bak"
		err = os.Rename(filename, bakfile)
		if err != nil {
			return fmt.Errorf("can't create backup file %s: %w", bakfile, err)
		}
	}
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, perms)
	if err != nil {
		return fmt.Errorf("can't open %s for writing: %w", filename, err)
	}
	defer file.Close()

	err = a.writeHeader(file, format)
	if err != nil {
		return fmt.Errorf("error writing %s: %w", filename, err)
	}
	_, err = file.Write(a.Body)
	if err != nil {
		return fmt.Errorf("error writing body to %s: %w", filename, err)
	}

	if exists {
		err = os.Remove(bakfile)
		if err != nil {
			return fmt.Errorf("failed to remove backup file %s: %w", bakfile, err)
		}
	}
	return nil
}

func ToStringSlice(x interface{}) []string {
	ss := []string{}
	switch v := x.(type) {
	case []string:
		ss = append(ss, v...)
	case []interface{}:
		for _, str := range v {
			ss = append(ss, str.(string))
		}
	case string:
		ss = append(ss, v)
	case interface{}:
		ss = append(ss, v.(string))
	}
	return ss
}
