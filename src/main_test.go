package main

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fileaxe/src/conf"
	"fileaxe/src/fileaxe"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strconv"
	"testing"
	"time"

	"github.com/triole/logseal"
)

func TestMainProcessor(t *testing.T) {
	fol := "../testdata/tmp"
	generateTestLogFiles(fol, 9)

	cnf := conf.InitTestConf("rt", fol, "\\.log$")
	lg = logseal.Init("debug")
	la := fileaxe.Init(cnf, lg)
	la.Run()

	files := la.Find(fol, "\\.log$", 0, 0, time.Now())
	verifyFiles(files, "d41d8cd98f00b204e9800998ecf8427e", 9, t)

	cnf = conf.InitTestConf("tn", fol, "\\.log$")
	la = fileaxe.Init(cnf, lg)
	la.Run()

	cnf = conf.InitTestConf("rm", fol, "\\.gz$")
	la = fileaxe.Init(cnf, lg)
	la.Run()
}

func verifyFiles(files fileaxe.FileInfos, hash string, amount int, t *testing.T) {
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

func generateTestLogFiles(fol string, j int) {
	os.MkdirAll(fol, 0755)
	for i := 1; i <= j; i++ {
		createFile(fmt.Sprintf(path.Join(fol, "log%v.log"), i))
	}
	time.Sleep(1000)
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
