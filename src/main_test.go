package main

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/triole/logseal"
)

func TestMainProcessor(t *testing.T) {
	lg = logseal.Init("trace")
	conf := initConf(
		"../testdata", "\\.log$", "1s", false, false,
	)
	time.Sleep(1500)

	generateTestLogFiles(0, 9)
	runProcessor(conf)

}

func generateTestLogFiles(i, j int) {
	for i := 1; i <= j; i++ {
		createFile(fmt.Sprintf("../testdata/log%v.log", i))
	}
}

func createFile(target string) {
	f, err := os.Create(target)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for _, line := range makeLines() {
		_, err := f.WriteString(line + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func makeLines() (lines []string) {
	for i := 1; i <= 999; i++ {
		lines = append(lines, makeLine(i))
	}
	return
}

func makeLine(i int) string {
	return fmt.Sprintf(
		"%d --- %s +++ %s +++ %s",
		i,
		sha512(strconv.Itoa(i)),
		sha512(strconv.Itoa(i+1000)),
		sha512(strconv.Itoa(i+2000)),
	)
}

func sha512(str string) string {
	hasher := sha1.New()
	hasher.Write([]byte(str))
	return hex.EncodeToString(hasher.Sum(nil))
}

func md5SumOfFile(filePath string) (string, error) {
	var returnMD5String string
	file, err := os.Open(filePath)
	if err != nil {
		return returnMD5String, err
	}
	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return returnMD5String, err
	}
	hashInBytes := hash.Sum(nil)[:16]
	returnMD5String = hex.EncodeToString(hashInBytes)
	return returnMD5String, nil
}
