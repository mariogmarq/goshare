package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
	"io/ioutil"
)

//Creates a random key for encryption of size int, it should be either 16, 24 or 32
//Is created from a secure source
func MakeKey(size int) ([]byte, error) {
	//Reads size bytes from math random and return
	key := make([]byte, size)
	n, err := rand.Read(key)
	if err != nil {
		return nil, err
	}

	if n != size {
		return nil, errors.New("Couldn't create key")
	}

	return key, nil
}

//Encrypts with AES standard the content of data
func Encrypt(key []byte, data io.Reader) ([]byte, error) {
	//Create the cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	//Creates a galois counter mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	//Read the data into a []byte
	rawdata, err := ioutil.ReadAll(data)
	if err != nil {
		return nil, err
	}
	//Create a random Nonce
	nonce := make([]byte, gcm.NonceSize())
	rand.Read(nonce)

	encrypted := gcm.Seal(nonce, nonce, rawdata, nil)

	return encrypted, nil
}

//Decrypts data with AES standard, key should be the same that used in encryption
func Decrypt(key []byte, data io.Reader) ([]byte, error) {
	//read the cifrated data
	ciphertext, err := ioutil.ReadAll(data)
	if err != nil {
		return nil, err
	}

	//Create the corresponding aes block and GCM
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcmDecrypt, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcmDecrypt.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, err
	}

	//Separate the message from the nonce
	nonce, encryptedMessage := ciphertext[:nonceSize], ciphertext[nonceSize:]

	//Decrypt
	decrypted, err := gcmDecrypt.Open(nil, nonce, encryptedMessage, nil)
	if err != nil {
		return nil, err
	}

	return decrypted, nil
}
