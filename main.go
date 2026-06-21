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

// version, commit, and date are stamped at build time by goreleaser via -X ldflags.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
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

// defaultSensitivePatterns lists path patterns whose file CONTENTS are redacted by
// default. The files still appear in the tree view and included-paths list — only
// their body is replaced with a placeholder so secrets never reach the output.
// Disable with -no-default-redact; force full inclusion with -unsafe-include-secrets.
var defaultSensitivePatterns = []string{
	`(^|/)\.env(\..*)?$`,
	`(^|/)\.envrc$`,
	`\.pem$`,
	`\.key$`,
	`\.pfx$`,
	`\.p12$`,
	`\.keystore$`,
	`(^|/)id_(rsa|dsa|ecdsa|ed25519)(\.pub)?$`,
	`(^|/)\.npmrc$`,
	`(^|/)\.pypirc$`,
	`(^|/)\.netrc$`,
	`(^|/)credentials(\..*)?$`,
	`(^|/)\.aws/`,
	`(^|/)\.ssh/`,
	`\.tfstate$`,
	`\.tfvars$`,
	`(^|/)secrets?\.(ya?ml|json|toml|ini)$`,
}

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
			os.Exit(0)
		}
	}

	if cfg.showVersion {
		fmt.Printf("CodeWeaver %s\ncommit: %s\nbuilt:  %s\n", version, commit, date)
		os.Exit(0)
	}

	logger := setupLogging(cfg)
	ignoreMatchers, includeMatchers, err := compileMatchers(cfg, logger)
	if err != nil {
		logger.Fatalf("%sError compiling regex patterns: %v%s", colorRed, err, colorReset)
	}

	sensitiveMatchers, err := compileSensitiveMatchers(cfg)
	if err != nil {
		logger.Fatalf("%sError compiling sensitive patterns: %v%s", colorRed, err, colorReset)
	}
	if !cfg.unsafeIncludeSecrets && len(sensitiveMatchers) > 0 {
		logger.Printf("%sSensitive content redaction active (%d patterns). Use -unsafe-include-secrets to disable.%s",
			colorYellow, len(sensitiveMatchers), colorReset)
	}

	finalMarkdownString, pathsForIncludedFile, pathsForExcludedFile, err := generateMarkdown(cfg, ignoreMatchers, includeMatchers, sensitiveMatchers, logger)
	if err != nil {
		logger.Fatalf("%sError generating markdown: %v%s", colorRed, err, colorReset)
	}

	err = writeOutput(cfg, finalMarkdownString, pathsForIncludedFile, pathsForExcludedFile, logger)
	if err != nil {
		logger.Fatalf("%sError writing output: %v%s", colorRed, err, colorReset)
	}

	byteCount := len(finalMarkdownString)
	logger.Printf("Output: %d bytes (~%d tokens)", byteCount, byteCount/4)
	logger.Printf("%sCodeWeaver finished successfully.%s", colorGreen, colorReset)
}

// isFlagHelpError checks if the error returned by flag.Parse is due to a request for help.
func isFlagHelpError(err error) bool {
	return err != nil && err.Error() == "flag: help requested"
}

// config holds the runtime configuration options parsed from command-line arguments.
type config struct {
	inputDirOriginal    string
	inputDirAbs         string
	outputFile          string
	ignorePatterns      []string
	includePatterns     []string
	includedPathsFile   string
	excludedPathsFile   string
	instruction         string
	addToClipboard      bool
	showVersion         bool
	redactPatterns      []string
	noDefaultRedact     bool
	unsafeIncludeSecrets bool
	maxFileSizeBytes    int64
	rootMarker          string
	here                bool
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

	redactStr := flag.String("redact", "", "Comma-separated extra regex patterns whose matching files have contents redacted.")
	flag.BoolVar(&cfg.noDefaultRedact, "no-default-redact", false, "Disables the built-in sensitive-file redaction list.")
	flag.BoolVar(&cfg.unsafeIncludeSecrets, "unsafe-include-secrets", false, "Embeds ALL file contents including sensitive files. Use with caution.")
	flag.Int64Var(&cfg.maxFileSizeBytes, "max-file-size", 5_000_000, "Skip file contents larger than this many bytes (0 = no limit).")
	flag.StringVar(&cfg.rootMarker, "root-marker", ".git,go.mod,package.json,pyproject.toml,Cargo.toml", "Comma-separated filenames/dirnames; walk up from -input until one is found and use that dir as root.")
	flag.BoolVar(&cfg.here, "here", false, "Skip root-marker detection and use -input (or CWD) as-is.")

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
	if !cfg.here && cfg.rootMarker != "" {
		markers := strings.Split(cfg.rootMarker, ",")
		if found := findProjectRoot(cfg.inputDirAbs, markers); found != "" {
			cfg.inputDirAbs = found
		}
	}
	if *ignoreStr != "" {
		cfg.ignorePatterns = strings.Split(*ignoreStr, ",")
	}
	if *includeStr != "" {
		cfg.includePatterns = strings.Split(*includeStr, ",")
	}
	if *redactStr != "" {
		cfg.redactPatterns = strings.Split(*redactStr, ",")
	}
	return cfg, nil
}

// findProjectRoot walks up from dir until it finds a directory containing one of the
// marker filenames/dirnames, returning that ancestor directory. Returns "" if not found.
func findProjectRoot(dir string, markers []string) string {
	for {
		for _, marker := range markers {
			m := strings.TrimSpace(marker)
			if m == "" {
				continue
			}
			if _, err := os.Stat(filepath.Join(dir, m)); err == nil {
				return dir
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

// setupLogging initializes the logger and prints the startup configuration.
// When output is stdout ("-"), logs go to stderr so they don't corrupt the Markdown stream.
func setupLogging(cfg *config) *log.Logger {
	logDest := os.Stdout
	if cfg.outputFile == "-" {
		logDest = os.Stderr
	}
	logger := log.New(logDest, "", 0)
	logger.Println("Starting CodeWeaver...")
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

// compileSensitiveMatchers builds the list of regexes used for content redaction.
// It merges the built-in default list (unless -no-default-redact) with any
// user-supplied -redact patterns.
func compileSensitiveMatchers(cfg *config) ([]*regexp.Regexp, error) {
	if cfg.unsafeIncludeSecrets {
		return nil, nil
	}
	var patterns []string
	if !cfg.noDefaultRedact {
		patterns = append(patterns, defaultSensitivePatterns...)
	}
	patterns = append(patterns, cfg.redactPatterns...)
	compiled := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		trimmed := strings.TrimSpace(p)
		if trimmed == "" {
			continue
		}
		r, err := regexp.Compile(trimmed)
		if err != nil {
			return nil, fmt.Errorf("sensitive pattern '%s': %w", trimmed, err)
		}
		compiled = append(compiled, r)
	}
	return compiled, nil
}

// isSensitive returns true if the relative path matches any sensitive content pattern.
func isSensitive(pathRelToInput string, sensitiveMatchers []*regexp.Regexp) (bool, string) {
	for _, pattern := range sensitiveMatchers {
		if pattern != nil && pattern.MatchString(pathRelToInput) {
			return true, pattern.String()
		}
	}
	return false, ""
}

// --- Markdown Generation ---

// generateMarkdown orchestrates the creation of the tree view and content sections.
func generateMarkdown(cfg *config, ignoreMatchers, includeMatchers, sensitiveMatchers []*regexp.Regexp, logger *log.Logger) (
	finalMarkdown string, processedPathsForFile []string, excludedPathsForFile []string, err error) {

	var fullMarkdown strings.Builder

	logger.Println("Processing paths and building content section...")
	contentBuilder := newContentBuilder(
		cfg.inputDirAbs, cfg.includedPathsFile, cfg.excludedPathsFile,
		ignoreMatchers, includeMatchers, sensitiveMatchers,
		cfg.maxFileSizeBytes,
		logger,
	)
	contentMarkdown, processedPaths, excludedPaths, err := contentBuilder.buildContentString()
	if err != nil {
		return "", nil, nil, fmt.Errorf("failed to build code content: %w", err)
	}

	if cfg.instruction != "" {
		fullMarkdown.WriteString(cfg.instruction)
		fullMarkdown.WriteString("\n\n")
	}

	logger.Println("Building tree view...")
	fullMarkdown.WriteString(mdTreeViewHeader)
	fullMarkdown.WriteString(filepath.ToSlash(cfg.inputDirOriginal) + "\n")

	processedPathsSet := make(map[string]struct{}, len(processedPaths))
	for _, p := range processedPaths {
		processedPathsSet[p] = struct{}{}
	}

	treeBuilder := newTreeBuilder(cfg.inputDirAbs, processedPathsSet)
	treeString, err := treeBuilder.buildTreeString()
	if err != nil {
		return "", nil, nil, fmt.Errorf("failed to build codebase tree: %w", err)
	}
	fullMarkdown.WriteString(treeString)
	fullMarkdown.WriteString(mdCodeBlockEnd)

	fullMarkdown.WriteString(mdContentHeader)
	fullMarkdown.WriteString(contentMarkdown)
	logger.Println()

	return fullMarkdown.String(), processedPaths, excludedPaths, nil
}

// --- Tree Builder ---

// treeBuilder is responsible for generating the visual directory tree structure.
type treeBuilder struct {
	rootAbsPath       string
	processedPathsSet map[string]struct{}
	output            strings.Builder
	depthOpen         map[int]bool
}

// newTreeBuilder creates a new tree builder instance.
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
	ignoreMatchers, includeMatchers, sensitiveMatchers []*regexp.Regexp
	maxFileSizeBytes                                   int64
	logger                                            *log.Logger
}

// newContentBuilder initializes a new contentBuilder.
func newContentBuilder(
	rootAbsPath, includedPathsFile, excludedPathsFile string,
	ignoreMatchers, includeMatchers, sensitiveMatchers []*regexp.Regexp,
	maxFileSizeBytes int64,
	logger *log.Logger,
) *contentBuilder {
	return &contentBuilder{
		rootAbsPath:       rootAbsPath,
		includedPathsFile: includedPathsFile,
		excludedPathsFile: excludedPathsFile,
		ignoreMatchers:    ignoreMatchers,
		includeMatchers:   includeMatchers,
		sensitiveMatchers: sensitiveMatchers,
		maxFileSizeBytes:  maxFileSizeBytes,
		logger:            logger,
	}
}

// buildContentString scans the directory and returns:
// 1. markdownContent: The combined markdown string for file contents.
// 2. allProcessedPaths: List of files that passed filters (used for tree building).
// 3. excludedPaths: List of files/dirs that failed filters.
// 4. err: Any error encountered.
func (cb *contentBuilder) buildContentString() (
	markdownContent string, allProcessedPaths []string, excludedPaths []string, err error) {

	var contentSB strings.Builder

	var localProcessedPaths []string
	var localExcludedPaths []string

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
		}

		if !shouldProcess(pathRelToInput, cb.ignoreMatchers, cb.includeMatchers) {
			if cb.excludedPathsFile == "" {
				cb.logger.Printf("%s%s %s%s\n", colorRed, logPrefixExclude, pathRelToInput, colorReset)
			}
			localExcludedPaths = append(localExcludedPaths, pathRelToInput)
			if !d.IsDir() {
				ext := strings.ToLower(filepath.Ext(pathRelToInput))
				if ext == "" {
					ext = "(no ext)"
				}
				excludedExtensions[ext] = struct{}{}
			}
			return nil
		}

		// S2: Symlink escape guard — prevent reading files whose real path is outside
		// the input root. WalkDir uses lstat so symlinked dirs aren't traversed, but
		// symlinked files would be silently read by os.ReadFile.
		if d.Type()&fs.ModeSymlink != 0 {
			realPath, symlinkErr := filepath.EvalSymlinks(currentWalkPath)
			if symlinkErr != nil {
				cb.logger.Printf("%sWarning: Cannot resolve symlink %s: %v — skipping%s\n",
					colorYellow, pathRelToInput, symlinkErr, colorReset)
				localExcludedPaths = append(localExcludedPaths, pathRelToInput)
				return nil
			}
			relReal, relErr := filepath.Rel(cb.rootAbsPath, realPath)
			if relErr != nil || strings.HasPrefix(filepath.ToSlash(relReal), "..") {
				cb.logger.Printf("%sWarning: Skipping symlink %s — target escapes root (%s)%s\n",
					colorYellow, pathRelToInput, realPath, colorReset)
				localExcludedPaths = append(localExcludedPaths, pathRelToInput)
				return nil
			}
		}

		if cb.includedPathsFile == "" {
			if d.IsDir() {
				cb.logger.Printf("%s %s\n", logPrefixInclude, pathRelToInput)
			} else {
				cb.logger.Printf("%s%s %s%s\n", colorGreen, logPrefixInclude, pathRelToInput, colorReset)
			}
		}

		if d.IsDir() {
			return nil
		}

		localProcessedPaths = append(localProcessedPaths, pathRelToInput)
		ext := strings.ToLower(filepath.Ext(pathRelToInput))
		if ext == "" {
			ext = "(no ext)"
		}
		includedExtensions[ext] = struct{}{}

		// S4: Skip oversized files rather than loading them into memory.
		if cb.maxFileSizeBytes > 0 {
			info, infoErr := d.Info()
			if infoErr == nil && info.Size() > cb.maxFileSizeBytes {
				cb.logger.Printf("%sSkipping large file %s (%d bytes > -max-file-size %d)%s\n",
					colorYellow, pathRelToInput, info.Size(), cb.maxFileSizeBytes, colorReset)
				contentSB.WriteString(fmt.Sprintf("%s%s\n\n[content skipped — file size %d bytes exceeds -max-file-size %d]\n\n",
					mdFileHeaderStart, pathRelToInput, info.Size(), cb.maxFileSizeBytes))
				return nil
			}
		}

		// S1: Redact contents of sensitive files. The file still appears in the tree
		// and included-paths list (already appended above); only its body is suppressed.
		if sensitive, matchedPattern := isSensitive(pathRelToInput, cb.sensitiveMatchers); sensitive {
			cb.logger.Printf("%sRedacting sensitive file: %s (matched: %s)%s\n",
				colorYellow, pathRelToInput, matchedPattern, colorReset)
			contentSB.WriteString(fmt.Sprintf("%s%s\n\n[content redacted — matched sensitive pattern: %s]\n\n",
				mdFileHeaderStart, pathRelToInput, matchedPattern))
			return nil
		}

		fileContent, readErr := os.ReadFile(currentWalkPath)
		if readErr != nil {
			cb.logger.Printf("%sWarning: Failed to read file %s: %v%s\n", colorRed, currentWalkPath, readErr, colorReset)
			contentSB.WriteString(fmt.Sprintf("%s%s\n```\nError reading file: %v\n```\n\n", mdFileHeaderStart, pathRelToInput, readErr))
			return nil
		}

		if len(fileContent) == 0 {
			return nil
		}

		if isBinary(fileContent) {
			cb.logger.Printf("Skipping binary file content: %s\n", pathRelToInput)
			return nil
		}

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
func isBinary(content []byte) bool {
	const maxBytesToCheck = 1024
	checkLen := len(content)
	if checkLen > maxBytesToCheck {
		checkLen = maxBytesToCheck
	}
	return bytes.IndexByte(content[:checkLen], 0) != -1
}

// countMaxBackticks calculates the maximum number of consecutive backticks in the byte slice.
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

// writeOutput writes the generated markdown to the output file (or stdout when "-"),
// saves included/excluded path lists, and optionally copies the result to the clipboard.
func writeOutput(cfg *config, markdownContent string, includedPaths, excludedPaths []string, logger *log.Logger) error {
	if cfg.outputFile == "-" {
		if _, err := fmt.Fprint(os.Stdout, markdownContent); err != nil {
			return fmt.Errorf("writing markdown to stdout: %w", err)
		}
	} else {
		outputFileSlash := filepath.ToSlash(cfg.outputFile)
		logger.Printf("Writing output to %s...", outputFileSlash)
		// 0600: aggregate may contain sensitive data; restrict to owner only (S3).
		err := os.WriteFile(cfg.outputFile, []byte(markdownContent), 0600)
		if err != nil {
			logger.Printf("%sError writing to output file %s: %v%s", colorRed, outputFileSlash, err, colorReset)
			return fmt.Errorf("writing output file %s: %w", cfg.outputFile, err)
		}
		logger.Printf("Markdown content written to %s", outputFileSlash)
	}

	if cfg.includedPathsFile != "" {
		includedFileSlash := filepath.ToSlash(cfg.includedPathsFile)
		if err := savePathsToFile(cfg.includedPathsFile, includedPaths, logger); err != nil {
			logger.Printf("%sWarning: Error saving included paths to %s: %v%s", colorRed, includedFileSlash, err, colorReset)
		}
	}
	if cfg.excludedPathsFile != "" {
		excludedFileSlash := filepath.ToSlash(cfg.excludedPathsFile)
		if err := savePathsToFile(cfg.excludedPathsFile, excludedPaths, logger); err != nil {
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
	sort.Strings(paths)
	var sb strings.Builder
	for _, p := range paths {
		sb.WriteString(p)
		sb.WriteString("\n")
	}
	// 0600: restrict to owner; these files list every path in the scanned tree (S3).
	err := os.WriteFile(filename, []byte(sb.String()), 0600)
	if err == nil {
		logger.Printf("Paths saved to %s", filepath.ToSlash(filename))
	} else {
		return fmt.Errorf("writing paths file %s: %w", filepath.ToSlash(filename), err)
	}
	return nil
}

// printHelp displays the application's usage information, available flags, and examples.
func printHelp() {
	fmt.Fprintf(os.Stderr, "%s%sCodeWeaver%s: Generate Markdown Documentation from Your Codebase.\n\n", colorBold, colorGreen, colorReset)

	fmt.Fprintf(os.Stderr, "%sUsage:%s\n", colorCyan, colorReset)
	fmt.Fprintf(os.Stderr, "  codeweaver [flags]\n\n")

	fmt.Fprintf(os.Stderr, "%sFlags:%s\n", colorCyan, colorReset)
	flag.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(os.Stderr, "  %s-%-20s%s\n", colorGreen, f.Name, colorReset)
		fmt.Fprintf(os.Stderr, "      %s", f.Usage)

		if f.DefValue != "" {
			if f.Name == "clipboard" && f.DefValue == "false" {
				// skip
			} else if f.Name == "version" && f.DefValue == "false" {
				// skip
			} else if f.Name == "no-default-redact" && f.DefValue == "false" {
				// skip
			} else if f.Name == "unsafe-include-secrets" && f.DefValue == "false" {
				// skip
			} else {
				fmt.Fprintf(os.Stderr, " %s(Default: %s)%s", colorYellow, f.DefValue, colorReset)
			}
		}
		fmt.Fprintf(os.Stderr, "\n")
	})
	fmt.Fprintln(os.Stderr)

	fmt.Fprintf(os.Stderr, "%sExamples:%s\n", colorCyan, colorReset)
	fmt.Fprintf(os.Stderr, "  %scodeweaver%s\n", colorGreen, colorReset)
	fmt.Fprintf(os.Stderr, "    Process current directory, output to codebase.md\n\n")

	fmt.Fprintf(os.Stderr, "  %scodeweaver -input src -output source_docs.md%s\n", colorGreen, colorReset)
	fmt.Fprintf(os.Stderr, "    Specify input directory and output filename\n\n")

	fmt.Fprintf(os.Stderr, "  %scodeweaver -ignore \"build/,vendor/\" -include \"\\.go$,\\.md$\"%s\n", colorGreen, colorReset)
	fmt.Fprintf(os.Stderr, "    Exclude 'build' and 'vendor' folders, but ONLY include .go and .md files\n\n")

	fmt.Fprintf(os.Stderr, "  %scodeweaver -instruction \"Analyze this code\" -clipboard%s\n", colorGreen, colorReset)
	fmt.Fprintf(os.Stderr, "    Prepend instruction and copy result to clipboard\n\n")

	fmt.Fprintf(os.Stderr, "  %scodeweaver -output -%s\n", colorGreen, colorReset)
	fmt.Fprintf(os.Stderr, "    Stream Markdown to stdout (logs go to stderr); useful for piping to LLMs\n\n")

	fmt.Fprintf(os.Stderr, "  %scodeweaver -here%s\n", colorGreen, colorReset)
	fmt.Fprintf(os.Stderr, "    Use CWD as-is, skipping root-marker detection\n\n")

	fmt.Fprintf(os.Stderr, "%sSensitive File Handling:%s\n", colorCyan, colorReset)
	fmt.Fprintf(os.Stderr, "  By default, files matching common secret patterns (.env, *.pem, id_rsa, etc.)\n")
	fmt.Fprintf(os.Stderr, "  appear in the tree but have their CONTENTS replaced with a redaction notice.\n")
	fmt.Fprintf(os.Stderr, "  Use %s-no-default-redact%s to disable built-in patterns.\n", colorLiteRed, colorReset)
	fmt.Fprintf(os.Stderr, "  Use %s-unsafe-include-secrets%s to embed all contents (use with caution).\n\n", colorLiteRed, colorReset)

	fmt.Fprintf(os.Stderr, "%sHow Filters Work:%s\n", colorCyan, colorReset)
	fmt.Fprintf(os.Stderr, "  1. %s-ignore%s (Blacklist): Matches are excluded. Checked first.\n", colorLiteRed, colorReset)
	fmt.Fprintf(os.Stderr, "  2. %s-include%s (Whitelist): If specified, ONLY matches are included.\n", colorGreen, colorReset)
	fmt.Fprintf(os.Stderr, "     If a path matches both (unlikely given logic), ignore takes precedence.\n")
	fmt.Fprintf(os.Stderr, "     Directories matching -ignore are skipped entirely.\n\n")

	fmt.Fprintf(os.Stderr, "%sRegex Notes:%s\n", colorCyan, colorReset)
	fmt.Fprintf(os.Stderr, "  - Patterns are Go regular expressions.\n")
	fmt.Fprintf(os.Stderr, "  - Use forward slashes '/' for paths (e.g., \"dir/file.txt\").\n")
	fmt.Fprintf(os.Stderr, "  - Match a file extension: \"\\.go$\"\n")
	fmt.Fprintf(os.Stderr, "  - Match a directory: \"^vendor/\" or \"/vendor/\"\n")
}
