package signr

import (
	"encoding/hex"
	"github.com/minio/sha256-simd"
	"github.com/mleku/ec/schnorr"
	"github.com/mleku/ec/secp"
	"github.com/mleku/signr/pkg/nostr"
	"github.com/pkg/errors"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
)

func (s *Signr) GetKeyPairNames() (list []string, err error) {

	keyMap := make(map[string]int)

	err = filepath.Walk(s.DataDir,
		func(path string, info fs.FileInfo, err error) (e error) {

			if info.IsDir() {
				return
			}

			filename := filepath.Base(path)
			if strings.HasSuffix(filename, ConfigExt) {

				return
			}

			if strings.HasSuffix(filename, DeletedExt) {

				return
			}

			splitted := strings.Split(filename, ".")
			if len(splitted) == 1 || splitted[1] == PubExt {

				keyMap[splitted[0]]++
			}
			return
		},
	)
	if err != nil {
		err = errors.Wrap(err, "failed while walking data directory")
		return
	}

	for i := range keyMap {

		if keyMap[i] == 2 {

			list = append(list, i)
		}
	}

	sort.Strings(list)

	return
}

func (s *Signr) GetList(g [][]string) (grid [][]string,
	encrypted map[string]struct{},
	err error) {

	grid = g

	var keySlice []string
	keySlice, err = s.GetKeyPairNames()
	if err != nil {

		s.Fatal("error reading in keychain data '%s'\n", err)
	}

	var data []byte
	encrypted = make(map[string]struct{})
	for i := range keySlice {

		pubFilename := keySlice[i] + "." + PubExt

		data, err = s.ReadFile(pubFilename)
		if err != nil {

			s.PrintErr("error reading file %s: %v\n", pubFilename, err)
			continue
		}

		var secData []byte
		if secData, err = s.ReadFile(keySlice[i]); err != nil {

			s. PrintErr("error reading file '%s': %v\n", keySlice[i], err)
			continue
		}

		for j, sb := range secData {

			if sb == ' ' {

				if len(secData) >= 64 && secData[j+1] == '*' {

					secData = secData[:64]
					encrypted[keySlice[i]] = struct{}{}

					break
				}
			}
		}

		var pk *secp256k1.PublicKey
		key := strings.TrimSpace(string(data))
		if pk, err = nostr.NpubToPublicKey(key); err != nil {

			s.PrintErr("error decoding key '%s' %s: %v\n",
				keySlice[i], pk, err)
			continue
		}

		spk := schnorr.SerializePubKey(pk)
		fingerprint := sha256.Sum256(spk)

		grid = append(grid,
			[]string{

				keySlice[i],
				"@" + hex.EncodeToString(fingerprint[:8]),
			},
		)
	}
	return
}
