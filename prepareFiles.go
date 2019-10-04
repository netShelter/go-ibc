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

func initDownload() (bsdir, dir string) {
	dir = prepareDir("")
	path := (download(dir))
	bsdir = unzip(dir, path)
	return
}

func prepareDir(basedir string) (dir string) {
	dir, err := ioutil.TempDir(basedir, "ibc-temp")
	evalErr(err)
	return
}

func download(dir string) (path string) {
	httpClient := http.Client{}
	resp, err := httpClient.Get(RepoZipURL)
	evalErr(err)
	defer resp.Body.Close()
	file, err := os.Create(filepath.Join(dir, "master.zip"))
	evalErr(err)
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	evalErr(err)
	path = file.Name()
	return
}

func unzip(dir, path string) (tmpdir string) {
	files, err := zip.OpenReader(path)
	evalErr(err)
	defer files.Close()
	tmpdir = prepareDir(dir)

	for _, file := range files.File {
		if (strings.HasSuffix(file.Name, ".netset") || strings.HasSuffix(file.Name, ".ipset")) && !file.FileInfo().IsDir() {
			tmppath := filepath.Join(tmpdir, filepath.Base(file.Name))
			if !strings.HasPrefix(tmppath, filepath.Clean(tmpdir)+string(os.PathSeparator)) {
				log.Fatalln("Error: Blocking Relative Path, which  is included in Zip !")
			}
			zippedFile, err := file.Open()
			evalErr(err)
			defer zippedFile.Close()
			fs, err := os.Create(tmppath)
			evalErr(err)
			_, err = io.Copy(fs, zippedFile)
			evalErr(err)
		}
	}
	return
}
