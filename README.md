# Laverna
Download Your Comics
For now, only works with http://xoxocomics.com links

Pull Requests encouraged!

## Build

```bash
# Just build all the binaries
scripts/build.sh

# Build binaries and upload them to gitlab for release tagging
cp .env.example .env # Fill out .env with private token and project id

source .env && scripts/build.sh upload

```

## Run
```bash
cd ~/go/src/
git@gitlab.com:PaperStreetHouse/laverna.git
cd laverna/
go build -o ~/go/bin/laverna main.go
~/go/bin/laverna
```

## Usage
When prompted
```
// Commands
get [url]

// Example
get http://xoxocomics.com/comic/miles-morales-ultimate-spider-man
```
