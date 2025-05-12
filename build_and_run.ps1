go build .

git describe --tags --abbrev=0

./CodeWeaver -clipboard -ignore="\.git.*,.+\.exe,codebase.md,excluded_paths.txt"