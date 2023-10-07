package cmd

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/mleku/ec/schnorr"
	secp "github.com/mleku/ec/secp"
	"github.com/mleku/signr/pkg/nostr"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

// listkeysCmd represents the listkeys command
var listkeysCmd = &cobra.Command{
	Use:   "listkeys",
	Short: "List the keys in the keychain",
	Long: `List the keys in the keychain with the name and fingerprint.
`,
	Run: func(cmd *cobra.Command, args []string) {

		keyMap := make(map[string]int)

		filepath.Walk(dataDir,
			func(path string, info fs.FileInfo, err error) (e error) {
				if info.IsDir() {
					return
				}
				filename := filepath.Base(path)
				if strings.HasSuffix(filename, configExt) {
					return
				}
				splitted := strings.Split(filename, ".")
				if len(splitted) == 1 || splitted[1] == pubExt {
					keyMap[splitted[0]]++
				}
				return
			})
		var keySlice []string
		for i := range keyMap {
			if keyMap[i] == 2 {
				keySlice = append(keySlice, i)
			}
		}
		sort.Strings(keySlice)
		var data []byte
		var err error
		encrypted := make(map[string]struct{})
		grid := [][]string{{"name", "fingerprint"}}
		for i := range keySlice {
			pubFilename := keySlice[i] + "." + pubExt
			data, err = os.ReadFile(filepath.Join(dataDir, pubFilename))
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr,
					"error reading file %s: %v\n", pubFilename, err)
				continue
			}
			var secData []byte
			secData, err = os.ReadFile(filepath.Join(dataDir, keySlice[i]))
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr,
					"error reading file %s: %v\n", secData, err)
				continue
			}
			if string(secData[0]) == "e" {
				encrypted[keySlice[i]] = struct{}{}
			}
			key := strings.TrimSpace(string(data))

			var pk *secp.PublicKey
			pk, err = nostr.DecodePublicKey(key)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr,
					"error decoding key '%s' %s: %v\n",
					keySlice[i], pk, err)
				continue
			}
			spk := schnorr.SerializePubKey(pk)
			fingerprint := sha256.Sum256(spk)
			grid = append(grid, []string{
				keySlice[i],
				"@" + hex.EncodeToString(fingerprint[:8]),
				// hex.EncodeToString(spk),
			})
		}
		var maxLen1, maxLen2 int
		for i := range grid {
			l := len(grid[i][0])
			if l > maxLen1 {
				maxLen1 = l
			}
			l = len(grid[i][1])
			if l > maxLen2 {
				maxLen2 = l
			}
		}
		header, tail := grid[0], grid[1:]
		grid = append([][]string{header},
			[]string{
				strings.Repeat("-", maxLen1) + " ",
				strings.Repeat("-", maxLen2),
			})
		grid = append(grid, tail...)
		maxLen1++
		fmt.Print("keys in keychain: (* = password protected)\n\n")
		for i := range grid {
			isDefault := "          "
			if grid[i][0] == defaultKey {
				isDefault = " (default)"
			}
			crypted := " "
			if _, ok := encrypted[grid[i][0]]; ok {
				crypted = "*"
			}
			grid[i][0] = grid[i][0] + strings.Repeat(" ",
				maxLen1-len(grid[i][0]))
			fmt.Printf("  %s %s%s\n", crypted, grid[i][0], grid[i][1]+isDefault)
		}

	},
}

func init() {
	rootCmd.AddCommand(listkeysCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listkeysCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listkeysCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
