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
	fol := "../testdata"
	lg = logseal.Init("trace")
	conf := initConf(fol, "\\.log$", "1s", false, false)
	time.Sleep(1500)

	generateTestLogFiles(0, 9)
	runProcessor(conf)

	checkFiles(fol, "\\.log$", "d41d8cd98f00b204e9800998ecf8427e", 9, t)
	checkFiles(fol, "\\.gz$", "0e93baf81315ce74e7484d374d550179", 9, t)
}

func checkFiles(fol, matcher, hash string, amount int, t *testing.T) {
	files := find(fol, matcher, 0, time.Now())
	if len(files) != amount {
		t.Errorf("test error, amount of files wrong: %d != %d", len(files), amount)
	}
	for _, fil := range files {
		md5sum, _ := md5SumOfFile(fil.Path)
		if md5sum != hash {
			t.Errorf("test error, log file hash wrong: %s != %s", md5sum, hash)
		}
	}
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
		sha512OfString(strconv.Itoa(i)),
		sha512OfString(strconv.Itoa(i+1000)),
		sha512OfString(strconv.Itoa(i+2000)),
	)
}

func sha512OfString(str string) string {
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
