go build .

git describe --tags --abbrev=0
# ./CodeWeaver -h
./CodeWeaver -clipboard -ignore="\.git.*,.+\.exe,codebase.md,excluded_paths.txt,DRAFT\.md,coverage$"