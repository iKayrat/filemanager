package random

import (
	"encoding/base64"
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(someString string, charset string) string {
	b := make([]byte, len(someString))
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return base64.StdEncoding.EncodeToString(b)
}

func String(someStr string) string {
	return StringWithCharset(someStr, charset)
}

// generate random name
// func genName(fname string) string {

// 	ext := filepath.Ext(fname)

// 	b := make([]byte, len(fname))
// 	_, err := rand.Read(b)
// 	if err != nil {
// 		panic(err)
// 	}

// 	var str strings.Builder

// 	str.WriteString(string(b))
// 	str.WriteString(".")
// 	str.WriteString(ext)

// 	log.Println("byteTostring:", string(b))
// 	log.Println("generated string:", str)

// 	return str.String()
// }
