# Tree View:
```
.
├─ LICENSE
├─ README.md
├─ build_and_run.ps1
├─ coverage.out
├─ go.mod
├─ go.sum
├─ goreleaser.yaml
├─ main.go
├─ main_test.go
├─ run_tests.ps1
└─ test_root
   ├─ File at root A.txt
   ├─ File at root B.md
   ├─ folder 01
   │  ├─ File at folder 01 I.txt
   │  ├─ File at folder 01 II.md
   │  └─ File at folder 01 III.csv
   └─ folder 02
      ├─ File at folder 02 I.txt
      ├─ File at folder 02 II.md
      ├─ File at folder 02 III.csv
      └─ folder 02 01
         ├─ File at folder 02 01 I.txt
         ├─ File at folder 02 01 II.md
         └─ File at folder 02 01 III.csv
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

./CodeWeaver -clipboard -ignore="\.git.*,.+\.exe,codebase.md,excluded_paths.txt,DRAFT\.md,coverage$"
```

## coverage.out
```out
mode: set
github.com/tesserato/CodeWeaver/main.go:32.13,34.16 2 0
github.com/tesserato/CodeWeaver/main.go:34.16,37.3 1 0
github.com/tesserato/CodeWeaver/main.go:39.2,39.21 1 0
github.com/tesserato/CodeWeaver/main.go:39.21,42.3 2 0
github.com/tesserato/CodeWeaver/main.go:44.2,44.18 1 0
github.com/tesserato/CodeWeaver/main.go:44.18,47.3 2 0
github.com/tesserato/CodeWeaver/main.go:49.2,54.33 5 0
github.com/tesserato/CodeWeaver/main.go:54.33,56.3 1 0
github.com/tesserato/CodeWeaver/main.go:57.2,57.33 1 0
github.com/tesserato/CodeWeaver/main.go:57.33,59.3 1 0
github.com/tesserato/CodeWeaver/main.go:60.2,60.24 1 0
github.com/tesserato/CodeWeaver/main.go:60.24,62.3 1 0
github.com/tesserato/CodeWeaver/main.go:63.2,66.16 3 0
github.com/tesserato/CodeWeaver/main.go:66.16,68.3 1 0
github.com/tesserato/CodeWeaver/main.go:69.2,70.16 2 0
github.com/tesserato/CodeWeaver/main.go:70.16,72.3 1 0
github.com/tesserato/CodeWeaver/main.go:73.2,84.16 7 0
github.com/tesserato/CodeWeaver/main.go:84.16,86.3 1 0
github.com/tesserato/CodeWeaver/main.go:87.2,94.16 6 0
github.com/tesserato/CodeWeaver/main.go:94.16,96.3 1 0
github.com/tesserato/CodeWeaver/main.go:97.2,101.16 3 0
github.com/tesserato/CodeWeaver/main.go:101.16,103.3 1 0
github.com/tesserato/CodeWeaver/main.go:104.2,107.33 2 0
github.com/tesserato/CodeWeaver/main.go:107.33,108.87 1 0
github.com/tesserato/CodeWeaver/main.go:108.87,110.4 1 0
github.com/tesserato/CodeWeaver/main.go:112.2,112.33 1 0
github.com/tesserato/CodeWeaver/main.go:112.33,113.87 1 0
github.com/tesserato/CodeWeaver/main.go:113.87,115.4 1 0
github.com/tesserato/CodeWeaver/main.go:119.2,119.24 1 0
github.com/tesserato/CodeWeaver/main.go:119.24,120.42 1 0
github.com/tesserato/CodeWeaver/main.go:120.42,122.4 1 0
github.com/tesserato/CodeWeaver/main.go:122.9,125.4 2 0
github.com/tesserato/CodeWeaver/main.go:142.36,157.37 13 0
github.com/tesserato/CodeWeaver/main.go:157.37,160.3 1 0
github.com/tesserato/CodeWeaver/main.go:162.2,164.16 3 0
github.com/tesserato/CodeWeaver/main.go:164.16,166.3 1 0
github.com/tesserato/CodeWeaver/main.go:169.2,170.16 2 0
github.com/tesserato/CodeWeaver/main.go:170.16,171.25 1 0
github.com/tesserato/CodeWeaver/main.go:171.25,173.4 1 0
github.com/tesserato/CodeWeaver/main.go:174.3,174.91 1 0
github.com/tesserato/CodeWeaver/main.go:176.2,176.19 1 0
github.com/tesserato/CodeWeaver/main.go:176.19,178.3 1 0
github.com/tesserato/CodeWeaver/main.go:180.2,180.22 1 0
github.com/tesserato/CodeWeaver/main.go:180.22,182.3 1 0
github.com/tesserato/CodeWeaver/main.go:183.2,183.23 1 0
github.com/tesserato/CodeWeaver/main.go:183.23,185.3 1 0
github.com/tesserato/CodeWeaver/main.go:187.2,187.17 1 0
github.com/tesserato/CodeWeaver/main.go:190.114,191.24 1 1
github.com/tesserato/CodeWeaver/main.go:191.24,193.3 1 1
github.com/tesserato/CodeWeaver/main.go:195.2,196.29 2 1
github.com/tesserato/CodeWeaver/main.go:196.29,198.27 2 1
github.com/tesserato/CodeWeaver/main.go:198.27,199.12 1 1
github.com/tesserato/CodeWeaver/main.go:201.3,203.17 3 1
github.com/tesserato/CodeWeaver/main.go:203.17,206.4 1 1
github.com/tesserato/CodeWeaver/main.go:207.3,207.51 1 1
github.com/tesserato/CodeWeaver/main.go:210.2,210.32 1 1
github.com/tesserato/CodeWeaver/main.go:210.32,212.3 1 0
github.com/tesserato/CodeWeaver/main.go:213.2,213.30 1 1
github.com/tesserato/CodeWeaver/main.go:218.98,220.41 1 1
github.com/tesserato/CodeWeaver/main.go:220.41,221.60 1 1
github.com/tesserato/CodeWeaver/main.go:221.60,223.4 1 1
github.com/tesserato/CodeWeaver/main.go:227.2,227.30 1 1
github.com/tesserato/CodeWeaver/main.go:227.30,229.43 2 1
github.com/tesserato/CodeWeaver/main.go:229.43,230.61 1 1
github.com/tesserato/CodeWeaver/main.go:230.61,232.10 2 1
github.com/tesserato/CodeWeaver/main.go:235.3,235.22 1 1
github.com/tesserato/CodeWeaver/main.go:235.22,237.4 1 1
github.com/tesserato/CodeWeaver/main.go:240.2,240.13 1 1
github.com/tesserato/CodeWeaver/main.go:251.104,258.2 1 1
github.com/tesserato/CodeWeaver/main.go:260.58,263.2 2 1
github.com/tesserato/CodeWeaver/main.go:265.83,267.16 2 1
github.com/tesserato/CodeWeaver/main.go:267.16,269.3 1 1
github.com/tesserato/CodeWeaver/main.go:271.2,272.32 2 1
github.com/tesserato/CodeWeaver/main.go:272.32,275.17 3 1
github.com/tesserato/CodeWeaver/main.go:275.17,277.4 1 0
github.com/tesserato/CodeWeaver/main.go:278.3,280.75 2 1
github.com/tesserato/CodeWeaver/main.go:280.75,282.4 1 1
github.com/tesserato/CodeWeaver/main.go:286.2,286.50 1 1
github.com/tesserato/CodeWeaver/main.go:286.50,288.3 1 1
github.com/tesserato/CodeWeaver/main.go:290.2,290.40 1 1
github.com/tesserato/CodeWeaver/main.go:290.40,294.20 3 1
github.com/tesserato/CodeWeaver/main.go:294.20,297.18 3 1
github.com/tesserato/CodeWeaver/main.go:297.18,299.5 1 0
github.com/tesserato/CodeWeaver/main.go:302.2,302.12 1 1
github.com/tesserato/CodeWeaver/main.go:305.78,307.29 2 1
github.com/tesserato/CodeWeaver/main.go:307.29,308.22 1 1
github.com/tesserato/CodeWeaver/main.go:308.22,310.4 1 1
github.com/tesserato/CodeWeaver/main.go:310.9,312.4 1 1
github.com/tesserato/CodeWeaver/main.go:315.2,315.12 1 1
github.com/tesserato/CodeWeaver/main.go:315.12,317.3 1 1
github.com/tesserato/CodeWeaver/main.go:317.8,319.3 1 1
github.com/tesserato/CodeWeaver/main.go:321.2,323.29 3 1
github.com/tesserato/CodeWeaver/main.go:336.177,345.2 1 1
github.com/tesserato/CodeWeaver/main.go:347.84,352.103 4 1
github.com/tesserato/CodeWeaver/main.go:352.103,353.17 1 1
github.com/tesserato/CodeWeaver/main.go:353.17,361.4 2 1
github.com/tesserato/CodeWeaver/main.go:363.3,364.20 2 1
github.com/tesserato/CodeWeaver/main.go:364.20,368.4 2 0
github.com/tesserato/CodeWeaver/main.go:369.3,373.28 2 1
github.com/tesserato/CodeWeaver/main.go:373.28,375.4 1 1
github.com/tesserato/CodeWeaver/main.go:377.3,377.76 1 1
github.com/tesserato/CodeWeaver/main.go:377.76,378.34 1 1
github.com/tesserato/CodeWeaver/main.go:378.34,380.5 1 1
github.com/tesserato/CodeWeaver/main.go:381.4,382.17 2 1
github.com/tesserato/CodeWeaver/main.go:382.17,384.5 1 1
github.com/tesserato/CodeWeaver/main.go:385.4,385.14 1 1
github.com/tesserato/CodeWeaver/main.go:389.3,389.17 1 1
github.com/tesserato/CodeWeaver/main.go:389.17,390.34 1 1
github.com/tesserato/CodeWeaver/main.go:390.34,392.5 1 1
github.com/tesserato/CodeWeaver/main.go:393.4,396.22 3 1
github.com/tesserato/CodeWeaver/main.go:396.22,402.5 4 0
github.com/tesserato/CodeWeaver/main.go:404.4,408.36 5 1
github.com/tesserato/CodeWeaver/main.go:410.3,410.13 1 1
github.com/tesserato/CodeWeaver/main.go:413.2,413.16 1 1
github.com/tesserato/CodeWeaver/main.go:413.16,415.3 1 1
github.com/tesserato/CodeWeaver/main.go:416.2,416.60 1 1
github.com/tesserato/CodeWeaver/main.go:419.81,420.21 1 1
github.com/tesserato/CodeWeaver/main.go:420.21,425.3 1 1
github.com/tesserato/CodeWeaver/main.go:427.2,428.26 2 1
github.com/tesserato/CodeWeaver/main.go:428.26,431.3 2 1
github.com/tesserato/CodeWeaver/main.go:433.2,434.16 2 1
github.com/tesserato/CodeWeaver/main.go:434.16,436.3 1 1
github.com/tesserato/CodeWeaver/main.go:437.2,437.12 1 1
github.com/tesserato/CodeWeaver/main.go:440.18,455.2 14 1

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
      - -X main.version={{.Version}}
      - -X main.commit={{.Commit}}
      - -X main.date={{.Date}}

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

const (
	colorRed       = "\033[31m"
	colorLiteRed   = "\033[91m"
	colorGreen     = "\033[32m"
	colorLiteGreen = "\033[92m"
	colorReset     = "\033[0m"
)

var (
	// Populated by goreleaser during build
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cfg, err := parseFlags()
	if err != nil {
		// parseFlags handles help and version, so other errors are actual problems.
		log.Fatalf("Error parsing flags: %v", err)
	}

	if cfg.showVersion {
		fmt.Printf("CodeWeaver version %s\ncommit %s\nbuilt at %s\n", version, commit, date)
		return
	}

	if cfg.showHelp {
		printHelp()
		return
	}

	logger := log.New(os.Stdout, "", 0) // Simple logger for progress messages

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

	ignoreMatchers, err := compileRegexPatterns(cfg.ignorePatterns, colorLiteRed, "- RGX:", logger)
	if err != nil {
		log.Fatalf("Error compiling ignore patterns: %v", err)
	}
	includeMatchers, err := compileRegexPatterns(cfg.includePatterns, colorLiteGreen, "+ RGX:", logger)
	if err != nil {
		log.Fatalf("Error compiling include patterns: %v", err)
	}
	logger.Println()

	var markdownContent strings.Builder

	// --- Build Tree View ---
	markdownContent.WriteString("# Tree View:\n```\n")
	// Display the original input path as the root, not the absolute one, for user-friendliness
	markdownContent.WriteString(filepath.ToSlash(cfg.inputDirOriginal) + "\n")

	treeBuilder := newTreeBuilder(cfg.inputDirAbs, ignoreMatchers, includeMatchers)
	treeString, err := treeBuilder.buildTreeString()
	if err != nil {
		log.Fatalf("Error building codebase tree: %v", err)
	}
	markdownContent.WriteString(treeString)
	markdownContent.WriteString("```\n")

	// --- Build Content Section ---
	markdownContent.WriteString("\n# Content:\n")
	contentBuilder := newContentBuilder(cfg.inputDirAbs, cfg.includedPathsFile, cfg.excludedPathsFile, ignoreMatchers, includeMatchers, logger)
	contentString, includedPaths, excludedPaths, err := contentBuilder.buildContentString()
	if err != nil {
		log.Fatalf("Error writing code content: %v", err)
	}
	markdownContent.WriteString(contentString)

	// --- Write to Output File ---
	err = os.WriteFile(cfg.outputFile, []byte(markdownContent.String()), 0644)
	if err != nil {
		log.Fatalf("Error writing to output file %s: %v", cfg.outputFile, err)
	}
	logger.Printf("Markdown content written to %s\n", cfg.outputFile)

	// --- Save Included/Excluded Paths ---
	if cfg.includedPathsFile != "" {
		if err := savePathsToFile(cfg.includedPathsFile, includedPaths, logger); err != nil {
			logger.Printf("%sWarning: Error saving included paths to %s: %v%s\n", colorRed, cfg.includedPathsFile, err, colorReset)
		}
	}
	if cfg.excludedPathsFile != "" {
		if err := savePathsToFile(cfg.excludedPathsFile, excludedPaths, logger); err != nil {
			logger.Printf("%sWarning: Error saving excluded paths to %s: %v%s\n", colorRed, cfg.excludedPathsFile, err, colorReset)
		}
	}

	// --- Copy to Clipboard ---
	if cfg.addToClipboard {
		if err := clipboard.Init(); err != nil {
			logger.Printf("%sWarning: Could not initialize clipboard: %v%s\n", colorRed, err, colorReset)
		} else {
			clipboard.Write(clipboard.FmtText, []byte(markdownContent.String()))
			logger.Println("Markdown content copied to clipboard.")
		}
	}
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
	showHelp          bool
	showVersion       bool
}

func parseFlags() (*config, error) {
	cfg := &config{}
	flag.StringVar(&cfg.inputDirOriginal, "input", ".", "The root directory to scan.")
	flag.StringVar(&cfg.outputFile, "output", "codebase.md", "The name of the output Markdown file.")
	ignoreStr := flag.String("ignore", `\.git.*`, "Comma-separated list of regular expressions for paths to *exclude* (relative to input directory).")
	includeStr := flag.String("include", "", "Comma-separated list of regular expressions. *Only* paths matching these are *included* (relative to input directory).")
	flag.StringVar(&cfg.includedPathsFile, "included-paths-file", "", "Saves the list of *included* paths to this file.")
	flag.StringVar(&cfg.excludedPathsFile, "excluded-paths-file", "", "Saves the list of *excluded* paths to this file.")
	flag.BoolVar(&cfg.addToClipboard, "clipboard", false, "Copies the generated Markdown to the clipboard.")
	flag.BoolVar(&cfg.showVersion, "version", false, "Displays the version and exits.")
	flag.BoolVar(&cfg.showHelp, "help", false, "Displays help message and exits.")

	flag.Usage = printHelp // Override default usage
	flag.Parse()

	if cfg.showHelp || cfg.showVersion {
		// Let main handle printing help/version
		return cfg, nil
	}

	var err error
	cfg.inputDirAbs, err = filepath.Abs(cfg.inputDirOriginal)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path for input directory '%s': %w", cfg.inputDirOriginal, err)
	}

	// Check if input directory exists and is a directory
	info, err := os.Stat(cfg.inputDirAbs)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("input directory '%s' does not exist", cfg.inputDirAbs)
		}
		return nil, fmt.Errorf("error accessing input directory '%s': %w", cfg.inputDirAbs, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("input path '%s' is not a directory", cfg.inputDirAbs)
	}

	if *ignoreStr != "" {
		cfg.ignorePatterns = strings.Split(*ignoreStr, ",")
	}
	if *includeStr != "" {
		cfg.includePatterns = strings.Split(*includeStr, ",")
	}

	return cfg, nil
}

func compileRegexPatterns(patterns []string, color, prefix string, logger *log.Logger) ([]*regexp.Regexp, error) {
	if len(patterns) == 0 {
		return nil, nil
	}
	// Pre-allocate close to the expected size, but use append for flexibility
	compiledMatchers := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		trimmedPattern := strings.TrimSpace(p)
		if trimmedPattern == "" {
			continue // Skip empty patterns that might result from trailing commas
		}
		logger.Printf("%s%s %s%s\n", color, prefix, trimmedPattern, colorReset)
		rgx, err := regexp.Compile(trimmedPattern)
		if err != nil {
			// Fail completely if any pattern is invalid
			return nil, fmt.Errorf("invalid regex pattern '%s': %w", trimmedPattern, err)
		}
		compiledMatchers = append(compiledMatchers, rgx) // Append the valid regex
	}
	// Return nil if no valid patterns were actually found after trimming/skipping
	if len(compiledMatchers) == 0 {
		return nil, nil
	}
	return compiledMatchers, nil
}

// shouldProcess determines if a path should be processed based on include and ignore patterns.
// The path argument must be relative to the input directory and use forward slashes.
func shouldProcess(pathRelToInput string, ignoreMatchers, includeMatchers []*regexp.Regexp) bool {
	// Check exclusion first
	for _, pattern := range ignoreMatchers {
		if pattern != nil && pattern.MatchString(pathRelToInput) {
			return false // Excluded
		}
	}

	// If include patterns are defined, path must match at least one
	if len(includeMatchers) > 0 {
		matchedInclude := false
		for _, pattern := range includeMatchers {
			if pattern != nil && pattern.MatchString(pathRelToInput) {
				matchedInclude = true
				break
			}
		}
		if !matchedInclude {
			return false // Not in include list
		}
	}

	return true // Included (or not excluded if no include list)
}

type treeBuilder struct {
	rootAbsPath     string
	ignoreMatchers  []*regexp.Regexp
	includeMatchers []*regexp.Regexp
	output          strings.Builder
	depthOpen       map[int]bool // Tracks open branches for │ character
}

func newTreeBuilder(rootAbsPath string, ignoreMatchers, includeMatchers []*regexp.Regexp) *treeBuilder {
	return &treeBuilder{
		rootAbsPath:     rootAbsPath,
		ignoreMatchers:  ignoreMatchers,
		includeMatchers: includeMatchers,
		depthOpen:       make(map[int]bool),
	}
}

func (tb *treeBuilder) buildTreeString() (string, error) {
	err := tb.printTreeRecursive(tb.rootAbsPath, 0)
	return tb.output.String(), err
}

func (tb *treeBuilder) printTreeRecursive(currentDirPath string, depth int) error {
	entries, err := os.ReadDir(currentDirPath)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", currentDirPath, err)
	}

	var filteredEntries []fs.DirEntry
	for _, entry := range entries {
		fullEntryPath := filepath.Join(currentDirPath, entry.Name())
		pathRelToInput, err := filepath.Rel(tb.rootAbsPath, fullEntryPath)
		if err != nil {
			return fmt.Errorf("failed to make path %s relative to %s: %w", fullEntryPath, tb.rootAbsPath, err)
		}
		pathRelToInput = filepath.ToSlash(pathRelToInput)

		if shouldProcess(pathRelToInput, tb.ignoreMatchers, tb.includeMatchers) {
			filteredEntries = append(filteredEntries, entry)
		}
	}

	// Sort entries alphabetically for consistent output
	sort.Slice(filteredEntries, func(i, j int) bool {
		return filteredEntries[i].Name() < filteredEntries[j].Name()
	})

	for i, entry := range filteredEntries {
		isLastEntry := (i == len(filteredEntries)-1)
		tb.printEntry(entry, depth, isLastEntry)

		if entry.IsDir() {
			tb.depthOpen[depth] = !isLastEntry
			err := tb.printTreeRecursive(filepath.Join(currentDirPath, entry.Name()), depth+1)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (tb *treeBuilder) printEntry(entry fs.DirEntry, depth int, isLast bool) {
	var prefix strings.Builder
	for i := 0; i < depth; i++ {
		if tb.depthOpen[i] {
			prefix.WriteString("│  ")
		} else {
			prefix.WriteString("   ") // Was "  " - need three spaces to align with "└─ " or "├─ "
		}
	}

	if isLast {
		prefix.WriteString("└─ ")
	} else {
		prefix.WriteString("├─ ")
	}

	tb.output.WriteString(prefix.String())
	tb.output.WriteString(entry.Name())
	tb.output.WriteString("\n")
}

type contentBuilder struct {
	rootAbsPath       string
	includedPathsFile string
	excludedPathsFile string
	ignoreMatchers    []*regexp.Regexp
	includeMatchers   []*regexp.Regexp
	logger            *log.Logger
}

func newContentBuilder(
	rootAbsPath string, includedPathsFile string, excludedPathsFile string, ignoreMatchers []*regexp.Regexp, includeMatchers []*regexp.Regexp, logger *log.Logger) *contentBuilder {
	return &contentBuilder{
		rootAbsPath:       rootAbsPath,
		includedPathsFile: includedPathsFile,
		excludedPathsFile: excludedPathsFile,
		ignoreMatchers:    ignoreMatchers,
		includeMatchers:   includeMatchers,
		logger:            logger,
	}
}

func (cb *contentBuilder) buildContentString() (string, []string, []string, error) {
	var content strings.Builder
	var includedPaths []string
	var excludedPaths []string

	err := filepath.WalkDir(cb.rootAbsPath, func(currentWalkPath string, d fs.DirEntry, err error) error {
		if err != nil {
			// If WalkDir encounters an error accessing a path (like permission denied or root not existing),
			// report it and return the error to stop the walk unless it's skippable.
			cb.logger.Printf("%sWarning: Error accessing %s: %v%s\n", colorRed, currentWalkPath, err, colorReset)
			// Allow skipping directories if the error is related to reading its contents,
			// but critical errors (like the starting path itself not existing) should halt.
			// Returning the error is generally safer.
			return err
		}

		pathRelToInput, relErr := filepath.Rel(cb.rootAbsPath, currentWalkPath)
		if relErr != nil {
			// This should ideally not happen if WalkDir starts from rootAbsPath
			cb.logger.Printf("%sWarning: Could not make path %s relative to %s: %v%s\n", colorRed, currentWalkPath, cb.rootAbsPath, relErr, colorReset)
			return nil // Skip this entry
		}
		pathRelToInput = filepath.ToSlash(pathRelToInput)

		// Don't process the root directory itself as a "file" or for primary filtering here;
		// WalkDir handles recursion into it. Filtering applies to its children.
		if pathRelToInput == "." {
			return nil // Continue walking
		}

		if !shouldProcess(pathRelToInput, cb.ignoreMatchers, cb.includeMatchers) {
			if cb.excludedPathsFile == "" {
				cb.logger.Printf("%s- %s%s\n", colorRed, pathRelToInput, colorReset)
			}
			excludedPaths = append(excludedPaths, pathRelToInput) // Store relative path
			if d.IsDir() {
				return filepath.SkipDir // Skip entire directory
			}
			return nil // Skip this file
		}

		// If we reach here, the path is included.
		if !d.IsDir() {
			if cb.includedPathsFile == "" {
				cb.logger.Printf("%s+ %s%s\n", colorGreen, pathRelToInput, colorReset)
			}
			includedPaths = append(includedPaths, pathRelToInput) // Store relative path

			fileContent, readErr := os.ReadFile(currentWalkPath)
			if readErr != nil {
				cb.logger.Printf("%sWarning: Failed to read file %s: %v%s\n", colorRed, currentWalkPath, readErr, colorReset)
				// Optionally, add a placeholder to the markdown for unreadable files
				content.WriteString(fmt.Sprintf("## %s\n", pathRelToInput))
				content.WriteString(fmt.Sprintf("```\nError reading file: %v\n```\n\n", readErr))
				return nil // Continue with next file
			}

			extension := strings.TrimPrefix(strings.ToLower(filepath.Ext(currentWalkPath)), ".")
			content.WriteString(fmt.Sprintf("## %s\n", pathRelToInput)) // Use relative path for header
			content.WriteString(fmt.Sprintf("```%s\n", extension))
			content.Write(fileContent) // Write bytes directly
			content.WriteString("\n```\n\n")
		}
		return nil
	})

	if err != nil {
		return "", nil, nil, fmt.Errorf("error walking directory %s: %w", cb.rootAbsPath, err)
	}
	return content.String(), includedPaths, excludedPaths, nil
}

func savePathsToFile(filename string, paths []string, logger *log.Logger) error {
	if len(paths) == 0 {
		// Optionally, create an empty file or just skip
		// logger.Printf("No paths to save to %s.", filename)
		// return os.WriteFile(filename, []byte{}, 0644)
		return nil // Do nothing if no paths
	}

	var sb strings.Builder
	for _, p := range paths {
		sb.WriteString(p)
		sb.WriteString("\n")
	}

	err := os.WriteFile(filename, []byte(sb.String()), 0644)
	if err == nil {
		logger.Printf("Paths saved to %s\n", filename)
	}
	return err
}

func printHelp() {
	fmt.Println("CodeWeaver: Generate Markdown Documentation from Your Codebase.")
	fmt.Printf("Version: %s, Commit: %s, Date: %s\n\n", version, commit, date)
	fmt.Println("Usage: codeweaver [options]")
	fmt.Println("\nOptions:")
	flag.PrintDefaults()
	fmt.Println("\nExamples:")
	fmt.Println("  codeweaver                               # Process current directory, output to codebase.md")
	fmt.Println("  codeweaver -input my_project -output docs.md")
	fmt.Println(`  codeweaver -ignore "build/,vendor/" -include "\.go$,\.md$"`)
	fmt.Println("  codeweaver -clipboard -excluded-paths-file ignored.txt")
	fmt.Println("\nNotes on patterns:")
	fmt.Println("  - Patterns are Go regular expressions.")
	fmt.Println("  - Paths for filtering are relative to the input directory (e.g., \"src/main.go\", not \"./src/main.go\" or \"/path/to/project/src/main.go\").")
	fmt.Println("  - Use forward slashes '/' in patterns for cross-platform compatibility (e.g., \"data/images/\").")
}

```

## main_test.go
```go
package main

import (
	"flag" // For testing parseFlags
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// --- Helpers (mustCompileRegex, normalizeNewlines, createTestFS, equalStringSlices) remain the same ---

func mustCompileRegex(pattern string) *regexp.Regexp {
	if pattern == "" {
		return nil
	}
	r, err := regexp.Compile(pattern)
	if err != nil {
		log.Fatalf("Failed to compile regex '%s': %v", pattern, err)
	}
	return r
}

func normalizeNewlines(s string) string {
	return strings.ReplaceAll(s, "\r\n", "\n")
}

func createTestFS(t *testing.T) (string, func()) {
	t.Helper()
	rootDir, err := os.MkdirTemp("", "codeweaver_test_fs_")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Define structure: map[relativePath]content ("" content for directory)
	structure := map[string]string{
		"file1.txt":                     "content of file1",
		"script.go":                     "package main\nfunc main() {}",
		"README.md":                     "# Test Readme",
		"data/":                         "", // Directory marker
		"data/image.png":                "fake png data",
		"data/config.yaml":              "key: value",
		"build/":                        "", // Directory marker
		"build/output.exe":              "binary data",
		"build/tmp/":                    "", // Directory marker
		"build/tmp/log.txt":             "log entry",
		".git/":                         "", // Directory marker
		".git/HEAD":                     "ref: refs/heads/main",
		"node_modules/":                 "", // Directory marker
		"node_modules/dep/":             "", // Directory marker
		"node_modules/dep/package.json": "{}",
		"empty_dir/":                    "", // Explicitly empty directory
		"docs/sub_docs/file_in_sub.txt": "nested doc content",
		"other.log":                     "another log",
	}

	for relPath, content := range structure {
		absPath := filepath.Join(rootDir, relPath)
		if strings.HasSuffix(relPath, "/") || content == "" && !strings.Contains(relPath, ".") { // Treat as directory if ends with / OR has no extension/dot and empty content
			if err := os.MkdirAll(absPath, 0755); err != nil {
				os.RemoveAll(rootDir) // Cleanup on error
				t.Fatalf("Failed to create dir %s: %v", absPath, err)
			}
		} else { // It's a file
			parentDir := filepath.Dir(absPath)
			if err := os.MkdirAll(parentDir, 0755); err != nil {
				os.RemoveAll(rootDir)
				t.Fatalf("Failed to create parent dir %s for file %s: %v", parentDir, absPath, err)
			}
			if err := os.WriteFile(absPath, []byte(content), 0644); err != nil {
				os.RemoveAll(rootDir)
				t.Fatalf("Failed to write file %s: %v", absPath, err)
			}
		}
	}

	cleanup := func() {
		os.RemoveAll(rootDir)
	}
	return rootDir, cleanup
}

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	// Assumes slices are sorted
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// --- Test Suite ---

func TestShouldProcess(t *testing.T) {
	// --- Test cases from previous version (verified) ---
	testCases := []struct {
		name            string
		path            string // Path relative to input dir, using forward slashes
		ignoreMatchers  []*regexp.Regexp
		includeMatchers []*regexp.Regexp
		expected        bool
	}{
		{"NoFilters_Allow", "file.txt", nil, nil, true},
		{"EmptyIgnoreList", "file.txt", []*regexp.Regexp{}, nil, true},
		{"EmptyIncludeList_MeansAllow", "file.txt", nil, []*regexp.Regexp{}, true},
		{"IgnoreMatch_Exact", "skip.txt", []*regexp.Regexp{mustCompileRegex(`^skip\.txt$`)}, nil, false},
		{"IgnoreMatch_DirPrefix", "skip/file.txt", []*regexp.Regexp{mustCompileRegex(`^skip/`)}, nil, false},
		{"IgnoreNoMatch", "keep/file.txt", []*regexp.Regexp{mustCompileRegex(`^skip/`)}, nil, true},
		{"IncludeMatch_Extension", "src/main.go", nil, []*regexp.Regexp{mustCompileRegex(`\.go$`)}, true},
		{"IncludeNoMatch_Extension", "src/main.txt", nil, []*regexp.Regexp{mustCompileRegex(`\.go$`)}, false},
		{"IncludeDir_Exact", "data/", nil, []*regexp.Regexp{mustCompileRegex(`^data/?$`)}, true}, // Note: ? for optional trailing slash
		{"IncludeDir_NoMatch", "other/", nil, []*regexp.Regexp{mustCompileRegex(`^data/?$`)}, false},
		{"IgnoreTakesPrecedence_PathMatchesBoth", "vendor/lib.go", []*regexp.Regexp{mustCompileRegex(`^vendor/`)}, []*regexp.Regexp{mustCompileRegex(`\.go$`)}, false},
		{"IncludeSatisfied_NotIgnored", "src/main.go", []*regexp.Regexp{mustCompileRegex(`^vendor/`)}, []*regexp.Regexp{mustCompileRegex(`\.go$`)}, true},
		{"MultipleIgnore_FirstMatch", "config/dev.log", []*regexp.Regexp{mustCompileRegex(`\.log$`), mustCompileRegex(`^config/`)}, nil, false},
		{"MultipleIgnore_SecondMatch", "build/app.exe", []*regexp.Regexp{mustCompileRegex(`\.log$`), mustCompileRegex(`\.exe$`)}, nil, false},
		{"MultipleInclude_OneMatchSufficient", "image.png", nil, []*regexp.Regexp{mustCompileRegex(`\.jpg$`), mustCompileRegex(`\.png$`)}, true},
		{"MultipleInclude_NoMatch", "archive.zip", nil, []*regexp.Regexp{mustCompileRegex(`\.jpg$`), mustCompileRegex(`\.png$`)}, false},
		{"PathIsEmptyString_Root_NoFilters", "", nil, nil, true}, // Special case handled elsewhere, but func should allow
		{"PathIsEmptyString_Root_Ignored", "", []*regexp.Regexp{mustCompileRegex(`^$`)}, nil, false},
		{"PathIsEmptyString_Root_Included", "", nil, []*regexp.Regexp{mustCompileRegex(`^$`)}, true},
		{"PathIsEmptyString_Root_NotIncluded", "", nil, []*regexp.Regexp{mustCompileRegex(`\.txt$`)}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := shouldProcess(tc.path, tc.ignoreMatchers, tc.includeMatchers)
			if actual != tc.expected {
				t.Errorf("shouldProcess(%q) with ignore/include patterns = %v; want %v", tc.path, actual, tc.expected)
			}
		})
	}
}

// Test parseFlags by simulating command-line arguments
func TestParseFlags(t *testing.T) {
	// Helper to run parseFlags with specific args
	runParse := func(args []string) (*config, error) {
		// Create a new flag set for testing to avoid interfering with global state
		testFlags := flag.NewFlagSet("test", flag.ContinueOnError)
		testFlags.SetOutput(io.Discard) // Suppress flag errors during testing

		cfg := &config{}
		originalArgs := os.Args                   // Backup original args
		defer func() { os.Args = originalArgs }() // Restore original args

		// Define flags on the test set
		testFlags.StringVar(&cfg.inputDirOriginal, "input", ".", "The root directory to scan.")
		testFlags.StringVar(&cfg.outputFile, "output", "codebase.md", "The name of the output Markdown file.")
		ignoreStr := testFlags.String("ignore", `\.git.*`, "Comma-separated list of regex patterns to exclude.")
		includeStr := testFlags.String("include", "", "Comma-separated list of regex patterns to include.")
		testFlags.StringVar(&cfg.includedPathsFile, "included-paths-file", "", "File to save included paths.")
		testFlags.StringVar(&cfg.excludedPathsFile, "excluded-paths-file", "", "File to save excluded paths.")
		testFlags.BoolVar(&cfg.addToClipboard, "clipboard", false, "Copies the generated Markdown to the clipboard.")
		testFlags.BoolVar(&cfg.showVersion, "version", false, "Displays the version and exits.")
		testFlags.BoolVar(&cfg.showHelp, "help", false, "Displays help message and exits.")

		// Simulate os.Args for parsing
		os.Args = append([]string{"cmd"}, args...)
		testFlags.Parse(args) // Parse the simulated args

		// Manually perform the logic from parseFlags that happens *after* testFlags.Parse()
		if cfg.showHelp || cfg.showVersion {
			return cfg, nil
		}

		var err error
		cfg.inputDirAbs, err = filepath.Abs(cfg.inputDirOriginal)
		if err != nil {
			return nil, fmt.Errorf("failed to get absolute path for input directory '%s': %w", cfg.inputDirOriginal, err)
		}
		info, err := os.Stat(cfg.inputDirAbs)
		if err != nil {
			if os.IsNotExist(err) {
				return nil, fmt.Errorf("input directory '%s' does not exist", cfg.inputDirAbs)
			}
			return nil, fmt.Errorf("error accessing input directory '%s': %w", cfg.inputDirAbs, err)
		}
		if !info.IsDir() {
			return nil, fmt.Errorf("input path '%s' is not a directory", cfg.inputDirAbs)
		}
		if *ignoreStr != "" {
			cfg.ignorePatterns = strings.Split(*ignoreStr, ",")
		}
		if *includeStr != "" {
			cfg.includePatterns = strings.Split(*includeStr, ",")
		}

		return cfg, nil
	}

	// Test Cases
	t.Run("Defaults", func(t *testing.T) {
		cfg, err := runParse([]string{})
		if err != nil {
			t.Fatalf("Expected no error for defaults, got %v", err)
		}
		cwd, _ := os.Getwd()
		absCwd, _ := filepath.Abs(cwd)
		if cfg.inputDirAbs != absCwd {
			t.Errorf("Default inputDirAbs mismatch: got %s, want %s", cfg.inputDirAbs, absCwd)
		}
		if cfg.outputFile != "codebase.md" {
			t.Errorf("Default outputFile mismatch: got %s", cfg.outputFile)
		}
		if len(cfg.ignorePatterns) != 1 || cfg.ignorePatterns[0] != `\.git.*` {
			t.Errorf("Default ignorePatterns mismatch: got %v", cfg.ignorePatterns)
		}
		if len(cfg.includePatterns) != 0 {
			t.Errorf("Default includePatterns should be empty, got %v", cfg.includePatterns)
		}
		if cfg.addToClipboard != false {
			t.Errorf("Default clipboard mismatch")
		}
	})

	t.Run("SetValues", func(t *testing.T) {
		// Use a known existing directory for input
		tmpDir, cleanup := createTestFS(t)
		defer cleanup()

		args := []string{
			"-input", tmpDir,
			"-output", "out.md",
			"-ignore", "a,b",
			"-include", "c,d",
			"-included-paths-file", "inc.txt",
			"-excluded-paths-file", "exc.txt",
			"-clipboard",
		}
		cfg, err := runParse(args)
		if err != nil {
			t.Fatalf("Expected no error setting values, got %v", err)
		}
		absTmpDir, _ := filepath.Abs(tmpDir)
		if cfg.inputDirAbs != absTmpDir {
			t.Errorf("Set inputDirAbs mismatch: got %s, want %s", cfg.inputDirAbs, absTmpDir)
		}
		if cfg.outputFile != "out.md" {
			t.Errorf("Set outputFile mismatch")
		}
		if !equalStringSlices(cfg.ignorePatterns, []string{"a", "b"}) {
			t.Errorf("Set ignorePatterns mismatch")
		}
		if !equalStringSlices(cfg.includePatterns, []string{"c", "d"}) {
			t.Errorf("Set includePatterns mismatch")
		}
		if cfg.includedPathsFile != "inc.txt" {
			t.Errorf("Set includedPathsFile mismatch")
		}
		if cfg.excludedPathsFile != "exc.txt" {
			t.Errorf("Set excludedPathsFile mismatch")
		}
		if cfg.addToClipboard != true {
			t.Errorf("Set clipboard mismatch")
		}
	})

	t.Run("HelpFlag", func(t *testing.T) {
		cfg, err := runParse([]string{"-help"})
		if err != nil {
			t.Fatalf("Expected no error for -help, got %v", err)
		}
		if !cfg.showHelp {
			t.Errorf("Expected showHelp to be true")
		}
	})

	t.Run("VersionFlag", func(t *testing.T) {
		cfg, err := runParse([]string{"-version"})
		if err != nil {
			t.Fatalf("Expected no error for -version, got %v", err)
		}
		if !cfg.showVersion {
			t.Errorf("Expected showVersion to be true")
		}
	})

	t.Run("Error_InputNotExist", func(t *testing.T) {
		nonExistentPath := filepath.Join(os.TempDir(), "codeweaver_non_existent_dir_abc123")
		os.Remove(nonExistentPath) // Ensure it doesn't exist
		_, err := runParse([]string{"-input", nonExistentPath})
		if err == nil {
			t.Fatalf("Expected error for non-existent input path, got nil")
		}
		if !strings.Contains(err.Error(), "does not exist") {
			t.Errorf("Expected 'does not exist' error, got: %v", err)
		}
	})

	t.Run("Error_InputIsFile", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "codeweaver_test_file_*.txt")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		filePath := tmpFile.Name()
		tmpFile.Close()
		defer os.Remove(filePath)

		_, err = runParse([]string{"-input", filePath})
		if err == nil {
			t.Fatalf("Expected error for file input path, got nil")
		}
		if !strings.Contains(err.Error(), "is not a directory") {
			t.Errorf("Expected 'is not a directory' error, got: %v", err)
		}
	})
}

func TestCompileRegexPatterns(t *testing.T) {
	testLogger := log.New(io.Discard, "", 0)

	t.Run("ValidPatterns", func(t *testing.T) {
		patterns := []string{"\\.go$", "^src/"}
		matchers, err := compileRegexPatterns(patterns, "", "", testLogger)
		if err != nil {
			t.Fatalf("Expected no error for valid patterns, got %v", err)
		}
		if len(matchers) != 2 {
			t.Fatalf("Expected 2 matchers, got %d", len(matchers))
		}
		if !matchers[0].MatchString("main.go") {
			t.Error("Matcher 0 failed")
		}
		if !matchers[1].MatchString("src/main.go") {
			t.Error("Matcher 1 failed")
		}
	})

	t.Run("InvalidPattern", func(t *testing.T) {
		patterns := []string{"\\.go$", "["} // Invalid regex
		_, err := compileRegexPatterns(patterns, "", "", testLogger)
		if err == nil {
			t.Fatal("Expected error for invalid pattern, got nil")
		}
		if !strings.Contains(err.Error(), "invalid regex pattern") {
			t.Errorf("Expected 'invalid regex pattern' error, got: %v", err)
		}
	})

	t.Run("EmptyInput", func(t *testing.T) {
		patterns := []string{}
		matchers, err := compileRegexPatterns(patterns, "", "", testLogger)
		if err != nil {
			t.Fatalf("Expected no error for empty input, got %v", err)
		}
		if matchers != nil {
			t.Errorf("Expected nil matchers for empty input, got %v", matchers)
		}
	})

	t.Run("MixValidAndEmptyStrings", func(t *testing.T) {
		// From comma-separated string: "  \\.go$  , , ^src/  "
		patterns := []string{"  \\.go$  ", "", " ^src/  "}
		matchers, err := compileRegexPatterns(patterns, "", "", testLogger)
		if err != nil {
			t.Fatalf("Expected no error for mix, got %v", err)
		}
		// Should skip the empty one ""
		count := 0
		for _, m := range matchers {
			if m != nil {
				count++
			}
		}
		if count != 2 {
			t.Fatalf("Expected 2 non-nil matchers, got %d from %v", count, matchers)
		}
	})
}

func TestTreeBuilder(t *testing.T) {
	// --- Existing Test Cases from previous version (verified & updated) ---
	rootDir, cleanup := createTestFS(t)
	defer cleanup()

	testCases := []struct {
		name            string
		ignoreMatchers  []*regexp.Regexp
		includeMatchers []*regexp.Regexp
		expectedTree    string
	}{
		// ... (use the corrected expectedTree values from the previous response) ...
		{
			name: "NoFilters",
			expectedTree: normalizeNewlines(`├─ .git
│  └─ HEAD
├─ README.md
├─ build
│  ├─ output.exe
│  └─ tmp
│     └─ log.txt
├─ data
│  ├─ config.yaml
│  └─ image.png
├─ docs
│  └─ sub_docs
│     └─ file_in_sub.txt
├─ empty_dir
├─ file1.txt
├─ node_modules
│  └─ dep
│     └─ package.json
├─ other.log
└─ script.go
`),
		},
		{
			name:           "IgnoreDotGitAndNodeModulesDirs",
			ignoreMatchers: []*regexp.Regexp{mustCompileRegex(`^\.git/?`), mustCompileRegex(`^node_modules/?`)},
			expectedTree: normalizeNewlines(`├─ README.md
├─ build
│  ├─ output.exe
│  └─ tmp
│     └─ log.txt
├─ data
│  ├─ config.yaml
│  └─ image.png
├─ docs
│  └─ sub_docs
│     └─ file_in_sub.txt
├─ empty_dir
├─ file1.txt
├─ other.log
└─ script.go
`), // Added other.log which wasn't ignored
		},
		{
			name:            "IncludeOnlyGoAndMdFiles",
			includeMatchers: []*regexp.Regexp{mustCompileRegex(`\.go$`), mustCompileRegex(`\.md$`)},
			expectedTree: normalizeNewlines(`├─ README.md
└─ script.go
`),
		},
		{
			name:            "IncludeYamlFilesAndDataDir",
			includeMatchers: []*regexp.Regexp{mustCompileRegex(`\.yaml$`), mustCompileRegex(`^data/?`)},
			expectedTree: normalizeNewlines(`└─ data
   ├─ config.yaml
   └─ image.png
`),
		},
		{
			name:            "IgnoreExe_IncludeDocsDirAndSubDocsFile",
			ignoreMatchers:  []*regexp.Regexp{mustCompileRegex(`\.exe$`)},
			includeMatchers: []*regexp.Regexp{mustCompileRegex(`^docs(/.*)?`), mustCompileRegex(`file_in_sub\.txt$`)},
			expectedTree: normalizeNewlines(`└─ docs
   └─ sub_docs
      └─ file_in_sub.txt
`),
		},
		{
			name:            "EmptyDir_Included",
			includeMatchers: []*regexp.Regexp{mustCompileRegex(`^empty_dir/?`)},
			expectedTree: normalizeNewlines(`└─ empty_dir
`),
		},
		{
			name:           "EmptyResult_IgnoreAll",
			ignoreMatchers: []*regexp.Regexp{mustCompileRegex(".")},
			expectedTree:   ``,
		},
		{
			name:            "IncludeNonExistentPattern",
			includeMatchers: []*regexp.Regexp{mustCompileRegex(`non_existent_pattern`)},
			expectedTree:    ``,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := newTreeBuilder(rootDir, tc.ignoreMatchers, tc.includeMatchers)
			actualTree, err := builder.buildTreeString()
			if err != nil {
				t.Fatalf("buildTreeString() failed: %v", err)
			}
			actualTree = normalizeNewlines(actualTree)

			if actualTree != tc.expectedTree {
				t.Errorf("Tree mismatch for '%s':\nExpected:\n%s\nGot:\n%s", tc.name, tc.expectedTree, actualTree)
			}
		})
	}

	// --- Error Case ---
	t.Run("Error_InputNotExist", func(t *testing.T) {
		nonExistentPath := filepath.Join(os.TempDir(), "codeweaver_non_existent_dir_tree_abc123")
		os.Remove(nonExistentPath) // Ensure non-existence
		builder := newTreeBuilder(nonExistentPath, nil, nil)
		_, err := builder.buildTreeString() // Should fail when trying to ReadDir
		if err == nil {
			t.Fatalf("Expected error for non-existent input path, got nil")
		}
		// Check if the error is about reading the directory (os-specific message check is brittle)
		if !strings.Contains(err.Error(), "failed to read directory") {
			t.Errorf("Expected error related to reading directory, got: %v", err)
		}
	})

	// --- Empty Directory Case ---
	t.Run("EmptyInputDir", func(t *testing.T) {
		emptyDir, cleanupEmpty := createTestFS(t) // Create our standard structure
		defer cleanupEmpty()
		// Create an actual empty dir *within* the temp structure to use as input
		actualEmptyDir := filepath.Join(emptyDir, "really_empty")
		err := os.Mkdir(actualEmptyDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create empty test dir: %v", err)
		}

		builder := newTreeBuilder(actualEmptyDir, nil, nil)
		actualTree, err := builder.buildTreeString()
		if err != nil {
			t.Fatalf("buildTreeString() failed for empty dir: %v", err)
		}
		if actualTree != "" {
			t.Errorf("Expected empty tree output for empty dir, got: %s", actualTree)
		}
	})
}

func TestContentBuilder(t *testing.T) {
	// --- Existing Test Cases from previous version (verified & updated) ---
	rootDir, cleanup := createTestFS(t)
	defer cleanup()

	testLogger := log.New(io.Discard, "", 0) // Suppress log output during tests

	testCases := []struct {
		name                  string
		ignoreMatchers        []*regexp.Regexp
		includeMatchers       []*regexp.Regexp
		expectedContentSubstr string
		expectedIncludedPaths []string // Sorted
		expectedExcludedPaths []string // Sorted
	}{
		// ... (use the corrected test cases from the previous response) ...
		{
			name:                  "NoFilters",
			expectedContentSubstr: "## file1.txt\n```txt\ncontent of file1\n```",
			expectedIncludedPaths: []string{
				".git/HEAD", "README.md", "build/output.exe", "build/tmp/log.txt",
				"data/config.yaml", "data/image.png", "docs/sub_docs/file_in_sub.txt",
				"file1.txt", "node_modules/dep/package.json", "other.log", "script.go", // Added other.log
			},
			expectedExcludedPaths: []string{},
		},
		{
			name:                  "IgnoreDotGitDirAndExeFiles",
			ignoreMatchers:        []*regexp.Regexp{mustCompileRegex(`^\.git/?`), mustCompileRegex(`\.exe$`)},
			expectedContentSubstr: "## README.md\n```md\n# Test Readme\n```",
			expectedIncludedPaths: []string{
				"README.md", "build/tmp/log.txt", "data/config.yaml", "data/image.png",
				"docs/sub_docs/file_in_sub.txt", "file1.txt", "node_modules/dep/package.json",
				"other.log", "script.go", // Added other.log
			},
			expectedExcludedPaths: []string{".git", "build/output.exe"},
		},
		{
			name:                  "IncludeOnlyGoAndMdFiles",
			includeMatchers:       []*regexp.Regexp{mustCompileRegex(`\.go$`), mustCompileRegex(`\.md$`)},
			expectedContentSubstr: "## script.go\n```go\npackage main",
			expectedIncludedPaths: []string{"README.md", "script.go"},
			expectedExcludedPaths: []string{
				".git", "build", "data", "docs", "empty_dir", "file1.txt", "node_modules", "other.log", // Added other.log
			},
		},
		{
			name:                  "ContentOfSpecificFile_Yaml_StrictInclude",
			includeMatchers:       []*regexp.Regexp{mustCompileRegex(`config\.yaml$`)},
			expectedContentSubstr: "",
			expectedIncludedPaths: []string{},
			expectedExcludedPaths: []string{
				".git", "README.md", "build", "data", "docs", "empty_dir", "file1.txt", "node_modules", "other.log", "script.go", // Added other.log
			},
		},
		{
			name:                  "ContentOfSpecificFile_Yaml_PermissiveIncludeDir",
			includeMatchers:       []*regexp.Regexp{mustCompileRegex(`config\.yaml$`), mustCompileRegex(`^data/?`)},
			expectedContentSubstr: normalizeNewlines("## data/config.yaml\n```yaml\nkey: value\n```\n\n"),
			expectedIncludedPaths: []string{"data/config.yaml", "data/image.png"},
			expectedExcludedPaths: []string{
				".git", "README.md", "build", "docs", "empty_dir", "file1.txt", "node_modules", "other.log", "script.go", // Added other.log
			},
		},
		{
			name:                  "EmptyDir_NotIncludedInContent",
			includeMatchers:       []*regexp.Regexp{mustCompileRegex(`^empty_dir/?`)},
			expectedContentSubstr: "",
			expectedIncludedPaths: []string{},
			expectedExcludedPaths: []string{
				".git", "README.md", "build", "data", "docs", "file1.txt", "node_modules", "other.log", "script.go",
			},
		},
		// --- New Case: Ignore takes precedence over include ---
		{
			name:                  "IgnoreLog_IncludeAll",
			ignoreMatchers:        []*regexp.Regexp{mustCompileRegex(`\.log$`)}, // Ensure this is the original pattern
			includeMatchers:       []*regexp.Regexp{mustCompileRegex(".")},      // Include everything not ignored
			expectedContentSubstr: "## README.md",                               // Check some non-log file exists
			expectedIncludedPaths: []string{ // All except other.log
				".git/HEAD", "README.md", "build/output.exe", "build/tmp/log.txt", // <-- CORRECT: build/tmp/log.txt IS included
				"data/config.yaml", "data/image.png", "docs/sub_docs/file_in_sub.txt",
				"file1.txt", "node_modules/dep/package.json", "script.go",
			},
			expectedExcludedPaths: []string{"other.log"}, // <-- CORRECT: Only other.log is excluded by the pattern
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sort.Strings(tc.expectedIncludedPaths)
			sort.Strings(tc.expectedExcludedPaths)

			builder := newContentBuilder(rootDir, "", "", tc.ignoreMatchers, tc.includeMatchers, testLogger)
			actualContentStr, actualIncludedPaths, actualExcludedPaths, err := builder.buildContentString()
			if err != nil {
				t.Fatalf("buildContentString() failed for '%s': %v", tc.name, err)
			}
			actualContentStr = normalizeNewlines(actualContentStr)
			sort.Strings(actualIncludedPaths)
			sort.Strings(actualExcludedPaths)

			if tc.expectedContentSubstr != "" && !strings.Contains(actualContentStr, tc.expectedContentSubstr) {
				t.Errorf("'%s': Generated content does not contain expected substring.\nExpected to find:\n%s\n-----\nActual Content:\n%s-----", tc.name, tc.expectedContentSubstr, actualContentStr)
			}
			if tc.expectedContentSubstr == "" && actualContentStr != "" {
				t.Errorf("'%s': Expected empty content, but got content:\n%s", tc.name, actualContentStr)
			}
			if !equalStringSlices(actualIncludedPaths, tc.expectedIncludedPaths) {
				t.Errorf("'%s': Included paths mismatch.\nExpected: %v\nGot:      %v", tc.name, tc.expectedIncludedPaths, actualIncludedPaths)
			}
			if !equalStringSlices(actualExcludedPaths, tc.expectedExcludedPaths) {
				t.Errorf("'%s': Excluded paths mismatch.\nExpected: %v\nGot:      %v", tc.name, tc.expectedExcludedPaths, actualExcludedPaths)
			}
		})
	}

	// --- Error Case ---
	t.Run("Error_InputNotExist", func(t *testing.T) {
		nonExistentPath := filepath.Join(os.TempDir(), "codeweaver_non_existent_dir_content_abc123")
		os.Remove(nonExistentPath) // Ensure it doesn't exist first
		builder := newContentBuilder(nonExistentPath, "", "", nil, nil, testLogger)
		_, _, _, err := builder.buildContentString() // Should fail on WalkDir
		if err == nil {
			t.Fatal("Expected error for non-existent input path, got nil")
		}
		// Fallback: Check if the error message contains the non-existent path string,
		// as os.ErrNotExist checking seems unreliable across platforms/versions here.
		if !strings.Contains(err.Error(), nonExistentPath) {
			t.Errorf("Expected error message related to path '%s', but got: %v", nonExistentPath, err)
		}
	})

	// --- Test Log Warning on Unreadable File (Difficult to do reliably cross-platform) ---
	// This often requires setting specific permissions which might not work or might
	// require elevated privileges. Skipping direct test for unreadable file warning,
	// but ensuring the error path in WalkDir's func is covered by InputNotExist.

	// --- Test Empty Input Dir ---
	t.Run("EmptyInputDir", func(t *testing.T) {
		emptyDir, cleanupEmpty := createTestFS(t)
		defer cleanupEmpty()
		actualEmptyDir := filepath.Join(emptyDir, "really_empty")
		os.Mkdir(actualEmptyDir, 0755)

		builder := newContentBuilder(actualEmptyDir, "", "", nil, nil, testLogger)
		content, included, excluded, err := builder.buildContentString()
		if err != nil {
			t.Fatalf("Failed processing empty dir: %v", err)
		}
		if content != "" {
			t.Errorf("Expected empty content string, got: %s", content)
		}
		if len(included) != 0 {
			t.Errorf("Expected zero included paths, got: %v", included)
		}
		if len(excluded) != 0 {
			t.Errorf("Expected zero excluded paths, got: %v", excluded)
		}
	})

}

func TestSavePathsToFile(t *testing.T) {
	testLogger := log.New(io.Discard, "", 0)

	t.Run("StandardSave", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "test_paths_*.txt")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		tmpFilePath := tmpFile.Name()
		tmpFile.Close()
		defer os.Remove(tmpFilePath)

		paths := []string{"path/to/file1.txt", "another/path.go", "root_file.md"}
		expectedContent := "path/to/file1.txt\nanother/path.go\nroot_file.md\n"

		err = savePathsToFile(tmpFilePath, paths, testLogger)
		if err != nil {
			t.Fatalf("savePathsToFile failed: %v", err)
		}

		contentBytes, err := os.ReadFile(tmpFilePath)
		if err != nil {
			t.Fatalf("Failed to read back saved paths file: %v", err)
		}
		if normalizeNewlines(string(contentBytes)) != normalizeNewlines(expectedContent) {
			t.Errorf("Content mismatch.\nExpected:\n%s\nGot:\n%s", expectedContent, string(contentBytes))
		}
	})

	t.Run("EmptyPaths", func(t *testing.T) {
		emptyPathsFile := filepath.Join(os.TempDir(), "empty_paths_save_test.txt")
		os.Remove(emptyPathsFile) // Ensure it doesn't exist
		defer os.Remove(emptyPathsFile)

		err := savePathsToFile(emptyPathsFile, []string{}, testLogger)
		if err != nil {
			t.Fatalf("savePathsToFile with empty paths failed: %v", err)
		}
		if _, err := os.Stat(emptyPathsFile); !os.IsNotExist(err) {
			t.Errorf("Expected file not to be created for empty paths, but it was (or other error: %v)", err)
		}
	})

	t.Run("Error_PathIsDirectory", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "test_save_dir_")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		err = savePathsToFile(tmpDir, []string{"a", "b"}, testLogger)
		if err == nil {
			t.Fatalf("Expected error when saving to a directory path, got nil")
		}
	})

	// Note: Testing write permission errors is hard cross-platform.
	// Testing saving to a non-existent parent directory:
	t.Run("Error_ParentDirNotExist", func(t *testing.T) {
		invalidPath := filepath.Join(os.TempDir(), "non_existent_save_parent", "paths.txt")
		os.RemoveAll(filepath.Dir(invalidPath)) // Ensure parent doesn't exist
		defer os.RemoveAll(filepath.Dir(invalidPath))

		err := savePathsToFile(invalidPath, []string{"a", "b"}, testLogger)
		if err == nil {
			t.Fatalf("Expected error when saving to non-existent parent dir, got nil")
		}
	})
}

// Test printHelp by capturing stdout
func TestPrintHelp(t *testing.T) {
	// Capture standard output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() { os.Stdout = oldStdout }() // Restore stdout

	printHelp() // Call the function that prints

	// Close the write end and read from the read end
	w.Close()
	outBytes, _ := io.ReadAll(r)
	output := string(outBytes)

	// Basic check for some expected content
	if !strings.Contains(output, "Usage: codeweaver") {
		t.Error("Help output did not contain 'Usage: codeweaver'")
	}
	if !strings.Contains(output, "Options:") {
		t.Error("Help output did not contain 'Options:'")
	}
	if !strings.Contains(output, "-input") {
		t.Error("Help output did not contain '-input'")
	}
}

```

## run_tests.ps1
```ps1
# go test -cover -v

go test -coverprofile="coverage.out"
go tool cover -func="coverage.out"
go tool cover -html="coverage.out"
```

## test_root/File at root A.txt
```txt

```

## test_root/File at root B.md
```md

```

## test_root/folder 01/File at folder 01 I.txt
```txt

```

## test_root/folder 01/File at folder 01 II.md
```md

```

## test_root/folder 01/File at folder 01 III.csv
```csv

```

## test_root/folder 02/File at folder 02 I.txt
```txt

```

## test_root/folder 02/File at folder 02 II.md
```md

```

## test_root/folder 02/File at folder 02 III.csv
```csv

```

## test_root/folder 02/folder 02 01/File at folder 02 01 I.txt
```txt

```

## test_root/folder 02/folder 02 01/File at folder 02 01 II.md
```md

```

## test_root/folder 02/folder 02 01/File at folder 02 01 III.csv
```csv

```

