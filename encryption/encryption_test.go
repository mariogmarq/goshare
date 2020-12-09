package encryption

import (
	"bytes"
	crand "crypto/rand"
	"math/rand"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	//the test consists in the same test 1000 times
	for i := 0; i < 1000; i++ {
		//Creates a 32bytes random key
		key, err := MakeKey(32)
		if err != nil {
			t.Errorf("Got error creating key %v in test %d", err, i)
		}

		//Creates random byte slice of size between 100 and 10000
		size := rand.Intn(9900)
		size = size + 100
		testData := make([]byte, size)

		//Fills data with random bytes
		writtenSize, err := crand.Read(testData)
		if writtenSize != size || err != nil {
			t.Errorf("Got error creating random slice in test %d", i)
		}

		//Create reader for data
		readerU := bytes.NewReader(testData)

		//encrypt data
		encryptedData, err := Encrypt(key, readerU)
		if err != nil {
			t.Errorf("Error encrypting in test %d", i)
		}

		//reader for encrypted data
		readerE := bytes.NewReader(encryptedData)

		//Decrypt data
		decryptedData, err := Decrypt(key, readerE)
		if err != nil {
			t.Errorf("Error decrypting in test %d", i)
		}

		//Compare data
		if bytes.Equal(testData, decryptedData) == false {
			t.Errorf("Decrypted and original data aren't the same in test %d", i)
		}

	}
}
