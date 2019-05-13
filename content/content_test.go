package content

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	md := `+++
tags = ["one", "two"]
title = "Integers"
categories = ["alpha", "beta"]
+++

Dull body text
`

	art, err := Parse([]byte(md))
	if err != nil {
		t.Errorf("test parse failed: %v", err)
	}
	if bytes.Compare(art.Body, []byte("\nDull body text\n")) != 0 {
		t.Errorf("body not preserved")
	}
	if !compareSlices(art.Tags.Slice(), []string{"one", "two"}) {
		t.Errorf("tag parse/fetch failed, got %#v", art.Tags.Slice())
	}
	if !compareSlices(art.Categories.Slice(), []string{"alpha", "beta"}) {
		t.Errorf("category parse/fetch failed, got %#v", art.Categories.Slice())
	}
	art.Tags.AddAll([]string{"three", "four"})
	art.Categories.AddAll([]string{"gamma"})
	// Write the file out and read it back in
	tmpfile, err := ioutil.TempFile("", "hugoutil_content_test")
	if err != nil {
		t.Errorf("can't create temporary file: %v", err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()
	art.Write(tmpfile.Name(), TOML)
	txt, err := ioutil.ReadFile(tmpfile.Name())
	if err != nil {
		t.Errorf("can't read temporary file %s: %v", tmpfile.Name(), err)
	}
	newart, err := Parse(txt)
	if err != nil {
		t.Errorf("new file bad: %v", err)
	}
	if !compareSlices(newart.Categories.Slice(), []string{"alpha", "beta", "gamma"}) {
		t.Errorf("new file had wrong categories")
	}
	if !compareSlices(newart.Tags.Slice(), []string{"one", "two", "three", "four"}) {
		t.Errorf("new file had wrong tags")
	}
}
