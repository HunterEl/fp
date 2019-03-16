package cmd

import (
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

	// TODO: Test the case where the repo url is fucked up
	repoThatDoesntExist := "https://github.com/HunterTestRepo"

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

	defer removeCachesDir()

	err = initCachesDir()
	if err != nil {
		t.Error("Could not create the .caches directory")
	}

	repoToFetch := "https://github.com/HunterEl/fp-test-scripts/"

	err = removeRepo(repoToFetch)
	if err != nil {
		t.Errorf("Could not remove the existing repo %s", repoToFetch)
	}

	_, err = fetchRepo(repoToFetch)
	defer removeRepo(repoToFetch)

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
	repoToFetch := "https://github.com/HunterEl/fp-test-scripts/"
	expectedPath := "github.com/HunterEl/fp-test-scripts"

	actualPath, err := repoURLToPathName(repoToFetch)
	if err != nil {
		t.Error(err)
	}

	equal := strings.EqualFold(expectedPath, actualPath)
	if !equal {
		t.Errorf("Expected %s did not match actual %s", expectedPath, actualPath)
	}
}
