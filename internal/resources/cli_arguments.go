package resources

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Arguments defines the structure to hold parsed arguments
type Arguments struct {
	Language  string
	ModelPath string
	Port      int
}
type ParsedArguments struct {
	Language  int32
	ModelPath string
	Port      int
}

type LanguageMap map[string]string

func processLanguageAndCode(args *Arguments) (int32, error) {
	// Read the language map from JSON file
	jsonFile, err := os.Open("languageMap.json")
	if err != nil {
		return 0x6E65, fmt.Errorf("error opening language map: %w", err) // Wrap error for context
	}
	defer jsonFile.Close()

	byteData, err := io.ReadAll(jsonFile)
	if err != nil {
		return 0x6E65, fmt.Errorf("error reading language map: %w", err)
	}

	var languageMap LanguageMap
	err = json.Unmarshal(byteData, &languageMap)
	if err != nil {
		return 0x6E65, fmt.Errorf("error parsing language map: %w", err)
	}

	hexCode, ok := languageMap[strings.ToLower(args.Language)]
	if !ok {
		return 0x6E65, fmt.Errorf("unsupported language: %s", args.Language)
	}

	languageCode, err := strconv.ParseInt(hexCode, 0, 32)
	if err != nil {
		return 0x6E65, fmt.Errorf("error converting hex code: %w", err)
	}

	return int32(languageCode), nil
}

// ParseFlags parses command line arguments and returns an Arguments struct
func ParseFlags() (*ParsedArguments, error) {
	args := &Arguments{}

	flag.StringVar(&args.Language, "l", "", "Language to be processed")
	flag.StringVar(&args.Language, "language", "", "Language to be processed") // Optional: Redundant to demonstrate
	flag.StringVar(&args.ModelPath, "m", "", "Path to the model file (required)")
	flag.StringVar(&args.ModelPath, "modelPath", "", "Path to the model file (required)") // Optional: Redundant
	flag.IntVar(&args.Port, "p", 3031, "Port to start the server on")
	flag.IntVar(&args.Port, "port", 3031, "Port to start the server on") // Optional: Redundant

	flag.Usage = func() {
		fmt.Println("Usage: your_program [OPTIONS]")
		fmt.Println("Options:")
		flag.PrintDefaults() // Print default values for all flags
	}

	// Parsing flags
	flag.Parse()

	args.Language = strings.ToLower(args.Language)

	if args.ModelPath == "" {
		return nil, fmt.Errorf("modelPath argument is required")
	}

	languageCode, err := processLanguageAndCode(args)
	if err != nil {
		fmt.Println("Error setting language, defaulting to English:", err)
		// Use default language code directly as the result here
	}

	return &ParsedArguments{
		Language:  languageCode,
		ModelPath: args.ModelPath,
		Port:      args.Port,
	}, nil
}
