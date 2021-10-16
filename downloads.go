package main

import (
	"io"
	"net/http"
)

func download(tmpdir string) (zippath string) {
	httpClient := http.Client{}
	resp, err := httpClient.Get(repoZipURL)
	evalErr(err, repoZipURL)

	defer resp.Body.Close()

	file := createFile(tmpdir, "master.zip")

	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	evalErr(err, file.Name())
	zippath = file.Name()

	return
}
