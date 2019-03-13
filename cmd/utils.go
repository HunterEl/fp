package cmd

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

var cachesDir string

// TODO: This uses relative pathing for caches dir and it might be smart to decide on
// an absolute location
func getCachesDir() string {
	if cachesDir != "" {
		return cachesDir
	}

	workingDir, _ := os.Getwd()
	cachesDir := fmt.Sprintf("%s/../.caches", workingDir)
	return cachesDir
}

// TODO: Of course use git hashes/checksums to validate that we have the most up to date versions
// NOTE: We will use '.caches' for the time being to store all remote repositories
func localRepoExists(repo string) (bool, error) {
	canonicalizedRepoPath, err := repoURLToPathName(repo)
	if err != nil {
		return false, err
	}

	repoLocation := fmt.Sprintf("%s/%s/", getCachesDir(), canonicalizedRepoPath)
	repoExists := true
	_, err = os.Stat(repoLocation)
	log.Print(err)
	if os.IsNotExist(err) {
		repoExists = false
	}

	return repoExists, nil
}

// NOTE: We will use '.caches' for the time being to store all remote repositories
func initCachesDir() (err error) {
	_, err = os.Stat(getCachesDir())
	if os.IsNotExist(err) {
		err = os.Mkdir(getCachesDir(), 0777)
	}

	return err
}

func removeCachesDir() (err error) {
	caches := getCachesDir()

	_, err = os.Stat(caches)
	if !os.IsNotExist(err) {
		err = os.RemoveAll(caches)
	}

	return err
}

func removeRepo(repoURL string) (err error) {
	caches := getCachesDir()
	canonicalizedRepoPath, err := repoURLToPathName(repoURL)
	if err != nil {
		return err
	}

	repoLocation := fmt.Sprintf("%s/%s", caches, canonicalizedRepoPath)

	_, err = os.Stat(repoLocation)
	if os.IsNotExist(err) {
		return nil
	}

	err = os.RemoveAll(repoLocation)
	return err
}

func repoURLToPathName(repoURL string) (pathName string, err error) {
	parsedURL, err := url.Parse(repoURL)
	if err != nil {
		return "", err
	}

	host := parsedURL.Host
	urlPath := parsedURL.Path

	if host == "" {
		return "", errors.New("Host name in the repo URL is empty")
	}

	// TODO: Is it an error if the path is also empty?
	// TODO: Someone could be fetching from a hosted git repo which was at the root

	filePath := []string{}

	pathParts := strings.Split(urlPath, "/")
	for _, part := range pathParts {
		if part == "" {
			continue
		}

		filePath = append(filePath, part)
	}

	joinedPath := strings.Join(filePath, "/")
	fullUnixPath := fmt.Sprintf("%s/%s", host, joinedPath)

	log.Printf("Full unix path name for %s is %s", repoURL, fullUnixPath)
	return fullUnixPath, err
}

func fetchRepo(repoURL string) (commandOutput string, err error) {
	err = initCachesDir()
	if err != nil {
		return "", err
	}

	err = removeRepo(repoURL)
	if err != nil {
		return "", err
	}

	caches := getCachesDir()
	canonicalizedRepoPath, err := repoURLToPathName(repoURL)
	if err != nil {
		return "", err
	}

	repoLocation := fmt.Sprintf("%s/%s", caches, canonicalizedRepoPath)

	fetchCommand := []string{"clone", repoURL, repoLocation}
	cmdVar := exec.Command("git", fetchCommand...)
	cmdVar.Env = os.Environ()

	if cmdOut, err = cmdVar.Output(); err != nil {
		fmt.Fprint(os.Stderr, err)
		return "", err
	}

	result := string(cmdOut)

	return result, err
}
