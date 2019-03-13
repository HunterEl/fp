package cmd

import (
	"fmt"
	"log"
	"os"
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
func localRepoExists(repo string) bool {
	repoLocation := fmt.Sprintf("%s/%s/", getCachesDir(), repo)
	log.Printf("Looking for repo dir %s", repoLocation)
	repoExists := true
	_, err := os.Stat(repoLocation)
	log.Print(err)
	if os.IsNotExist(err) {
		log.Printf("Repo %s does not exist in %s folder...", repoLocation, getCachesDir())
		repoExists = false
	}

	return repoExists
}

// NOTE: We will use '.caches' for the time being to store all remote repositories
func initCachesDir() (err error) {
	_, err = os.Stat(getCachesDir())
	if os.IsNotExist(err) {
		log.Printf("%s does not exist. creating now...", getCachesDir())
		err = os.Mkdir(getCachesDir(), 0777)
	}

	if err != nil {
		log.Printf("Could not create the %s directory...", getCachesDir())
	}

	return err
}

func removeCachesDir() (err error) {
	caches := getCachesDir()
	log.Printf("Checking caches dir %s", caches)
	_, err = os.Stat(caches)
	if !os.IsNotExist(err) {
		log.Printf("Removing %s", getCachesDir())
		err = os.RemoveAll(caches)
	}

	return err
}
