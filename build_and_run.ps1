# mit the symbol table, debug information and the DWARF symbol table by passing -s and -w go build -ldflags="-s -w" .
go build -ldflags="-s -w" .

./codeweaver -h

./codeweaver -ignore="\.git.*,.+\.exe" -excluded-paths-file="excluded_paths.txt"