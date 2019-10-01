package main

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func initDownload() {
	dir := prepareDir("")
	//defer os.RemoveAll(dir)
	path := (download(dir))
	unzip(dir, path)
}

func prepareDir(basedir string) (dir string) {
	dir, err := ioutil.TempDir(basedir, "ibc-temp")
	if err != nil {
		log.Fatalln("Error: ", err)
	}
	return
}

func download(dir string) (path string) {
	httpClient := http.Client{}
	resp, err := httpClient.Get(RepoZipURL)
	if err != nil {
		log.Fatalln("Error: ", err)
	}
	defer resp.Body.Close()
	file, err := os.Create(filepath.Join(dir, "master.zip"))
	if err != nil {
		log.Fatalln("Error: ", err)
	}
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Fatalln("Error: ", err)
	}
	path = file.Name()
	return
}

func unzip(dir, path string) {
	files, err := zip.OpenReader(path)
	if err != nil {
		log.Fatalln("Error: ", err)
	}
	defer files.Close()
	tmpdir := prepareDir(dir)

	for _, file := range files.File {
		if !file.FileInfo().IsDir() {
			err := os.MkdirAll(filepath.Dir(file.Name), os.ModePerm)
			if err != nil {
				log.Fatalln("Error8: ", err)
			}
		} else if strings.HasSuffix(file.Name, ".ipset") && !file.FileInfo().IsDir() {
			tmppath := filepath.Join(tmpdir, file.Name)
			if !strings.HasPrefix(tmppath, filepath.Clean(tmpdir)+string(os.PathSeparator)) {
				log.Fatalln("Error: Blocking Relative Path, which  is included in Zip !")
			}
			zippedFile, err := file.Open()
			if err != nil {
				log.Fatalln("Error5: ", err)
			}
			defer zippedFile.Close()
			fs, err := os.Create(tmppath)
			if err != nil {
				log.Fatalln("Error6: ", err)
			}
			_, err = io.Copy(fs, zippedFile)
			if err != nil {
				log.Fatalln("Error7: ", err)
			}
		}
	}
}
