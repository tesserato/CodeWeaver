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
		// parseFlags now handles printing help if -help is used.
		// If flag.Parse() itself returned an error (e.g., unknown flag), exit.
		// flag.ErrHelp is not expected here anymore if Usage is set correctly.
		fmt.Fprintf(os.Stderr, "%sError parsing flags: %v%s\n", colorRed, err, colorReset)
		// Provide context for help flag usage
		if !errors.Is(err, flag.ErrHelp) { // flag.ErrHelp is less likely now, but check just in case
			fmt.Fprintln(os.Stderr, "Use --help for usage information.")
		}
		os.Exit(2) // Standard exit code for command line usage errors
	}
	for _, arg := range os.Args[1:] { // Check args excluding the program name
		arg_lower := strings.ToLower(arg)
		if arg_lower == "-help" || arg_lower == "--help" || arg_lower == "h" || arg_lower == "help" {
			printHelp()
			os.Exit(0) // Exit successfully after showing help
		}
	}

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
		logger.Fatalf("%sError writing output: %v%s", colorRed, err, colorReset)
	}

	logger.Printf("%sCodeWeaver finished successfully.%s", colorGreen, colorReset)
}

// --- Configuration & Setup ---

// config holds the application's configuration values derived from flags.
type config struct {
	inputDirOriginal  string
	inputDirAbs       string
	outputFile        string
	ignorePatterns    []string
	includePatterns   []string
	includedPathsFile string
	excludedPathsFile string
	addToClipboard    bool
	showHelp          bool // Keep for potential future use, though Usage handles exit
	showVersion       bool
}

// parseFlags defines, parses, and validates command-line flags using the global flag package.
// It returns a config struct or an error.
func parseFlags() (*config, error) {
	cfg := &config{}

	// --- Define flags directly on the global 'flag' package ---
	flag.StringVar(&cfg.inputDirOriginal, "input", ".", "The root directory to scan.")
	flag.StringVar(&cfg.outputFile, "output", "codebase.md", "The name of the output Markdown file.")
	ignoreStr := flag.String("ignore", `\.git.*`, "Comma-separated list of regular expressions for paths to *exclude*.")
	includeStr := flag.String("include", "", "Comma-separated list of regular expressions. *Only* paths matching these are *included*.")
	flag.StringVar(&cfg.includedPathsFile, "included-paths-file", "", "Saves the list of *included* paths to this file.")
	flag.StringVar(&cfg.excludedPathsFile, "excluded-paths-file", "", "Saves the list of *excluded* paths to this file.")
	flag.BoolVar(&cfg.addToClipboard, "clipboard", false, "Copies the generated Markdown to the clipboard.")
	flag.BoolVar(&cfg.showVersion, "version", false, "Displays the version and exits.")
	// The -help flag is automatically handled by the flag package if Usage is set.
	// We still define cfg.showHelp in case we need the value later, but don't need to explicitly check it for exit.
	flag.BoolVar(&cfg.showHelp, "help", false, "Displays help message and exits.")

	// --- Set custom Usage ---
	// This function will be called by the flag package if -help is used or if parsing fails.
	flag.Usage = func() {
		printHelp()
		os.Exit(0) // Exit successfully after showing help
	}

	// --- Parse using global flag package ---
	flag.Parse()

	// Handle version flag exit *after* parsing is successful
	if cfg.showVersion {
		// Return the config state; main will handle printing version and exiting.
		return cfg, nil
	}

	// --- Post-parsing Validation (only if not showing help/version) ---
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

type treeBuilder struct {
	rootAbsPath     string
	ignoreMatchers  []*regexp.Regexp
	includeMatchers []*regexp.Regexp
	output          strings.Builder
	depthOpen       map[int]bool
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
		return fmt.Errorf("read directory %s: %w", currentDirPath, err)
	}

	var filteredEntries []fs.DirEntry
	for _, entry := range entries {
		fullEntryPath := filepath.Join(currentDirPath, entry.Name())
		pathRelToInput, err := tb.getRelativePath(fullEntryPath)
		if err != nil {
			return err
		}
		if shouldProcess(pathRelToInput, tb.ignoreMatchers, tb.includeMatchers) {
			filteredEntries = append(filteredEntries, entry)
		}
	}

	sort.Slice(filteredEntries, func(i, j int) bool {
		return strings.ToLower(filteredEntries[i].Name()) < strings.ToLower(filteredEntries[j].Name())
	})

	for i, entry := range filteredEntries {
		isLastEntry := (i == len(filteredEntries)-1)
		tb.printEntryLine(entry, depth, isLastEntry)

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
	tb.output.WriteString(prefix.String())
	tb.output.WriteString(entry.Name())
	tb.output.WriteString("\n")
}

func (tb *treeBuilder) getRelativePath(fullPath string) (string, error) {
	pathRelToInput, err := filepath.Rel(tb.rootAbsPath, fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to make path %s relative to %s: %w", fullPath, tb.rootAbsPath, err)
	}
	return filepath.ToSlash(pathRelToInput), nil
}

// --- Content Builder ---

type contentBuilder struct {
	rootAbsPath       string
	includedPathsFile string
	excludedPathsFile string
	ignoreMatchers    []*regexp.Regexp
	includeMatchers   []*regexp.Regexp
	logger            *log.Logger
}

func newContentBuilder(rootAbsPath string, includedPathsFile string, excludedPathsFile string, ignoreMatchers []*regexp.Regexp, includeMatchers []*regexp.Regexp, logger *log.Logger) *contentBuilder {
	return &contentBuilder{
		rootAbsPath:       rootAbsPath,
		includedPathsFile: includedPathsFile,
		excludedPathsFile: excludedPathsFile,
		ignoreMatchers:    ignoreMatchers,
		includeMatchers:   includeMatchers,
		logger:            logger,
	}
}

func (cb *contentBuilder) buildContentString() (content string, includedPaths []string, excludedPaths []string, err error) {
	var contentSB strings.Builder

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

		process := shouldProcess(pathRelToInput, cb.ignoreMatchers, cb.includeMatchers)

		if !process {
			if cb.excludedPathsFile == "" {
				cb.logger.Printf("%s%s %s%s\n", colorRed, logPrefixExclude, pathRelToInput, colorReset)
			}
			excludedPaths = append(excludedPaths, pathRelToInput)
			return nil
		}

		if d.IsDir() {
			return nil
		}

		if cb.includedPathsFile == "" {
			cb.logger.Printf("%s%s %s%s\n", colorGreen, logPrefixInclude, pathRelToInput, colorReset)
		}
		includedPaths = append(includedPaths, pathRelToInput)

		fileContent, readErr := os.ReadFile(currentWalkPath)
		if readErr != nil {
			cb.logger.Printf("%sWarning: Failed to read file %s: %v%s\n", colorRed, currentWalkPath, readErr, colorReset)
			contentSB.WriteString(fmt.Sprintf("%s%s\n", mdFileHeaderStart, pathRelToInput))
			contentSB.WriteString(fmt.Sprintf("%s\nError reading file: %v%s", mdCodeBlockStart, readErr, mdCodeBlockEndNL))
			return nil
		}

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
	return contentSB.String(), includedPaths, excludedPaths, nil
}

// --- Path Filtering Logic ---

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

// --- Output Handling ---

func writeOutput(cfg *config, markdownContent string, includedPaths, excludedPaths []string, logger *log.Logger) error {
	logger.Printf("Writing output to %s...", cfg.outputFile)
	err := os.WriteFile(cfg.outputFile, []byte(markdownContent), 0644)
	if err != nil {
		logger.Printf("%sError writing to output file %s: %v%s", colorRed, cfg.outputFile, err, colorReset)
		return fmt.Errorf("writing output file %s: %w", cfg.outputFile, err)
	}
	logger.Printf("Markdown content written to %s", cfg.outputFile)

	if cfg.includedPathsFile != "" {
		if err := savePathsToFile(cfg.includedPathsFile, includedPaths, logger); err != nil {
			logger.Printf("%sWarning: Error saving included paths to %s: %v%s", colorRed, cfg.includedPathsFile, err, colorReset)
		}
	}
	if cfg.excludedPathsFile != "" {
		if err := savePathsToFile(cfg.excludedPathsFile, excludedPaths, logger); err != nil {
			logger.Printf("%sWarning: Error saving excluded paths to %s: %v%s", colorRed, cfg.excludedPathsFile, err, colorReset)
		}
	}

	if cfg.addToClipboard {
		logger.Println("Attempting to copy to clipboard...")
		err := clipboard.Init()
		if err != nil {
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
		return fmt.Errorf("writing paths file %s: %w", filename, err)
	}
	return nil
}

// --- Help Message ---

// printHelp displays the command-line help message.
func printHelp() {
	// Use os.Stderr for help message output
	fmt.Fprintf(os.Stderr, "CodeWeaver: Generate Markdown Documentation from Your Codebase.\n")
	fmt.Fprintf(os.Stderr, "Usage: codeweaver [options]\n")
	fmt.Fprintf(os.Stderr, "\nOptions:\n")
	flag.PrintDefaults() // This will now print the flags defined globally
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
