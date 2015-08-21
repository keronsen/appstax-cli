package hosting

import (
	"appstax-cli/appstax/apiclient"
	"appstax-cli/appstax/fail"
	"appstax-cli/appstax/log"
	"archive/tar"
	"bufio"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func UploadStatic(archivePath string, progressWriter io.Writer) error {
	_, _, err := apiclient.PostFile(apiclient.Url("/appstax/hosting/static"), archivePath, progressWriter)
	return err
}

func UploadServer(archivePath string, progressWriter io.Writer) error {
	_, _, err := apiclient.PostFile(apiclient.Url("/appstax/hosting/server/code"), archivePath, progressWriter)
	return err
}

func CreateServer() error {
	_, _, err := apiclient.Post(apiclient.Url("/appstax/hosting/server"), "")
	return err
}

func DeleteServer() error {
	_, _, err := apiclient.Delete(apiclient.Url("/appstax/hosting/server"))
	return err
}

func PrepareArchive(rootPath string) (string, int64, error) {
	file, err := ioutil.TempFile("", "")
	fail.Handle(err)
	defer file.Close()
	fileWriter := bufio.NewWriter(file)
	defer fileWriter.Flush()
	gzipWriter, err := gzip.NewWriterLevel(fileWriter, gzip.BestCompression)
	fail.Handle(err)
	defer gzipWriter.Close()
	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	fullRootPath, err := filepath.Abs(rootPath)
	fail.Handle(err)
	err = addAllToArchive(fullRootPath, tarWriter)
	if err != nil {
		return "", 0, err
	}

	tarWriter.Close()
	gzipWriter.Close()
	fileWriter.Flush()
	file.Close()

	fileInfo, err := os.Stat(file.Name())
	fail.Handle(err)
	return file.Name(), fileInfo.Size(), nil
}

func addAllToArchive(fullRootPath string, tarWriter *tar.Writer) error {
	log.Debugf("Creating archive by walking from root path %s", fullRootPath)
	return filepath.Walk(fullRootPath, func(path string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !fileInfo.IsDir() && fileInfo.Name()[:1] != "." {
			err := addFileToArchive(path, path[len(fullRootPath+"/"):], tarWriter, fileInfo)
			if err != nil {
				return err
			}
		} else {
			log.Debugf("Ignoring path %s", path)
		}
		return nil
	})
}

func addFileToArchive(filePath string, addPath string, tarWriter *tar.Writer, fileInfo os.FileInfo) error {
	addPath = filepath.ToSlash(addPath)
	log.Debugf("Adding file %s from %s", addPath, filePath)
	fileReader, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer fileReader.Close()

	header := new(tar.Header)
	header.Name = addPath
	header.Size = fileInfo.Size()
	header.Mode = int64(fileInfo.Mode())
	header.ModTime = fileInfo.ModTime()

	err = tarWriter.WriteHeader(header)
	fail.Handle(err)
	_, err = io.Copy(tarWriter, fileReader)
	return err
}
