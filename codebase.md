**Before generating any code, state the features added and the features removed (hopefully none) by the code you generated - when in doubt, output this part and ask for confirmation before generating code!**

1. When implementing the requested changes, generate the complete modified files for easy copy and paste;
2. Change the code as little as possible;
3. Do not Introduce regressions or arbitrary simplifications: keep comments, checks, asserts, etc;
4. Generate professional and standard code
5. Do not add ephemerous comments, like Changed, Fix Start, Removed, etc. Always generate a final, professional version of the codebase;
6. Do not add the path at the top of the file.

# Tree View:
```
.
├── build_and_run.ps1
├── build_and_show_help.ps1
├── go.mod
├── goreleaser.yaml
├── LICENSE
├── main.go
├── main_test.go
├── README.md
└── run_tests.ps1

```

# Content:

## LICENSE

```
MIT License

Copyright (c) 2024 Carlos Tarjano

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

```


## README.md

````md
# CodeWeaver: Generate Markdown Documentation from Your Codebase

[![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

CodeWeaver is a command-line tool that transforms your codebase into a single, navigable Markdown document. It recursively scans a directory, creating a tree-like representation of your project's file structure and embedding the content of each file within markdown code blocks.  This simplifies codebase sharing, documentation, and integration with AI/ML tools by providing a consolidated, readable Markdown output.

The output for the current repository can be found [here](https://github.com/tesserato/CodeWeaver/blob/main/codebase.md).

## Key Features

*   **Comprehensive Codebase Documentation:** Generates a Markdown file outlining your project's directory and file structure in a clear, tree-like format.
*   **Code Content Inclusion:** Embeds the *complete* content of each file within the Markdown document, using code blocks based on file extensions.
*   **Flexible Path Filtering:** Uses regular expressions to define `include` and / or `ignore` patterns, giving you precise control over which files are included.
*   **Optional Path Logging:**  Saves lists of included and excluded file paths to separate files for detailed tracking.
*   **Clipboard Integration:**  Optionally copies the generated Markdown to the clipboard for easy pasting.
*   **Simple CLI:** A straightforward command-line interface with intuitive options.

## Installation

**Using `go install` (Recommended):**

Requires Go 1.18 or later.

```bash
go install github.com/tesserato/CodeWeaver@latest
```

To install a specific version:

```bash
go install github.com/tesserato/CodeWeaver@vX.Y.Z  # Replace X.Y.Z with the desired version
```

**From Pre-built Executables:**

Download the appropriate executable for your operating system from the [releases page](https://github.com/tesserato/CodeWeaver/releases).

After downloading, make the executable:

```bash
chmod +x codeweaver  # On Linux/macOS
```

## Usage

```bash
codeweaver [options]
```

For help:

```bash
codeweaver -h
```

**Options:**

| Option                            | Description                                                                                                     | Default Value           |
| :-------------------------------- | :-------------------------------------------------------------------------------------------------------------- | :---------------------- |
| `-input <directory>`              | The root directory to scan.                                                                                     | `.` (current directory) |
| `-output <filename>`              | The name of the output Markdown file.                                                                           | `codebase.md`           |
| `-ignore "<regex patterns>"`      | Comma-separated list of regular expressions for paths to *exclude*.  Example: `\.git.*,node_modules,*.log`      | `\.git.*`               |
| `-include "<regex patterns>"`     | Comma-separated list of regular expressions. *Only* paths matching these are *included*. Example: `\.go$,\.md$` | None                    |
| `-included-paths-file <filename>` | Saves the list of *included* paths to this file.                                                                | None                    |
| `-excluded-paths-file <filename>` | Saves the list of *excluded* paths to this file.                                                                | None                    |
| `-clipboard`                      | Copies the generated Markdown to the clipboard.                                                                | `false`                 |
| `-version`                        | Displays the version and exits.                                                                                 |                         |
| `-help`                           | Displays this help message and exits.                                                                           |                         |

**Understanding `-include` and `-ignore`**

These flags control which files and directories are included in the generated documentation.

*   **`-ignore` (Blacklist):**  Excludes files/directories matching *any* of the provided regular expressions.
*   **`-include` (Whitelist):**  *Only* includes files/directories matching *at least one* of the provided regular expressions.  If `-include` is used, everything else is *excluded* by default.

**Behavior Table:**

| `-ignore` | `-include` | Behavior                                                                                                                                                                       |
| :-------- | :--------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| No        | No         | Includes all files/directories except the input directory itself (`.`).                                                                                                        |
| Yes       | No         | Excludes files/directories matching `-ignore`; includes everything else.                                                                                                       |
| No        | Yes        | *Only* includes files/directories matching `-include`. Everything else is excluded.                                                                                            |
| Yes       | Yes        | Includes files/directories that match *at least one* `-include` pattern AND do *not* match *any* `-ignore` pattern.  `-include` creates a whitelist, and `-ignore` filters it. |

## Examples

**1. Basic Usage:**

```bash
codeweaver
```

Creates `codebase.md` in the current directory, documenting the structure and content (excluding paths matching the default ignore pattern `\.git.*`).

**2. Different Input/Output:**

```bash
codeweaver -input=my_project -output=project_docs.md
```

Processes `my_project` and saves the output to `project_docs.md`.

**3. Ignoring Files/Directories:**

```bash
codeweaver -ignore="\.log,temp,build"
```

Excludes files/directories named `.log`, `temp`, or `build`.

**4. Including Only Specific Files:**

```bash
codeweaver -include="\.go$,\.md$"
```

Includes *only* Go (`.go`) and Markdown (`.md`) files.

**5. Combining `include` and `ignore`:**

```bash
codeweaver -include="\.go$,\.md$" -ignore="vendor,test"
```

Includes Go and Markdown files, *except* those with "vendor" or "test" in their paths.

**6. Saving Included/Excluded Paths:**

```bash
codeweaver -ignore="node_modules" -included-paths-file=included.txt -excluded-paths-file=excluded.txt
```

Creates `codebase.md`, saves included paths to `included.txt`, and excluded paths to `excluded.txt`.

**7. Copying to Clipboard:**

```bash
codeweaver -clipboard
```

Creates `codebase.md` and copies its content to the clipboard.

**8. Regex Examples:**

*   `.`: Matches any single character.
*   `*`: Matches zero or more of the preceding character.
*   `+`: Matches one or more of the preceding character.
*   `?`: Matches zero or one of the preceding character.
*   `[abc]`: Matches any one of the characters inside the brackets.
*   `[^abc]`: Matches any character *not* inside the brackets.
*   `[a-z]`: Matches any character in the range a-z.
*   `^`: Matches the beginning of the string.
*   `$`: Matches the end of the string.
*   `\.`: Matches a literal dot (.). You need to escape it because `.` has special meaning in regex.
*   `\|`: Used for alternation (OR).  e.g., `a\|b` matches either "a" or "b".
* `.*\.py[cod]$`: matches python files that end with pyc, pyd or pyo.
* `.*\.pdf`: matches PDF files.
* `(dir1\|dir2)`: matches `dir1` or `dir2`

**9. Complete example:**
```bash
codeweaver -input=. -output=codebase.md -ignore="\.git.*,.+\.exe,codebase.md,excluded_paths.txt" -include="\.go$,\.md$,\.ps1$,\.yaml$,\.txt$,\.csv$" -excluded-paths-file="excluded_paths.txt" -clipboard
```
This command will:

* Process the current directory (`.`).
* Generate documentation and save it in `codebase.md`.
* Exclude files matching `.git.*`, `.+\.exe`, the output file (`codebase.md`), and the file where the excluded paths will be saved.
* Include *only* files with the extensions .go, .md, .ps1, .yaml, .txt, and .csv.
* Save the list of excluded files in a file named `excluded_paths.txt`.
* Copy the generated Markdown to the system clipboard.

## Contributing

Contributions are welcome!  Please open an issue or submit a pull request on the project's GitHub repository.

## License

CodeWeaver is released under the [MIT License](LICENSE).

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=tesserato/CodeWeaver&type=Date)](https://star-history.com/#tesserato/CodeWeaver&Date)

## Alternatives

This section lists tools with similar or overlapping functionality.

**GitHub Repositories**

| Project                                                                                  | Stars                                                                                                                                                                        |
| :--------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [ai-context](https://github.com/tanq16/ai-context)                                       | [![GitHub stars](https://img.shields.io/github/stars/tanq16/ai-context?style=social)](https://github.com/tanq16/ai-context)                                                  |
| [bundle-codebases](https://github.com/manfrin/bundle-codebases)                          | [![GitHub stars](https://img.shields.io/github/stars/manfrin/bundle-codebases?style=social)](https://github.com/manfrin/bundle-codebases)                                    |
| [code2prompt](https://github.com/mufeedvh/code2prompt)                                   | [![GitHub stars](https://img.shields.io/github/stars/mufeedvh/code2prompt?style=social)](https://github.com/mufeedvh/code2prompt)                                            |
| [code2text](https://github.com/forrest321/code2text)                                     | [![GitHub stars](https://img.shields.io/github/stars/forrest321/code2text?style=social)](https://github.com/forrest321/code2text)                                            |
| [codefetch](https://github.com/regenrek/codefetch)                                       | [![GitHub stars](https://img.shields.io/github/stars/regenrek/codefetch?style=social)](https://github.com/regenrek/codefetch)                                                |
| [copcon](https://github.com/kasperjunge/copcon)                                          | [![GitHub stars](https://img.shields.io/github/stars/kasperjunge/copcon?style=social)](https://github.com/kasperjunge/copcon)                                                |
| [describe](https://github.com/rodlaf/describe)                                           | [![GitHub stars](https://img.shields.io/github/stars/rodlaf/describe?style=social)](https://github.com/rodlaf/describe)                                                      |
| [feed-llm](https://github.com/nahco314/feed-llm)                                         | [![GitHub stars](https://img.shields.io/github/stars/nahco314/feed-llm?style=social)](https://github.com/nahco314/feed-llm)                                                  |
| [files-to-prompt](https://github.com/simonw/files-to-prompt)                             | [![GitHub stars](https://img.shields.io/github/stars/simonw/files-to-prompt?style=social)](https://github.com/simonw/files-to-prompt)                                        |
| [ggrab](https://github.com/keizo/ggrab)                                                  | [![GitHub stars](https://img.shields.io/github/stars/keizo/ggrab?style=social)](https://github.com/keizo/ggrab)                                                              |
| [gitingest](https://gitingest.com/)                                                      | [![GitHub stars](https://img.shields.io/github/stars/cyclotruc/gitingest?style=social)](https://github.com/cyclotruc/gitingest)                                              |
| [gitpodcast](https://gitpodcast.com)                                                     | [![GitHub stars](https://img.shields.io/github/stars/BandarLabs/gitpodcast?style=social)](https://github.com/BandarLabs/gitpodcast)                                          |
| [globcat.sh](https://github.com/jzombie/globcat.sh)                                      | [![GitHub stars](https://img.shields.io/github/stars/jzombie/globcat.sh?style=social)](https://github.com/jzombie/globcat.sh)                                                |
| [grimoire](https://github.com/foresturquhart/grimoire)                                   | [![GitHub stars](https://img.shields.io/github/stars/foresturquhart/grimoire?style=social)](https://github.com/foresturquhart/grimoire)                                      |
| [llmcat](https://github.com/azer/llmcat)                                                 | [![GitHub stars](https://img.shields.io/github/stars/azer/llmcat?style=social)](https://github.com/azer/llmcat)                                                              |
| [RepoMix](https://github.com/yamadashy/repomix)                                          | [![GitHub stars](https://img.shields.io/github/stars/yamadashy/repomix?style=social)](https://github.com/yamadashy/repomix)                                                  |
| [techdocs](https://github.com/thesurlydev/techdocs)                                      | [![GitHub stars](https://img.shields.io/github/stars/thesurlydev/techdocs?style=social)](https://github.com/thesurlydev/techdocs)                                            |
| [thisismy](https://github.com/franzenzenhofer/thisismy)                                  | [![GitHub stars](https://img.shields.io/github/stars/franzenzenhofer/thisismy?style=social)](https://github.com/franzenzenhofer/thisismy)                                    |
| [yek](https://github.com/bodo-run/yek)                                                   | [![GitHub stars](https://img.shields.io/github/stars/bodo-run/yek?style=social)](https://github.com/bodo-run/yek)                                                            |
| [your-source-to-prompt](https://github.com/Dicklesworthstone/your-source-to-prompt.html) | [![GitHub stars](https://img.shields.io/github/stars/Dicklesworthstone/your-source-to-prompt.html?style=social)](https://github.com/Dicklesworthstone/your-source-to-prompt) |
| [ingest](https://github.com/sammcj/ingest)                                               | [![GitHub stars](https://img.shields.io/github/stars/sammcj/ingest?style=social)](https://github.com/sammcj/ingest)                                                          |
| [onefilellm](https://github.com/jimmc414/onefilellm)                                     | [![GitHub stars](https://img.shields.io/github/stars/jimmc414/onefilellm?style=social)](https://github.com/jimmc414/onefilellm)                                              |
| [repo2file](https://github.com/artkulak/repo2file)                                       | [![GitHub stars](https://img.shields.io/github/stars/artkulak/repo2file?style=social)](https://github.com/artkulak/repo2file)                                                |
| [clipsource](https://github.com/strizzo/clipsource)                                      | [![GitHub stars](https://img.shields.io/github/stars/strizzo/clipsource?style=social)](https://github.com/strizzo/clipsource)                                                |

**Other Tools**

*   **r2md:**  A Rust crate ([https://crates.io/crates/r2md](https://crates.io/crates/r2md)).
*   **repo2txt:** A web-based tool ([https://chathub.gg/repo2txt](https://chathub.gg/repo2txt) and [https://repo2txt.simplebasedomain.com/local.html](https://repo2txt.simplebasedomain.com/local.html)).
*  **repoprompt:** A web service ([https://www.repoprompt.com](https://www.repoprompt.com)).

**VSCode Extensions**

*   **Codebase to Markdown:** ([https://marketplace.visualstudio.com/items?itemName=DVYIO.combine-open-files](https://marketplace.visualstudio.com/items?itemName=DVYIO.combine-open-files))
````


## build_and_run.ps1

```ps1
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
```


## build_and_show_help.ps1

```ps1
go build .
./CodeWeaver -h

./CodeWeaver --help



```


## go.mod

```mod
module github.com/tesserato/CodeWeaver

go 1.23.0

require golang.design/x/clipboard v0.7.0

require (
	golang.org/x/exp v0.0.0-20190731235908-ec7cb31e5a56 // indirect
	golang.org/x/image v0.6.0 // indirect
	golang.org/x/mobile v0.0.0-20230301163155-e0f57694e12c // indirect
	golang.org/x/sys v0.5.0 // indirect
)

```


## goreleaser.yaml

```yaml
# vim: set ts=2 sw=2 tw=0 fo=cnqoj

version: 2

builds:
  - env:
      - CGO_ENABLED=0
    goos:
      - linux
      - windows
      - darwin
    ldflags:
      - -s -w
      - -X github.com/tesserato/CodeWeaver.version={{.Version}}
      - -X github.com/tesserato/CodeWeaver.commit={{.Commit}}
      - -X github.com/tesserato/CodeWeaver.date={{.Date}}

archives:
  - formats: [tar.gz]
    # this name template makes the OS and Arch compatible with the results of `uname`.
    name_template: >-
      {{ .ProjectName }}_
      {{- title .Os }}_
      {{- if eq .Arch "amd64" }}x86_64
      {{- else if eq .Arch "386" }}i386
      {{- else }}{{ .Arch }}{{ end }}
      {{- if .Arm }}v{{ .Arm }}{{ end }}
    # use zip for windows archives
    format_overrides:
      - goos: windows
        formats: [ 'zip' ]

changelog:
  sort: asc
  filters:
    exclude:
      - "^docs:"
      - "^test:"

release:
  footer: >-

    ---

    Released by [GoReleaser](https://github.com/goreleaser/goreleaser).
```


## main.go

````go
// Package main implements CodeWeaver, a command-line tool that transforms a codebase
// into a single, navigable Markdown document. It generates a tree view of the file
// structure and embeds the content of each file within markdown code blocks.
package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"golang.design/x/clipboard"
)

// --- Constants ---

// ANSI color codes for terminal output formatting.
const (
	colorRed       = "\033[31m"
	colorLiteRed   = "\033[91m"
	colorGreen     = "\033[32m"
	colorLiteGreen = "\033[92m"
	colorYellow    = "\033[33m"
	colorCyan      = "\033[36m"
	colorBold      = "\033[1m"
	colorReset     = "\033[0m"
)

// Markdown formatting constants.
const (
	mdTreeViewHeader  = "# Tree View:\n```\n"
	mdCodeBlockEnd    = "\n```\n"
	mdContentHeader   = "\n# Content:\n"
	mdFileHeaderStart = "\n## "
	// mdCodeBlockStart and mdCodeBlockEndNL are replaced by dynamic generation logic.
)

// Log prefixes for pattern matching output.
const (
	rgxLogPrefixIgnore  = "- RGX:"
	rgxLogPrefixInclude = "+ RGX:"
)

// Log prefixes for file processing status.
const (
	logPrefixExclude = "-"
	logPrefixInclude = "+"
)

// --- Main Execution ---

// main is the entry point of the application. It orchestrates the configuration parsing,
// logging setup, regex compilation, markdown generation, and output writing.
func main() {

	// 1. Parse Configuration
	cfg, err := parseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%sError: %v%s\n", colorRed, err, colorReset)
		if !isFlagHelpError(err) {
			fmt.Fprintln(os.Stderr, "Use -h for usage information.")
		}
		os.Exit(2)
	}

	// Check args excluding the program name for help flags explicitly
	for _, arg := range os.Args[1:] {
		arg_lower := strings.ToLower(arg)
		if arg_lower == "-help" || arg_lower == "--help" || arg_lower == "h" || arg_lower == "-h" || arg_lower == "help" {
			printHelp()
			os.Exit(0) // Exit successfully after showing help
		}
	}

	logger := setupLogging(cfg)
	ignoreMatchers, includeMatchers, err := compileMatchers(cfg, logger)
	if err != nil {
		logger.Fatalf("%sError compiling regex patterns: %v%s", colorRed, err, colorReset)
	}

	// generateMarkdown now orchestrates content and tree generation based on a single processing pass
	finalMarkdownString, pathsForIncludedFile, pathsForExcludedFile, err := generateMarkdown(cfg, ignoreMatchers, includeMatchers, logger)
	if err != nil {
		logger.Fatalf("%sError generating markdown: %v%s", colorRed, err, colorReset)
	}

	err = writeOutput(cfg, finalMarkdownString, pathsForIncludedFile, pathsForExcludedFile, logger)
	if err != nil {
		logger.Fatalf("%sError writing output: %v%s", colorRed, err, colorReset)
	}

	logger.Printf("%sCodeWeaver finished successfully.%s", colorGreen, colorReset)
}

// isFlagHelpError checks if the error returned by flag.Parse is due to a request for help.
func isFlagHelpError(err error) bool {
	return err != nil && err.Error() == "flag: help requested"
}

// config holds the runtime configuration options parsed from command-line arguments.
type config struct {
	inputDirOriginal  string
	inputDirAbs       string
	outputFile        string
	ignorePatterns    []string
	includePatterns   []string
	includedPathsFile string
	excludedPathsFile string
	instruction       string
	addToClipboard    bool
	showVersion       bool
}

// parseFlags defines and parses the command-line flags. It validates the input directory
// and returns a populated config struct or an error.
func parseFlags() (*config, error) {
	cfg := &config{}
	flag.StringVar(&cfg.inputDirOriginal, "input", ".", "The root directory to scan.")
	flag.StringVar(&cfg.outputFile, "output", "codebase.md", "The name of the output Markdown file.")
	ignoreStr := flag.String("ignore", `\.git.*`, "Comma-separated list of regular expressions for paths to *exclude*.")
	includeStr := flag.String("include", "", "Comma-separated list of regular expressions. *Only* paths matching these are *included*.")
	flag.StringVar(&cfg.includedPathsFile, "included-paths-file", "", "Saves the list of *included* paths to this file.")
	flag.StringVar(&cfg.excludedPathsFile, "excluded-paths-file", "", "Saves the list of *excluded* paths to this file.")
	flag.StringVar(&cfg.instruction, "instruction", "", "Optional text to prepend to the generated Markdown file.")
	flag.BoolVar(&cfg.addToClipboard, "clipboard", false, "Copies the generated Markdown to the clipboard.")
	flag.BoolVar(&cfg.showVersion, "version", false, "Displays the version and exits.")

	// Override default Usage to print custom help
	flag.Usage = func() { printHelp(); os.Exit(0) }
	err := flag.CommandLine.Parse(os.Args[1:])
	if err != nil {
		return nil, err
	}
	if cfg.showVersion {
		return cfg, nil
	}

	var validationErr error
	cfg.inputDirAbs, validationErr = filepath.Abs(cfg.inputDirOriginal)
	if validationErr != nil {
		return nil, fmt.Errorf("abs path for '%s': %w", cfg.inputDirOriginal, validationErr)
	}
	info, validationErr := os.Stat(cfg.inputDirAbs)
	if validationErr != nil {
		if os.IsNotExist(validationErr) {
			return nil, fmt.Errorf("input dir '%s' not exist", cfg.inputDirAbs)
		}
		return nil, fmt.Errorf("accessing input dir '%s': %w", cfg.inputDirAbs, validationErr)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("input path '%s' not a dir", cfg.inputDirAbs)
	}
	if *ignoreStr != "" {
		cfg.ignorePatterns = strings.Split(*ignoreStr, ",")
	}
	if *includeStr != "" {
		cfg.includePatterns = strings.Split(*includeStr, ",")
	}
	return cfg, nil
}

// setupLogging initializes the logger and prints the startup configuration.
func setupLogging(cfg *config) *log.Logger {
	logger := log.New(os.Stdout, "", 0)
	logger.Println("Starting CodeWeaver...")
	// Use filepath.ToSlash to ensure consistent path separators in logs
	logger.Println("Input directory:", filepath.ToSlash(cfg.inputDirAbs))
	logger.Println("Output file:", filepath.ToSlash(cfg.outputFile))
	if cfg.includedPathsFile != "" {
		logger.Println("Included paths will be saved to:", filepath.ToSlash(cfg.includedPathsFile))
	}
	if cfg.excludedPathsFile != "" {
		logger.Println("Excluded paths will be saved to:", filepath.ToSlash(cfg.excludedPathsFile))
	}
	if cfg.instruction != "" {
		logger.Println("Instruction text provided.")
	}
	if cfg.addToClipboard {
		logger.Println("Result will be copied to clipboard.")
	}
	logger.Println()
	return logger
}

// compileMatchers compiles the provided string patterns into regular expressions.
func compileMatchers(cfg *config, logger *log.Logger) (ignore, include []*regexp.Regexp, err error) {
	logger.Println("Compiling ignore patterns:")
	ignore, err = compileRegexList(cfg.ignorePatterns, colorLiteRed, rgxLogPrefixIgnore, logger)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid ignore pattern: %w", err)
	}
	logger.Println("Compiling include patterns:")
	include, err = compileRegexList(cfg.includePatterns, colorLiteGreen, rgxLogPrefixInclude, logger)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid include pattern: %w", err)
	}
	logger.Println()
	return
}

// compileRegexList is a helper that compiles a slice of regex strings, logging each one.
func compileRegexList(patterns []string, color, prefix string, logger *log.Logger) ([]*regexp.Regexp, error) {
	if len(patterns) == 0 {
		logger.Printf("  (No patterns provided)")
		return nil, nil
	}
	compiled := make([]*regexp.Regexp, 0, len(patterns))
	validFound := false
	for _, p := range patterns {
		trimmed := strings.TrimSpace(p)
		if trimmed == "" {
			continue
		}
		logger.Printf("%s  %s %s%s\n", color, prefix, trimmed, colorReset)
		rgx, rgxErr := regexp.Compile(trimmed)
		if rgxErr != nil {
			return nil, fmt.Errorf("pattern '%s': %w", trimmed, rgxErr)
		}
		compiled = append(compiled, rgx)
		validFound = true
	}
	if !validFound {
		logger.Printf("  (No valid patterns found)")
		return nil, nil
	}
	return compiled, nil
}

// --- Markdown Generation ---

// generateMarkdown orchestrates the creation of the tree view and content sections.
// It adopts a two-pass approach:
//  1. Content Generation: Scans files, applies filters, and builds the content markdown.
//     This step identifies exactly which files are "included".
//  2. Tree Generation: Builds the directory tree using ONLY the files identified in step 1.
//     This ensures that directories which become empty due to filtering are not shown.
func generateMarkdown(cfg *config, ignoreMatchers, includeMatchers []*regexp.Regexp, logger *log.Logger) (
	finalMarkdown string, processedPathsForFile []string, excludedPathsForFile []string, err error) {

	var fullMarkdown strings.Builder

	// --- Build Content Section FIRST to get processedPaths ---
	logger.Println("Processing paths and building content section...")
	contentBuilder := newContentBuilder(cfg.inputDirAbs, cfg.includedPathsFile, cfg.excludedPathsFile, ignoreMatchers, includeMatchers, logger)
	// contentMarkdown is the markdown string of file contents
	// processedPaths contains ALL files that passed shouldProcess (directories are stripped for cleaner tree)
	// excludedPaths contains all files/dirs that failed shouldProcess
	contentMarkdown, processedPaths, excludedPaths, err := contentBuilder.buildContentString()
	if err != nil {
		return "", nil, nil, fmt.Errorf("failed to build code content: %w", err)
	}

	// --- Prepend Instruction if present ---
	if cfg.instruction != "" {
		fullMarkdown.WriteString(cfg.instruction)
		fullMarkdown.WriteString("\n\n")
	}

	// --- Build Tree View using processedPaths ---
	logger.Println("Building tree view...")
	fullMarkdown.WriteString(mdTreeViewHeader)
	fullMarkdown.WriteString(filepath.ToSlash(cfg.inputDirOriginal) + "\n")

	processedPathsSet := make(map[string]struct{}, len(processedPaths))
	for _, p := range processedPaths {
		processedPathsSet[p] = struct{}{}
	}

	treeBuilder := newTreeBuilder(cfg.inputDirAbs, processedPathsSet) // Modified constructor
	treeString, err := treeBuilder.buildTreeString()
	if err != nil {
		return "", nil, nil, fmt.Errorf("failed to build codebase tree: %w", err)
	}
	fullMarkdown.WriteString(treeString)
	fullMarkdown.WriteString(mdCodeBlockEnd)

	// Append content string AFTER tree
	fullMarkdown.WriteString(mdContentHeader)
	fullMarkdown.WriteString(contentMarkdown) // The content part from contentBuilder
	logger.Println()

	// processedPaths will be used for the -included-paths-file
	return fullMarkdown.String(), processedPaths, excludedPaths, nil
}

// --- Tree Builder ---

// treeBuilder is responsible for generating the visual directory tree structure.
type treeBuilder struct {
	rootAbsPath       string
	processedPathsSet map[string]struct{} // Set of paths that passed filters (from contentBuilder)
	output            strings.Builder
	depthOpen         map[int]bool
}

// newTreeBuilder creates a new tree builder instance, taking the set of processed paths to filter the tree.
func newTreeBuilder(rootAbsPath string, processedPathsSet map[string]struct{}) *treeBuilder {
	return &treeBuilder{
		rootAbsPath:       rootAbsPath,
		processedPathsSet: processedPathsSet,
		depthOpen:         make(map[int]bool),
	}
}

// buildTreeString initiates the recursive tree building process and returns the result string.
func (tb *treeBuilder) buildTreeString() (string, error) {
	err := tb.printTreeRecursive(tb.rootAbsPath, 0)
	return tb.output.String(), err
}

// printTreeRecursive traverses the directory structure. It only includes entries that
// are present in processedPathsSet or are directories containing such entries.
func (tb *treeBuilder) printTreeRecursive(currentDirPath string, depth int) error {
	entries, err := os.ReadDir(currentDirPath)
	if err != nil {
		return fmt.Errorf("read directory %s: %w", currentDirPath, err)
	}

	var displayableEntries []fs.DirEntry
	for _, entry := range entries {
		fullEntryPath := filepath.Join(currentDirPath, entry.Name())
		pathRelToInput, relErr := tb.getRelativePath(fullEntryPath)
		if relErr != nil {
			return relErr
		}

		// Determine if this entry should be part of the tree:
		// 1. Itself is in processedPathsSet OR
		// 2. It's a directory and any processedPath is a descendant of it.
		shouldDisplay := false
		if _, isProcessed := tb.processedPathsSet[pathRelToInput]; isProcessed {
			shouldDisplay = true
		} else if entry.IsDir() {
			dirPathPrefix := pathRelToInput + "/"
			for processedPath := range tb.processedPathsSet {
				if strings.HasPrefix(processedPath, dirPathPrefix) {
					shouldDisplay = true
					break
				}
			}
		}

		if shouldDisplay {
			displayableEntries = append(displayableEntries, entry)
		}
	}

	sort.Slice(displayableEntries, func(i, j int) bool {
		return strings.ToLower(displayableEntries[i].Name()) < strings.ToLower(displayableEntries[j].Name())
	})

	for i, entry := range displayableEntries {
		isLastEntry := (i == len(displayableEntries)-1)
		tb.printEntryLine(entry, depth, isLastEntry)
		if entry.IsDir() {
			tb.depthOpen[depth] = !isLastEntry
			if err := tb.printTreeRecursive(filepath.Join(currentDirPath, entry.Name()), depth+1); err != nil {
				return err
			}
		}
	}
	return nil
}

// printEntryLine formats and appends a single line of the tree view (e.g., ├── filename).
func (tb *treeBuilder) printEntryLine(entry fs.DirEntry, depth int, isLast bool) {
	var prefix strings.Builder
	for i := 0; i < depth; i++ {
		if tb.depthOpen[i] {
			prefix.WriteString("│   ")
		} else {
			prefix.WriteString("    ")
		}
	}
	if isLast {
		prefix.WriteString("└── ")
	} else {
		prefix.WriteString("├── ")
	}
	tb.output.WriteString(prefix.String() + entry.Name() + "\n")
}

// getRelativePath calculates the path relative to the root input directory and ensures forward slashes.
func (tb *treeBuilder) getRelativePath(fullPath string) (string, error) {
	pathRelToInput, err := filepath.Rel(tb.rootAbsPath, fullPath)
	if err != nil {
		return "", fmt.Errorf("rel path for %s to %s: %w", fullPath, tb.rootAbsPath, err)
	}
	return filepath.ToSlash(pathRelToInput), nil
}

// --- Content Builder ---

// contentBuilder handles scanning the directory, filtering files, reading content,
// and formatting it into Markdown.
type contentBuilder struct {
	rootAbsPath, includedPathsFile, excludedPathsFile string
	ignoreMatchers, includeMatchers                   []*regexp.Regexp
	logger                                            *log.Logger
}

// newContentBuilder initializes a new contentBuilder.
func newContentBuilder(rootAbsPath, includedPathsFile, excludedPathsFile string, ignoreMatchers, includeMatchers []*regexp.Regexp, logger *log.Logger) *contentBuilder {
	return &contentBuilder{rootAbsPath, includedPathsFile, excludedPathsFile, ignoreMatchers, includeMatchers, logger}
}

// buildContentString scans the directory and returns:
// 1. markdownContent: The combined markdown string for file contents.
// 2. allProcessedPaths: List of files that passed filters (used for tree building).
// 3. excludedPaths: List of files/dirs that failed filters.
// 4. err: Any error encountered.
func (cb *contentBuilder) buildContentString() (
	markdownContent string, allProcessedPaths []string, excludedPaths []string, err error) {

	var contentSB strings.Builder // For actual file content markdown

	// Slices to store paths
	var localProcessedPaths []string
	var localExcludedPaths []string

	// Maps to store extensions
	includedExtensions := make(map[string]struct{})
	excludedExtensions := make(map[string]struct{})

	walkErr := filepath.WalkDir(cb.rootAbsPath, func(currentWalkPath string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			cb.logger.Printf("%sWarning: Error accessing %s: %v%s\n", colorRed, currentWalkPath, walkErr, colorReset)
			if d != nil && d.IsDir() && errors.Is(walkErr, fs.ErrPermission) {
				return fs.SkipDir
			}
			return walkErr
		}
		pathRelToInput, relErr := filepath.Rel(cb.rootAbsPath, currentWalkPath)
		if relErr != nil {
			cb.logger.Printf("%sWarning: Could not make path %s relative to %s: %v%s\n", colorRed, currentWalkPath, cb.rootAbsPath, relErr, colorReset)
			return nil
		}
		pathRelToInput = filepath.ToSlash(pathRelToInput)
		if pathRelToInput == "." {
			return nil
		} // Skip root itself from lists

		if !shouldProcess(pathRelToInput, cb.ignoreMatchers, cb.includeMatchers) {
			// Log exclusion if not saving to file (to avoid verbose output)
			if cb.excludedPathsFile == "" {
				cb.logger.Printf("%s%s %s%s\n", colorRed, logPrefixExclude, pathRelToInput, colorReset)
			}
			localExcludedPaths = append(localExcludedPaths, pathRelToInput)
			// Track extension for excluded file
			if !d.IsDir() {
				ext := strings.ToLower(filepath.Ext(pathRelToInput))
				if ext == "" {
					ext = "(no ext)"
				}
				excludedExtensions[ext] = struct{}{}
			}
			return nil // Path excluded, continue walk
		}

		// Log inclusion if not saving to file.
		// Directories are logged without color to distinguish from files.
		if cb.includedPathsFile == "" {
			if d.IsDir() {
				cb.logger.Printf("%s %s\n", logPrefixInclude, pathRelToInput)
			} else {
				cb.logger.Printf("%s%s %s%s\n", colorGreen, logPrefixInclude, pathRelToInput, colorReset)
			}
		}

		// If it's a directory, don't add to processedPaths.
		// This effectively removes "empty folders" (or folders with no included files) from the Tree View.
		if d.IsDir() {
			return nil // Continue into directory
		}

		// Path passed filters and is a file, add to processedPaths
		localProcessedPaths = append(localProcessedPaths, pathRelToInput)
		// Track extension for included file
		ext := strings.ToLower(filepath.Ext(pathRelToInput))
		if ext == "" {
			ext = "(no ext)"
		}
		includedExtensions[ext] = struct{}{}

		// Process file for content inclusion
		fileContent, readErr := os.ReadFile(currentWalkPath)
		if readErr != nil {
			cb.logger.Printf("%sWarning: Failed to read file %s: %v%s\n", colorRed, currentWalkPath, readErr, colorReset)
			// For error cases, we'll default to standard fencing (3 backticks)
			contentSB.WriteString(fmt.Sprintf("%s%s\n```\nError reading file: %v\n```\n\n", mdFileHeaderStart, pathRelToInput, readErr))
			return nil // File processed (passed filters), but content not added
		}

		if len(fileContent) == 0 {
			// Empty file: it's in localProcessedPaths, but no markdown content for it.
			return nil
		}

		// Check for binary content
		if isBinary(fileContent) {
			cb.logger.Printf("Skipping binary file content: %s\n", pathRelToInput)
			return nil
		}

		// Add file content to markdown with Dynamic Fencing
		extension := strings.TrimPrefix(strings.ToLower(filepath.Ext(currentWalkPath)), ".")
		maxBackticks := countMaxBackticks(fileContent)
		fenceLen := 3
		if maxBackticks >= 3 {
			fenceLen = maxBackticks + 1
		}
		fence := strings.Repeat("`", fenceLen)

		contentSB.WriteString(fmt.Sprintf("%s%s\n", mdFileHeaderStart, pathRelToInput))
		contentSB.WriteString(fmt.Sprintf("\n%s%s\n", fence, extension))
		contentSB.Write(fileContent)
		contentSB.WriteString(fmt.Sprintf("\n%s\n\n", fence))

		return nil
	})

	// Print Extension Summary
	if walkErr == nil {
		cb.logger.Println()
		printExtensionSummary(includedExtensions, colorGreen, "Included extensions:", cb.logger)
		printExtensionSummary(excludedExtensions, colorRed, "Excluded extensions:", cb.logger)
	}

	if walkErr != nil {
		return "", nil, nil, fmt.Errorf("walking directory %s: %w", cb.rootAbsPath, walkErr)
	}
	return contentSB.String(), localProcessedPaths, localExcludedPaths, nil
}

// printExtensionSummary logs a sorted list of file extensions found during processing.
func printExtensionSummary(extMap map[string]struct{}, color, label string, logger *log.Logger) {
	if len(extMap) == 0 {
		return
	}
	exts := make([]string, 0, len(extMap))
	for ext := range extMap {
		exts = append(exts, ext)
	}
	sort.Strings(exts)
	logger.Printf("%s%s %s%s\n", color, label, strings.Join(exts, ", "), colorReset)
}

// isBinary checks if the content contains a null byte in the first 1024 bytes.
// This is a standard heuristic to detect binary files.
func isBinary(content []byte) bool {
	const maxBytesToCheck = 1024
	checkLen := len(content)
	if checkLen > maxBytesToCheck {
		checkLen = maxBytesToCheck
	}
	return bytes.IndexByte(content[:checkLen], 0) != -1
}

// countMaxBackticks calculates the maximum number of consecutive backticks in the byte slice.
// This is used to determine the length of the markdown code block fence.
func countMaxBackticks(content []byte) int {
	maxCount := 0
	currentCount := 0
	for _, b := range content {
		if b == '`' {
			currentCount++
		} else {
			if currentCount > maxCount {
				maxCount = currentCount
			}
			currentCount = 0
		}
	}
	if currentCount > maxCount {
		maxCount = currentCount
	}
	return maxCount
}

// shouldProcess determines whether a file path should be processed based on the
// provided ignore and include regex matchers.
func shouldProcess(pathRelToInput string, ignoreMatchers, includeMatchers []*regexp.Regexp) bool {
	for _, pattern := range ignoreMatchers {
		if pattern != nil && pattern.MatchString(pathRelToInput) {
			return false
		}
	}
	if len(includeMatchers) > 0 {
		matchedInclude := false
		for _, pattern := range includeMatchers {
			if pattern != nil && pattern.MatchString(pathRelToInput) {
				matchedInclude = true
				break
			}
		}
		if !matchedInclude {
			return false
		}
	}
	return true
}

// writeOutput writes the generated markdown to the output file, saves included/excluded path lists,
// and optionally copies the result to the clipboard.
func writeOutput(cfg *config, markdownContent string, includedPaths, excludedPaths []string, logger *log.Logger) error {
	outputFileSlash := filepath.ToSlash(cfg.outputFile)
	logger.Printf("Writing output to %s...", outputFileSlash)
	err := os.WriteFile(cfg.outputFile, []byte(markdownContent), 0644)
	if err != nil {
		logger.Printf("%sError writing to output file %s: %v%s", colorRed, outputFileSlash, err, colorReset)
		return fmt.Errorf("writing output file %s: %w", cfg.outputFile, err)
	}
	logger.Printf("Markdown content written to %s", outputFileSlash)

	if cfg.includedPathsFile != "" {
		includedFileSlash := filepath.ToSlash(cfg.includedPathsFile)
		if err := savePathsToFile(cfg.includedPathsFile, includedPaths, logger); err != nil { // Pass `includedPaths` from generateMarkdown
			logger.Printf("%sWarning: Error saving included paths to %s: %v%s", colorRed, includedFileSlash, err, colorReset)
		}
	}
	if cfg.excludedPathsFile != "" {
		excludedFileSlash := filepath.ToSlash(cfg.excludedPathsFile)
		if err := savePathsToFile(cfg.excludedPathsFile, excludedPaths, logger); err != nil { // Pass `excludedPaths` from generateMarkdown
			logger.Printf("%sWarning: Error saving excluded paths to %s: %v%s", colorRed, excludedFileSlash, err, colorReset)
		}
	}
	if cfg.addToClipboard {
		logger.Println("Attempting to copy to clipboard...")
		if err := clipboard.Init(); err != nil {
			logger.Printf("%sWarning: Could not initialize clipboard: %v%s", colorRed, err, colorReset)
		} else {
			clipboard.Write(clipboard.FmtText, []byte(markdownContent))
			logger.Println("Markdown content copied to clipboard.")
		}
	}
	return nil
}

// savePathsToFile writes a list of file paths to the specified file, sorted alphabetically.
func savePathsToFile(filename string, paths []string, logger *log.Logger) error {
	if len(paths) == 0 {
		logger.Printf("No paths to save to %s.", filepath.ToSlash(filename))
		return nil
	}
	sort.Strings(paths) // Sort for consistent output
	var sb strings.Builder
	for _, p := range paths {
		sb.WriteString(p)
		sb.WriteString("\n")
	}
	err := os.WriteFile(filename, []byte(sb.String()), 0644)
	if err == nil {
		logger.Printf("Paths saved to %s", filepath.ToSlash(filename))
	} else {
		return fmt.Errorf("writing paths file %s: %w", filepath.ToSlash(filename), err)
	}
	return nil
}

// printHelp displays the application's usage information, available flags, and examples.
func printHelp() {
	// Header
	fmt.Fprintf(os.Stderr, "%s%sCodeWeaver%s: Generate Markdown Documentation from Your Codebase.\n\n", colorBold, colorGreen, colorReset)

	// Usage
	fmt.Fprintf(os.Stderr, "%sUsage:%s\n", colorCyan, colorReset)
	fmt.Fprintf(os.Stderr, "  codeweaver [flags]\n\n")

	// Flags
	fmt.Fprintf(os.Stderr, "%sFlags:%s\n", colorCyan, colorReset)
	flag.VisitAll(func(f *flag.Flag) {
		// Format: -flag
		//         Description (Default: value)
		// Flag name in Green, Description in default (white/reset)
		fmt.Fprintf(os.Stderr, "  %s-%-20s%s\n", colorGreen, f.Name, colorReset)
		fmt.Fprintf(os.Stderr, "      %s", f.Usage)

		// Print default value if not empty
		if f.DefValue != "" {
			// Don't print defaults for boolean flags that are false (cleaner output)
			if f.Name == "clipboard" && f.DefValue == "false" {
				// skip
			} else if f.Name == "version" && f.DefValue == "false" {
				// skip
			} else {
				fmt.Fprintf(os.Stderr, " %s(Default: %s)%s", colorYellow, f.DefValue, colorReset)
			}
		}
		fmt.Fprintf(os.Stderr, "\n")
	})
	fmt.Fprintln(os.Stderr)

	// Examples
	fmt.Fprintf(os.Stderr, "%sExamples:%s\n", colorCyan, colorReset)
	fmt.Fprintf(os.Stderr, "  %scodeweaver%s\n", colorGreen, colorReset)
	fmt.Fprintf(os.Stderr, "    Process current directory, output to codebase.md\n\n")

	fmt.Fprintf(os.Stderr, "  %scodeweaver -input src -output source_docs.md%s\n", colorGreen, colorReset)
	fmt.Fprintf(os.Stderr, "    Specify input directory and output filename\n\n")

	fmt.Fprintf(os.Stderr, "  %scodeweaver -ignore \"build/,vendor/\" -include \"\\.go$,\\.md$\"%s\n", colorGreen, colorReset)
	fmt.Fprintf(os.Stderr, "    Exclude 'build' and 'vendor' folders, but ONLY include .go and .md files\n\n")

	fmt.Fprintf(os.Stderr, "  %scodeweaver -instruction \"Analyze this code\" -clipboard%s\n", colorGreen, colorReset)
	fmt.Fprintf(os.Stderr, "    Prepend instruction and copy result to clipboard\n\n")

	// Filter Logic
	fmt.Fprintf(os.Stderr, "%sHow Filters Work:%s\n", colorCyan, colorReset)
	fmt.Fprintf(os.Stderr, "  1. %s-ignore%s (Blacklist): Matches are excluded. Checked first.\n", colorLiteRed, colorReset)
	fmt.Fprintf(os.Stderr, "  2. %s-include%s (Whitelist): If specified, ONLY matches are included.\n", colorGreen, colorReset)
	fmt.Fprintf(os.Stderr, "     If a path matches both (unlikely given logic), ignore takes precedence.\n")
	fmt.Fprintf(os.Stderr, "     Directories matching -ignore are skipped entirely.\n\n")

	// Regex Notes
	fmt.Fprintf(os.Stderr, "%sRegex Notes:%s\n", colorCyan, colorReset)
	fmt.Fprintf(os.Stderr, "  - Patterns are Go regular expressions.\n")
	fmt.Fprintf(os.Stderr, "  - Use forward slashes '/' for paths (e.g., \"dir/file.txt\").\n")
	fmt.Fprintf(os.Stderr, "  - Match a file extension: \"\\.go$\"\n")
	fmt.Fprintf(os.Stderr, "  - Match a directory: \"^vendor/\" or \"/vendor/\"\n")
}

````


## main_test.go

`````go
package main

import (
	"bufio"
	"bytes" // For capturing command output
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec" // Needed to run commands
	"path/filepath"
	"regexp"
	"runtime" // Needed for OS-specific details
	"sort"
	"strings"
	"testing"
)

// --- Test Helpers (keep previous helpers like mustCompileRegex, normalizeNewlines, etc.) ---

// mustCompileRegex compiles a regex string and panics on error.
// Useful for initializing static regexes in test cases.
func mustCompileRegex(pattern string) *regexp.Regexp {
	if pattern == "" {
		return nil
	}
	r, err := regexp.Compile(pattern)
	if err != nil {
		log.Fatalf("Test setup: Failed to compile regex '%s': %v", pattern, err)
	}
	return r
}

// normalizeNewlines replaces carriage returns with newlines to ensure consistent comparisons across OS.
func normalizeNewlines(s string) string { return strings.ReplaceAll(s, "\r\n", "\n") }

// createTestFS creates a temporary file system structure for testing.
// It returns the root directory path and a cleanup function.
func createTestFS(t *testing.T) (string, func()) {
	t.Helper()
	rootDir := t.TempDir()
	structure := map[string]string{
		"file1.txt": "content of file1",
		"script.go": "package main\nfunc main() {}",
		"README.md": "# Test Readme",
		"data/":     "", "data/image.png": "fake png data", "data/config.yaml": "key: value",
		"build/": "", "build/output.exe": "binary\x00data", // Binary file with null byte
		"build/tmp/": "", "build/tmp/log.txt": "log entry",
		".git/": "", ".git/HEAD": "ref: refs/heads/main",
		"node_modules/": "", "node_modules/dep/": "", "node_modules/dep/package.json": "{}",
		"empty_dir/":                    "",
		"docs/sub_docs/file_in_sub.txt": "nested doc content",
		"other.log":                     "another log",
		"empty_file.txt":                "",                                                             // Explicitly empty file
		"doc_with_ticks.md":             "Here is a code block:\n```go\nfmt.Println(\"Hi\")\n```\nEnd.", // File with 3 backticks
	}
	for relPath, content := range structure {
		absPath := filepath.Join(rootDir, relPath)
		parentDir := filepath.Dir(absPath)
		if err := os.MkdirAll(parentDir, 0755); err != nil {
			t.Fatalf("Failed to create parent dir %s: %v", parentDir, err)
		}
		isDirectory := strings.HasSuffix(relPath, "/") || (content == "" && !strings.Contains(filepath.Base(relPath), ".") && relPath != "empty_file.txt")
		if isDirectory {
			if err := os.MkdirAll(absPath, 0755); err != nil && !errors.Is(err, os.ErrExist) {
				t.Fatalf("Failed to create dir %s: %v", absPath, err)
			}
		} else {
			if err := os.WriteFile(absPath, []byte(content), 0644); err != nil {
				t.Fatalf("Failed to write file %s: %v", absPath, err)
			}
		}
	}
	return rootDir, func() {}
}

// equalStringSlices checks if two string slices contain the same elements (order independent).
func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	sort.Strings(a) // Ensure sorted for comparison
	sort.Strings(b)
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// runMainLogic simulates the execution of the main program for integration testing.
// It sets up flags, captures logs, and avoids exiting the process on error.
func runMainLogic(args []string, outputDir string) (string, error) {
	originalArgs := os.Args
	os.Args = append([]string{"codeweaver"}, args...)
	defer func() { os.Args = originalArgs }()

	var logBuf bytes.Buffer
	originalLoggerOutput := log.Writer()
	testRunLogger := log.New(&logBuf, "", 0)
	log.SetOutput(&logBuf)
	log.SetFlags(0)
	defer func() { log.SetOutput(originalLoggerOutput); log.SetFlags(log.LstdFlags) }()

	// Use the actual parseFlags from main.go
	// We need to temporarily set flag.Usage to avoid os.Exit(0) if -h is present in args
	originalUsage := flag.Usage
	// var parseErr error
	// flag.Usage = func() { parseErr = flag.ErrHelp /* Mark that help was requested */ }
	defer func() { flag.Usage = originalUsage }() // Restore original usage

	// Reset CommandLine flags before each run to avoid pollution between test cases
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError) // Use ContinueOnError for tests
	flag.CommandLine.SetOutput(io.Discard)                               // Suppress default error output from flag parsing

	cfg, err := parseFlags() // Call the actual parseFlags
	if err != nil {
		if errors.Is(err, flag.ErrHelp) { // If flag.Parse() itself returns ErrHelp for -h
			printHelp() // Manually call printHelp to get its output in logBuf for assertion
		}
		// For other parse errors or if parseFlags helper returned ErrHelp
		return logBuf.String(), err
	}

	// Handle version/help flags first (they stop execution)
	// if cfg.showVersion {
	// 	testRunLogger.Printf("CodeWeaver version %s\ncommit %s\nbuilt at %s\n", version, commit, date)
	// 	return logBuf.String(), nil
	// }
	// Note: -help is now handled before parseFlags in main.go, -h by flag.Usage
	// The test for `runMainLogic([]string{"-help"}, ...)` in TestMainExecutionFlows will check this.

	// Perform post-parsing steps (validation) from original parseFlags
	// (This block is now mostly handled by parseFlags itself, but we keep inputDirAbs for logging)
	if cfg.inputDirAbs == "" { // If parseFlags didn't set it (e.g., due to version/help early exit)
		absPath, pathErr := filepath.Abs(cfg.inputDirOriginal)
		if pathErr != nil {
			return logBuf.String(), fmt.Errorf("error getting abs path in test: %w", pathErr)
		}
		cfg.inputDirAbs = absPath
	}

	testRunLogger.Println("Starting CodeWeaver...")
	testRunLogger.Println("Input directory:", filepath.ToSlash(cfg.inputDirAbs))
	testRunLogger.Println("Output file:", filepath.ToSlash(filepath.Join(outputDir, cfg.outputFile)))
	if cfg.includedPathsFile != "" {
		testRunLogger.Println("Included paths will be saved to:", filepath.ToSlash(filepath.Join(outputDir, cfg.includedPathsFile)))
	}
	if cfg.excludedPathsFile != "" {
		testRunLogger.Println("Excluded paths will be saved to:", filepath.ToSlash(filepath.Join(outputDir, cfg.excludedPathsFile)))
	}
	if cfg.instruction != "" {
		testRunLogger.Println("Instruction text provided.")
	}
	if cfg.addToClipboard {
		testRunLogger.Println("Result will be copied to clipboard.")
	}
	testRunLogger.Println()

	ignoreMatchers, includeMatchers, err := compileMatchers(cfg, testRunLogger)
	if err != nil {
		err = fmt.Errorf("Error compiling regex patterns: %w", err)
		testRunLogger.Println(colorRed + err.Error() + colorReset)
		return logBuf.String(), err
	}

	// Call the refactored generateMarkdown
	finalMarkdownString, pathsForIncludedFile, pathsForExcludedFile, err := generateMarkdown(cfg, ignoreMatchers, includeMatchers, testRunLogger)
	if err != nil {
		err = fmt.Errorf("Error generating markdown: %w", err)
		testRunLogger.Println(colorRed + err.Error() + colorReset)
		return logBuf.String(), err
	}

	fullOutputPath := filepath.Join(outputDir, cfg.outputFile)
	fullIncludedPath := ""
	if cfg.includedPathsFile != "" {
		fullIncludedPath = filepath.Join(outputDir, cfg.includedPathsFile)
	}
	fullExcludedPath := ""
	if cfg.excludedPathsFile != "" {
		fullExcludedPath = filepath.Join(outputDir, cfg.excludedPathsFile)
	}

	writeCfg := *cfg
	writeCfg.outputFile = fullOutputPath
	writeCfg.includedPathsFile = fullIncludedPath
	writeCfg.excludedPathsFile = fullExcludedPath

	originalClipboardState := writeCfg.addToClipboard
	if writeCfg.addToClipboard {
		writeCfg.addToClipboard = false
	}

	err = writeOutput(&writeCfg, finalMarkdownString, pathsForIncludedFile, pathsForExcludedFile, testRunLogger)
	if err != nil {
		return logBuf.String(), err
	}

	if originalClipboardState {
		testRunLogger.Println("Markdown content copied to clipboard (simulated).")
	}

	return logBuf.String(), nil
}

// --- Test Suite ---

// TestShouldProcess validates the file filtering logic based on ignore and include patterns.
func TestShouldProcess(t *testing.T) {
	// This test remains the same as it tests the standalone filtering logic.
	testCases := []struct {
		name            string
		path            string
		ignoreMatchers  []*regexp.Regexp
		includeMatchers []*regexp.Regexp
		expected        bool
	}{
		{"NoFilters_Allow", "file.txt", nil, nil, true},
		{"IgnoreMatch_Exact", "skip.txt", []*regexp.Regexp{mustCompileRegex(`^skip\.txt$`)}, nil, false},
		{"IncludeMatch_Extension", "src/main.go", nil, []*regexp.Regexp{mustCompileRegex(`\.go$`)}, true},
		{"IgnoreTakesPrecedence", "vendor/lib.go", []*regexp.Regexp{mustCompileRegex(`^vendor/`)}, []*regexp.Regexp{mustCompileRegex(`\.go$`)}, false},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := shouldProcess(tc.path, tc.ignoreMatchers, tc.includeMatchers)
			if actual != tc.expected {
				t.Errorf("shouldProcess(%q) = %v; want %v", tc.path, actual, tc.expected)
			}
		})
	}
}

// TestParseFlags validates command-line flag parsing, including defaults and error conditions.
func TestParseFlags(t *testing.T) {
	runParse := func(t *testing.T, args []string) (string, error) {
		t.Helper()
		testOutputDir := t.TempDir() // Need output dir for runMainLogic

		// Determine if a default valid input dir is needed
		needsDefaultInput := true
		for _, arg := range args {
			if arg == "-input" {
				needsDefaultInput = false // Input explicitly provided
				break
			}
			// Assume tests providing non-existent or file paths *intend* to test that behavior
			if strings.Contains(arg, "non_existent") || strings.Contains(arg, "test_file_") {
				needsDefaultInput = false
				break
			}
		}

		finalArgs := args
		var cleanupInput func() = func() {} // No-op cleanup initially
		if needsDefaultInput {
			// Provide a default valid input directory if -input isn't specified
			// and it doesn't look like an input error test.
			validTempInputDir, cleanup := createTestFS(t)
			finalArgs = append([]string{"-input", validTempInputDir}, args...)
			cleanupInput = cleanup // Assign real cleanup if dir was created
		}
		defer cleanupInput() // Defer cleanup for the default input dir if created

		// Run the main logic simulation
		logOutput, err := runMainLogic(finalArgs, testOutputDir)
		return logOutput, err
	}

	// Test Cases
	t.Run("Defaults", func(t *testing.T) {
		logOutput, err := runParse(t, []string{})
		if err != nil && !errors.Is(err, flag.ErrHelp) {
			t.Fatalf("Expected no error for defaults, got %v. Log:\n%s", err, logOutput)
		}
		if !strings.Contains(logOutput, "Output file:") || !strings.Contains(logOutput, "codebase.md") {
			t.Errorf("Default output file logging mismatch. Log:\n%s", logOutput)
		}
		if !strings.Contains(logOutput, "- RGX: \\.git.*") {
			t.Errorf("Default ignore pattern logging mismatch. Log:\n%s", logOutput)
		}
	})

	t.Run("SetValues", func(t *testing.T) {
		args := []string{
			// -input is added by runParse helper if needed for this non-error test
			"-output", "out.md",
			"-ignore", "a,b",
			"-include", "c,d",
			"-included-paths-file", "inc.txt",
			"-excluded-paths-file", "exc.txt",
			"-instruction", "Hello Instruction",
			"-clipboard",
		}
		logOutput, err := runParse(t, args)
		if err != nil && !errors.Is(err, flag.ErrHelp) {
			t.Fatalf("Expected no error setting values, got %v. Log:\n%s", err, logOutput)
		}
		if !strings.Contains(logOutput, "Output file:") || !strings.Contains(logOutput, "out.md") {
			t.Errorf("Set outputFile logging mismatch. Log:\n%s", logOutput)
		}
		if !strings.Contains(logOutput, "- RGX: a") || !strings.Contains(logOutput, "- RGX: b") {
			t.Errorf("Set ignorePatterns logging mismatch. Log:\n%s", logOutput)
		}
		if !strings.Contains(logOutput, "+ RGX: c") || !strings.Contains(logOutput, "+ RGX: d") {
			t.Errorf("Set includePatterns logging mismatch. Log:\n%s", logOutput)
		}
		if !strings.Contains(logOutput, "Included paths will be saved to:") || !strings.Contains(logOutput, "inc.txt") {
			t.Errorf("Set includedPathsFile logging mismatch. Log:\n%s", logOutput)
		}
		if !strings.Contains(logOutput, "Excluded paths will be saved to:") || !strings.Contains(logOutput, "exc.txt") {
			t.Errorf("Set excludedPathsFile logging mismatch. Log:\n%s", logOutput)
		}
		if !strings.Contains(logOutput, "Instruction text provided.") {
			t.Errorf("Set instruction logging mismatch. Log:\n%s", logOutput)
		}
		if !strings.Contains(logOutput, "Result will be copied to clipboard.") {
			t.Errorf("Set clipboard logging mismatch. Log:\n%s", logOutput)
		}
	})

	t.Run("Error_InputNotExist", func(t *testing.T) {
		nonExistentPath := filepath.Join(t.TempDir(), "codeweaver_non_existent_dir_abc123") // Use test's tempdir
		// runParse will *not* add default -input because path contains "non_existent"
		_, err := runParse(t, []string{"-input", nonExistentPath})
		if err == nil {
			t.Fatalf("Expected error for non-existent input path, got nil")
		}
		if !strings.Contains(err.Error(), "not exist") {
			t.Errorf("Expected 'not exist' error, got err: %v", err)
		}
	})

	// --- FIX: Test case for Error_InputIsFile ---
	t.Run("Error_InputIsFile", func(t *testing.T) {
		inputFileDir := t.TempDir() // Dir to hold the file
		filePath := filepath.Join(inputFileDir, "test_file_for_input.txt")
		if err := os.WriteFile(filePath, []byte("I am a file"), 0644); err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}

		// Pass the file path explicitly with the -input flag
		_, err := runParse(t, []string{"-input", filePath})
		if err == nil {
			t.Fatalf("Expected error for file input path, got nil")
		}
		if !strings.Contains(err.Error(), "not a dir") {
			t.Errorf("Expected 'is not a directory' error, got err: %v", err)
		}
	})
}

// TestCompileRegexPatterns verifies that valid regex patterns are compiled correctly
// and invalid ones return appropriate errors.
func TestCompileRegexPatterns(t *testing.T) {
	// This test remains largely the same as it tests regex compilation.
	testLogger := log.New(io.Discard, "", 0)
	t.Run("ValidPatterns", func(t *testing.T) {
		patterns := []string{"\\.go$", "^src/"}
		matchers, err := compileRegexList(patterns, "", "", testLogger)
		if err != nil {
			t.Fatalf("no error expected, got %v", err)
		}
		if len(matchers) != 2 {
			t.Fatalf("expected 2 matchers, got %d", len(matchers))
		}
	})
	t.Run("InvalidPattern", func(t *testing.T) {
		patterns := []string{"["}
		_, err := compileRegexList(patterns, "", "", testLogger)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

// TestTreeBuilder verifies the tree construction logic, ensuring that filters
// applied during content generation are correctly reflected in the tree structure
// (e.g., hiding empty directories).
func TestTreeBuilder(t *testing.T) {
	rootDir, cleanup := createTestFS(t) // Create our standard test file system
	defer cleanup()

	testCases := []struct {
		name              string
		processedPaths    []string // Paths that contentBuilder determined should be processed
		expectedTreeLines []string // Expected lines in the tree output
	}{
		{
			name: "FullTree_AllPathsProcessed",
			processedPaths: []string{ // Simulate all paths being processed by contentBuilder (FILES ONLY now)
				".git/HEAD",
				"README.md",
				"build/output.exe", "build/tmp/log.txt",
				"data/config.yaml", "data/image.png",
				"doc_with_ticks.md", // New file
				"docs/sub_docs/file_in_sub.txt",
				// "empty_dir", // Dirs are not in processedPaths
				"file1.txt", "empty_file.txt",
				"node_modules/dep/package.json",
				"other.log",
				"script.go",
			},
			expectedTreeLines: []string{
				"├── .git", "│   └── HEAD",
				"├── README.md",
				"├── build", "│   ├── output.exe", "│   └── tmp", "│       └── log.txt",
				"├── data", "│   ├── config.yaml", "│   └── image.png",
				"├── doc_with_ticks.md",
				"├── docs", "│   └── sub_docs", "│       └── file_in_sub.txt",
				// "├── empty_dir", // Should NOT appear
				"├── empty_file.txt",
				"├── file1.txt",
				"├── node_modules", "│   └── dep", "│       └── package.json",
				"├── other.log",
				"└── script.go",
			},
		},
		{
			name: "PartialTree_OnlyGoAndMdFilesProcessed",
			processedPaths: []string{ // Only .go and .md files
				"README.md",
				"script.go",
			},
			expectedTreeLines: []string{
				"├── README.md",
				"└── script.go",
			},
		},
		{
			name: "PartialTree_SpecificFilesAndTheirDirs",
			processedPaths: []string{
				"docs/sub_docs/file_in_sub.txt", // file
				"data/config.yaml",              // file
			},
			expectedTreeLines: []string{
				"├── data", "│   └── config.yaml",
				"└── docs", "    └── sub_docs", "        └── file_in_sub.txt",
			},
		},
		{
			name:           "EmptyDir_WhenProcessed",
			processedPaths: []string{
				// "empty_dir", // Empty dir would not produce a processed path for a file
			},
			expectedTreeLines: []string{}, // Expect empty tree
		},
		{
			name:              "NoPathsProcessed",
			processedPaths:    []string{},
			expectedTreeLines: []string{}, // Expect an empty tree
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			processedSet := make(map[string]struct{})
			for _, p := range tc.processedPaths {
				processedSet[p] = struct{}{}
			}

			builder := newTreeBuilder(rootDir, processedSet)
			actualTree, err := builder.buildTreeString()
			if err != nil {
				t.Fatalf("buildTreeString() failed: %v", err)
			}
			actualTree = normalizeNewlines(actualTree)

			scanner := bufio.NewScanner(strings.NewReader(actualTree))
			actualLinesSet := make(map[string]bool)
			actualLineCount := 0
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text()) // Trim spaces for more robust comparison
				if line != "" {
					actualLinesSet[line] = true
					actualLineCount++
				}
			}
			if err := scanner.Err(); err != nil {
				t.Fatalf("Error scanning actual tree output: %v", err)
			}

			missingLines := []string{}
			for _, expectedLine := range tc.expectedTreeLines {
				if !actualLinesSet[strings.TrimSpace(expectedLine)] { // Compare trimmed lines
					missingLines = append(missingLines, expectedLine)
				}
			}

			if len(missingLines) > 0 {
				t.Errorf("Tree mismatch. Missing expected lines:\n%s\nActual Tree (trimmed lines):\n%s",
					strings.Join(missingLines, "\n"), actualTree)
			}
			if actualLineCount != len(tc.expectedTreeLines) && len(tc.expectedTreeLines) > 0 { // Avoid error for expected empty tree
				t.Errorf("Tree mismatch. Expected %d non-empty lines, got %d.\nActual Tree:\n%s",
					len(tc.expectedTreeLines), actualLineCount, actualTree)
			}
			if len(tc.expectedTreeLines) == 0 && actualLineCount != 0 {
				t.Errorf("Tree mismatch. Expected empty tree, got %d lines.\nActual Tree:\n%s", actualLineCount, actualTree)
			}
		})
	}
}

// TestContentBuilder verifies the content generation pass.
// It checks correct path filtering, exclusion of binary/empty files, and correct markdown formatting.
func TestContentBuilder(t *testing.T) {
	rootDir, cleanup := createTestFS(t)
	defer cleanup()
	testLogger := log.New(io.Discard, "", 0)

	testCases := []struct {
		name                             string
		ignoreMatchers                   []*regexp.Regexp
		includeMatchers                  []*regexp.Regexp
		expectedContentSubstr            string   // Substring to find in generated markdown content
		expectedProcessedPaths           []string // All FILES that passed filters (no directories)
		expectedExcludedPaths            []string // All files AND DIRS that failed filters
		expectEmptyFileSkippedInContent  bool     // If an empty file should be processed but not in content markdown
		expectBinaryFileSkippedInContent bool     // If a binary file should be processed but not in content markdown
		dynamicFenceCheck                bool     // If true, check for dynamic fencing on doc_with_ticks.md
	}{
		{
			name:                  "NoFilters_AllProcessed_ContentForAllNonEmpty",
			expectedContentSubstr: "## file1.txt\n\n```txt\ncontent of file1\n```",
			expectedProcessedPaths: []string{ // FILES ONLY
				".git/HEAD", "README.md",
				"build/output.exe", "build/tmp/log.txt",
				"data/config.yaml", "data/image.png",
				"doc_with_ticks.md",
				"docs/sub_docs/file_in_sub.txt",
				"empty_file.txt", "file1.txt",
				"node_modules/dep/package.json",
				"other.log", "script.go",
			},
			expectedExcludedPaths:            []string{},
			expectEmptyFileSkippedInContent:  true,
			expectBinaryFileSkippedInContent: true,
			dynamicFenceCheck:                true,
		},
		{
			name:                   "IncludeOnlyGoAndMdFiles",
			includeMatchers:        []*regexp.Regexp{mustCompileRegex(`\.go$`), mustCompileRegex(`\.md$`)},
			expectedContentSubstr:  "## script.go\n\n```go\npackage main",                   // README content also present
			expectedProcessedPaths: []string{"README.md", "doc_with_ticks.md", "script.go"}, // Only these files pass
			expectedExcludedPaths: []string{ // All other files and dirs
				".git", ".git/HEAD",
				"build", "build/output.exe", "build/tmp", "build/tmp/log.txt",
				"data", "data/config.yaml", "data/image.png",
				"docs", "docs/sub_docs", "docs/sub_docs/file_in_sub.txt",
				"empty_dir", "empty_file.txt", "file1.txt",
				"node_modules", "node_modules/dep", "node_modules/dep/package.json",
				"other.log",
			},
		},
		{
			name:                   "IgnoreGitDir_IncludeTxtFiles",
			ignoreMatchers:         []*regexp.Regexp{mustCompileRegex(`^\.git(/.*)?$`)}, // Ignore .git dir and its contents
			includeMatchers:        []*regexp.Regexp{mustCompileRegex(`\.txt$`)},
			expectedContentSubstr:  "## file1.txt\n\n```txt\ncontent of file1\n```",
			expectedProcessedPaths: []string{"build/tmp/log.txt", "docs/sub_docs/file_in_sub.txt", "empty_file.txt", "file1.txt"},
			expectedExcludedPaths: []string{
				".git", ".git/HEAD", // Explicitly ignored
				"README.md", "script.go", "other.log", // Not .txt
				"build", "build/output.exe", "build/tmp", // Dirs not .txt, output.exe not .txt
				"data", "data/config.yaml", "data/image.png", // Not .txt
				"doc_with_ticks.md",     // Not .txt
				"docs", "docs/sub_docs", // Dirs not .txt
				"empty_dir",                                                         // Dir not .txt
				"node_modules", "node_modules/dep", "node_modules/dep/package.json", // Not .txt
			},
			expectEmptyFileSkippedInContent: true,
		},
		{
			name:                   "IncludeEmptyFile_VerifyNoContentGenerated",
			includeMatchers:        []*regexp.Regexp{mustCompileRegex(`empty_file\.txt$`)},
			expectedProcessedPaths: []string{"empty_file.txt"},
			expectedExcludedPaths: []string{
				".git", ".git/HEAD", "README.md", "script.go", "other.log",
				"build", "build/output.exe", "build/tmp", "build/tmp/log.txt",
				"data", "data/config.yaml", "data/image.png",
				"doc_with_ticks.md",
				"docs", "docs/sub_docs", "docs/sub_docs/file_in_sub.txt",
				"empty_dir", "file1.txt",
				"node_modules", "node_modules/dep", "node_modules/dep/package.json",
			},
			expectEmptyFileSkippedInContent: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newContentBuilder(rootDir, "", "", tc.ignoreMatchers, tc.includeMatchers, testLogger)
			actualContentStr, actualProcessedPaths, actualExcludedPaths, err := builder.buildContentString()
			if err != nil {
				t.Fatalf("buildContentString() failed: %v", err)
			}
			actualContentStr = normalizeNewlines(actualContentStr)

			// Validate markdown content substring
			if tc.expectedContentSubstr != "" && !strings.Contains(actualContentStr, tc.expectedContentSubstr) {
				t.Errorf("Content mismatch: Substring not found.\nExpected to find:\n%s\nActual Content:\n%s-----", tc.expectedContentSubstr, actualContentStr)
			}

			// Validate if empty file content was correctly skipped from markdown
			if tc.expectEmptyFileSkippedInContent {
				if strings.Contains(actualContentStr, "## empty_file.txt") {
					t.Errorf("Empty file 'empty_file.txt' was found in markdown content, but should have been skipped.")
				}
			}

			// Validate if binary file content was correctly skipped from markdown
			if tc.expectBinaryFileSkippedInContent {
				if strings.Contains(actualContentStr, "## build/output.exe") {
					t.Errorf("Binary file 'build/output.exe' was found in markdown content, but should have been skipped.")
				}
			}

			// Validate dynamic fencing for doc_with_ticks.md
			if tc.dynamicFenceCheck {
				// doc_with_ticks.md has 3 backticks, so fence should be 4 backticks
				expectedFenceStart := "\n````md\n"
				if !strings.Contains(actualContentStr, expectedFenceStart) {
					t.Errorf("Dynamic fencing mismatch. Expected 4 backticks for doc_with_ticks.md.\nActual Content fragment:\n%s", actualContentStr)
				}
			}

			// Validate processed paths
			if !equalStringSlices(actualProcessedPaths, tc.expectedProcessedPaths) {
				t.Errorf("Processed paths mismatch.\nExpected: %v\nGot:      %v", tc.expectedProcessedPaths, actualProcessedPaths)
			}

			// Validate excluded paths
			if !equalStringSlices(actualExcludedPaths, tc.expectedExcludedPaths) {
				t.Errorf("Excluded paths mismatch.\nExpected: %v\nGot:      %v", tc.expectedExcludedPaths, actualExcludedPaths)
			}
		})
	}
}

// TestSavePathsToFile verifies the helper function for saving path lists.
func TestSavePathsToFile(t *testing.T) {
	// This test remains largely the same.
	testLogger := log.New(io.Discard, "", 0)
	t.Run("StandardSave", func(t *testing.T) {
		testDir := t.TempDir()
		tmpFilePath := filepath.Join(testDir, "test_paths.txt")
		paths := []string{"b", "a", "c"} // unsorted
		expectedContent := "a\nb\nc\n"   // savePathsToFile sorts them
		err := savePathsToFile(tmpFilePath, paths, testLogger)
		if err != nil {
			t.Fatalf("savePathsToFile failed: %v", err)
		}
		contentBytes, err := os.ReadFile(tmpFilePath)
		if err != nil {
			t.Fatalf("Failed to read back file: %v", err)
		}
		if normalizeNewlines(string(contentBytes)) != normalizeNewlines(expectedContent) {
			t.Errorf("Content mismatch.\nExpected:\n%sGot:\n%s", expectedContent, string(contentBytes))
		}
	})
}

// TestPrintHelp tests the output of the compiled binary when called with -help or -h.
func TestPrintHelp(t *testing.T) {
	// This test remains the same as it tests the compiled binary.
	tempDir := t.TempDir()
	binaryName := "codeweaver_test_binary"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}
	binaryPath := filepath.Join(tempDir, binaryName)

	buildCommand := exec.Command("go", "build", "-o", binaryPath, ".")
	buildOutput, err := buildCommand.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to build: %v\nOutput:\n%s", err, string(buildOutput))
	}

	helpArgsToTest := []string{"-h", "-help"}
	for _, helpArg := range helpArgsToTest {
		t.Run(fmt.Sprintf("WithArg_%s", helpArg), func(t *testing.T) {
			runCommand := exec.Command(binaryPath, helpArg)
			var stderrBuf bytes.Buffer
			runCommand.Stderr = &stderrBuf
			err := runCommand.Run() // Expect exit code 0 for help
			if err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() != 0 {
					t.Errorf("Binary with '%s' exited %d. Stderr:\n%s", helpArg, exitErr.ExitCode(), stderrBuf.String())
				} else if !ok {
					t.Fatalf("Failed to run binary with '%s': %v", helpArg, err)
				}
			}
			output := stderrBuf.String()
			if !strings.Contains(output, "Usage:") {
				t.Errorf("Help for '%s' missing 'Usage:'. Output:\n%s", helpArg, output)
			}
			if !strings.Contains(output, "Flags:") {
				t.Errorf("Help for '%s' missing 'Flags:'. Output:\n%s", helpArg, output)
			}
			expectedFlags := []string{"-input", "-output", "-ignore", "-clipboard", "-instruction"}
			for _, flagSig := range expectedFlags {
				if !strings.Contains(output, " "+flagSig) {
					t.Errorf("Help for '%s' missing flag '%s'. Output:\n%s", helpArg, flagSig, output)
				}
			}
		})
	}
}

// TestMainExecutionFlows verifies full application flows, checking output files,
// tree structure, and content integrity.
func TestMainExecutionFlows(t *testing.T) {
	baseInputDir, cleanupInput := createTestFS(t)
	defer cleanupInput()

	t.Run("SuccessfulRun_Basic_CheckTreeAndContent", func(t *testing.T) {
		testOutputDir := t.TempDir()
		outputFileName := "basic_run.md"
		args := []string{"-input", baseInputDir, "-output", outputFileName}
		logOutput, err := runMainLogic(args, testOutputDir)
		if err != nil {
			t.Fatalf("Run failed: %v. Log:\n%s", err, logOutput)
		}

		expectedOutputPath := filepath.ToSlash(filepath.Join(testOutputDir, outputFileName))
		if !strings.Contains(logOutput, "Markdown content written to "+expectedOutputPath) {
			t.Errorf("Missing 'Markdown content written...' log. Expected path '%s'. Got:\n%s", expectedOutputPath, logOutput)
		}
		if _, statErr := os.Stat(filepath.Join(testOutputDir, outputFileName)); statErr != nil {
			t.Errorf("Expected output file '%s' to exist, stat failed: %v", expectedOutputPath, statErr)
		}

		// Read the generated markdown and verify tree and content parts
		generatedMdBytes, readErr := os.ReadFile(filepath.Join(testOutputDir, outputFileName))
		if readErr != nil {
			t.Fatalf("Failed to read generated markdown file: %v", readErr)
		}
		generatedMd := string(generatedMdBytes)

		// Check tree view presence (basic check)
		if !strings.Contains(generatedMd, "# Tree View:") {
			t.Errorf("Generated markdown missing Tree View header.")
		}
		if strings.Contains(generatedMd, "├── .git") {
			t.Errorf("Generated markdown tree should NOT contain ignored entry '.git'. Output:\n%s", generatedMd)
		}
		if !strings.Contains(generatedMd, "├── build") { // Check for a non-ignored directory
			t.Errorf("Generated markdown tree missing expected entry 'build'. Output:\n%s", generatedMd)
		}
		if !strings.Contains(generatedMd, "└── script.go") { // Example tree entry
			t.Errorf("Generated markdown tree missing expected entry 'script.go'.")
		}

		// Check content section presence
		if !strings.Contains(generatedMd, "# Content:") {
			t.Errorf("Generated markdown missing Content header.")
		}
		if !strings.Contains(generatedMd, "## file1.txt") { // Example content entry
			t.Errorf("Generated markdown content missing expected file 'file1.txt'.")
		}
		if !strings.Contains(generatedMd, "content of file1") {
			t.Errorf("Generated markdown content for 'file1.txt' incorrect.")
		}
		// Check that empty_file.txt content is NOT present
		if strings.Contains(generatedMd, "## empty_file.txt") {
			t.Errorf("Generated markdown should NOT contain content section for 'empty_file.txt'.")
		}
	})

	t.Run("Run_WithIncludeAndExclude_CheckPathFiles", func(t *testing.T) {
		testOutputDir := t.TempDir()
		includedFile := "inc.log"
		excludedFile := "exc.log"
		args := []string{
			"-input", baseInputDir,
			"-output", "filtered_run.md",
			"-ignore", `^\.git(/.*)?$,build/output\.exe$`, // Ignore .git dir and specific exe
			"-include", `\.txt$,\.md$`, // Include only .txt and .md files
			"-included-paths-file", includedFile,
			"-excluded-paths-file", excludedFile,
		}
		logOutput, err := runMainLogic(args, testOutputDir)
		if err != nil {
			t.Fatalf("Run failed: %v. Log:\n%s", err, logOutput)
		}

		// Check extension summary logic
		if !strings.Contains(logOutput, "Included extensions: .md, .txt") {
			t.Errorf("Missing/Incorrect included extension summary. Log:\n%s", logOutput)
		}
		// Excluded should contain .go (script.go), .log (other.log), .exe (build/output.exe), .yaml (data/config.yaml), .png (data/image.png), .json (node_modules/dep/package.json)
		// .git dir contents are ignored via regex "^\.git", but directory handling might affect if files inside are reached?
		// Code logic: if excluded, we return nil. So if .git matches, we don't scan inside. Thus no extensions from inside .git.
		// However, files like script.go are definitely excluded.
		if !strings.Contains(logOutput, "Excluded extensions:") || !strings.Contains(logOutput, ".go") || !strings.Contains(logOutput, ".exe") {
			t.Errorf("Missing/Incorrect excluded extension summary. Log:\n%s", logOutput)
		}

		fullIncludedPath := filepath.Join(testOutputDir, includedFile)
		fullExcludedPath := filepath.Join(testOutputDir, excludedFile)

		// Verify included paths file content
		if _, statErr := os.Stat(fullIncludedPath); statErr != nil {
			t.Errorf("Included paths file '%s' not found.", fullIncludedPath)
		} else {
			incBytes, _ := os.ReadFile(fullIncludedPath)
			incContent := string(incBytes)
			expectedIncludes := []string{"README.md", "build/tmp/log.txt", "docs/sub_docs/file_in_sub.txt", "empty_file.txt", "file1.txt"}
			for _, p := range expectedIncludes {
				if !strings.Contains(incContent, p+"\n") {
					t.Errorf("Included paths file missing '%s'. Content:\n%s", p, incContent)
				}
			}
			if strings.Contains(incContent, ".git/HEAD") {
				t.Errorf("Included paths file should NOT contain ignored '.git/HEAD'. Content:\n%s", incContent)
			}
			if strings.Contains(incContent, "script.go") {
				t.Errorf("Included paths file should NOT contain non-txt/md 'script.go'. Content:\n%s", incContent)
			}
		}

		// Verify excluded paths file content
		if _, statErr := os.Stat(fullExcludedPath); statErr != nil {
			t.Errorf("Excluded paths file '%s' not found.", fullExcludedPath)
		} else {
			excBytes, _ := os.ReadFile(fullExcludedPath)
			excContent := string(excBytes)
			expectedExcludes := []string{".git", ".git/HEAD", "build/output.exe", "script.go", "other.log"} // Files and dirs
			for _, p := range expectedExcludes {
				if !strings.Contains(excContent, p+"\n") {
					t.Errorf("Excluded paths file missing '%s'. Content:\n%s", p, excContent)
				}
			}
			if strings.Contains(excContent, "README.md") {
				t.Errorf("Excluded paths file should NOT contain included 'README.md'. Content:\n%s", excContent)
			}
		}
	})

	t.Run("Run_WithInstruction", func(t *testing.T) {
		testOutputDir := t.TempDir()
		outputFileName := "instruction_run.md"
		instructionText := "This is a prompt instruction."
		args := []string{
			"-input", baseInputDir,
			"-output", outputFileName,
			"-instruction", instructionText,
		}
		logOutput, err := runMainLogic(args, testOutputDir)
		if err != nil {
			t.Fatalf("Run failed: %v. Log:\n%s", err, logOutput)
		}

		expectedOutputPath := filepath.Join(testOutputDir, outputFileName)
		generatedMdBytes, readErr := os.ReadFile(expectedOutputPath)
		if readErr != nil {
			t.Fatalf("Failed to read generated markdown file: %v", readErr)
		}
		generatedMd := string(generatedMdBytes)

		// Check if instruction is at the beginning
		if !strings.HasPrefix(generatedMd, instructionText+"\n\n# Tree View:") {
			t.Errorf("Instruction text not found at the beginning of the file.\nFile Start:\n%s...", generatedMd[:min(len(generatedMd), 100)])
		}
	})
}

// min returns the smaller of two integers.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

`````


## run_tests.ps1

```ps1
# go test -cover -v

go test -coverprofile="coverage.out"
go tool cover -func="coverage.out"
go tool cover -html="coverage.out"
```

