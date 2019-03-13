package cmd

import (
	"log"
	"os"
	"testing"
)

func TestInitCachesDir(t *testing.T) {
	err := os.RemoveAll("../.caches/")
	if err != nil {
		log.Fatal("Could not remove the .caches directory....")
		t.Error("Could not remove the .caches directory...", err)
	}

	err = initCachesDir()

	if err != nil {
		log.Fatal("Could not create the .caches/ directory...")
		t.Error("Could not create the .caches/ directory", err)
	}

	_, err = os.Stat("../.caches")
	if os.IsNotExist(err) {
		t.Error(".caches/ directory does not exist!")
	}
}

func TestCheckForLocalRepo(t *testing.T) {
	err := initCachesDir()
	if err != nil {
		t.Error("Could not create the .caches/ directory")
	}

	repoThatDoesntExist := "HunterTestRepo"

	log.Printf("Looking for repo %s", repoThatDoesntExist)

	repoExists := localRepoExists(repoThatDoesntExist)
	if repoExists == true {
		t.Errorf("Repo %s should not exist, but dooes!", repoThatDoesntExist)
	}
}