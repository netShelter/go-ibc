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
	files0, err0 := ioutil.ReadDir(os.TempDir())
	evalErr(err0)
	for _, file0 := range files0 {
		if file0.IsDir() && strings.HasPrefix(filepath.Base(filepath.Join(os.TempDir(), file0.Name())), prefix) {
			bsdir = filepath.Join(os.TempDir(), file0.Name())
			files1, err1 := ioutil.ReadDir(bsdir)
			evalErr(err1)
			for _, file1 := range files1 {
				if file1.IsDir() && strings.HasPrefix(filepath.Base(filepath.Join(bsdir, file1.Name())), prefix) {
					files2, err2 := ioutil.ReadDir(filepath.Join(bsdir, file1.Name()))
					evalErr(err2)
					if len(files2) > 0 {
						dir = filepath.Join(bsdir, file1.Name())
						return
					}
				}
			}
		}
	}
	bsdir = prepareDir("")
	path := download(bsdir)
	dir = unzip(bsdir, path)

	return
}

func prepareDir(basedir string) (dir string) {
	dir, err := ioutil.TempDir(basedir, prefix)
	evalErr(err, prefix, basedir, dir)
	return
}

func download(dir string) (path string) {
	httpClient := http.Client{}
	resp, err := httpClient.Get(repoZipURL)
	evalErr(err, repoZipURL)
	defer resp.Body.Close()
	file, err := os.Create(filepath.Join(dir, "master.zip"))
	evalErr(err, file.Name(), dir)
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	evalErr(err, file.Name())
	path = file.Name()
	return
}

func unzip(dir, path string) (tmpdir string) {
	files, err0 := zip.OpenReader(path)
	evalErr(err0, path)
	defer files.Close()
	tmpdir = prepareDir(dir)

	for _, file := range files.File {
		if (strings.HasSuffix(file.Name, ".netset") || strings.HasSuffix(file.Name, ".ipset")) && !file.FileInfo().IsDir() {
			tmppath := filepath.Join(tmpdir, filepath.Base(file.Name))
			if !strings.HasPrefix(tmppath, filepath.Clean(tmpdir)+string(os.PathSeparator)) {
				log.Fatalln("Error: Blocking Relative Path, which  is included in Zip !")
			}
			zippedFile, err1 := file.Open()
			evalErr(err1, file.FileInfo().Name())
			fs, err2 := os.Create(tmppath)
			evalErr(err2, tmppath, fs.Name())
			_, err3 := io.Copy(fs, zippedFile)
			evalErr(err3)
			err4 := zippedFile.Close()
			evalErr(err4)
		}
	}
	defer os.Remove(path)

	return
}
