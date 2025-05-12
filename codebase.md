# Tree View:
```
.
├─ build_and_run.ps1
├─ build_and_show_help.ps1
├─ coverage.out
├─ go.mod
├─ go.sum
├─ goreleaser.yaml
├─ LICENSE
├─ main.go
├─ main_test.go
├─ README.md
└─ run_tests.ps1

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
```md
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
```


## build_and_run.ps1
```ps1
go build .

git describe --tags --abbrev=0
# ./CodeWeaver -h
./CodeWeaver -clipboard -ignore="\.git.*,.+\.exe,codebase.md,excluded_paths.txt,DRAFT\.md,coverage$"
```


## build_and_show_help.ps1
```ps1
go build .
./CodeWeaver -h

./CodeWeaver --help



```


## coverage.out
```out
mode: set
github.com/tesserato/CodeWeaver/main.go:44.13,48.16 2 0
github.com/tesserato/CodeWeaver/main.go:48.16,50.28 2 0
github.com/tesserato/CodeWeaver/main.go:50.28,52.4 1 0
github.com/tesserato/CodeWeaver/main.go:53.3,53.13 1 0
github.com/tesserato/CodeWeaver/main.go:56.2,56.34 1 0
github.com/tesserato/CodeWeaver/main.go:56.34,58.116 2 0
github.com/tesserato/CodeWeaver/main.go:58.116,61.4 2 0
github.com/tesserato/CodeWeaver/main.go:64.2,66.16 3 0
github.com/tesserato/CodeWeaver/main.go:66.16,68.3 1 0
github.com/tesserato/CodeWeaver/main.go:71.2,72.16 2 0
github.com/tesserato/CodeWeaver/main.go:72.16,74.3 1 0
github.com/tesserato/CodeWeaver/main.go:76.2,77.16 2 0
github.com/tesserato/CodeWeaver/main.go:77.16,79.3 1 0
github.com/tesserato/CodeWeaver/main.go:81.2,81.80 1 0
github.com/tesserato/CodeWeaver/main.go:84.38,86.2 1 0
github.com/tesserato/CodeWeaver/main.go:100.36,113.22 10 1
github.com/tesserato/CodeWeaver/main.go:113.22,113.49 2 0
github.com/tesserato/CodeWeaver/main.go:114.2,115.16 2 1
github.com/tesserato/CodeWeaver/main.go:115.16,117.3 1 0
github.com/tesserato/CodeWeaver/main.go:118.2,118.21 1 1
github.com/tesserato/CodeWeaver/main.go:118.21,120.3 1 0
github.com/tesserato/CodeWeaver/main.go:122.2,124.26 3 1
github.com/tesserato/CodeWeaver/main.go:124.26,126.3 1 0
github.com/tesserato/CodeWeaver/main.go:127.2,128.26 2 1
github.com/tesserato/CodeWeaver/main.go:128.26,129.35 1 1
github.com/tesserato/CodeWeaver/main.go:129.35,131.4 1 1
github.com/tesserato/CodeWeaver/main.go:132.3,132.89 1 0
github.com/tesserato/CodeWeaver/main.go:134.2,134.19 1 1
github.com/tesserato/CodeWeaver/main.go:134.19,136.3 1 1
github.com/tesserato/CodeWeaver/main.go:137.2,137.22 1 1
github.com/tesserato/CodeWeaver/main.go:137.22,139.3 1 1
github.com/tesserato/CodeWeaver/main.go:140.2,140.23 1 1
github.com/tesserato/CodeWeaver/main.go:140.23,142.3 1 1
github.com/tesserato/CodeWeaver/main.go:143.2,143.17 1 1
github.com/tesserato/CodeWeaver/main.go:146.44,151.33 5 0
github.com/tesserato/CodeWeaver/main.go:151.33,153.3 1 0
github.com/tesserato/CodeWeaver/main.go:154.2,154.33 1 0
github.com/tesserato/CodeWeaver/main.go:154.33,156.3 1 0
github.com/tesserato/CodeWeaver/main.go:157.2,157.24 1 0
github.com/tesserato/CodeWeaver/main.go:157.24,159.3 1 0
github.com/tesserato/CodeWeaver/main.go:160.2,161.15 2 0
github.com/tesserato/CodeWeaver/main.go:164.101,167.16 3 1
github.com/tesserato/CodeWeaver/main.go:167.16,169.3 1 0
github.com/tesserato/CodeWeaver/main.go:170.2,172.16 3 1
github.com/tesserato/CodeWeaver/main.go:172.16,174.3 1 0
github.com/tesserato/CodeWeaver/main.go:175.2,176.8 2 1
github.com/tesserato/CodeWeaver/main.go:179.110,180.24 1 1
github.com/tesserato/CodeWeaver/main.go:180.24,183.3 2 1
github.com/tesserato/CodeWeaver/main.go:184.2,186.29 3 1
github.com/tesserato/CodeWeaver/main.go:186.29,188.20 2 1
github.com/tesserato/CodeWeaver/main.go:188.20,189.12 1 0
github.com/tesserato/CodeWeaver/main.go:191.3,193.20 3 1
github.com/tesserato/CodeWeaver/main.go:193.20,195.4 1 1
github.com/tesserato/CodeWeaver/main.go:196.3,197.20 2 1
github.com/tesserato/CodeWeaver/main.go:199.2,199.17 1 1
github.com/tesserato/CodeWeaver/main.go:199.17,202.3 2 0
github.com/tesserato/CodeWeaver/main.go:203.2,203.22 1 1
github.com/tesserato/CodeWeaver/main.go:211.98,222.16 5 1
github.com/tesserato/CodeWeaver/main.go:222.16,224.3 1 0
github.com/tesserato/CodeWeaver/main.go:227.2,232.35 5 1
github.com/tesserato/CodeWeaver/main.go:232.35,234.3 1 1
github.com/tesserato/CodeWeaver/main.go:236.2,238.16 3 1
github.com/tesserato/CodeWeaver/main.go:238.16,240.3 1 0
github.com/tesserato/CodeWeaver/main.go:241.2,250.66 6 1
github.com/tesserato/CodeWeaver/main.go:263.93,269.2 1 1
github.com/tesserato/CodeWeaver/main.go:271.58,274.2 2 1
github.com/tesserato/CodeWeaver/main.go:276.83,278.16 2 1
github.com/tesserato/CodeWeaver/main.go:278.16,280.3 1 0
github.com/tesserato/CodeWeaver/main.go:282.2,283.32 2 1
github.com/tesserato/CodeWeaver/main.go:283.32,286.20 3 1
github.com/tesserato/CodeWeaver/main.go:286.20,288.4 1 0
github.com/tesserato/CodeWeaver/main.go:293.3,294.74 2 1
github.com/tesserato/CodeWeaver/main.go:294.74,296.4 1 1
github.com/tesserato/CodeWeaver/main.go:296.9,296.27 1 1
github.com/tesserato/CodeWeaver/main.go:296.27,298.52 2 1
github.com/tesserato/CodeWeaver/main.go:298.52,299.56 1 1
github.com/tesserato/CodeWeaver/main.go:299.56,301.11 2 1
github.com/tesserato/CodeWeaver/main.go:306.3,306.20 1 1
github.com/tesserato/CodeWeaver/main.go:306.20,308.4 1 1
github.com/tesserato/CodeWeaver/main.go:311.2,311.53 1 1
github.com/tesserato/CodeWeaver/main.go:311.53,313.3 1 1
github.com/tesserato/CodeWeaver/main.go:315.2,315.43 1 1
github.com/tesserato/CodeWeaver/main.go:315.43,318.20 3 1
github.com/tesserato/CodeWeaver/main.go:318.20,320.102 2 1
github.com/tesserato/CodeWeaver/main.go:320.102,322.5 1 0
github.com/tesserato/CodeWeaver/main.go:325.2,325.12 1 1
github.com/tesserato/CodeWeaver/main.go:328.82,330.29 2 1
github.com/tesserato/CodeWeaver/main.go:330.29,331.22 1 1
github.com/tesserato/CodeWeaver/main.go:331.22,333.4 1 1
github.com/tesserato/CodeWeaver/main.go:333.9,335.4 1 1
github.com/tesserato/CodeWeaver/main.go:337.2,337.12 1 1
github.com/tesserato/CodeWeaver/main.go:337.12,339.3 1 1
github.com/tesserato/CodeWeaver/main.go:339.8,341.3 1 1
github.com/tesserato/CodeWeaver/main.go:342.2,342.62 1 1
github.com/tesserato/CodeWeaver/main.go:345.73,347.16 2 1
github.com/tesserato/CodeWeaver/main.go:347.16,349.3 1 0
github.com/tesserato/CodeWeaver/main.go:350.2,350.46 1 1
github.com/tesserato/CodeWeaver/main.go:361.168,363.2 1 1
github.com/tesserato/CodeWeaver/main.go:371.89,379.111 4 1
github.com/tesserato/CodeWeaver/main.go:379.111,380.21 1 1
github.com/tesserato/CodeWeaver/main.go:380.21,382.69 2 0
github.com/tesserato/CodeWeaver/main.go:382.69,384.5 1 0
github.com/tesserato/CodeWeaver/main.go:385.4,385.18 1 0
github.com/tesserato/CodeWeaver/main.go:387.3,388.20 2 1
github.com/tesserato/CodeWeaver/main.go:388.20,391.4 2 0
github.com/tesserato/CodeWeaver/main.go:392.3,393.28 2 1
github.com/tesserato/CodeWeaver/main.go:393.28,395.4 1 1
github.com/tesserato/CodeWeaver/main.go:397.3,397.76 1 1
github.com/tesserato/CodeWeaver/main.go:397.76,399.34 1 1
github.com/tesserato/CodeWeaver/main.go:399.34,401.5 1 1
github.com/tesserato/CodeWeaver/main.go:402.4,403.14 2 1
github.com/tesserato/CodeWeaver/main.go:407.3,409.33 2 1
github.com/tesserato/CodeWeaver/main.go:409.33,411.4 1 1
github.com/tesserato/CodeWeaver/main.go:414.3,414.16 1 1
github.com/tesserato/CodeWeaver/main.go:414.16,416.4 1 1
github.com/tesserato/CodeWeaver/main.go:419.3,420.21 2 1
github.com/tesserato/CodeWeaver/main.go:420.21,424.4 3 0
github.com/tesserato/CodeWeaver/main.go:426.3,426.28 1 1
github.com/tesserato/CodeWeaver/main.go:426.28,429.4 1 1
github.com/tesserato/CodeWeaver/main.go:432.3,437.13 6 1
github.com/tesserato/CodeWeaver/main.go:440.2,440.20 1 1
github.com/tesserato/CodeWeaver/main.go:440.20,442.3 1 0
github.com/tesserato/CodeWeaver/main.go:443.2,443.73 1 1
github.com/tesserato/CodeWeaver/main.go:446.98,447.41 1 1
github.com/tesserato/CodeWeaver/main.go:447.41,448.60 1 1
github.com/tesserato/CodeWeaver/main.go:448.60,450.4 1 1
github.com/tesserato/CodeWeaver/main.go:452.2,452.30 1 1
github.com/tesserato/CodeWeaver/main.go:452.30,454.43 2 1
github.com/tesserato/CodeWeaver/main.go:454.43,455.61 1 1
github.com/tesserato/CodeWeaver/main.go:455.61,457.10 2 1
github.com/tesserato/CodeWeaver/main.go:460.3,460.22 1 1
github.com/tesserato/CodeWeaver/main.go:460.22,462.4 1 1
github.com/tesserato/CodeWeaver/main.go:464.2,464.13 1 1
github.com/tesserato/CodeWeaver/main.go:467.120,470.16 3 1
github.com/tesserato/CodeWeaver/main.go:470.16,473.3 2 0
github.com/tesserato/CodeWeaver/main.go:474.2,475.33 2 1
github.com/tesserato/CodeWeaver/main.go:475.33,476.87 1 1
github.com/tesserato/CodeWeaver/main.go:476.87,478.4 1 0
github.com/tesserato/CodeWeaver/main.go:480.2,480.33 1 1
github.com/tesserato/CodeWeaver/main.go:480.33,481.87 1 1
github.com/tesserato/CodeWeaver/main.go:481.87,483.4 1 0
github.com/tesserato/CodeWeaver/main.go:485.2,485.24 1 1
github.com/tesserato/CodeWeaver/main.go:485.24,487.42 2 0
github.com/tesserato/CodeWeaver/main.go:487.42,489.4 1 0
github.com/tesserato/CodeWeaver/main.go:489.9,492.4 2 0
github.com/tesserato/CodeWeaver/main.go:494.2,494.12 1 1
github.com/tesserato/CodeWeaver/main.go:497.81,498.21 1 1
github.com/tesserato/CodeWeaver/main.go:498.21,501.3 2 0
github.com/tesserato/CodeWeaver/main.go:502.2,504.26 3 1
github.com/tesserato/CodeWeaver/main.go:504.26,507.3 2 1
github.com/tesserato/CodeWeaver/main.go:508.2,509.16 2 1
github.com/tesserato/CodeWeaver/main.go:509.16,511.3 1 1
github.com/tesserato/CodeWeaver/main.go:511.8,513.3 1 0
github.com/tesserato/CodeWeaver/main.go:514.2,514.12 1 1
github.com/tesserato/CodeWeaver/main.go:517.18,533.2 15 0

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


## go.sum
```sum
github.com/BurntSushi/xgb v0.0.0-20160522181843-27f122750802/go.mod h1:IVnqGOEym/WlBOVXweHU+Q+/VP0lqqI8lqeDx9IjBqo=
github.com/yuin/goldmark v1.4.13/go.mod h1:6yULJ656Px+3vBD8DxQVa3kxgyrAnzto9xy5taEt/CY=
golang.design/x/clipboard v0.7.0 h1:4Je8M/ys9AJumVnl8m+rZnIvstSnYj1fvzqYrU3TXvo=
golang.design/x/clipboard v0.7.0/go.mod h1:PQIvqYO9GP29yINEfsEn5zSQKAz3UgXmZKzDA6dnq2E=
golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2/go.mod h1:djNgcEr1/C05ACkg1iLfiJU5Ep61QUkGW8qpdssI0+w=
golang.org/x/crypto v0.0.0-20190510104115-cbcb75029529/go.mod h1:yigFU9vqHzYiE8UmvKecakEJjdnWj3jj499lnFckfCI=
golang.org/x/crypto v0.0.0-20210921155107-089bfa567519/go.mod h1:GvvjBRRGRdwPK5ydBHafDWAxML/pGHZbMvKqRZ5+Abc=
golang.org/x/exp v0.0.0-20190731235908-ec7cb31e5a56 h1:estk1glOnSVeJ9tdEZZc5mAMDZk5lNJNyJ6DvrBkTEU=
golang.org/x/exp v0.0.0-20190731235908-ec7cb31e5a56/go.mod h1:JhuoJpWY28nO4Vef9tZUw9qufEGTyX1+7lmHxV5q5G4=
golang.org/x/image v0.0.0-20190227222117-0694c2d4d067/go.mod h1:kZ7UVZpmo3dzQBMxlp+ypCbDeSB+sBbTgSJuh5dn5js=
golang.org/x/image v0.6.0 h1:bR8b5okrPI3g/gyZakLZHeWxAR8Dn5CyxXv1hLH5g/4=
golang.org/x/image v0.6.0/go.mod h1:MXLdDR43H7cDJq5GEGXEVeeNhPgi+YYEQ2pC1byI1x0=
golang.org/x/mobile v0.0.0-20190312151609-d3739f865fa6/go.mod h1:z+o9i4GpDbdi3rU15maQ/Ox0txvL9dWGYEHz965HBQE=
golang.org/x/mobile v0.0.0-20230301163155-e0f57694e12c h1:Gk61ECugwEHL6IiyyNLXNzmu8XslmRP2dS0xjIYhbb4=
golang.org/x/mobile v0.0.0-20230301163155-e0f57694e12c/go.mod h1:aAjjkJNdrh3PMckS4B10TGS2nag27cbKR1y2BpUxsiY=
golang.org/x/mod v0.1.0/go.mod h1:0QHyrYULN0/3qlju5TqG8bIK38QM8yzMo5ekMj3DlcY=
golang.org/x/mod v0.6.0-dev.0.20220419223038-86c51ed26bb4/go.mod h1:jJ57K6gSWd91VN4djpZkiMVwK6gcyfeH4XE8wZrZaV4=
golang.org/x/mod v0.8.0/go.mod h1:iBbtSCu2XBx23ZKBPSOrRkjjQPZFPuis4dIYUhu/chs=
golang.org/x/net v0.0.0-20190311183353-d8887717615a/go.mod h1:t9HGtf8HONx5eT2rtn7q6eTqICYqUVnKs3thJo3Qplg=
golang.org/x/net v0.0.0-20190404232315-eb5bcb51f2a3/go.mod h1:t9HGtf8HONx5eT2rtn7q6eTqICYqUVnKs3thJo3Qplg=
golang.org/x/net v0.0.0-20190620200207-3b0461eec859/go.mod h1:z5CRVTTTmAJ677TzLLGU+0bjPO0LkuOLi4/5GtJWs/s=
golang.org/x/net v0.0.0-20210226172049-e18ecbb05110/go.mod h1:m0MpNAwzfU5UDzcl9v0D8zg8gWTRqZa9RBIspLL5mdg=
golang.org/x/net v0.0.0-20220722155237-a158d28d115b/go.mod h1:XRhObCWvk6IyKnWLug+ECip1KBveYUHfp+8e9klMJ9c=
golang.org/x/net v0.6.0/go.mod h1:2Tu9+aMcznHK/AK1HMvgo6xiTLG5rD5rZLDS+rp2Bjs=
golang.org/x/sync v0.0.0-20190423024810-112230192c58/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
golang.org/x/sync v0.0.0-20220722155255-886fb9371eb4/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
golang.org/x/sync v0.1.0/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
golang.org/x/sys v0.0.0-20190215142949-d0b11bdaac8a/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
golang.org/x/sys v0.0.0-20190412213103-97732733099d/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
golang.org/x/sys v0.0.0-20201119102817-f84b799fce68/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
golang.org/x/sys v0.0.0-20210615035016-665e8c7367d1/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
golang.org/x/sys v0.0.0-20220722155257-8c9f86f7a55f/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
golang.org/x/sys v0.5.0 h1:MUK/U/4lj1t1oPg0HfuXDN/Z1wv31ZJ/YcPiGccS4DU=
golang.org/x/sys v0.5.0/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
golang.org/x/term v0.0.0-20201126162022-7de9c90e9dd1/go.mod h1:bj7SfCRtBDWHUb9snDiAeCFNEtKQo2Wmx5Cou7ajbmo=
golang.org/x/term v0.0.0-20210927222741-03fcf44c2211/go.mod h1:jbD1KX2456YbFQfuXm/mYQcufACuNUgVhRMnK/tPxf8=
golang.org/x/term v0.5.0/go.mod h1:jMB1sMXY+tzblOD4FWmEbocvup2/aLOaQEp7JmGp78k=
golang.org/x/text v0.3.0/go.mod h1:NqM8EUOU14njkJ3fqMW+pc6Ldnwhi/IjpwHt7yyuwOQ=
golang.org/x/text v0.3.3/go.mod h1:5Zoc/QRtKVWzQhOtBMvqHzDpF6irO9z98xDceosuGiQ=
golang.org/x/text v0.3.7/go.mod h1:u+2+/6zg+i71rQMx5EYifcz6MCKuco9NR6JIITiCfzQ=
golang.org/x/text v0.7.0/go.mod h1:mrYo+phRRbMaCq/xk9113O4dZlRixOauAjOtrjsXDZ8=
golang.org/x/text v0.8.0/go.mod h1:e1OnstbJyHTd6l/uOt8jFFHp6TRDWZR/bV3emEE/zU8=
golang.org/x/tools v0.0.0-20180917221912-90fa682c2a6e/go.mod h1:n7NCudcB/nEzxVGmLbDWY5pfWTLqBcC2KZ6jyYvM4mQ=
golang.org/x/tools v0.0.0-20190312151545-0bb0c0a6e846/go.mod h1:LCzVGOaR6xXOjkQ3onu1FJEFr0SW1gC7cKk1uF8kGRs=
golang.org/x/tools v0.0.0-20191119224855-298f0cb1881e/go.mod h1:b+2E5dAYhXwXZwtnZ6UAqBI28+e2cm9otk0dWdXHAEo=
golang.org/x/tools v0.1.12/go.mod h1:hNGJHUnrk76NpqgfD5Aqm5Crs+Hm0VOH/i9J2+nxYbc=
golang.org/x/tools v0.6.0/go.mod h1:Xwgl3UAJ/d3gWutnCtw505GrjyAbvKui8lOU390QaIU=
golang.org/x/xerrors v0.0.0-20190717185122-a985d3407aa7/go.mod h1:I/5z698sn9Ka8TeJc9MKroUUfqBBauWjQqLJ2OPfmY0=

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
```go
package main

import (
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
const (
	colorRed       = "\033[31m"
	colorLiteRed   = "\033[91m"
	colorGreen     = "\033[32m"
	colorLiteGreen = "\033[92m"
	colorReset     = "\033[0m"
)
const (
	mdTreeViewHeader  = "# Tree View:\n```\n"
	mdCodeBlockEnd    = "\n```\n"
	mdContentHeader   = "\n# Content:\n"
	mdFileHeaderStart = "\n## "
	mdCodeBlockStart  = "\n```"
	mdCodeBlockEndNL  = "\n```\n\n"
)
const (
	rgxLogPrefixIgnore  = "- RGX:"
	rgxLogPrefixInclude = "+ RGX:"
)
const (
	logPrefixExclude = "-"
	logPrefixInclude = "+"
)

// --- Main Execution ---
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

	for _, arg := range os.Args[1:] { // Check args excluding the program name
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

	// generateMarkdown now orchestrates content and tree generation based on single processing pass
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

func isFlagHelpError(err error) bool {
	return err != nil && err.Error() == "flag: help requested"
}

type config struct {
	inputDirOriginal  string
	inputDirAbs       string
	outputFile        string
	ignorePatterns    []string
	includePatterns   []string
	includedPathsFile string
	excludedPathsFile string
	addToClipboard    bool
	showVersion       bool
}

func parseFlags() (*config, error) {
	cfg := &config{}
	flag.StringVar(&cfg.inputDirOriginal, "input", ".", "The root directory to scan.")
	flag.StringVar(&cfg.outputFile, "output", "codebase.md", "The name of the output Markdown file.")
	ignoreStr := flag.String("ignore", `\.git.*`, "Comma-separated list of regular expressions for paths to *exclude*.")
	includeStr := flag.String("include", "", "Comma-separated list of regular expressions. *Only* paths matching these are *included*.")
	flag.StringVar(&cfg.includedPathsFile, "included-paths-file", "", "Saves the list of *included* paths to this file.")
	flag.StringVar(&cfg.excludedPathsFile, "excluded-paths-file", "", "Saves the list of *excluded* paths to this file.")
	flag.BoolVar(&cfg.addToClipboard, "clipboard", false, "Copies the generated Markdown to the clipboard.")
	flag.BoolVar(&cfg.showVersion, "version", false, "Displays the version and exits.")
	// var helpFlag bool
	// flag.BoolVar(&helpFlag, "h", false, "Displays help message and exits.")

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

func setupLogging(cfg *config) *log.Logger {
	logger := log.New(os.Stdout, "", 0)
	logger.Println("Starting CodeWeaver...")
	logger.Println("Input directory:", cfg.inputDirAbs)
	logger.Println("Output file:", cfg.outputFile)
	if cfg.includedPathsFile != "" {
		logger.Println("Included paths will be saved to:", cfg.includedPathsFile)
	}
	if cfg.excludedPathsFile != "" {
		logger.Println("Excluded paths will be saved to:", cfg.excludedPathsFile)
	}
	if cfg.addToClipboard {
		logger.Println("Result will be copied to clipboard.")
	}
	logger.Println()
	return logger
}

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
// It now calls contentBuilder first, then uses its results for the treeBuilder.
func generateMarkdown(cfg *config, ignoreMatchers, includeMatchers []*regexp.Regexp, logger *log.Logger) (
	finalMarkdown string, processedPathsForFile []string, excludedPathsForFile []string, err error) {

	var fullMarkdown strings.Builder

	// --- Build Content Section FIRST to get processedPaths ---
	logger.Println("Processing paths and building content section...")
	contentBuilder := newContentBuilder(cfg.inputDirAbs, cfg.includedPathsFile, cfg.excludedPathsFile, ignoreMatchers, includeMatchers, logger)
	// contentMarkdown is the markdown string of file contents
	// processedPaths contains ALL files AND directories that passed shouldProcess
	// excludedPaths contains all files/dirs that failed shouldProcess
	contentMarkdown, processedPaths, excludedPaths, err := contentBuilder.buildContentString()
	if err != nil {
		return "", nil, nil, fmt.Errorf("failed to build code content: %w", err)
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

type treeBuilder struct {
	rootAbsPath       string
	processedPathsSet map[string]struct{} // Set of paths that passed filters (from contentBuilder)
	output            strings.Builder
	depthOpen         map[int]bool
}

// newTreeBuilder creates a new tree builder instance, now taking processedPathsSet.
func newTreeBuilder(rootAbsPath string, processedPathsSet map[string]struct{}) *treeBuilder {
	return &treeBuilder{
		rootAbsPath:       rootAbsPath,
		processedPathsSet: processedPathsSet,
		depthOpen:         make(map[int]bool),
	}
}

func (tb *treeBuilder) buildTreeString() (string, error) {
	err := tb.printTreeRecursive(tb.rootAbsPath, 0)
	return tb.output.String(), err
}

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

func (tb *treeBuilder) printEntryLine(entry fs.DirEntry, depth int, isLast bool) {
	var prefix strings.Builder
	for i := 0; i < depth; i++ {
		if tb.depthOpen[i] {
			prefix.WriteString("│  ")
		} else {
			prefix.WriteString("   ")
		}
	}
	if isLast {
		prefix.WriteString("└─ ")
	} else {
		prefix.WriteString("├─ ")
	}
	tb.output.WriteString(prefix.String() + entry.Name() + "\n")
}

func (tb *treeBuilder) getRelativePath(fullPath string) (string, error) {
	pathRelToInput, err := filepath.Rel(tb.rootAbsPath, fullPath)
	if err != nil {
		return "", fmt.Errorf("rel path for %s to %s: %w", fullPath, tb.rootAbsPath, err)
	}
	return filepath.ToSlash(pathRelToInput), nil
}

// --- Content Builder ---

type contentBuilder struct {
	rootAbsPath, includedPathsFile, excludedPathsFile string
	ignoreMatchers, includeMatchers                   []*regexp.Regexp
	logger                                            *log.Logger
}

func newContentBuilder(rootAbsPath, includedPathsFile, excludedPathsFile string, ignoreMatchers, includeMatchers []*regexp.Regexp, logger *log.Logger) *contentBuilder {
	return &contentBuilder{rootAbsPath, includedPathsFile, excludedPathsFile, ignoreMatchers, includeMatchers, logger}
}

// buildContentString now returns:
// 1. markdown string for file contents
// 2. allProcessedPaths (files AND dirs that passed shouldProcess)
// 3. excludedPaths (files AND dirs that failed shouldProcess)
// 4. error
func (cb *contentBuilder) buildContentString() (
	markdownContent string, allProcessedPaths []string, excludedPaths []string, err error) {

	var contentSB strings.Builder // For actual file content markdown

	// Slices to store paths
	var localProcessedPaths []string
	var localExcludedPaths []string

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
			return nil // Path excluded, continue walk
		}

		// Path passed filters, add to processedPaths
		localProcessedPaths = append(localProcessedPaths, pathRelToInput)
		// Log inclusion if not saving to file
		if cb.includedPathsFile == "" {
			cb.logger.Printf("%s%s %s%s\n", colorGreen, logPrefixInclude, pathRelToInput, colorReset)
		}

		// If it's a directory or an empty/unreadable file, don't add its content to markdown
		if d.IsDir() {
			return nil // Continue into directory
		}

		// Process file for content inclusion
		fileContent, readErr := os.ReadFile(currentWalkPath)
		if readErr != nil {
			cb.logger.Printf("%sWarning: Failed to read file %s: %v%s\n", colorRed, currentWalkPath, readErr, colorReset)
			contentSB.WriteString(fmt.Sprintf("%s%s\n%s\nError reading file: %v%s", mdFileHeaderStart, pathRelToInput, mdCodeBlockStart, readErr, mdCodeBlockEndNL))
			return nil // File processed (passed filters), but content not added
		}

		if len(fileContent) == 0 {
			// Empty file: it's in localProcessedPaths, but no markdown content for it.
			return nil
		}

		// Add file content to markdown
		extension := strings.TrimPrefix(strings.ToLower(filepath.Ext(currentWalkPath)), ".")
		contentSB.WriteString(fmt.Sprintf("%s%s", mdFileHeaderStart, pathRelToInput))
		contentSB.WriteString(fmt.Sprintf("%s%s\n", mdCodeBlockStart, extension))
		contentSB.Write(fileContent)
		contentSB.WriteString(mdCodeBlockEndNL)
		return nil
	})

	if walkErr != nil {
		return "", nil, nil, fmt.Errorf("walking directory %s: %w", cb.rootAbsPath, walkErr)
	}
	return contentSB.String(), localProcessedPaths, localExcludedPaths, nil
}

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

func writeOutput(cfg *config, markdownContent string, includedPaths, excludedPaths []string, logger *log.Logger) error {
	logger.Printf("Writing output to %s...", cfg.outputFile)
	err := os.WriteFile(cfg.outputFile, []byte(markdownContent), 0644)
	if err != nil {
		logger.Printf("%sError writing to output file %s: %v%s", colorRed, cfg.outputFile, err, colorReset)
		return fmt.Errorf("writing output file %s: %w", cfg.outputFile, err)
	}
	logger.Printf("Markdown content written to %s", cfg.outputFile)
	if cfg.includedPathsFile != "" {
		if err := savePathsToFile(cfg.includedPathsFile, includedPaths, logger); err != nil { // Pass `includedPaths` from generateMarkdown
			logger.Printf("%sWarning: Error saving included paths to %s: %v%s", colorRed, cfg.includedPathsFile, err, colorReset)
		}
	}
	if cfg.excludedPathsFile != "" {
		if err := savePathsToFile(cfg.excludedPathsFile, excludedPaths, logger); err != nil { // Pass `excludedPaths` from generateMarkdown
			logger.Printf("%sWarning: Error saving excluded paths to %s: %v%s", colorRed, cfg.excludedPathsFile, err, colorReset)
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

func savePathsToFile(filename string, paths []string, logger *log.Logger) error {
	if len(paths) == 0 {
		logger.Printf("No paths to save to %s.", filename)
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
		logger.Printf("Paths saved to %s", filename)
	} else {
		return fmt.Errorf("writing paths file %s: %w", filename, err)
	}
	return nil
}

func printHelp() {
	fmt.Fprintf(os.Stderr, "CodeWeaver: Generate Markdown Documentation from Your Codebase.\n")
	fmt.Fprintf(os.Stderr, "Usage: codeweaver [options]\nFor help, use -h or --help.\n")
	fmt.Fprintf(os.Stderr, "\nOptions:\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\nExamples:\n")
	fmt.Fprintf(os.Stderr, "  codeweaver                                   # Process current directory, output to codebase.md\n")
	fmt.Fprintf(os.Stderr, "  codeweaver -input my_project -output docs.md # Specify input and output\n")
	fmt.Fprintf(os.Stderr, `  codeweaver -ignore "build/,vendor/" -include "\.go$,\.md$" # Filter paths`+"\n")
	fmt.Fprintf(os.Stderr, "  codeweaver -clipboard -excluded-paths-file ignored.txt # Copy & log excluded\n")
	fmt.Fprintf(os.Stderr, "\nNotes on patterns:\n")
	fmt.Fprintf(os.Stderr, "  - Patterns are Go regular expressions (https://pkg.go.dev/regexp/syntax).\n")
	fmt.Fprintf(os.Stderr, "  - Paths for filtering are relative to the input directory (e.g., \"src/main.go\").\n")
	fmt.Fprintf(os.Stderr, "  - Use forward slashes '/' in patterns for cross-platform compatibility (e.g., \"data/images/\").\n")
	fmt.Fprintf(os.Stderr, "  - Use 'path/to/dir/?' to match a directory itself (with or without a trailing slash).\n")
	fmt.Fprintf(os.Stderr, "  - Use 'path/to/dir(/.*)?' to match a directory AND its contents.\n")
}

```


## main_test.go
```go
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

func normalizeNewlines(s string) string { return strings.ReplaceAll(s, "\r\n", "\n") }

func createTestFS(t *testing.T) (string, func()) {
	t.Helper()
	rootDir := t.TempDir()
	structure := map[string]string{
		"file1.txt": "content of file1",
		"script.go": "package main\nfunc main() {}",
		"README.md": "# Test Readme",
		"data/":     "", "data/image.png": "fake png data", "data/config.yaml": "key: value",
		"build/": "", "build/output.exe": "binary data", "build/tmp/": "", "build/tmp/log.txt": "log entry",
		".git/": "", ".git/HEAD": "ref: refs/heads/main",
		"node_modules/": "", "node_modules/dep/": "", "node_modules/dep/package.json": "{}",
		"empty_dir/":                    "",
		"docs/sub_docs/file_in_sub.txt": "nested doc content",
		"other.log":                     "another log",
		"empty_file.txt":                "", // Explicitly empty file
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

// runMainLogic - Keep this helper as it's useful for testing *internal* logic flows
// and error handling paths without the overhead of compilation for every case.
// ... (runMainLogic implementation needs minor adjustment for how it calls generateMarkdown if its signature changed, but the core flag parsing remains the same) ...
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
	testRunLogger.Println("Input directory:", cfg.inputDirAbs)
	testRunLogger.Println("Output file:", filepath.Join(outputDir, cfg.outputFile))
	if cfg.includedPathsFile != "" {
		testRunLogger.Println("Included paths will be saved to:", filepath.Join(outputDir, cfg.includedPathsFile))
	}
	if cfg.excludedPathsFile != "" {
		testRunLogger.Println("Excluded paths will be saved to:", filepath.Join(outputDir, cfg.excludedPathsFile))
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

// TestTreeBuilder now focuses on tree construction given a set of processed paths.
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
			processedPaths: []string{ // Simulate all paths being processed by contentBuilder
				".git", ".git/HEAD",
				"README.md",
				"build", "build/output.exe", "build/tmp", "build/tmp/log.txt",
				"data", "data/config.yaml", "data/image.png",
				"docs", "docs/sub_docs", "docs/sub_docs/file_in_sub.txt",
				"empty_dir",
				"file1.txt", "empty_file.txt",
				"node_modules", "node_modules/dep", "node_modules/dep/package.json",
				"other.log",
				"script.go",
			},
			expectedTreeLines: []string{
				"├─ .git", "│  └─ HEAD",
				"├─ README.md",
				"├─ build", "│  ├─ output.exe", "│  └─ tmp", "│     └─ log.txt",
				"├─ data", "│  ├─ config.yaml", "│  └─ image.png",
				"├─ docs", "│  └─ sub_docs", "│     └─ file_in_sub.txt",
				"├─ empty_dir",
				"├─ empty_file.txt",
				"├─ file1.txt",
				"├─ node_modules", "│  └─ dep", "│     └─ package.json",
				"├─ other.log",
				"└─ script.go",
			},
		},
		{
			name: "PartialTree_OnlyGoAndMdFilesProcessed",
			processedPaths: []string{ // Only .go and .md files (and their parent dirs for tree structure)
				"README.md",
				"script.go",
				// Implicitly, parent directories like "." are needed for WalkDir to start
				// but treeBuilder logic ensures parent dirs of processed files are shown.
			},
			expectedTreeLines: []string{
				"├─ README.md",
				"└─ script.go",
			},
		},
		{
			name: "PartialTree_SpecificFilesAndTheirDirs",
			processedPaths: []string{
				"docs", "docs/sub_docs", "docs/sub_docs/file_in_sub.txt", // file + its parent dirs
				"data", "data/config.yaml", // file + its parent dir
			},
			expectedTreeLines: []string{
				"├─ data", "│  └─ config.yaml",
				"└─ docs", "   └─ sub_docs", "      └─ file_in_sub.txt",
			},
		},
		{
			name: "EmptyDir_WhenProcessed",
			processedPaths: []string{
				"empty_dir",
			},
			expectedTreeLines: []string{
				"└─ empty_dir",
			},
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

// TestContentBuilder focuses on the output of contentBuilder: markdown content, processed paths, and excluded paths.
func TestContentBuilder(t *testing.T) {
	rootDir, cleanup := createTestFS(t)
	defer cleanup()
	testLogger := log.New(io.Discard, "", 0)

	testCases := []struct {
		name                            string
		ignoreMatchers                  []*regexp.Regexp
		includeMatchers                 []*regexp.Regexp
		expectedContentSubstr           string   // Substring to find in generated markdown content
		expectedProcessedPaths          []string // All files AND DIRS that passed filters
		expectedExcludedPaths           []string // All files AND DIRS that failed filters
		expectEmptyFileSkippedInContent bool     // If an empty file should be processed but not in content markdown
	}{
		{
			name:                  "NoFilters_AllProcessed_ContentForAllNonEmpty",
			expectedContentSubstr: "## file1.txt\n```txt\ncontent of file1\n```",
			expectedProcessedPaths: []string{ // Includes dirs now
				".git", ".git/HEAD", "README.md",
				"build", "build/output.exe", "build/tmp", "build/tmp/log.txt",
				"data", "data/config.yaml", "data/image.png",
				"docs", "docs/sub_docs", "docs/sub_docs/file_in_sub.txt",
				"empty_dir", "empty_file.txt", "file1.txt",
				"node_modules", "node_modules/dep", "node_modules/dep/package.json",
				"other.log", "script.go",
			},
			expectedExcludedPaths:           []string{},
			expectEmptyFileSkippedInContent: true,
		},
		{
			name:                   "IncludeOnlyGoAndMdFiles",
			includeMatchers:        []*regexp.Regexp{mustCompileRegex(`\.go$`), mustCompileRegex(`\.md$`)},
			expectedContentSubstr:  "## script.go\n```go\npackage main", // README content also present
			expectedProcessedPaths: []string{"README.md", "script.go"},  // Only these files pass
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
			expectedContentSubstr:  "## file1.txt\n```txt\ncontent of file1\n```",
			expectedProcessedPaths: []string{"build/tmp/log.txt", "docs/sub_docs/file_in_sub.txt", "empty_file.txt", "file1.txt"},
			expectedExcludedPaths: []string{
				".git", ".git/HEAD", // Explicitly ignored
				"README.md", "script.go", "other.log", // Not .txt
				"build", "build/output.exe", "build/tmp", // Dirs not .txt, output.exe not .txt
				"data", "data/config.yaml", "data/image.png", // Not .txt
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
			if !strings.Contains(output, "Usage: codeweaver [options]") {
				t.Errorf("Help for '%s' missing 'Usage:'. Output:\n%s", helpArg, output)
			}
			if !strings.Contains(output, "Options:") {
				t.Errorf("Help for '%s' missing 'Options:'. Output:\n%s", helpArg, output)
			}
			expectedFlags := []string{"-input string", "-output string", "-ignore string", "-clipboard"}
			for _, flagSig := range expectedFlags {
				if !strings.Contains(output, " "+flagSig) {
					t.Errorf("Help for '%s' missing flag '%s'. Output:\n%s", helpArg, flagSig, output)
				}
			}
		})
	}
}

// TestMainExecutionFlows needs careful review for path assertions.
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

		expectedOutputPath := filepath.Join(testOutputDir, outputFileName)
		if !strings.Contains(logOutput, "Markdown content written to "+expectedOutputPath) {
			t.Errorf("Missing 'Markdown content written...' log. Expected path '%s'. Got:\n%s", expectedOutputPath, logOutput)
		}
		if _, statErr := os.Stat(expectedOutputPath); statErr != nil {
			t.Errorf("Expected output file '%s' to exist, stat failed: %v", expectedOutputPath, statErr)
		}

		// Read the generated markdown and verify tree and content parts
		generatedMdBytes, readErr := os.ReadFile(expectedOutputPath)
		if readErr != nil {
			t.Fatalf("Failed to read generated markdown file: %v", readErr)
		}
		generatedMd := string(generatedMdBytes)

		// Check tree view presence (basic check)
		if !strings.Contains(generatedMd, "# Tree View:") {
			t.Errorf("Generated markdown missing Tree View header.")
		}
		if strings.Contains(generatedMd, "├─ .git") {
			t.Errorf("Generated markdown tree should NOT contain ignored entry '.git'. Output:\n%s", generatedMd)
		}
		if !strings.Contains(generatedMd, "├─ build") { // Check for a non-ignored directory
			t.Errorf("Generated markdown tree missing expected entry 'build'. Output:\n%s", generatedMd)
		}
		if !strings.Contains(generatedMd, "└─ script.go") { // Example tree entry
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

	// ... (other TestMainExecutionFlows like VersionFlag, HelpFlag, Error_InvalidPatterns, etc. can remain similar,
	//      just ensure they use testOutputDir for runMainLogic) ...
	// t.Run("VersionFlag", func(t *testing.T) {
	// 	testOutputDir := t.TempDir()
	// 	logOutput, err := runMainLogic([]string{"-version"}, testOutputDir)
	// 	if err != nil {
	// 		t.Fatalf("runMainLogic with -version failed: %v", err)
	// 	}
	// 	expected := fmt.Sprintf("CodeWeaver version %s", version)
	// 	if !strings.Contains(logOutput, expected) {
	// 		t.Errorf("Expected '%s' in output, got:\n%s", expected, logOutput)
	// 	}
	// })

	// t.Run("HelpFlag_via_runMainLogic", func(t *testing.T) { // Differentiate from binary test
	// 	testOutputDir := t.TempDir()
	// 	// Use -h for runMainLogic as it relies on flag package's default handling for Usage
	// 	logOutput, err := runMainLogic([]string{"-h"}, testOutputDir)
	// 	if !errors.Is(err, flag.ErrHelp) {
	// 		t.Fatalf("runMainLogic with -h did not return flag.ErrHelp, got err: %v. Log:\n%s", err, logOutput)
	// 	}
	// 	// Check that printHelp was invoked (which is part of flag.Usage)
	// 	if !strings.Contains(logOutput, "Usage: codeweaver [options]") {
	// 		t.Errorf("Expected 'Usage:' in help output from runMainLogic, got:\n%s", logOutput)
	// 	}
	// })

}

```


## run_tests.ps1
```ps1
# go test -cover -v

go test -coverprofile="coverage.out"
go tool cover -func="coverage.out"
go tool cover -html="coverage.out"
```

