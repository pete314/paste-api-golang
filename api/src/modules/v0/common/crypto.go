//Author: Peter Nagy <https://peternagy.ie>
//Since: 06, 2017
//Description: collection of crypto functions
package common

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"hash/crc32"
	"time"
)

var (
	crc32q = crc32.MakeTable(0xD5828281)
)

func NewUUID() string {
	return uuid.NewV4().String()
}

func GenRandomHash() string {
	bytes := GetRandomBytes(32)
	if len(bytes) == 0 {
		//Fail over
		bytes = []byte(string(time.Now().Unix()))
	}
	hasher := sha256.New()
	hasher.Write(bytes)

	return hex.EncodeToString(hasher.Sum(nil))
}

//CreateBcryptPassword - create bcrypt password
func CreateBcryptPassword(uPassword string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(uPassword), 7)
	if err == nil {
		return hash, nil
	}

	CheckError("Common/crypto - Error with bcrypt", err, false)
	return nil, err

}

//ValidateBcryptPassword - compare hash and password
func ValidateBcryptPassword(dPassword, uPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(dPassword), []byte(uPassword))
	if err == nil {
		return true
	}

	CheckError("Common/crypto - Error with bcrypt", err, false)
	return false

}

//GetRandomBytes - Get random bytes from kernel
func GetRandomBytes(size int) []byte {
	data := make([]byte, size)
	_, err := rand.Read(data)
	CheckError("Common/crypto - Error generating random bytes:", err, false)

	return data
}

//GetStringHash - sha1
func GetStringHash(d string) string {
	csum := crc32.Checksum([]byte(d), crc32q)

	hasher := sha256.New()
	hasher.Write([]byte(fmt.Sprintf("%d:%s", csum, d)))

	return hex.EncodeToString(hasher.Sum(nil))
}
