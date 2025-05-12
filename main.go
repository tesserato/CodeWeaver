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