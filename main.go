package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	defaultDropitPath = "./.dropit"
	dropitPathEnv     = "DROPITPATH"
)

func dropitPath() string {
	path := os.Getenv(dropitPathEnv)
	if strings.TrimSpace(path) == "" {
		return defaultDropitPath
	}
	return path
}

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println("missing file path")
		return
	}

	dropFile(args[0])

}

func dropFile(path string) error {
	creationDate := time.Now().UTC()
	destFileName := creationDate.Format("20060102T150405") + "-" + path //TODO: handle absolute path
	destFilePath := filepath.Join(dropitPath(), destFileName)

	err := copyFile(path, destFilePath)
	if err != nil {
		return err
	}

	return nil
}

//TODO: review the whole function
func copyFile(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	_, err = os.Stat(dst)
	if err == nil {
		return fmt.Errorf("file %s already exists", dst)
	}

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	buf := make([]byte, 4096)
	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if _, err := destination.Write(buf[:n]); err != nil {
			return err
		}
	}

	return nil
}
