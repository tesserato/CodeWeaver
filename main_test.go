package main

import (
	"io/ioutil" // For ioutil.Discard
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// mustCompileRegex is a helper for tests to compile regex patterns.
// It panics if compilation fails, as this indicates a test setup error.
func mustCompileRegex(pattern string) *regexp.Regexp {
	if pattern == "" {
		return nil // Allow empty patterns to result in nil, which is handled by shouldProcess
	}
	r, err := regexp.Compile(pattern)
	if err != nil {
		log.Fatalf("Failed to compile regex '%s': %v", pattern, err)
	}
	return r
}

// normalizeNewlines ensures consistent line endings for string comparisons.
func normalizeNewlines(s string) string {
	return strings.ReplaceAll(s, "\r\n", "\n")
}

func TestShouldProcess(t *testing.T) {
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
		{"IncludeDir_Exact", "data/", nil, []*regexp.Regexp{mustCompileRegex(`^data/$`)}, true},
		{"IncludeDir_NoMatch", "other/", nil, []*regexp.Regexp{mustCompileRegex(`^data/$`)}, false},
		{"IgnoreTakesPrecedence_PathMatchesBoth", "vendor/lib.go", []*regexp.Regexp{mustCompileRegex(`^vendor/`)}, []*regexp.Regexp{mustCompileRegex(`\.go$`)}, false},
		{"IncludeSatisfied_NotIgnored", "src/main.go", []*regexp.Regexp{mustCompileRegex(`^vendor/`)}, []*regexp.Regexp{mustCompileRegex(`\.go$`)}, true},
		{"MultipleIgnore_FirstMatch", "config/dev.log", []*regexp.Regexp{mustCompileRegex(`\.log$`), mustCompileRegex(`^config/`)}, nil, false},
		{"MultipleIgnore_SecondMatch", "build/app.exe", []*regexp.Regexp{mustCompileRegex(`\.log$`), mustCompileRegex(`\.exe$`)}, nil, false},
		{"MultipleInclude_OneMatchSufficient", "image.png", nil, []*regexp.Regexp{mustCompileRegex(`\.jpg$`), mustCompileRegex(`\.png$`)}, true},
		{"MultipleInclude_NoMatch", "archive.zip", nil, []*regexp.Regexp{mustCompileRegex(`\.jpg$`), mustCompileRegex(`\.png$`)}, false},
		{"PathIsEmptyString_Root_NoFilters", "", nil, nil, true},
		{"PathIsEmptyString_Root_Ignored", "", []*regexp.Regexp{mustCompileRegex(`^$`)}, nil, false}, // Regex `^$` matches empty string
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

// createTestFS creates a temporary directory structure for integration-like tests.
// It returns the absolute path to the root of the test structure and a cleanup function.
func createTestFS(t *testing.T) (string, func()) {
	t.Helper()
	rootDir, err := ioutil.TempDir("", "codeweaver_test_fs_")
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
	}

	for relPath, content := range structure {
		absPath := filepath.Join(rootDir, relPath)
		if strings.HasSuffix(relPath, "/") { // It's a directory
			if err := os.MkdirAll(absPath, 0755); err != nil {
				os.RemoveAll(rootDir) // Cleanup on error
				t.Fatalf("Failed to create dir %s: %v", absPath, err)
			}
		} else { // It's a file
			// Ensure parent directory exists first
			parentDir := filepath.Dir(absPath)
			if err := os.MkdirAll(parentDir, 0755); err != nil {
				os.RemoveAll(rootDir)
				t.Fatalf("Failed to create parent dir %s for file %s: %v", parentDir, absPath, err)
			}
			if err := ioutil.WriteFile(absPath, []byte(content), 0644); err != nil {
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

func TestTreeBuilder(t *testing.T) {
	rootDir, cleanup := createTestFS(t)
	defer cleanup()

	testCases := []struct {
		name            string
		ignoreMatchers  []*regexp.Regexp
		includeMatchers []*regexp.Regexp
		expectedTree    string
	}{
		{
			name: "NoFilters",
			// Sorted and prefixed
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
└─ script.go
`),
		},
		{
			name:           "IgnoreDotGitAndNodeModulesDirs",
			ignoreMatchers: []*regexp.Regexp{mustCompileRegex(`^\.git/?`), mustCompileRegex(`^node_modules/?`)},
			// Sorted and prefixed
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
└─ script.go
`),
		},
		{
			name:            "IncludeOnlyGoAndMdFiles",
			includeMatchers: []*regexp.Regexp{mustCompileRegex(`\.go$`), mustCompileRegex(`\.md$`)},
			// Sorted and prefixed
			expectedTree: normalizeNewlines(`├─ README.md
└─ script.go
`),
		},
		{
			name: "IncludeYamlFilesAndDataDir",
			// Include data dir itself (and its contents if they also match an include rule or if no other file-specific include rule prevents them)
			// and .yaml files anywhere.
			// ^data/? will match "data", "data/config.yaml", "data/image.png"
			// \.yaml$ will match "data/config.yaml"
			// So, 'data' directory is listed.
			// Inside 'data', 'config.yaml' matches both. 'image.png' matches '^data/?'.
			includeMatchers: []*regexp.Regexp{mustCompileRegex(`\.yaml$`), mustCompileRegex(`^data/?`)},
			// Sorted and prefixed
			expectedTree: normalizeNewlines(`└─ data
   ├─ config.yaml
   └─ image.png
`),
		},
		{
			name:            "IgnoreExe_IncludeDocsDirAndSubDocsFile",
			ignoreMatchers:  []*regexp.Regexp{mustCompileRegex(`\.exe$`)},
			includeMatchers: []*regexp.Regexp{mustCompileRegex(`^docs(/.*)?`), mustCompileRegex(`file_in_sub\.txt$`)},
			// Only "docs" dir and its specified contents should appear. "docs" is the only top-level item.
			// Sorted and prefixed
			expectedTree: normalizeNewlines(`└─ docs
   └─ sub_docs
      └─ file_in_sub.txt
`),
		},
		{
			name:            "EmptyDir_Included",
			includeMatchers: []*regexp.Regexp{mustCompileRegex(`^empty_dir/?`)},
			// Sorted and prefixed
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
}

func TestContentBuilder(t *testing.T) {
	rootDir, cleanup := createTestFS(t)
	defer cleanup()

	testLogger := log.New(ioutil.Discard, "", 0)

	testCases := []struct {
		name                  string
		ignoreMatchers        []*regexp.Regexp
		includeMatchers       []*regexp.Regexp
		expectedContentSubstr string
		expectedIncludedPaths []string // Sorted
		expectedExcludedPaths []string // Sorted
	}{
		{
			name:                  "NoFilters",
			expectedContentSubstr: "## file1.txt\n```txt\ncontent of file1\n```",
			expectedIncludedPaths: []string{
				".git/HEAD", "README.md", "build/output.exe", "build/tmp/log.txt",
				"data/config.yaml", "data/image.png", "docs/sub_docs/file_in_sub.txt",
				"file1.txt", "node_modules/dep/package.json", "script.go",
			},
			expectedExcludedPaths: []string{},
		},
		{
			name:                  "IgnoreDotGitDirAndExeFiles",
			ignoreMatchers:        []*regexp.Regexp{mustCompileRegex(`^\.git/?`), mustCompileRegex(`\.exe$`)},
			expectedContentSubstr: "## README.md\n```md\n# Test Readme\n```",
			expectedIncludedPaths: []string{
				"README.md", "build/tmp/log.txt", "data/config.yaml", "data/image.png",
				"docs/sub_docs/file_in_sub.txt", "file1.txt", "node_modules/dep/package.json", "script.go",
			},
			// .git dir is skipped, output.exe file is skipped.
			expectedExcludedPaths: []string{".git", "build/output.exe"},
		},
		{
			name:                  "IncludeOnlyGoAndMdFiles",
			includeMatchers:       []*regexp.Regexp{mustCompileRegex(`\.go$`), mustCompileRegex(`\.md$`)},
			expectedContentSubstr: "## script.go\n```go\npackage main",
			expectedIncludedPaths: []string{"README.md", "script.go"},
			// Directories not matching .go or .md are skipped.
			// Files not matching .go or .md are individually excluded if their parent dir was traversed (not the case here for most).
			expectedExcludedPaths: []string{
				".git", "build", "data", "docs", "empty_dir", "file1.txt", "node_modules",
			},
		},
		{
			name: "ContentOfSpecificFile_Yaml_StrictInclude", // Test name reflects strict include for dirs
			// With only "config.yaml$" as include, "data" dir itself won't match and will be skipped.
			includeMatchers:       []*regexp.Regexp{mustCompileRegex(`config\.yaml$`)},
			expectedContentSubstr: "", // No content because 'data' dir is skipped
			expectedIncludedPaths: []string{},
			// All top-level items are excluded as they don't match "config.yaml$"
			// and "data" dir gets skipped, so its children are not processed.
			expectedExcludedPaths: []string{
				".git", "README.md", "build", "data", "docs", "empty_dir", "file1.txt", "node_modules", "script.go",
			},
		},
		{
			name: "ContentOfSpecificFile_Yaml_PermissiveIncludeDir",
			// To include data/config.yaml, data dir must also be allowed by an include rule.
			includeMatchers:       []*regexp.Regexp{mustCompileRegex(`config\.yaml$`), mustCompileRegex(`^data/?`)},
			expectedContentSubstr: normalizeNewlines("## data/config.yaml\n```yaml\nkey: value\n```\n\n"),
			// data/image.png also gets included because `^data/?` matches it and it's not ignored.
			expectedIncludedPaths: []string{"data/config.yaml", "data/image.png"},
			expectedExcludedPaths: []string{
				".git", "README.md", "build", "docs", "empty_dir", "file1.txt", "node_modules", "script.go",
			},
		},

		{
			name:                  "EmptyDir_NotIncludedInContent",
			includeMatchers:       []*regexp.Regexp{mustCompileRegex(`^empty_dir/?`)},
			expectedContentSubstr: "",
			expectedIncludedPaths: []string{},
			expectedExcludedPaths: []string{
				".git", "README.md", "build", "data", "docs", "file1.txt", "node_modules", "script.go",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Sort expected paths for reliable comparison
			sort.Strings(tc.expectedIncludedPaths)
			sort.Strings(tc.expectedExcludedPaths)

			builder := newContentBuilder(rootDir,"", "", tc.ignoreMatchers, tc.includeMatchers, testLogger)
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

			// Check included paths
			if !equalStringSlices(actualIncludedPaths, tc.expectedIncludedPaths) {
				t.Errorf("'%s': Included paths mismatch.\nExpected: %v\nGot:      %v", tc.name, tc.expectedIncludedPaths, actualIncludedPaths)
			}

			// Check excluded paths
			if !equalStringSlices(actualExcludedPaths, tc.expectedExcludedPaths) {
				t.Errorf("'%s': Excluded paths mismatch.\nExpected: %v\nGot:      %v", tc.name, tc.expectedExcludedPaths, actualExcludedPaths)
			}
		})
	}
}

// Helper function to compare string slices
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

func TestSavePathsToFile(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "test_paths_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpFilePath := tmpFile.Name()
	tmpFile.Close() // Close immediately, savePathsToFile will recreate/truncate
	defer os.Remove(tmpFilePath)

	testLogger := log.New(ioutil.Discard, "", 0)

	paths := []string{"path/to/file1.txt", "another/path.go", "root_file.md"}
	expectedContent := "path/to/file1.txt\nanother/path.go\nroot_file.md\n"

	err = savePathsToFile(tmpFilePath, paths, testLogger)
	if err != nil {
		t.Fatalf("savePathsToFile failed: %v", err)
	}

	contentBytes, err := ioutil.ReadFile(tmpFilePath)
	if err != nil {
		t.Fatalf("Failed to read back saved paths file: %v", err)
	}

	if normalizeNewlines(string(contentBytes)) != normalizeNewlines(expectedContent) {
		t.Errorf("Content mismatch for saved paths file.\nExpected:\n%s\nGot:\n%s", expectedContent, string(contentBytes))
	}

	// Test with empty paths
	emptyPathsFile := filepath.Join(filepath.Dir(tmpFilePath), "empty_paths.txt")
	defer os.Remove(emptyPathsFile)
	err = savePathsToFile(emptyPathsFile, []string{}, testLogger)
	if err != nil {
		t.Fatalf("savePathsToFile with empty paths failed: %v", err)
	}
	// Check if file exists and is empty (or doesn't exist, current logic returns nil if no paths)
	// Current logic returns nil and does nothing if paths is empty. So file shouldn't be created.
	if _, err := os.Stat(emptyPathsFile); !os.IsNotExist(err) {
		t.Errorf("Expected empty_paths.txt to not be created for empty paths, but it was (or other error: %v)", err)
	}
}

// Note: Testing parseFlags is more involved as it interacts with os.Args and os.Exit.
// It's often tested by running the binary with different arguments in an integration test setup,
// or by refactoring it to be more testable (e.g., taking args as a slice).
// For now, the core logic functions (shouldProcess, treeBuilder, contentBuilder) are prioritized.
