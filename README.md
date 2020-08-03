
# hugoutil

This is a command-line utility to aid with tagging and categorizing [Hugo](https://gohugo.io/) posts.
Like Hugo, it is written in Go for speed.

## Features

 - Batch add and remove tags and categories from the command line.
 - Automatically extract metadata from article text using [IBM Watson][watson] [Natural Language Understanding][nlu],
 and either apply it to the source files as is, or interactively select from the identified tags and categories.
 - Convert Hugo frontmatter between TOML and YAML formats, including batch conversion. (Hugo [can also do this][convert], 
 but I don't find it as convenient.)

IBM Watson functionality requires an IBM Cloud account. You can run this program using an 
[IBM Cloud Lite][cloud] account, no credit card required.

[nlu]: https://www.ibm.com/cloud/watson-natural-language-understanding/resources
[cloud]: https://www.ibm.com/cloud/free/
[watson]: https://www.ibm.com/watson
[convert]: https://gohugo.io/commands/hugo_convert/

## Disclaimer

Because this program can update many files in a single run, it can cause massive damage to your data if you aren't careful.
Please make sure you have a backup of your content files before use. I recommend a Git repository, then you can use 
`git diff` to check the changes before committing them. As per the license, no warranty is offered.
A bug could eat your latest posting if you don't make a copy of it first.

## Installation

```
git clone https://github.com/lpar/hugoutil
cd hugoutil
go build
```

## Example command line use

```
% hugoutil --untag 'US' --tag 'politics,USA' --uncategorize war --categorize 'Civil War' 04/*.md
```

See `--help` for a description of the supported options.

## Example interactive Watson session

```
% hugoutil -i 04/14.md
Updating 04/14.md (Lincoln assassinated)
0: John Wilkes Booth
1: President Abraham Lincoln
2: US
3: White House
4: Virginia
5: 12 days
6: unrest and war
7: assault
8: government
a: 11th April
b: stage actor John Wilkes Booth
c: President Abraham Lincoln
d: White House
e: Confederate sympathizer
f: April
g: important officials of the US government
h: northern Virginia
i: Booth
j: conspirators
k: assassination
Select categories by number and any number of keywords by letter
> 01k
```

## Notes

 - Uncategorizing/untagging happens before categorizing/tagging. If you remove a tag and add the same tag in a single run, it will end up set.

 - You can use the tagging and categorization functionality without needing any kind of cloud account and without sending any data to the cloud.

 - This utility is a not an official IBM product. It's a utility I wrote for my own use, and a demonstration of how you can call IBM Watson using the [official IBM Watson Cloud SDK for Go](https://github.com/watson-developer-cloud/go-sdk).
 
 - If you're interested in deploying Go web applications to IBM Cloud, check out [IBM-Cloud/get-started-go][getstarted].
 
[getstarted]: https://github.com/IBM-Cloud/get-started-go

## Copyright

Copyright © IBM Corporation 2019-2020. Apache License 2.0.
 
