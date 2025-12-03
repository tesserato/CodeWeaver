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
