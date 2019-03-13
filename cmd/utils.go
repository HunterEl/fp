package cmd

import (
	"fmt"
	"log"
	"os"
)

// TODO: Of course use git hashes/checksums to validate that we have the most up to date versions
// NOTE: We will use '.caches' for the time being to store all remote repositories
func localRepoExists(repo string) bool {
	repoLocation := fmt.Sprintf("../.caches/%s/", repo)
	log.Printf("Looking for repo dir %s", repoLocation)
	repoExists := true
	_, err := os.Stat(repoLocation)
	log.Print(err)
	if os.IsNotExist(err) {
		log.Printf("Repo %s does not exist in .caches/ folder...", repoLocation)
		repoExists = false
	}

	return repoExists
}

// NOTE: We will use '.caches' for the time being to store all remote repositories
func initCachesDir() (err error) {
	_, err = os.Stat("../.caches")
	if os.IsNotExist(err) {
		log.Print(".caches/ does not exist. creating now...")
		err = os.Mkdir("../.caches", 0777)
	}

	if err != nil {
		log.Printf("Could not create the .caches/ directory...")
	}

	return err
}

func removeCachesDir() (err error) {
	_, err = os.Stat("../.caches")

	if !os.IsNotExist(err) {
		log.Printf(".caches/ does not exist...")
		err = os.RemoveAll("../caches/")
	}

	return err
}
