// Copyright (C) 2017 Minio Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package crypto

import (
	"compress/gzip"
	"encoding/csv"

	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/minio/sio"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/hkdf"

	"opensignauxfaibles/tools/altares/pkg/minoos"
	"opensignauxfaibles/tools/altares/pkg/utils"
	"opensignauxfaibles/tools/altares/test"
)

func Test_ExampleEncrypt(t *testing.T) {
	mc := minoos.New(test.NewS3ForTest(t), test.FakeBucketName())
	stock := test.GenerateStockCSV(50)

	remoteFileName := "altares.csv.gz.x"

	key := buildKey()

	reader, writer := io.Pipe()
	gzw := gzip.NewWriter(writer)
	var wg sync.WaitGroup

	go func() {
		defer cclose(writer, "fermeture du writer encrypté ")
		slog.Info("encrypte", slog.String("status", "start"))
		wg.Add(1)
		_, err := sio.Encrypt(gzw, stock, sio.Config{Key: key[:]})
		utils.ManageError(err, "erreur à l'encryption")
		err = gzw.Flush()
		utils.ManageError(err, "erreur au vidage")
		wg.Done()
		slog.Info("encrypte", slog.String("status", "end"))
	}()

	slog.Info("depose le fichier sur OOS", slog.String("status", "start"))
	size, lastModified := mc.PutAltaresFile(remoteFileName, reader)
	slog.Info(
		"depose le fichier sur OOS",
		slog.String("status", "end"),
		slog.Any("wrote", size),
		slog.Time("lastModified", lastModified),
	)
	cclose(reader, "fermeture du pipe reader")
	wg.Wait()

	files := mc.ListAltaresFiles()
	slog.Info("liste les fichiers sur le bucket", slog.Any("files", files))
	assert.Contains(t, files, "altares/"+remoteFileName)

	// fetch then decompress then uncrypt

	localTarget, err := os.Create(filepath.Join(os.TempDir(), remoteFileName))
	require.NoError(t, err)
	slog.Info("copie locale créée", slog.String("path", localTarget.Name()))
	defer cclose(localTarget, "fermeture de la copie locale")

	slog.Info("récupère le fichier sur OOS", slog.String("status", "start"))
	remote := mc.GetAltaresFile(remoteFileName)
	slog.Info("récupère le fichier sur OOS", slog.String("status", "end"))
	gzr, err := gzip.NewReader(remote)
	require.NoError(t, err)

	var decrypted int64
	decryptReader, err := sio.DecryptReader(gzr, sio.Config{Key: key[:]})
	require.NoError(t, err)
	csvReader := csv.NewReader(decryptReader)
	record, err := csvReader.Read()
	require.NoError(t, err)
	slog.Info("une ligne", slog.Any("record", record))
	if decrypted, err = sio.Decrypt(localTarget, remote, sio.Config{Key: key[:]}); err != nil {
		if _, ok := err.(sio.Error); ok {
			utils.ManageError(err, "fichier avec erreur d'encryption") // add error handling - here we know that the data is malformed/not authentic.
		}
		utils.ManageError(err, "erreur de déchiffrage des données") // add error handling
	}
	slog.Debug("fichier décrypté", slog.Any("decrypted", decrypted))
	cclose(remote, "fermeture du fichier remote")
	require.NoError(t, err)
}

func buildKey() [32]byte {
	// the master key used to derive encryption keys
	// this key must be keep secret
	// generate a random nonce to derive an encryption key from the master key
	// this nonce must be saved to be able to decrypt the data again - it is not
	// required to keep it secret
	// derive an encryption key from the master key and the nonce
	masterkey, err := hex.DecodeString("000102030405060708090A0B0C0D0E0FF0E0D0C0B0A090807060504030201000") // use your own key here
	utils.ManageError(err, "erreur au décadocage de la clé hexadécimale")

	var nonce [32]byte
	_, err = io.ReadFull(rand.Reader, nonce[:])
	utils.ManageError(err, "erreur à la lecture de données random")

	var key [32]byte
	kdf := hkdf.New(sha256.New, masterkey, nonce[:], nil)
	_, err = io.ReadFull(kdf, key[:])
	utils.ManageError(err, "erreur à la dérication de la clé d'encryption")
	return key
}

func cclose(o io.Closer, s string) {
	err := o.Close()
	slog.Debug(s)
	utils.ManageError(err, "erreur "+s)
}

func ExampleDecrypt() {
	// the master key used to derive encryption keys
	masterkey, err := hex.DecodeString("000102030405060708090A0B0C0D0E0FF0E0D0C0B0A090807060504030201000") // use your own key here
	if err != nil {
		fmt.Printf("Cannot decode hex key: %v", err) // add error handling
		return
	}

	// the nonce used to derive the encryption key
	nonce, err := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000001") // use your generated nonce here
	if err != nil {
		fmt.Printf("Cannot decode hex key: %v", err) // add error handling
		return
	}

	// derive the encryption key from the master key and the nonce
	var key [32]byte
	kdf := hkdf.New(sha256.New, masterkey, nonce, nil)
	if _, err = io.ReadFull(kdf, key[:]); err != nil {
		fmt.Printf("Failed to derive encryption key: %v", err) // add error handling
		return
	}

	input := os.Stdin   // customize for your needs - the encrypted data
	output := os.Stdout // customize from your needs - the decrypted output

	if _, err = sio.Decrypt(output, input, sio.Config{Key: key[:]}); err != nil {
		if _, ok := err.(sio.Error); ok {
			fmt.Printf("Malformed encrypted data: %v", err) // add error handling - here we know that the data is malformed/not authentic.
			return
		}
		fmt.Printf("Failed to decrypt data: %v", err) // add error handling
		return
	}
}

func ExampleEncryptReader() {
	// the master key used to derive encryption keys
	// this key must be keep secret
	masterkey, err := hex.DecodeString("000102030405060708090A0B0C0D0E0FF0E0D0C0B0A090807060504030201000") // use your own key here
	if err != nil {
		fmt.Printf("Cannot decode hex key: %v", err) // add error handling
		return
	}

	// generate a random nonce to derive an encryption key from the master key
	// this nonce must be saved to be able to decrypt the data again - it is not
	// required to keep it secret
	var nonce [32]byte
	if _, err = io.ReadFull(rand.Reader, nonce[:]); err != nil {
		fmt.Printf("Failed to read random data: %v", err) // add error handling
		return
	}

	// derive an encryption key from the master key and the nonce
	var key [32]byte
	kdf := hkdf.New(sha256.New, masterkey, nonce[:], nil)
	if _, err = io.ReadFull(kdf, key[:]); err != nil {
		fmt.Printf("Failed to derive encryption key: %v", err) // add error handling
		return
	}

	input := os.Stdin // customize for your needs - the plaintext input
	encrypted, err := sio.EncryptReader(input, sio.Config{Key: key[:]})
	if err != nil {
		fmt.Printf("Failed to encrypted reader: %v", err) // add error handling
		return
	}

	// the encrypted io.Reader can be used like every other reader - e.g. for copying
	if _, err := io.Copy(os.Stdout, encrypted); err != nil {
		fmt.Printf("Failed to copy data: %v", err) // add error handling
		return
	}
}

func ExampleEncryptWriter() {
	// the master key used to derive encryption keys
	// this key must be keep secret
	masterkey, err := hex.DecodeString("000102030405060708090A0B0C0D0E0FF0E0D0C0B0A090807060504030201000") // use your own key here
	if err != nil {
		fmt.Printf("Cannot decode hex key: %v", err) // add error handling
		return
	}

	// generate a random nonce to derive an encryption key from the master key
	// this nonce must be saved to be able to decrypt the data again - it is not
	// required to keep it secret
	var nonce [32]byte
	if _, err = io.ReadFull(rand.Reader, nonce[:]); err != nil {
		fmt.Printf("Failed to read random data: %v", err) // add error handling
		return
	}

	// derive an encryption key from the master key and the nonce
	var key [32]byte
	kdf := hkdf.New(sha256.New, masterkey, nonce[:], nil)
	if _, err = io.ReadFull(kdf, key[:]); err != nil {
		fmt.Printf("Failed to derive encryption key: %v", err) // add error handling
		return
	}

	output := os.Stdout // customize for your needs - the encrypted output
	encrypted, err := sio.EncryptWriter(output, sio.Config{Key: key[:]})
	if err != nil {
		fmt.Printf("Failed to encrypted writer: %v", err) // add error handling
		return
	}

	// the encrypted io.Writer can be used now but it MUST be closed at the end to
	// finalize the encryption.
	if _, err = io.Copy(encrypted, os.Stdin); err != nil {
		fmt.Printf("Failed to copy data: %v", err) // add error handling
		return
	}
	if err = encrypted.Close(); err != nil {
		fmt.Printf("Failed to finalize encryption: %v", err) // add error handling
		return
	}
}
