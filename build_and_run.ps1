go build .
git describe --tags --abbrev=0

$instruction = @"
**Before generating any code, state the features added and the features removed (hopefully none) by the code you generated - when in doubt, output this part and ask for confirmation before generating code!**

1. When implementing the requested changes, generate the complete modified files for easy copy and paste;
2. Change the code as little as possible;
3. Do not Introduce regressions or arbitrary simplifications: keep comments, checks, asserts, etc;
4. Generate professional and standard code
5. Do not add ephemerous comments, like `Changed`, `Fix Start`, `Removed`, etc. Always generate a final, professional version of the codebase;
6. Do not add the path at the top of the file.
"@


./CodeWeaver -clipboard -instruction $instruction -ignore="\.git.*,.+\.exe,codebase.md,excluded_paths.txt,DRAFT\.md,coverage.*,go\.sum"