package cmd

import (
	"encoding/hex"
	"github.com/minio/sha256-simd"
	"github.com/mleku/ec/schnorr"
	secp "github.com/mleku/ec/secp"
	"github.com/mleku/signr/pkg/nostr"
	"github.com/pkg/errors"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
)

func GetList(g [][]string) (grid [][]string, encrypted map[string]struct{},
	err error) {

	grid = g
	keyMap := make(map[string]int)

	err = filepath.Walk(dataDir,
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
	if err != nil {
		err = errors.Wrap(err,
			"failed wile walking data directory")
		return
	}
	var keySlice []string
	for i := range keyMap {
		if keyMap[i] == 2 {
			keySlice = append(keySlice, i)
		}
	}
	sort.Strings(keySlice)
	var data []byte
	encrypted = make(map[string]struct{})
	for i := range keySlice {
		pubFilename := keySlice[i] + "." + pubExt
		data, err = ReadFile(pubFilename)
		if err != nil {
			printErr(
				"error reading file %s: %v\n", pubFilename, err)
			continue
		}
		var secData []byte
		secData, err = ReadFile(keySlice[i])
		if err != nil {
			printErr(
				"error reading file '%s': %v\n", keySlice[i], err)
			continue
		}
		if string(secData[0]) == "e" {
			encrypted[keySlice[i]] = struct{}{}
		}
		key := strings.TrimSpace(string(data))

		var pk *secp.PublicKey
		pk, err = nostr.DecodePublicKey(key)
		if err != nil {
			printErr(
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
	return
}
