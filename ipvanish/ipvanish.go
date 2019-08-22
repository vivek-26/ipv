package ipvanish

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"

	"github.com/mitchellh/go-homedir"
	"github.com/vivek-26/ipv/reporter"
)

const (
	ipvanishServersURL  = "https://www.ipvanish.com/software/configs/configs.zip"
	clientConnInfoURL   = "https://www.ipvanish.com/api/get-location.php"
	configDirName       = ".ipv"
	credentialsFileName = "/.credentials.txt"
)

// GetServers returns pointer to a slice of ipvanish servers.
// Pointer is returned as slice holds many elements.
func GetServers(countryCode string) *[]IPVServer {
	var servers []IPVServer

	resp, err := http.Get(ipvanishServersURL)
	if err != nil {
		reporter.Error(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		reporter.Error(err)
	}

	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		reporter.Error(err)
	}

	// Slice of zip files of selected country
	countryZipFiles := make([]*zip.File, 0)

	// Read all config files from zip file
	for _, zipFile := range zipReader.File {
		// Filter based on country
		if strings.HasPrefix(zipFile.Name, "ipvanish-"+countryCode) {
			countryZipFiles = append(countryZipFiles, zipFile)
		}
	}

	// Get server info from zip files concurrently using
	// worker pool.
	numJobs := len(countryZipFiles)
	numWorkers := getNumberOfWorkers()

	jobs := make(chan *zip.File, 10)
	results := make(chan *IPVServer, 10)

	var wg sync.WaitGroup
	// Start worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go getServerInfoFromZipFile(jobs, results, &wg)
	}

	// Goroutine to push jobs
	wg.Add(1)
	go pushJobs(countryZipFiles, jobs, &wg)

	// Process results
	for i := 0; i < numJobs; i++ {
		server := <-results
		if server != nil {
			servers = append(servers, *server)
		}
	}

	wg.Wait()
	return &servers
}

// getNumberOfWorkers returns the number of worker goroutines
// to be created based on number of CPUs available.
func getNumberOfWorkers() int {
	if runtime.NumCPU() < 4 {
		return 25
	}
	return 50
}

// pushJobs ranges over zip files and puts them on jobs channel
func pushJobs(countryZipFiles []*zip.File, jobs chan<- *zip.File, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := range countryZipFiles {
		jobs <- countryZipFiles[i]
	}
	close(jobs)
}

// getServerInfoFromZipFile itertates over jobs channel to get zip file
// and extracts hostname from it. In case of an error, it writes nil
// to results channel.
func getServerInfoFromZipFile(jobs <-chan *zip.File, results chan<- *IPVServer, wg *sync.WaitGroup) {
	defer wg.Done()

	for zipFile := range jobs {
		f, err := zipFile.Open()
		if err != nil {
			results <- nil
			continue
		}

		dataBytes, err := ioutil.ReadAll(f)
		if err != nil {
			results <- nil
			continue
		}

		if err = f.Close(); err != nil {
			reporter.Warn("Could not close zip file")
		}

		vpnConfig := string(dataBytes) // OpenVPN config

		// Get hostname from config
		re := regexp.MustCompile(`(?m)remote(?P<Hostname>.*)?443`)
		matches := re.FindStringSubmatch(vpnConfig)
		if len(matches) != 2 {
			reporter.Warn(
				fmt.Sprintf("Cannot find server hostname from config file %v", zipFile.Name),
			)
			results <- nil
			continue
		}
		serverHostname := strings.TrimSpace(matches[1])

		// Get IP address of server hostname
		addr, err := net.LookupIP(serverHostname)
		if err != nil {
			results <- nil
			continue
		}

		results <- &IPVServer{
			IP:       addr[0].String(),
			Hostname: serverHostname,
			Latency:  0,
		}
	}
}

// GetClientInfo returns current connection info of user
func GetClientInfo() ClientInfo {
	var clientInfo ClientInfo

	res, err := http.Get(clientConnInfoURL)
	if err != nil {
		reporter.Error(err)
	}
	defer res.Body.Close()

	if err = json.NewDecoder(res.Body).Decode(&clientInfo); err != nil {
		reporter.Error(err)
	}

	return clientInfo
}

// getConfigDirPath returns config directory path
func getConfigDirPath() string {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		reporter.Error(err)
	}

	// Config directory path
	return filepath.Join(home, configDirName)
}

// getCredentialsFilepath returns credentials file path
func getCredentialsFilepath() string {
	return getConfigDirPath() + credentialsFileName
}

// CreateCredentials creates a new credentials file.
// It replaces the previous file if it exists.
func CreateCredentials(username, password string) {
	credentialsFile := getCredentialsFilepath()

	credentials := fmt.Sprintf("%v\n%v\n", username, password)
	if err := ioutil.WriteFile(credentialsFile, []byte(credentials), 0644); err != nil {
		reporter.Error(err)
	}
}
