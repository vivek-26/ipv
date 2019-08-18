package ipvanish

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"strings"

	"github.com/vivek-26/ipv/reporter"
)

const (
	ipvanishServersURL = "https://www.ipvanish.com/software/configs/configs.zip"
	clientConnInfoURL  = "https://www.ipvanish.com/api/get-location.php"
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

	// Read all config files from zip file
	for _, zipFile := range zipReader.File {
		// Filter based on country
		if strings.HasPrefix(zipFile.Name, "ipvanish-"+countryCode) {
			// Get server info from config file
			server := getServerInfoFromZipFile(zipFile)
			if server != nil {
				servers = append(servers, *server)
			}
		}
	}

	return &servers
}

func getServerInfoFromZipFile(zf *zip.File) *IPVServer {
	f, err := zf.Open()
	if err != nil {
		reporter.Error(err)
	}
	defer f.Close()

	dataBytes, err := ioutil.ReadAll(f)
	if err != nil {
		reporter.Error(err)
	}

	vpnConfig := string(dataBytes) // OpenVPN config

	// Get hostname from config
	re := regexp.MustCompile(`(?m)remote(?P<Hostname>.*)?443`)
	matches := re.FindStringSubmatch(vpnConfig)
	if len(matches) != 2 {
		reporter.Warn(
			fmt.Sprintf("Cannot find server hostname from config file %v", zf.Name),
		)
		return nil
	}
	serverHostname := strings.TrimSpace(matches[1])

	// Get IP address of server hostname
	addr, err := net.LookupIP(serverHostname)
	if err != nil {
		return nil
	}

	return &IPVServer{
		IP:       addr[0].String(),
		Hostname: serverHostname,
		Latency:  0,
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
