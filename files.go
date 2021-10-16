package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func getBlocklistFilesFromSource() (dir string) {
	if copyexist, dir := existLocalCopyOfSource(); copyexist {
		return dir
	}

	return getLocalCopyOfSource()
}

func getLocalCopyOfSource() (dir string) {
	tmpdir := createTempDir("")
	path := download(tmpdir)
	dir = unzip(tmpdir, path)

	return
}

func existLocalCopyOfSource() (copyexist bool, dir string) {
	files0, err0 := ioutil.ReadDir(os.TempDir())
	evalErr(err0)

	for _, file0 := range files0 {
		if file0.IsDir() && strings.HasPrefix(filepath.Base(filepath.Join(os.TempDir(), file0.Name())), prefix) {
			bsdir := filepath.Join(os.TempDir(), file0.Name())
			files1, err1 := ioutil.ReadDir(bsdir)
			evalErr(err1)

			for _, file1 := range files1 {
				if file1.IsDir() && strings.HasPrefix(filepath.Base(filepath.Join(bsdir, file1.Name())), prefix) {
					files2, err2 := ioutil.ReadDir(filepath.Join(bsdir, file1.Name()))
					evalErr(err2)

					if len(files2) > 0 {
						dir = filepath.Join(bsdir, file1.Name())
						return true, dir
					}
				}
			}
		}
	}

	return false, ""
}

func createTempDir(basedir string) (dir string) {
	dir, err := ioutil.TempDir(basedir, prefix)
	evalErr(err, prefix, basedir, dir)

	return
}

func createFile(path, filename string) *os.File {
	file, err := os.Create(filepath.Join(path, filename))
	evalErr(err, file.Name(), path)

	return file
}
