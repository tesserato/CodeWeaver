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
	matchers := make([]*regexp.Regexp, len(patterns))
	for i, p := range patterns {
		trimmedPattern := strings.TrimSpace(p)
		if trimmedPattern == "" {
			continue // Skip empty patterns that might result from trailing commas
		}
		logger.Printf("%s%s %s%s\n", color, prefix, trimmedPattern, colorReset)
		rgx, err := regexp.Compile(trimmedPattern)
		if err != nil {
			return nil, fmt.Errorf("invalid regex pattern '%s': %w", trimmedPattern, err)
		}
		matchers[i] = rgx
	}
	return matchers, nil
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
			// Report error and attempt to continue if possible, unless it's critical.
			cb.logger.Printf("%sWarning: Error accessing %s: %v%s\n", colorRed, currentWalkPath, err, colorReset)
			if d != nil && d.IsDir() { // If it's a directory error, might not be skippable
				return filepath.SkipDir // Try to skip this problematic directory
			}
			return nil // Skip this problematic file entry
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
