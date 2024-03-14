package resources

import (
	"flag"
	"fmt"
	"strings"
)

// Arguments defines the structure to hold parsed arguments
type Arguments struct {
	Language  string
	ModelPath string
}
type ParsedArguments struct {
	Language  int32
	ModelPath string
}

// ParseFlags parses command line arguments and returns an Arguments struct
func ParseFlags() (*ParsedArguments, error) {
	args := &Arguments{}

	flag.StringVar(&args.Language, "l", "", "Language to be processed")
	flag.StringVar(&args.Language, "language", "", "Language to be processed") // Optional: Redundant to demonstrate
	flag.StringVar(&args.ModelPath, "m", "", "Path to the model file (required)")
	flag.StringVar(&args.ModelPath, "modelPath", "", "Path to the model file (required)") // Optional: Redundant

	flag.Usage = func() {
		fmt.Println("Usage: your_program [OPTIONS]")
		fmt.Println("Options:")
		flag.PrintDefaults() // Print default values for all flags
	}

	// Parsing flags
	flag.Parse()

	args.Language = strings.ToLower(args.Language)

	var pickedCode int32
	// Validate against LanguageMap and get associated code
	if code, exists := LanguageMap[args.Language]; exists {
		fmt.Println("Language code:", code) // Use the code as needed
		pickedCode = code
	} else {
		fmt.Println("unsupported language: ", args.Language, " Defaulting to english")
		pickedCode = 0x6E65 // Default to english
	}
	// Check for required flags

	if args.ModelPath == "" {
		return nil, fmt.Errorf("modelPath argument is required")
	}

	return &ParsedArguments{
		Language:  pickedCode,
		ModelPath: args.ModelPath,
	}, nil
}
