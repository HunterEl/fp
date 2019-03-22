package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

var cachesDir string

// Config struct mapping the fp config json
type Config struct {
	Commands commands `json:"commands"`
	// TODO: ADD OTHER CONFIG PROPERTIES HERE AS WELL
}

type commands map[string]Command

// Command struct represents info relating to each command
type Command struct {
	Command        string   `json:"command"`
	Environment    string   `json:"environment"`
	Lang           string   `json:"lang"`
	RunCommands    []string `json:"runCommands"`
	InstallCommand string   `json:"install"`
}

type AppConfig struct {
	commandsRepo string
	// TODO: Add more commands here as we progress
}

// Iterate the app config lines and try to find the 'commandsRepo' argument
// returns the commandsRepo value (if it exists) or and error if something went _horribly_ wrong along the way
func readAppConfig() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	configLocation := fmt.Sprintf("%s/fp.rc", cwd)
	b, err := ioutil.ReadFile(configLocation)
	if err != nil {
		return "", err
	}

	str := string(b)

	lines := strings.Split(str, "\n")
	commandsRepo := ""
	for _, line := range lines {
		if line == "" {
			continue
		}

		lineParts := strings.Split(line, "=")
		if len(lineParts) < 2 {
			continue
		}

		if lineParts[0] == "commandsRepo" {
			commandsRepo = lineParts[1]
			break
		}
	}

	if commandsRepo == "" {
		return "", errors.New("Could not find 'commandsRepo' in the app config. It is required, you monkey")
	}

	return commandsRepo, nil
}

// TODO: This uses relative pathing for caches dir and it might be smart to decide on
// TODO: an absolute location
func getCachesDir() string {
	if cachesDir != "" {
		return cachesDir
	}

	workingDir, _ := os.Getwd()
	cachesDir := fmt.Sprintf("%s/.caches", workingDir)
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
	return fileExists(repoLocation)
}

func localRepoHasConfig(repo string) (bool, error) {
	canonicalizedRepoPath, err := repoURLToPathName(repo)
	if err != nil {
		return false, err
	}

	configLocation := fmt.Sprintf("%s/%s/config.json", getCachesDir(), canonicalizedRepoPath)
	return fileExists(configLocation)
}

func fullRepoPath(repo string) (string, error) {
	canonicalizedRepoPath, err := repoURLToPathName(repo)
	if err != nil {
		return "", err
	}

	cachesDir := getCachesDir()

	configLocation := fmt.Sprintf("%s/%s", cachesDir, canonicalizedRepoPath)
	return configLocation, nil
}

func loadRepoConfig(repo string) (Config, error) {
	canonicalizedRepoPath, err := repoURLToPathName(repo)
	if err != nil {
		return Config{}, err
	}

	configLocation := fmt.Sprintf("%s/%s/config.json", getCachesDir(), canonicalizedRepoPath)
	log.Printf("Config location of %s", configLocation)
	jsonFile, err := os.Open(configLocation)

	if err != nil {
		log.Printf("Could not read config file from %s", configLocation)
		return Config{}, err
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var config Config
	json.Unmarshal([]byte(byteValue), &config)

	return config, nil
}

func fileExists(path string) (bool, error) {
	repoExists := true
	_, err = os.Stat(path)
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
	if os.IsNotExist(err) {
		return nil
	}

	err = os.RemoveAll(caches)

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

	// TODO: make the command and exec parameters to it's easier to test and
	// TODO: generic to use in other areas
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
