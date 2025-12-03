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
