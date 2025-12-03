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
	instruction       string
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
	flag.StringVar(&cfg.instruction, "instruction", "", "Optional text to prepend to the generated Markdown file.")
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

		// Check for binary content
		if isBinary(fileContent) {
			cb.logger.Printf("Skipping binary file content: %s\n", pathRelToInput)
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
	fmt.Fprintf(os.Stderr, "  codeweaver -instruction \"Please analyze this code.\" # Add instruction at top\n")
	fmt.Fprintf(os.Stderr, "\nNotes on patterns:\n")
	fmt.Fprintf(os.Stderr, "  - Patterns are Go regular expressions (https://pkg.go.dev/regexp/syntax).\n")
	fmt.Fprintf(os.Stderr, "  - Paths for filtering are relative to the input directory (e.g., \"src/main.go\").\n")
	fmt.Fprintf(os.Stderr, "  - Use forward slashes '/' in patterns for cross-platform compatibility (e.g., \"data/images/\").\n")
	fmt.Fprintf(os.Stderr, "  - Use 'path/to/dir/?' to match a directory itself (with or without a trailing slash).\n")
	fmt.Fprintf(os.Stderr, "  - Use 'path/to/dir(/.*)?' to match a directory AND its contents.\n")
}
