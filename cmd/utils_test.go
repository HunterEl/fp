package cmd

import (
	"log"
	"os"
	"strings"
	"testing"
)

func TestInitCachesDir(t *testing.T) {
	err := os.RemoveAll(getCachesDir())
	if err != nil {
		t.Error("Could not remove the .caches directory...", err)
	}

	err = initCachesDir()
	defer removeCachesDir()

	if err != nil {
		t.Error("Could not create the .caches/ directory", err)
	}

	_, err = os.Stat(getCachesDir())
	if os.IsNotExist(err) {
		t.Error(".caches/ directory does not exist!")
	}
}

func TestCheckForLocalRepo(t *testing.T) {
	err := initCachesDir()
	if err != nil {
		t.Error("Could not create the .caches/ directory")
	}

	defer removeCachesDir()

	repoThatDoesntExist := "HunterTestRepo"

	log.Printf("Looking for repo %s", repoThatDoesntExist)

	repoExists, err := localRepoExists(repoThatDoesntExist)
	if err != nil {
		t.Error(err)
	}

	if repoExists == true {
		t.Errorf("Repo %s should not exist, but dooes!", repoThatDoesntExist)
	}
}

func TestFetchRepo(t *testing.T) {
	err := removeCachesDir()
	if err != nil {
		t.Error(err)
	}

	err = initCachesDir()
	if err != nil {
		t.Error("Could not create the .caches directory")
	}

	// defer removeCachesDir()

	repoToFetch := "https://github.com/blisspointmedia/Scripts"

	err = removeRepo(repoToFetch)
	if err != nil {
		t.Errorf("Could not remove the existing repo %s", repoToFetch)
	}

	result, err := fetchRepo(repoToFetch)
	log.Printf("Fetch command output: %s", result)
	// defer removeRepo(repoToFetch)

	if err != nil {
		t.Errorf("Could not fetch the repo %s", repoToFetch)
	}

	repoExists, err := localRepoExists(repoToFetch)
	if err != nil {
		t.Error(err)
	}

	if !repoExists {
		t.Errorf("%s did not exist after fetch. Something went wrong...", repoToFetch)
	}
}

func TestRepoURLConversion(t *testing.T) {
	repoToFetch := "https://github.com/blisspointmedia/Scripts"
	expectedPath := "github.com/blisspointmedia/Scripts"

	actualPath, err := repoURLToPathName(repoToFetch)
	if err != nil {
		t.Error(err)
	}

	equal := strings.EqualFold(expectedPath, actualPath)
	if !equal {
		t.Errorf("Expected %s did not match actual %s", expectedPath, actualPath)
	}
}
