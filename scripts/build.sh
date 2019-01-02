#!/bin/bash
gooses=( "linux" "windows" "darwin" )
goarches=( "amd64" "386" )

echo "Copy/Paste the below into the release notes to provide download links on gitlab"

for os in "${gooses[@]}"
do

    echo "# $os"

    for arch in "${goarches[@]}"
    do
        ext="bin"
        if [[ "$os" == "windows" ]]; then
            ext="exe"
        fi
        GOOS=$os GOARCH=$arch go build -gcflags "-trimpath=$PWD" -asmflags "-trimpath=$PWD" -o bin/laverna-$os-$arch.$ext

        if [[ "$1" == "upload" ]]; then
            curl -s --request "POST" --header "Private-Token: $TOKEN" --form "file=@bin/laverna-$os-$arch.$ext" "https://gitlab.com/api/v4/projects/$PROJECTID/uploads" | jq -r '.markdown'
            echo ""
        fi
    done;
done;
