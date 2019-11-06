package main

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func unzip(dir, path string) (tmpdir string) {
	files, err0 := zip.OpenReader(path)
	evalErr(err0, path)

	defer files.Close()

	tmpdir = createTempDir(dir)

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

	return tmpdir
}
