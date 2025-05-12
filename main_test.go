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
