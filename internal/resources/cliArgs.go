package resources

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

//go:embed languageMap.json
var languageMapData []byte // Embedded language map file as a byte slice

// Arguments holds the parsed CLI arguments
type Arguments struct {
	Language  string
	ModelPath string
	Port      int
}

// ParsedArguments holds the processed arguments
type ParsedArguments struct {
	Language  int32
	ModelPath string
	Port      int
}

// LanguageMap represents the mapping of languages to their hex codes
type LanguageMap map[string]string

func processLanguageAndCode(language string) (int32, error) {
    var languageMap LanguageMap
    err := json.Unmarshal(languageMapData, &languageMap)
    if err != nil {
        return 0x6E65, fmt.Errorf("error parsing language map: %w", err)
    }

    hexCode, ok := languageMap[strings.ToLower(language)]
    if !ok {
        return 0x6E65, fmt.Errorf("unsupported language")
    }

    fmt.Printf("Hex Code Found: %s\n", hexCode)

    languageCode, err := strconv.ParseInt(hexCode, 0, 32)
    if err != nil {
        return 0x6E65, fmt.Errorf("error converting hex code: %w", err)
    }

    return int32(languageCode), nil
}

func ApplyExitOnHelp(c *cobra.Command, exitCode int) {
	helpFunc := c.HelpFunc()
	c.SetHelpFunc(func(c *cobra.Command, s []string) {
		helpFunc(c, s)
		os.Exit(exitCode)
	})
}

func ParseFlags() (*ParsedArguments, error) {
    args := &Arguments{}

    var parsedArgs *ParsedArguments

    rootCmd := &cobra.Command{
        Use:   "whisper",
        Short: "Audio transcription using the OpenAI Whisper models",
        RunE: func(cmd *cobra.Command, _ []string) error {
            // Process language code with fallback
            languageCode, err := processLanguageAndCode(args.Language)
            if err != nil {
                fmt.Println("Error setting language, defaulting to English")
                // Default to English
                languageCode = 0x6E65
            }

            parsedArgs = &ParsedArguments{
                Language:  languageCode,
                ModelPath: args.ModelPath,
                Port:      args.Port,
            }
            return nil
        },
    }

    rootCmd.Flags().StringVarP(&args.Language, "language", "l", "", "Language to be processed")
    rootCmd.Flags().StringVarP(&args.ModelPath, "modelPath", "m", "ggml-medium.bin", "Path to the model file (required)")
    rootCmd.Flags().IntVarP(&args.Port, "port", "p", 3000, "Port to start the server on")

	ApplyExitOnHelp(rootCmd, 0)

    err := rootCmd.Execute()
    if err != nil {
        return nil, err
    }

    return parsedArgs, nil
}

