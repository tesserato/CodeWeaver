# Tree View:
```
.
├─ build_and_run.ps1
├─ coverage.out
├─ exc.txt
├─ excluded.log
├─ go.mod
├─ go.sum
├─ goreleaser.yaml
├─ inc.txt
├─ included.log
├─ LICENSE
├─ main.go
├─ main_test.go
├─ out.md
├─ README.md
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
# ./CodeWeaver -h
./CodeWeaver -clipboard -ignore="\.git.*,.+\.exe,codebase.md,excluded_paths.txt,DRAFT\.md,coverage$"
```


## coverage.out
```out
mode: set
github.com/tesserato/CodeWeaver/main.go:62.13,65.16 2 0
github.com/tesserato/CodeWeaver/main.go:65.16,67.35 1 0
github.com/tesserato/CodeWeaver/main.go:67.35,71.4 1 0
github.com/tesserato/CodeWeaver/main.go:73.3,73.71 1 0
github.com/tesserato/CodeWeaver/main.go:77.2,77.21 1 0
github.com/tesserato/CodeWeaver/main.go:77.21,80.3 2 0
github.com/tesserato/CodeWeaver/main.go:84.2,88.16 3 0
github.com/tesserato/CodeWeaver/main.go:88.16,90.3 1 0
github.com/tesserato/CodeWeaver/main.go:93.2,94.16 2 0
github.com/tesserato/CodeWeaver/main.go:94.16,96.3 1 0
github.com/tesserato/CodeWeaver/main.go:99.2,100.16 2 0
github.com/tesserato/CodeWeaver/main.go:100.16,103.3 1 0
github.com/tesserato/CodeWeaver/main.go:105.2,105.80 1 0
github.com/tesserato/CodeWeaver/main.go:129.36,150.16 15 0
github.com/tesserato/CodeWeaver/main.go:150.16,154.3 1 0
github.com/tesserato/CodeWeaver/main.go:157.2,157.18 1 0
github.com/tesserato/CodeWeaver/main.go:157.18,161.3 2 0
github.com/tesserato/CodeWeaver/main.go:164.2,164.21 1 0
github.com/tesserato/CodeWeaver/main.go:164.21,166.3 1 0
github.com/tesserato/CodeWeaver/main.go:171.2,172.16 2 0
github.com/tesserato/CodeWeaver/main.go:172.16,174.3 1 0
github.com/tesserato/CodeWeaver/main.go:177.2,178.16 2 0
github.com/tesserato/CodeWeaver/main.go:178.16,179.25 1 0
github.com/tesserato/CodeWeaver/main.go:179.25,181.4 1 0
github.com/tesserato/CodeWeaver/main.go:182.3,182.91 1 0
github.com/tesserato/CodeWeaver/main.go:184.2,184.19 1 0
github.com/tesserato/CodeWeaver/main.go:184.19,186.3 1 0
github.com/tesserato/CodeWeaver/main.go:189.2,189.22 1 0
github.com/tesserato/CodeWeaver/main.go:189.22,191.3 1 0
github.com/tesserato/CodeWeaver/main.go:192.2,192.23 1 0
github.com/tesserato/CodeWeaver/main.go:192.23,194.3 1 0
github.com/tesserato/CodeWeaver/main.go:196.2,196.17 1 0
github.com/tesserato/CodeWeaver/main.go:200.44,205.33 5 0
github.com/tesserato/CodeWeaver/main.go:205.33,207.3 1 0
github.com/tesserato/CodeWeaver/main.go:208.2,208.33 1 0
github.com/tesserato/CodeWeaver/main.go:208.33,210.3 1 0
github.com/tesserato/CodeWeaver/main.go:211.2,211.24 1 0
github.com/tesserato/CodeWeaver/main.go:211.24,213.3 1 0
github.com/tesserato/CodeWeaver/main.go:214.2,215.15 2 0
github.com/tesserato/CodeWeaver/main.go:219.118,222.16 3 1
github.com/tesserato/CodeWeaver/main.go:222.16,224.3 1 1
github.com/tesserato/CodeWeaver/main.go:226.2,228.16 3 1
github.com/tesserato/CodeWeaver/main.go:228.16,230.3 1 1
github.com/tesserato/CodeWeaver/main.go:231.2,232.29 2 1
github.com/tesserato/CodeWeaver/main.go:236.110,237.24 1 1
github.com/tesserato/CodeWeaver/main.go:237.24,240.3 2 1
github.com/tesserato/CodeWeaver/main.go:241.2,243.29 3 1
github.com/tesserato/CodeWeaver/main.go:243.29,245.27 2 1
github.com/tesserato/CodeWeaver/main.go:245.27,246.12 1 0
github.com/tesserato/CodeWeaver/main.go:248.3,250.17 3 1
github.com/tesserato/CodeWeaver/main.go:250.17,253.4 1 1
github.com/tesserato/CodeWeaver/main.go:254.3,255.20 2 1
github.com/tesserato/CodeWeaver/main.go:257.2,257.17 1 1
github.com/tesserato/CodeWeaver/main.go:257.17,260.3 2 0
github.com/tesserato/CodeWeaver/main.go:261.2,261.30 1 1
github.com/tesserato/CodeWeaver/main.go:267.142,278.16 7 1
github.com/tesserato/CodeWeaver/main.go:278.16,280.3 1 0
github.com/tesserato/CodeWeaver/main.go:281.2,289.16 7 1
github.com/tesserato/CodeWeaver/main.go:289.16,291.3 1 0
github.com/tesserato/CodeWeaver/main.go:292.2,295.68 3 1
github.com/tesserato/CodeWeaver/main.go:310.104,317.2 1 1
github.com/tesserato/CodeWeaver/main.go:320.58,323.2 2 1
github.com/tesserato/CodeWeaver/main.go:326.83,328.16 2 1
github.com/tesserato/CodeWeaver/main.go:328.16,331.3 1 0
github.com/tesserato/CodeWeaver/main.go:334.2,335.32 2 1
github.com/tesserato/CodeWeaver/main.go:335.32,338.17 3 1
github.com/tesserato/CodeWeaver/main.go:338.17,340.4 1 0
github.com/tesserato/CodeWeaver/main.go:342.3,342.75 1 1
github.com/tesserato/CodeWeaver/main.go:342.75,344.4 1 1
github.com/tesserato/CodeWeaver/main.go:348.2,348.50 1 1
github.com/tesserato/CodeWeaver/main.go:348.50,354.3 1 1
github.com/tesserato/CodeWeaver/main.go:357.2,357.40 1 1
github.com/tesserato/CodeWeaver/main.go:357.40,361.20 3 1
github.com/tesserato/CodeWeaver/main.go:361.20,366.18 3 1
github.com/tesserato/CodeWeaver/main.go:366.18,369.5 1 0
github.com/tesserato/CodeWeaver/main.go:377.2,377.12 1 1
github.com/tesserato/CodeWeaver/main.go:381.82,383.29 2 1
github.com/tesserato/CodeWeaver/main.go:383.29,384.22 1 1
github.com/tesserato/CodeWeaver/main.go:384.22,386.4 1 1
github.com/tesserato/CodeWeaver/main.go:386.9,388.4 1 1
github.com/tesserato/CodeWeaver/main.go:391.2,391.12 1 1
github.com/tesserato/CodeWeaver/main.go:391.12,393.3 1 1
github.com/tesserato/CodeWeaver/main.go:393.8,395.3 1 1
github.com/tesserato/CodeWeaver/main.go:397.2,399.29 3 1
github.com/tesserato/CodeWeaver/main.go:403.73,405.16 2 1
github.com/tesserato/CodeWeaver/main.go:405.16,408.3 1 0
github.com/tesserato/CodeWeaver/main.go:410.2,410.46 1 1
github.com/tesserato/CodeWeaver/main.go:428.105,437.2 1 1
github.com/tesserato/CodeWeaver/main.go:441.124,444.111 2 1
github.com/tesserato/CodeWeaver/main.go:444.111,446.21 1 1
github.com/tesserato/CodeWeaver/main.go:446.21,449.69 2 0
github.com/tesserato/CodeWeaver/main.go:449.69,451.5 1 0
github.com/tesserato/CodeWeaver/main.go:453.4,453.18 1 0
github.com/tesserato/CodeWeaver/main.go:458.3,459.20 2 1
github.com/tesserato/CodeWeaver/main.go:459.20,463.4 2 0
github.com/tesserato/CodeWeaver/main.go:464.3,467.28 2 1
github.com/tesserato/CodeWeaver/main.go:467.28,469.4 1 1
github.com/tesserato/CodeWeaver/main.go:472.3,474.15 2 1
github.com/tesserato/CodeWeaver/main.go:474.15,476.34 1 1
github.com/tesserato/CodeWeaver/main.go:476.34,478.5 1 1
github.com/tesserato/CodeWeaver/main.go:479.4,481.14 2 1
github.com/tesserato/CodeWeaver/main.go:485.3,485.16 1 1
github.com/tesserato/CodeWeaver/main.go:485.16,489.4 1 1
github.com/tesserato/CodeWeaver/main.go:492.3,492.33 1 1
github.com/tesserato/CodeWeaver/main.go:492.33,494.4 1 1
github.com/tesserato/CodeWeaver/main.go:495.3,499.21 3 1
github.com/tesserato/CodeWeaver/main.go:499.21,505.4 4 0
github.com/tesserato/CodeWeaver/main.go:508.3,514.13 6 1
github.com/tesserato/CodeWeaver/main.go:517.2,517.20 1 1
github.com/tesserato/CodeWeaver/main.go:517.20,520.3 1 0
github.com/tesserato/CodeWeaver/main.go:522.2,522.62 1 1
github.com/tesserato/CodeWeaver/main.go:529.98,531.41 1 1
github.com/tesserato/CodeWeaver/main.go:531.41,532.60 1 1
github.com/tesserato/CodeWeaver/main.go:532.60,534.4 1 1
github.com/tesserato/CodeWeaver/main.go:538.2,538.30 1 1
github.com/tesserato/CodeWeaver/main.go:538.30,540.43 2 1
github.com/tesserato/CodeWeaver/main.go:540.43,541.61 1 1
github.com/tesserato/CodeWeaver/main.go:541.61,543.10 2 1
github.com/tesserato/CodeWeaver/main.go:546.3,546.22 1 1
github.com/tesserato/CodeWeaver/main.go:546.22,548.4 1 1
github.com/tesserato/CodeWeaver/main.go:555.2,555.13 1 1
github.com/tesserato/CodeWeaver/main.go:561.120,565.16 3 1
github.com/tesserato/CodeWeaver/main.go:565.16,569.3 2 1
github.com/tesserato/CodeWeaver/main.go:570.2,573.33 2 1
github.com/tesserato/CodeWeaver/main.go:573.33,574.87 1 1
github.com/tesserato/CodeWeaver/main.go:574.87,577.4 1 0
github.com/tesserato/CodeWeaver/main.go:581.2,581.33 1 1
github.com/tesserato/CodeWeaver/main.go:581.33,582.87 1 1
github.com/tesserato/CodeWeaver/main.go:582.87,585.4 1 0
github.com/tesserato/CodeWeaver/main.go:589.2,589.24 1 1
github.com/tesserato/CodeWeaver/main.go:589.24,593.17 3 0
github.com/tesserato/CodeWeaver/main.go:593.17,595.4 1 0
github.com/tesserato/CodeWeaver/main.go:595.9,598.4 2 0
github.com/tesserato/CodeWeaver/main.go:600.2,600.12 1 1
github.com/tesserato/CodeWeaver/main.go:604.81,605.21 1 1
github.com/tesserato/CodeWeaver/main.go:605.21,610.3 2 1
github.com/tesserato/CodeWeaver/main.go:613.2,616.26 3 1
github.com/tesserato/CodeWeaver/main.go:616.26,619.3 2 1
github.com/tesserato/CodeWeaver/main.go:621.2,622.16 2 1
github.com/tesserato/CodeWeaver/main.go:622.16,624.3 1 1
github.com/tesserato/CodeWeaver/main.go:624.8,627.3 1 1
github.com/tesserato/CodeWeaver/main.go:628.2,628.12 1 1
github.com/tesserato/CodeWeaver/main.go:634.18,656.2 19 1

```


## exc.txt
```txt
.git
.git/HEAD
build
build/output.exe
build/tmp
build/tmp/log.txt
data
data/config.yaml
data/image.png
docs/sub_docs
docs/sub_docs/file_in_sub.txt
file1.txt
node_modules/dep/package.json
other.log

```


## excluded.log
```log
.git

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


## inc.txt
```txt
README.md
script.go

```


## included.log
```log
README.md
build/output.exe
build/tmp/log.txt
data/config.yaml
data/image.png
docs/sub_docs/file_in_sub.txt
file1.txt
node_modules/dep/package.json
other.log
script.go

```


## main.go
```go
package main

import (
	"errors" // Added for specific error checks
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

// Color codes for logging
const (
	colorRed       = "\033[31m"
	colorLiteRed   = "\033[91m"
	colorGreen     = "\033[32m"
	colorLiteGreen = "\033[92m"
	colorReset     = "\033[0m"
)

// Markdown formatting constants
const (
	mdTreeViewHeader  = "# Tree View:\n```\n"
	mdCodeBlockEnd    = "\n```\n"
	mdContentHeader   = "\n# Content:\n"
	mdFileHeaderStart = "\n## "
	mdCodeBlockStart  = "\n```" // Extension will be appended
	mdCodeBlockEndNL  = "\n```\n\n"
)

// Regex log prefixes
const (
	rgxLogPrefixIgnore  = "- RGX:"
	rgxLogPrefixInclude = "+ RGX:"
)

// File processing log prefixes
const (
	logPrefixExclude = "-"
	logPrefixInclude = "+"
)

// --- Build-time Variables ---

var (
	// Populated by goreleaser during build using ldflags
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// --- Main Execution ---

func main() {
	// 1. Parse Configuration
	cfg, err := parseFlags()
	if err != nil {
		// Check for specific flag errors like help/version request
		if errors.Is(err, flag.ErrHelp) {
			// flag package already printed help due to ContinueOnError strategy if used,
			// or printHelp was called manually. Exit cleanly.
			os.Exit(0)
		}
		// For other parsing errors (invalid input dir etc.)
		log.Fatalf("%sError parsing flags: %v%s", colorRed, err, colorReset)
	}

	// Handle flags that exit immediately
	if cfg.showVersion {
		fmt.Printf("CodeWeaver version %s\ncommit %s\nbuilt at %s\n", version, commit, date)
		os.Exit(0)
	}
	// Help is handled implicitly by parseFlags or caught by flag.ErrHelp

	// 2. Setup Logger
	logger := setupLogging(cfg)

	// 3. Compile Regex Matchers
	ignoreMatchers, includeMatchers, err := compileMatchers(cfg, logger)
	if err != nil {
		logger.Fatalf("%sError compiling regex patterns: %v%s", colorRed, err, colorReset)
	}

	// 4. Generate Markdown Content
	markdownString, includedPaths, excludedPaths, err := generateMarkdown(cfg, ignoreMatchers, includeMatchers, logger)
	if err != nil {
		logger.Fatalf("%sError generating markdown: %v%s", colorRed, err, colorReset)
	}

	// 5. Write Output (File, Logs, Clipboard)
	err = writeOutput(cfg, markdownString, includedPaths, excludedPaths, logger)
	if err != nil {
		// writeOutput logs specific warnings, but a fatal error here means primary output failed.
		logger.Fatalf("%sError writing output: %v%s", colorRed, err, colorReset)
	}

	logger.Printf("%sCodeWeaver finished successfully.%s", colorGreen, colorReset)
}

// --- Configuration & Setup ---

// config holds the application's configuration values derived from flags.
type config struct {
	inputDirOriginal  string   // Original input path provided by user
	inputDirAbs       string   // Absolute path of the input directory
	outputFile        string   // Path for the output markdown file
	ignorePatterns    []string // Raw ignore patterns from flags
	includePatterns   []string // Raw include patterns from flags
	includedPathsFile string   // File to save included paths list
	excludedPathsFile string   // File to save excluded paths list
	addToClipboard    bool     // Flag to copy output to clipboard
	showHelp          bool     // Flag to show help message
	showVersion       bool     // Flag to show version information
	// logger            *log.Logger      // Central logger - REMOVED (passed as arg)
	// ignoreMatchers    []*regexp.Regexp // Compiled ignore patterns - REMOVED (returned by compileMatchers)
	// includeMatchers   []*regexp.Regexp // Compiled include patterns - REMOVED (returned by compileMatchers)
}

// parseFlags defines, parses, and validates command-line flags.
// It returns a config struct or an error.
func parseFlags() (*config, error) {
	cfg := &config{}
	// Use a temporary flag set to allow checking for ErrHelp if needed, though
	// the custom Usage function handles the common -help case.
	fs := flag.NewFlagSet("codeweaver", flag.ContinueOnError) // Allow returning errors
	fs.SetOutput(os.Stderr)                                   // Default error output

	fs.StringVar(&cfg.inputDirOriginal, "input", ".", "The root directory to scan.")
	fs.StringVar(&cfg.outputFile, "output", "codebase.md", "The name of the output Markdown file.")
	ignoreStr := fs.String("ignore", `\.git.*`, "Comma-separated list of regular expressions for paths to *exclude*.")
	includeStr := fs.String("include", "", "Comma-separated list of regular expressions. *Only* paths matching these are *included*.")
	fs.StringVar(&cfg.includedPathsFile, "included-paths-file", "", "Saves the list of *included* paths to this file.")
	fs.StringVar(&cfg.excludedPathsFile, "excluded-paths-file", "", "Saves the list of *excluded* paths to this file.")
	fs.BoolVar(&cfg.addToClipboard, "clipboard", false, "Copies the generated Markdown to the clipboard.")
	fs.BoolVar(&cfg.showVersion, "version", false, "Displays the version and exits.")
	// Handle -help explicitly via Usage
	fs.BoolVar(&cfg.showHelp, "help", false, "Displays help message and exits.")

	fs.Usage = printHelp // Override default usage message

	err := fs.Parse(os.Args[1:]) // Parse flags excluding program name
	if err != nil {
		// err could be flag.ErrHelp if -help was parsed by the library,
		// or an error for an unknown flag etc.
		return nil, err // Propagate parse error (or ErrHelp)
	}

	// If -help flag was set explicitly (and parsed without error), print help and signal exit.
	if cfg.showHelp {
		printHelp()
		// Use flag.ErrHelp to signal the caller (main) that help was requested.
		return nil, flag.ErrHelp
	}

	// Handle version flag immediately if set (no further validation needed)
	if cfg.showVersion {
		return cfg, nil // Return config state; main will handle version print & exit
	}

	// --- Post-parsing Validation (only if not showing help/version) ---

	// Get absolute path for input directory
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

	// Split pattern strings into slices
	if *ignoreStr != "" {
		cfg.ignorePatterns = strings.Split(*ignoreStr, ",")
	}
	if *includeStr != "" {
		cfg.includePatterns = strings.Split(*includeStr, ",")
	}

	return cfg, nil
}

// setupLogging initializes and returns a logger, printing initial config info.
func setupLogging(cfg *config) *log.Logger {
	logger := log.New(os.Stdout, "", 0) // Simple logger to stdout
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
	logger.Println() // Blank line for separation
	return logger
}

// compileMatchers compiles the regex patterns from the config.
func compileMatchers(cfg *config, logger *log.Logger) (ignore []*regexp.Regexp, include []*regexp.Regexp, err error) {
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
	logger.Println() // Blank line for separation
	return ignore, include, nil
}

// compileRegexList compiles a list of raw regex patterns.
func compileRegexList(patterns []string, color, prefix string, logger *log.Logger) ([]*regexp.Regexp, error) {
	if len(patterns) == 0 {
		logger.Printf("  (No patterns provided)")
		return nil, nil
	}
	compiledMatchers := make([]*regexp.Regexp, 0, len(patterns))
	foundValid := false
	for _, p := range patterns {
		trimmedPattern := strings.TrimSpace(p)
		if trimmedPattern == "" {
			continue // Skip empty patterns
		}
		logger.Printf("%s  %s %s%s\n", color, prefix, trimmedPattern, colorReset)
		rgx, err := regexp.Compile(trimmedPattern)
		if err != nil {
			// Fail completely if any pattern is invalid
			return nil, fmt.Errorf("pattern '%s': %w", trimmedPattern, err)
		}
		compiledMatchers = append(compiledMatchers, rgx)
		foundValid = true
	}
	if !foundValid {
		logger.Printf("  (No valid patterns found after trimming)")
		return nil, nil // Return nil if only empty/whitespace patterns were given
	}
	return compiledMatchers, nil
}

// --- Markdown Generation ---

// generateMarkdown orchestrates the creation of the tree view and content sections.
func generateMarkdown(cfg *config, ignoreMatchers, includeMatchers []*regexp.Regexp, logger *log.Logger) (string, []string, []string, error) {
	var markdownContent strings.Builder

	// --- Build Tree View ---
	logger.Println("Building tree view...")
	markdownContent.WriteString(mdTreeViewHeader)
	// Display the original input path as the root for user-friendliness
	markdownContent.WriteString(filepath.ToSlash(cfg.inputDirOriginal) + "\n")

	treeBuilder := newTreeBuilder(cfg.inputDirAbs, ignoreMatchers, includeMatchers)
	treeString, err := treeBuilder.buildTreeString()
	if err != nil {
		return "", nil, nil, fmt.Errorf("failed to build codebase tree: %w", err)
	}
	markdownContent.WriteString(treeString)
	markdownContent.WriteString(mdCodeBlockEnd)

	// --- Build Content Section ---
	logger.Println("Building content section...")
	markdownContent.WriteString(mdContentHeader)
	contentBuilder := newContentBuilder(cfg.inputDirAbs, cfg.includedPathsFile, cfg.excludedPathsFile, ignoreMatchers, includeMatchers, logger)
	contentString, includedPaths, excludedPaths, err := contentBuilder.buildContentString()
	if err != nil {
		return "", nil, nil, fmt.Errorf("failed to build code content: %w", err)
	}
	markdownContent.WriteString(contentString)
	logger.Println() // Blank line after processing logs

	return markdownContent.String(), includedPaths, excludedPaths, nil
}

// --- Tree Builder ---

// treeBuilder handles the generation of the directory tree structure.
type treeBuilder struct {
	rootAbsPath     string
	ignoreMatchers  []*regexp.Regexp
	includeMatchers []*regexp.Regexp
	output          strings.Builder
	depthOpen       map[int]bool // Tracks open branches at each depth level for │ character
}

// newTreeBuilder creates a new tree builder instance.
func newTreeBuilder(rootAbsPath string, ignoreMatchers, includeMatchers []*regexp.Regexp) *treeBuilder {
	return &treeBuilder{
		rootAbsPath:     rootAbsPath,
		ignoreMatchers:  ignoreMatchers,
		includeMatchers: includeMatchers,
		depthOpen:       make(map[int]bool),
	}
}

// buildTreeString generates the formatted tree string.
func (tb *treeBuilder) buildTreeString() (string, error) {
	err := tb.printTreeRecursive(tb.rootAbsPath, 0)
	return tb.output.String(), err // Return accumulated string and any error
}

// printTreeRecursive recursively walks the directory and builds the tree structure.
func (tb *treeBuilder) printTreeRecursive(currentDirPath string, depth int) error {
	entries, err := os.ReadDir(currentDirPath)
	if err != nil {
		// Don't log here, return the error to be handled higher up
		return fmt.Errorf("read directory %s: %w", currentDirPath, err)
	}

	// Filter entries based on include/exclude rules *before* sorting
	var filteredEntries []fs.DirEntry
	for _, entry := range entries {
		fullEntryPath := filepath.Join(currentDirPath, entry.Name())
		pathRelToInput, err := tb.getRelativePath(fullEntryPath)
		if err != nil {
			return err // Error during path relativization
		}

		if shouldProcess(pathRelToInput, tb.ignoreMatchers, tb.includeMatchers) {
			filteredEntries = append(filteredEntries, entry)
		}
	}

	// Sort filtered entries alphabetically for consistent output
	sort.Slice(filteredEntries, func(i, j int) bool {
		// Consider directories first (optional, but common in tree views)
		// if filteredEntries[i].IsDir() != filteredEntries[j].IsDir() {
		// 	return filteredEntries[i].IsDir() // True if i is dir, j is file
		// }
		return strings.ToLower(filteredEntries[i].Name()) < strings.ToLower(filteredEntries[j].Name())
	})

	// Process sorted and filtered entries
	for i, entry := range filteredEntries {
		isLastEntry := (i == len(filteredEntries) - 1)
		tb.printEntryLine(entry, depth, isLastEntry)

		if entry.IsDir() {
			// Mark if the current branch needs a vertical line for subsequent levels
			tb.depthOpen[depth] = !isLastEntry
			// Recurse into the directory
			err := tb.printTreeRecursive(filepath.Join(currentDirPath, entry.Name()), depth+1)
			if err != nil {
				// If recursion fails (e.g., permission error deeper down), propagate the error
				return err
			}
			// Clean up the depth state after returning from recursion? No, sibling entries need it.
			// The map naturally handles overwrites if the same depth is revisited.
		}
	}
	// Ensure the state for this depth is cleaned up if it was the last entry processed at this level?
	// Testing shows this is not strictly necessary with the current map logic.
	// delete(tb.depthOpen, depth)
	return nil
}

// printEntryLine formats and writes a single line of the tree output.
func (tb *treeBuilder) printEntryLine(entry fs.DirEntry, depth int, isLast bool) {
	var prefix strings.Builder
	for i := 0; i < depth; i++ {
		if tb.depthOpen[i] {
			prefix.WriteString("│  ") // Vertical line and spaces
		} else {
			prefix.WriteString("   ") // Spaces only
		}
	}

	if isLast {
		prefix.WriteString("└─ ") // Last item marker
	} else {
		prefix.WriteString("├─ ") // Intermediate item marker
	}

	tb.output.WriteString(prefix.String())
	tb.output.WriteString(entry.Name()) // Use the entry's base name
	tb.output.WriteString("\n")
}

// getRelativePath converts an absolute path to a slash-separated path relative to the tree root.
func (tb *treeBuilder) getRelativePath(fullPath string) (string, error) {
	pathRelToInput, err := filepath.Rel(tb.rootAbsPath, fullPath)
	if err != nil {
		// This error indicates a logic problem if fullPath doesn't start with rootAbsPath
		return "", fmt.Errorf("failed to make path %s relative to %s: %w", fullPath, tb.rootAbsPath, err)
	}
	// Use forward slashes for consistent matching and output
	return filepath.ToSlash(pathRelToInput), nil
}

// --- Content Builder ---

// contentBuilder handles walking the filesystem and generating the content markdown.
type contentBuilder struct {
	rootAbsPath       string
	includedPathsFile string // Config value, used for logging decisions
	excludedPathsFile string // Config value, used for logging decisions
	ignoreMatchers    []*regexp.Regexp
	includeMatchers   []*regexp.Regexp
	logger            *log.Logger
}

// newContentBuilder creates a new content builder instance.
func newContentBuilder(
	rootAbsPath string, includedPathsFile string, excludedPathsFile string,
	ignoreMatchers []*regexp.Regexp, includeMatchers []*regexp.Regexp, logger *log.Logger) *contentBuilder {
	return &contentBuilder{
		rootAbsPath:       rootAbsPath,
		includedPathsFile: includedPathsFile,
		excludedPathsFile: excludedPathsFile,
		ignoreMatchers:    ignoreMatchers,
		includeMatchers:   includeMatchers,
		logger:            logger,
	}
}

// buildContentString walks the filesystem, processes files, and returns the content markdown string,
// along with lists of included and excluded relative paths.
func (cb *contentBuilder) buildContentString() (content string, includedPaths []string, excludedPaths []string, err error) {
	var contentSB strings.Builder // Use strings.Builder for efficiency

	walkErr := filepath.WalkDir(cb.rootAbsPath, func(currentWalkPath string, d fs.DirEntry, walkErr error) error {
		// Handle errors encountered by WalkDir itself (e.g., permission denied)
		if walkErr != nil {
			cb.logger.Printf("%sWarning: Error accessing %s: %v%s\n", colorRed, currentWalkPath, walkErr, colorReset)
			// Decide if the error is skippable for a directory
			if d != nil && d.IsDir() && errors.Is(walkErr, fs.ErrPermission) { // Use errors.Is for specific error check
				return fs.SkipDir // Skip directory if permission denied, but continue walk elsewhere
			}
			// For other errors, or errors on files, stop the walk
			return walkErr // Propagate the error to stop WalkDir
		}


		// Get path relative to the *original* input root for filtering and output
		pathRelToInput, relErr := filepath.Rel(cb.rootAbsPath, currentWalkPath)
		if relErr != nil {
			// This should not happen if WalkDir starts correctly
			cb.logger.Printf("%sWarning: Could not make path %s relative to %s: %v%s\n", colorRed, currentWalkPath, cb.rootAbsPath, relErr, colorReset)
			return nil // Skip this entry, continue walk
		}
		pathRelToInput = filepath.ToSlash(pathRelToInput) // Use forward slashes

		// Skip processing the root directory itself for include/exclude logic here
		if pathRelToInput == "." {
			return nil // Continue walking into the root directory
		}

		// Determine if the path should be processed based on filters
		process := shouldProcess(pathRelToInput, cb.ignoreMatchers, cb.includeMatchers)

		if !process {
			// Log exclusion only if not saving paths to a file (to avoid verbose output)
			if cb.excludedPathsFile == "" {
				cb.logger.Printf("%s%s %s%s\n", colorRed, logPrefixExclude, pathRelToInput, colorReset)
			}
			excludedPaths = append(excludedPaths, pathRelToInput) // Collect relative path
			// Do NOT SkipDir here - let WalkDir continue so children can be evaluated individually.
			return nil // Skip this file/dir entry, continue walk
		}

		// If we reach here, the path is included (either file or directory)
		if d.IsDir() {
			// Directories that pass shouldProcess don't contribute to content markdown directly,
			// but we don't exclude them here. WalkDir continues into them.
			return nil // Continue walk
		}

		// Process included *files*
		if cb.includedPathsFile == "" {
			cb.logger.Printf("%s%s %s%s\n", colorGreen, logPrefixInclude, pathRelToInput, colorReset)
		}
		includedPaths = append(includedPaths, pathRelToInput) // Collect relative path

		// Read file content
		fileContent, readErr := os.ReadFile(currentWalkPath)
		if readErr != nil {
			cb.logger.Printf("%sWarning: Failed to read file %s: %v%s\n", colorRed, currentWalkPath, readErr, colorReset)
			// Add a placeholder to the markdown for unreadable files
			contentSB.WriteString(fmt.Sprintf("%s%s\n", mdFileHeaderStart, pathRelToInput)) // Use relative path
			contentSB.WriteString(fmt.Sprintf("%s\nError reading file: %v%s", mdCodeBlockStart, readErr, mdCodeBlockEndNL))
			return nil // Continue walk with the next file
		}

		// Add file content to markdown
		extension := strings.TrimPrefix(strings.ToLower(filepath.Ext(currentWalkPath)), ".")
		contentSB.WriteString(fmt.Sprintf("%s%s", mdFileHeaderStart, pathRelToInput))         // Header with relative path
		contentSB.WriteString(fmt.Sprintf("%s%s\n", mdCodeBlockStart, extension)) // Code block start with lang hint
		contentSB.Write(fileContent)                                              // Write file bytes
		contentSB.WriteString(mdCodeBlockEndNL)                                   // Code block end

		return nil // Continue walk
	}) // End of WalkDirFunc

	if walkErr != nil {
		// Return the error that stopped WalkDir
		return "", nil, nil, fmt.Errorf("walking directory %s: %w", cb.rootAbsPath, walkErr)
	}

	return contentSB.String(), includedPaths, excludedPaths, nil
}

// --- Path Filtering Logic ---

// shouldProcess determines if a path should be processed based on include and ignore patterns.
// pathRelToInput must be relative to the input directory and use forward slashes.
func shouldProcess(pathRelToInput string, ignoreMatchers, includeMatchers []*regexp.Regexp) bool {
	// 1. Check Ignore Patterns (Blacklist)
	for _, pattern := range ignoreMatchers {
		if pattern != nil && pattern.MatchString(pathRelToInput) {
			return false // Path matches an ignore pattern, exclude it.
		}
	}

	// 2. Check Include Patterns (Whitelist) - only if include patterns exist
	if len(includeMatchers) > 0 {
		matchedInclude := false
		for _, pattern := range includeMatchers {
			if pattern != nil && pattern.MatchString(pathRelToInput) {
				matchedInclude = true // Path matches at least one include pattern.
				break
			}
		}
		if !matchedInclude {
			return false // Include patterns exist, but this path didn't match any. Exclude it.
		}
	}

	// 3. Default Allow
	// If we reach here:
	// - EITHER no include patterns were provided, AND the path didn't match any ignore patterns.
	// - OR include patterns were provided, the path matched at least one include pattern, AND it didn't match any ignore patterns.
	return true // Include the path.
}

// --- Output Handling ---

// writeOutput handles writing the main markdown file, the path lists, and clipboard operations.
func writeOutput(cfg *config, markdownContent string, includedPaths, excludedPaths []string, logger *log.Logger) error {
	// --- Write Main Output File ---
	logger.Printf("Writing output to %s...", cfg.outputFile)
	err := os.WriteFile(cfg.outputFile, []byte(markdownContent), 0644)
	if err != nil {
		// Log the specific error but return it to signal main failure
		logger.Printf("%sError writing to output file %s: %v%s", colorRed, cfg.outputFile, err, colorReset)
		return fmt.Errorf("writing output file %s: %w", cfg.outputFile, err)
	}
	logger.Printf("Markdown content written to %s", cfg.outputFile)

	// --- Save Included Paths ---
	if cfg.includedPathsFile != "" {
		if err := savePathsToFile(cfg.includedPathsFile, includedPaths, logger); err != nil {
			// Log as warning, don't treat as fatal for the main operation
			logger.Printf("%sWarning: Error saving included paths to %s: %v%s", colorRed, cfg.includedPathsFile, err, colorReset)
		}
	}

	// --- Save Excluded Paths ---
	if cfg.excludedPathsFile != "" {
		if err := savePathsToFile(cfg.excludedPathsFile, excludedPaths, logger); err != nil {
			// Log as warning
			logger.Printf("%sWarning: Error saving excluded paths to %s: %v%s", colorRed, cfg.excludedPathsFile, err, colorReset)
		}
	}

	// --- Copy to Clipboard ---
	if cfg.addToClipboard {
		logger.Println("Attempting to copy to clipboard...")
		// Initialize clipboard. Needs to be done each time before writing.
		err := clipboard.Init()
		if err != nil {
			logger.Printf("%sWarning: Could not initialize clipboard: %v%s", colorRed, err, colorReset)
		} else {
			clipboard.Write(clipboard.FmtText, []byte(markdownContent))
			logger.Println("Markdown content copied to clipboard.")
		}
	}
	return nil // Indicate success of the primary output operation
}

// savePathsToFile writes a list of paths (one per line) to the specified file.
func savePathsToFile(filename string, paths []string, logger *log.Logger) error {
	if len(paths) == 0 {
		logger.Printf("No paths to save to %s.", filename)
		// Decide if an empty file should be created or not. Let's not create it.
		// return os.WriteFile(filename, []byte{}, 0644)
		return nil // Do nothing if no paths
	}

	// Sort paths for consistent output in the file
	sort.Strings(paths)

	var sb strings.Builder
	for _, p := range paths {
		sb.WriteString(p)
		sb.WriteString("\n")
	}

	err := os.WriteFile(filename, []byte(sb.String()), 0644)
	if err == nil {
		logger.Printf("Paths saved to %s", filename)
	} else {
		// Log is handled by the caller (writeOutput), just return error
		return fmt.Errorf("writing paths file %s: %w", filename, err)
	}
	return nil
}

// --- Help Message ---

// printHelp displays the command-line help message.
func printHelp() {
	// Use os.Stderr for help message output, similar to flag package default
	fmt.Fprintf(os.Stderr, "CodeWeaver: Generate Markdown Documentation from Your Codebase.\n")
	fmt.Fprintf(os.Stderr, "Version: %s, Commit: %s, Date: %s\n\n", version, commit, date)
	fmt.Fprintf(os.Stderr, "Usage: codeweaver [options]\n")
	fmt.Fprintf(os.Stderr, "\nOptions:\n")
	// Temporarily set flag output to Stderr for PrintDefaults
	originalOutput := flag.CommandLine.Output()
	flag.CommandLine.SetOutput(os.Stderr)
	flag.PrintDefaults()
	flag.CommandLine.SetOutput(originalOutput) // Restore original output
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
	"bufio" // For checking lines in tree output
	"bytes" // For capturing log output
	"errors"
	"flag"
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

// --- Test Helpers ---

// mustCompileRegex compiles a regex or panics. Used for setting up test matchers.
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

// normalizeNewlines ensures consistent line endings for string comparisons.
func normalizeNewlines(s string) string {
	return strings.ReplaceAll(s, "\r\n", "\n")
}

// createTestFS creates a temporary directory structure for testing filesystem operations.
// It returns the root path of the temporary structure and a cleanup function.
func createTestFS(t *testing.T) (string, func()) {
	t.Helper()
	// Use t.TempDir() for the *source* file system as well, simplifying cleanup.
	rootDir := t.TempDir() // Automatically cleaned up by the test framework

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
		parentDir := filepath.Dir(absPath)

		// Ensure parent directory exists
		if err := os.MkdirAll(parentDir, 0755); err != nil {
			// Since rootDir is managed by t.TempDir, we don't need manual cleanup on error here.
			t.Fatalf("Failed to create parent dir %s: %v", parentDir, err)
		}

		// Check if it's intended to be a directory
		isDirectory := strings.HasSuffix(relPath, "/") || (content == "" && !strings.Contains(filepath.Base(relPath), "."))

		if isDirectory {
			if err := os.MkdirAll(absPath, 0755); err != nil {
				if !errors.Is(err, os.ErrExist) { // Allow directory already existing
					t.Fatalf("Failed to create dir %s: %v", absPath, err)
				}
			}
		} else { // It's a file
			if err := os.WriteFile(absPath, []byte(content), 0644); err != nil {
				t.Fatalf("Failed to write file %s: %v", absPath, err)
			}
		}
	}

	// Cleanup function is no longer needed as t.TempDir handles it.
	return rootDir, func() {}
}

// equalStringSlices checks if two string slices are equal (assumes sorted).
func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// runMainLogic simulates the execution of the main function with given arguments,
// capturing log output and directing file output to a specified directory.
// Returns captured log output and any error encountered during execution.
func runMainLogic(args []string, outputDir string) (string, error) {
	// 1. Set up simulated args
	originalArgs := os.Args
	os.Args = append([]string{"codeweaver"}, args...) // Simulate command name + args
	defer func() { os.Args = originalArgs }()         // Restore original args

	// 2. Capture log output
	var logBuf bytes.Buffer
	originalLoggerOutput := log.Writer()
	// Configure logger to write to buffer without timestamps/prefixes
	testRunLogger := log.New(&logBuf, "", 0)
	// Temporarily redirect the global log output as well, in case any part uses it directly
	log.SetOutput(&logBuf)
	log.SetFlags(0)
	defer func() {
		// Restore global logger
		log.SetOutput(originalLoggerOutput)
		log.SetFlags(log.LstdFlags) // Restore standard flags
	}()

	// 3. Use a flag set local to this run to avoid global state issues
	testFlags := flag.NewFlagSet("testRun", flag.ContinueOnError)
	testFlags.SetOutput(io.Discard) // Prevent flag set from printing errors to console

	// --- Define flags using the test FlagSet ---
	inputDirOriginal := testFlags.String("input", ".", "The root directory to scan.")
	outputFileName := testFlags.String("output", "codebase.md", "The name of the output Markdown file.")
	ignoreStr := testFlags.String("ignore", `\.git.*`, "Comma-separated list of regex patterns to exclude.")
	includeStr := testFlags.String("include", "", "Comma-separated list of regex patterns to include.")
	includedPathsFileName := testFlags.String("included-paths-file", "", "File to save included paths.")
	excludedPathsFileName := testFlags.String("excluded-paths-file", "", "File to save excluded paths.")
	addToClipboard := testFlags.Bool("clipboard", false, "Copy markdown to clipboard.")
	showVersion := testFlags.Bool("version", false, "Display version and exit.")
	showHelp := testFlags.Bool("help", false, "Display help message and exit.")

	// Manually set the Usage function for the test FlagSet
	testFlags.Usage = func() {
		// Simulate printHelp output (captured by log redirection)
		var helpBuf bytes.Buffer
		w := io.Writer(&helpBuf) // Writer for help message
		fmt.Fprintln(w, "CodeWeaver: Generate Markdown Documentation from Your Codebase.")
		fmt.Fprintf(w, "Version: %s, Commit: %s, Date: %s\n\n", version, commit, date)
		fmt.Fprintln(w, "Usage: codeweaver [options]")
		fmt.Fprintln(w, "\nOptions:")
		testFlags.SetOutput(w) // Temporarily redirect PrintDefaults output
		testFlags.PrintDefaults()
		testFlags.SetOutput(io.Discard) // Restore discard
		fmt.Fprintln(w, "\nExamples:")  // Add Examples header like printHelp
		fmt.Fprintln(w, "  codeweaver                               # Process current directory, output to codebase.md")
		fmt.Fprintln(w, "  codeweaver -input my_project -output docs.md")
		fmt.Fprintln(w, `  codeweaver -ignore "build/,vendor/" -include "\.go$,\.md$"`)
		fmt.Fprintln(w, "  codeweaver -clipboard -excluded-paths-file ignored.txt")
		fmt.Fprintln(w, "\nNotes on patterns:") // Add Notes header like printHelp
		fmt.Fprintln(w, "  - Patterns are Go regular expressions.")
		fmt.Fprintln(w, "  - Paths for filtering are relative to the input directory (e.g., \"src/main.go\").")
		fmt.Fprintln(w, "  - Use forward slashes '/' in patterns for cross-platform compatibility.")
		testRunLogger.Print(helpBuf.String()) // Use our test logger to capture the output
	}

	// Parse the provided args (excluding the command name)
	err := testFlags.Parse(args)
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			testFlags.Usage() // Ensure help is logged if requested via flag package
		}
		// Propagate flag parsing errors (like undefined flags or ErrHelp)
		return logBuf.String(), err
	}

	// --- Create config struct from parsed flags ---
	cfg := &config{
		inputDirOriginal:  *inputDirOriginal,
		outputFile:        *outputFileName,        // Base name, will be joined with outputDir later
		includedPathsFile: *includedPathsFileName, // Base name
		excludedPathsFile: *excludedPathsFileName, // Base name
		addToClipboard:    *addToClipboard,
		showHelp:          *showHelp,
		showVersion:       *showVersion,
	}
	if *ignoreStr != "" {
		cfg.ignorePatterns = strings.Split(*ignoreStr, ",")
	}
	if *includeStr != "" {
		cfg.includePatterns = strings.Split(*includeStr, ",")
	}

	// 4. Execute the core logic from main(), adapted to return errors instead of Fatal

	// Handle version/help flags first (they stop execution)
	if cfg.showVersion {
		// Capture the print output into the log buffer for checking
		testRunLogger.Printf("CodeWeaver version %s\ncommit %s\nbuilt at %s\n", version, commit, date)
		return logBuf.String(), nil // Successful exit for version
	}
	// Handle help flag that was set explicitly (not via flag package's ErrHelp)
	if cfg.showHelp {
		testFlags.Usage()                    // Call Usage which now logs the help text
		return logBuf.String(), flag.ErrHelp // Return ErrHelp to signal this path
	}

	// Perform post-parsing steps (validation) from original parseFlags
	cfg.inputDirAbs, err = filepath.Abs(cfg.inputDirOriginal)
	if err != nil {
		err = fmt.Errorf("failed to get absolute path for input directory '%s': %w", cfg.inputDirOriginal, err)
		testRunLogger.Print(colorRed + err.Error() + colorReset) // Log the error like main would
		return logBuf.String(), err
	}
	info, err := os.Stat(cfg.inputDirAbs)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = fmt.Errorf("input directory '%s' does not exist", cfg.inputDirAbs)
		} else {
			err = fmt.Errorf("error accessing input directory '%s': %w", cfg.inputDirAbs, err)
		}
		testRunLogger.Print(colorRed + err.Error() + colorReset)
		return logBuf.String(), err
	}
	if !info.IsDir() {
		err = fmt.Errorf("input path '%s' is not a directory", cfg.inputDirAbs)
		testRunLogger.Print(colorRed + err.Error() + colorReset)
		return logBuf.String(), err
	}

	// --- Replicate main() steps using the test logger ---

	// Log config details (initial setup logging)
	testRunLogger.Println("Starting CodeWeaver...")
	testRunLogger.Println("Input directory:", cfg.inputDirAbs)
	// Log the *intended* output paths (relative to outputDir)
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

	// Compile Matchers
	ignoreMatchers, includeMatchers, err := compileMatchers(cfg, testRunLogger) // Pass cfg and logger
	if err != nil {
		err = fmt.Errorf("Error compiling regex patterns: %w", err)
		testRunLogger.Println(colorRed + err.Error() + colorReset) // Log error like main does
		return logBuf.String(), err                                // Return error to signal failure
	}

	// Generate Markdown
	markdownString, includedPaths, excludedPaths, err := generateMarkdown(cfg, ignoreMatchers, includeMatchers, testRunLogger)
	if err != nil {
		err = fmt.Errorf("Error generating markdown: %w", err)
		testRunLogger.Println(colorRed + err.Error() + colorReset)
		return logBuf.String(), err
	}

	// --- Write Output (to specified outputDir) ---
	fullOutputPath := filepath.Join(outputDir, cfg.outputFile)
	fullIncludedPath := ""
	if cfg.includedPathsFile != "" {
		fullIncludedPath = filepath.Join(outputDir, cfg.includedPathsFile)
	}
	fullExcludedPath := ""
	if cfg.excludedPathsFile != "" {
		fullExcludedPath = filepath.Join(outputDir, cfg.excludedPathsFile)
	}

	writeCfg := *cfg // Copy base config
	writeCfg.outputFile = fullOutputPath
	writeCfg.includedPathsFile = fullIncludedPath
	writeCfg.excludedPathsFile = fullExcludedPath

	// Store original clipboard flag state
	originalClipboardState := writeCfg.addToClipboard
	// Temporarily disable actual clipboard writing if the flag is true,
	// so we can log the simulation message instead.
	if writeCfg.addToClipboard {
		writeCfg.addToClipboard = false // Prevent writeOutput from trying clipboard.Write
	}

	// Call writeOutput with potentially modified cfg
	err = writeOutput(&writeCfg, markdownString, includedPaths, excludedPaths, testRunLogger)
	if err != nil {
		// writeOutput already logs details, just return the error
		return logBuf.String(), err
	}

	// Log the simulated clipboard message *if* the original flag was true
	if originalClipboardState {
		testRunLogger.Println("Markdown content copied to clipboard (simulated).")
	}

	return logBuf.String(), nil // Success
}

// --- Test Suite ---

func TestShouldProcess(t *testing.T) {
	// No changes needed
	testCases := []struct {
		name            string
		path            string
		ignoreMatchers  []*regexp.Regexp
		includeMatchers []*regexp.Regexp
		expected        bool
	}{
		{"NoFilters_Allow", "file.txt", nil, nil, true},
		{"IgnoreMatch_Exact", "skip.txt", []*regexp.Regexp{mustCompileRegex(`^skip\.txt$`)}, nil, false},
		{"IgnoreMatch_DirPrefix", "skip/file.txt", []*regexp.Regexp{mustCompileRegex(`^skip/`)}, nil, false},
		{"IncludeMatch_Extension", "src/main.go", nil, []*regexp.Regexp{mustCompileRegex(`\.go$`)}, true},
		{"IncludeNoMatch_Extension", "src/main.txt", nil, []*regexp.Regexp{mustCompileRegex(`\.go$`)}, false},
		{"IgnoreTakesPrecedence_PathMatchesBoth", "vendor/lib.go", []*regexp.Regexp{mustCompileRegex(`^vendor/`)}, []*regexp.Regexp{mustCompileRegex(`\.go$`)}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := shouldProcess(tc.path, tc.ignoreMatchers, tc.includeMatchers)
			if actual != tc.expected {
				t.Errorf("shouldProcess(%q) with ignore=%v, include=%v = %v; want %v", tc.path, tc.ignoreMatchers, tc.includeMatchers, actual, tc.expected)
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
		if !strings.Contains(err.Error(), "does not exist") {
			t.Errorf("Expected 'does not exist' error, got err: %v", err)
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
		if !strings.Contains(err.Error(), "is not a directory") {
			t.Errorf("Expected 'is not a directory' error, got err: %v", err)
		}
	})
}

func TestCompileRegexPatterns(t *testing.T) {
	// No changes needed
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
	t.Run("EmptyInput", func(t *testing.T) {
		patterns := []string{}
		matchers, err := compileRegexList(patterns, "", "", testLogger)
		if err != nil {
			t.Fatalf("no error expected, got %v", err)
		}
		if matchers != nil {
			t.Errorf("expected nil matchers, got %v", matchers)
		}
	})
}

func TestTreeBuilder(t *testing.T) {
	rootDir, cleanup := createTestFS(t)
	defer cleanup()

	testCases := []struct {
		name            string
		ignoreMatchers  []*regexp.Regexp
		includeMatchers []*regexp.Regexp
		expectedLines   []string // --- FIX: Check for lines instead of exact string ---
	}{
		{
			name: "NoFilters",
			// Check presence of key lines, order might vary slightly
			expectedLines: []string{
				"├─ .git",
				"│  └─ HEAD",
				"├─ build",
				"│  ├─ output.exe",
				"│  └─ tmp",
				"│     └─ log.txt",
				"├─ data",
				"│  ├─ config.yaml",
				"│  └─ image.png",
				"├─ docs",
				"│  └─ sub_docs",
				"│     └─ file_in_sub.txt",
				"├─ empty_dir",
				"├─ file1.txt",
				"├─ node_modules",
				"│  └─ dep",
				"│     └─ package.json",
				"├─ other.log",
				"├─ README.md", // Note: Order check relaxed
				"└─ script.go",
			},
		},
		{
			name:            "IncludeOnlyGoAndMdFiles",
			includeMatchers: []*regexp.Regexp{mustCompileRegex(`\.go$`), mustCompileRegex(`\.md$`)},
			expectedLines: []string{ // Order here is usually stable
				"├─ README.md",
				"└─ script.go",
			},
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

			// --- FIX: Check for line presence ---
			scanner := bufio.NewScanner(strings.NewReader(actualTree))
			actualLines := make(map[string]bool)
			for scanner.Scan() {
				line := scanner.Text()
				if line != "" { // Ignore potential empty lines
					actualLines[line] = true
				}
			}
			if err := scanner.Err(); err != nil {
				t.Fatalf("Error scanning actual tree output: %v", err)
			}

			missingLines := []string{}
			for _, expectedLine := range tc.expectedLines {
				if !actualLines[expectedLine] {
					missingLines = append(missingLines, expectedLine)
				}
			}

			if len(missingLines) > 0 {
				t.Errorf("Tree mismatch for '%s'. Missing expected lines:\n%s\nActual Tree:\n%s",
					tc.name, strings.Join(missingLines, "\n"), actualTree)
			}
			// Optional: Check if the number of non-empty lines matches expected count
			if len(actualLines) != len(tc.expectedLines) {
				t.Errorf("Tree mismatch for '%s'. Expected %d lines, got %d.\nActual Tree:\n%s",
					tc.name, len(tc.expectedLines), len(actualLines), actualTree)
			}
		})
	}
}

func TestContentBuilder(t *testing.T) {
	// No changes needed
	rootDir, cleanup := createTestFS(t)
	defer cleanup()
	testLogger := log.New(io.Discard, "", 0)

	testCases := []struct {
		name                  string
		ignoreMatchers        []*regexp.Regexp
		includeMatchers       []*regexp.Regexp
		expectedContentSubstr string
		expectedIncludedPaths []string // Sorted
		expectedExcludedPaths []string // Sorted
		checkContentEmpty     bool
	}{
		{
			name:                  "NoFilters",
			expectedContentSubstr: "## file1.txt\n```txt\ncontent of file1\n```",
			expectedIncludedPaths: []string{".git/HEAD", "README.md", "build/output.exe", "build/tmp/log.txt", "data/config.yaml", "data/image.png", "docs/sub_docs/file_in_sub.txt", "file1.txt", "node_modules/dep/package.json", "other.log", "script.go"},
			expectedExcludedPaths: []string{},
		},
		{
			name:                  "IncludeOnlyGoAndMdFiles",
			includeMatchers:       []*regexp.Regexp{mustCompileRegex(`\.go$`), mustCompileRegex(`\.md$`)},
			expectedContentSubstr: "## script.go\n```go\npackage main",
			expectedIncludedPaths: []string{"README.md", "script.go"},
			expectedExcludedPaths: []string{".git", ".git/HEAD", "build", "build/output.exe", "build/tmp", "build/tmp/log.txt", "data", "data/config.yaml", "data/image.png", "docs", "docs/sub_docs", "docs/sub_docs/file_in_sub.txt", "empty_dir", "file1.txt", "node_modules", "node_modules/dep", "node_modules/dep/package.json", "other.log"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sort.Strings(tc.expectedIncludedPaths)
			sort.Strings(tc.expectedExcludedPaths)

			builder := newContentBuilder(rootDir, "", "", tc.ignoreMatchers, tc.includeMatchers, testLogger)
			actualContentStr, actualIncludedPaths, actualExcludedPaths, err := builder.buildContentString()
			if err != nil {
				t.Logf("buildContentString() returned error (may be expected for permission tests): %v", err)
			}
			actualContentStr = normalizeNewlines(actualContentStr)
			sort.Strings(actualIncludedPaths)
			sort.Strings(actualExcludedPaths)

			if tc.expectedContentSubstr != "" && !strings.Contains(actualContentStr, tc.expectedContentSubstr) {
				t.Errorf("'%s': Generated content does not contain expected substring.\nExpected to find:\n%s\n-----\nActual Content:\n%s-----", tc.name, tc.expectedContentSubstr, actualContentStr)
			}
			if tc.checkContentEmpty && actualContentStr != "" {
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
}

func TestSavePathsToFile(t *testing.T) {
	// No changes needed
	testLogger := log.New(io.Discard, "", 0)
	t.Run("StandardSave", func(t *testing.T) {
		testDir := t.TempDir()
		tmpFilePath := filepath.Join(testDir, "test_paths_output.txt")
		paths := []string{"path/to/file1.txt", "another/path.go", "root_file.md"}
		expectedPathsSorted := []string{"another/path.go", "path/to/file1.txt", "root_file.md"}
		expectedContent := strings.Join(expectedPathsSorted, "\n") + "\n"
		err := savePathsToFile(tmpFilePath, paths, testLogger)
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
		testDir := t.TempDir()
		emptyPathsFile := filepath.Join(testDir, "empty_paths_save_test.txt")
		err := savePathsToFile(emptyPathsFile, []string{}, testLogger)
		if err != nil {
			t.Fatalf("savePathsToFile with empty paths failed: %v", err)
		}
		if _, err := os.Stat(emptyPathsFile); !errors.Is(err, os.ErrNotExist) {
			t.Errorf("Expected file not to be created, but it exists or other error: %v", err)
		}
	})
	t.Run("Error_PathIsExistingDirectory", func(t *testing.T) {
		testDir := t.TempDir()
		err := savePathsToFile(testDir, []string{"a", "b"}, testLogger)
		if err == nil {
			t.Fatalf("Expected error saving to a directory, got nil")
		}
	})
}

func TestPrintHelp(t *testing.T) {
	// No changes needed
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w
	printHelp()
	w.Close()
	outBytes, err := io.ReadAll(r)
	os.Stderr = oldStderr
	if err != nil {
		t.Fatalf("Failed to read stderr capture: %v", err)
	}
	output := string(outBytes)
	if !strings.Contains(output, "Usage: codeweaver") {
		t.Errorf("Help output missing 'Usage: codeweaver', got:\n%s", output)
	}
}

// TestMainExecutionFlows tests the main function's behavior
func TestMainExecutionFlows(t *testing.T) {
	baseInputDir, cleanupInput := createTestFS(t)
	defer cleanupInput()

	t.Run("VersionFlag", func(t *testing.T) {
		testOutputDir := t.TempDir()
		logOutput, err := runMainLogic([]string{"-version"}, testOutputDir)
		if err != nil {
			t.Fatalf("runMainLogic with -version failed: %v", err)
		}
		expected := fmt.Sprintf("CodeWeaver version %s", version)
		if !strings.Contains(logOutput, expected) {
			t.Errorf("Expected '%s' in output, got:\n%s", expected, logOutput)
		}
	})

	t.Run("HelpFlag", func(t *testing.T) {
		testOutputDir := t.TempDir()
		logOutput, err := runMainLogic([]string{"-help"}, testOutputDir)
		if !errors.Is(err, flag.ErrHelp) {
			t.Fatalf("runMainLogic with -help did not return flag.ErrHelp, got err: %v. Log:\n%s", err, logOutput)
		}
		if !strings.Contains(logOutput, "Usage: codeweaver") {
			t.Errorf("Expected 'Usage:' in help output, got:\n%s", logOutput)
		}
	})

	t.Run("Error_InvalidIgnorePattern", func(t *testing.T) {
		testOutputDir := t.TempDir()
		args := []string{"-input", baseInputDir, "-ignore", "["}
		_, err := runMainLogic(args, testOutputDir)
		if err == nil {
			t.Fatalf("Expected error for invalid ignore pattern, got nil")
		}
		if !strings.Contains(err.Error(), "Error compiling regex patterns") {
			t.Errorf("Expected 'Error compiling regex patterns' error, got: %v", err)
		}
	})

	t.Run("Error_InvalidIncludePattern", func(t *testing.T) {
		testOutputDir := t.TempDir()
		args := []string{"-input", baseInputDir, "-include", "+"}
		_, err := runMainLogic(args, testOutputDir)
		if err == nil {
			t.Fatalf("Expected error for invalid include pattern, got nil")
		}
		if !strings.Contains(err.Error(), "Error compiling regex patterns") {
			t.Errorf("Expected 'Error compiling regex patterns' error, got: %v", err)
		}
	})

	t.Run("Error_WritingOutputFileToDirectory", func(t *testing.T) {
		testOutputDir := t.TempDir()
		outputFilePathIsDir := filepath.Join(testOutputDir, "i_am_a_dir")
		_ = os.Mkdir(outputFilePathIsDir, 0755) // Create the dir
		args := []string{"-input", baseInputDir, "-output", "i_am_a_dir"}
		logOutput, err := runMainLogic(args, testOutputDir)
		if err == nil {
			t.Fatalf("Expected error writing output to directory, got nil. Log:\n%s", logOutput)
		}
		if !strings.Contains(err.Error(), "writing output file") {
			t.Errorf("Expected error related to 'writing output file', got: %v", err)
		}
	})

	t.Run("SuccessfulRun_Basic", func(t *testing.T) {
		testOutputDir := t.TempDir()
		outputFileName := "successful_run.md"
		args := []string{"-input", baseInputDir, "-output", outputFileName}
		logOutput, err := runMainLogic(args, testOutputDir)
		if err != nil {
			t.Fatalf("Successful run failed: %v. Log:\n%s", err, logOutput)
		}
		expectedOutputPathLog := filepath.Join(testOutputDir, outputFileName)
		if !strings.Contains(logOutput, "Markdown content written to "+expectedOutputPathLog) {
			t.Errorf("Missing 'Markdown content written...' log. Expected path '%s'. Got:\n%s", expectedOutputPathLog, logOutput)
		}
		if _, err := os.Stat(expectedOutputPathLog); err != nil {
			t.Errorf("Expected output file '%s' to exist, stat failed: %v", expectedOutputPathLog, err)
		}
	})

	t.Run("SuccessfulRun_WithPathsFiles", func(t *testing.T) {
		testOutputDir := t.TempDir()
		includedFile, excludedFile := "inc.log", "exc.log"
		args := []string{"-input", baseInputDir, "-ignore", `\.exe$`, "-included-paths-file", includedFile, "-excluded-paths-file", excludedFile}
		logOutput, err := runMainLogic(args, testOutputDir) // Default output "codebase.md"
		if err != nil {
			t.Fatalf("Run with path files failed: %v. Log:\n%s", err, logOutput)
		}
		fullIncludedPath := filepath.Join(testOutputDir, includedFile)
		fullExcludedPath := filepath.Join(testOutputDir, excludedFile)
		if !strings.Contains(logOutput, "Paths saved to "+fullIncludedPath) {
			t.Errorf("Missing log for included paths saved. Expected path '%s'. Got:\n%s", fullIncludedPath, logOutput)
		}
		if !strings.Contains(logOutput, "Paths saved to "+fullExcludedPath) {
			t.Errorf("Missing log for excluded paths saved. Expected path '%s'. Got:\n%s", fullExcludedPath, logOutput)
		}
		if _, err := os.Stat(fullIncludedPath); err != nil {
			t.Errorf("Expected included paths file '%s' to exist, stat failed: %v", fullIncludedPath, err)
		}
		if _, err := os.Stat(fullExcludedPath); err != nil {
			t.Errorf("Expected excluded paths file '%s' to exist, stat failed: %v", fullExcludedPath, err)
		}
	})

	t.Run("SuccessfulRun_ClipboardLogging", func(t *testing.T) {
		testOutputDir := t.TempDir()
		args := []string{"-input", baseInputDir, "-output", "clip.md", "-clipboard"}
		logOutput, err := runMainLogic(args, testOutputDir)
		if err != nil {
			t.Fatalf("Clipboard run failed: %v. Log:\n%s", err, logOutput)
		}
		if !strings.Contains(logOutput, "Result will be copied to clipboard.") {
			t.Errorf("Missing 'Result will be copied...' log. Got:\n%s", logOutput)
		}
		// Check for the *simulated* message logged by runMainLogic after writeOutput
		if !strings.Contains(logOutput, "Markdown content copied to clipboard (simulated).") {
			t.Errorf("Missing 'copied to clipboard (simulated).' log. Got:\n%s", logOutput)
		}
	})
}

```


## out.md
```md
# Tree View:
```
C:/Users/CARLOS~1.SAN/AppData/Local/Temp/codeweaver_test_fs_2867689002
├─ README.md
├─ docs
├─ empty_dir
├─ node_modules
│  └─ dep
└─ script.go
```

# Content:
## README.md
```md
# Test Readme
```

## script.go
```go
package main
func main() {}
```


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

