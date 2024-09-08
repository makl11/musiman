package cmd

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"

	"github.com/makl11/musiman/context_keys"
	"github.com/makl11/musiman/data"
	"github.com/makl11/musiman/scanner"
)

var (
	minSizeStr  string
	ignorePaths []string
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:     "scan [directory]",
	Short:   "Recursively scan a directory for music files (defaults to current directory if not specified)",
	Args:    cobra.MaximumNArgs(1),
	PreRunE: data.InitDb,
	Run: func(cmd *cobra.Command, args []string) {
		db := cmd.Context().Value(context_keys.DB).(*sqlx.DB) // Never nil, InitDb returns error if it fails
		defer db.Close()

		dir := "."
		if len(args) > 0 {
			dir = args[0]
		}

		minSize, err := parseSize(minSizeStr)
		if err != nil {
			fmt.Println("Error parsing min-size:", err)
			os.Exit(1)
		}

		err = scanner.ScanDirForMusic(dir, minSize, ignorePaths)
		if err != nil {
			fmt.Println("Error scanning directory:", err)
			os.Exit(1)
		}
	},
}

func init() {
	scanCmd.Flags().StringVarP(&minSizeStr, "min-size", "m", "0B", "Minimum file size to include in the scan (e.g. 1KB, 1MB, 1MiB)")
	scanCmd.Flags().StringArrayVarP(&ignorePaths, "ignore", "i", []string{}, "Relative path[s] to ignore during the scan. Can be specified multiple times (e. g. -i a/path --ignore another/path)")
	rootCmd.AddCommand(scanCmd)
}

var (
	ErrMissingArgumentValue = errors.New("missing argument value")
	ErrMissingNumericValue  = errors.New("missing numeric value in size")
	ErrInvalidNumericValue  = errors.New("invalid numeric value in size")
	ErrUnknownSizeUnit      = errors.New("unknown size unit")
)

var unitMultipliers = map[string]uint64{
	"B":   1,
	"KB":  1000,
	"MB":  1000 * 1000,
	"GB":  1000 * 1000 * 1000,
	"KIB": 1024,
	"MIB": 1024 * 1024,
	"GIB": 1024 * 1024 * 1024,
}

func parseSize(sizeStr string) (uint64, error) {
	sizeStr = strings.TrimSpace(sizeStr)
	if sizeStr == "" {
		return 0, fmt.Errorf("%w: %s", ErrMissingArgumentValue, "min-size")
	}

	firstLetterIdx := strings.IndexFunc(sizeStr, func(r rune) bool {
		return !(r == ',' || r == '.' || r == ' ' || r == '_') && (r < '0' || r > '9')
	})

	if firstLetterIdx == 0 {
		return 0, fmt.Errorf("%w: %s", ErrMissingNumericValue, sizeStr)
	}

	if firstLetterIdx > 0 {
		unit := strings.ToUpper(strings.TrimSpace(sizeStr[firstLetterIdx:]))

		multiplier, ok := unitMultipliers[unit]

		if !ok {
			return 0, fmt.Errorf("%w: %s", ErrUnknownSizeUnit, unit)
		}

		valueStr := strings.TrimSpace(sizeStr[:firstLetterIdx])
		value, err := strconv.ParseUint(valueStr, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("%w: %s", ErrInvalidNumericValue, sizeStr)
		}
		return value * multiplier, nil
	}

	value, err := strconv.ParseUint(sizeStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%w: %s", ErrInvalidNumericValue, sizeStr)
	}
	return value, nil
}
