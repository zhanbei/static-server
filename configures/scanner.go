package configures

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

const (
	// The global site to fallthrough.
	DirGlobalSite = "default.sites"

	FileSiteConfigure = "site.json"
)

// Modular Sites: _.domain.com
// Real Sites
func ScanSites(rootDir string) (global *StaticSite, modular StaticSites, sites StaticSites, err error) {
	stat, err := os.Stat(rootDir)
	if os.IsNotExist(err) {
		return
	}
	if !stat.IsDir() {
		return
	}
	files, err := ioutil.ReadDir(rootDir)
	if err != nil {
		return
	}

	sitesWithoutConfigures := make([]string, 0)
	sitesWithConfigures := make([]string, 0)

	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		name := file.Name()
		if strings.HasPrefix(name, ".") || !strings.Contains(name, ".") {
			// 简单地检查域名的合法性！
			continue
		}
		conf, err := ReadSiteConfigure(path.Join(rootDir, name))
		if err != nil {
			fmt.Println("Failed to get/parse site configures:", name, err)
			continue
		}
		if conf == nil {
			sitesWithoutConfigures = append(sitesWithoutConfigures, name)
		} else {
			sitesWithConfigures = append(sitesWithConfigures, name)
		}

		site, special := NewSite(name, conf)
		if special {
			if site.Name == DirGlobalSite {
				global = site
			} else {
				modular = append(modular, site)
			}
		} else {
			sites = append(sites, site)
		}
	}

	if len(sitesWithoutConfigures) > 0 {
		fmt.Println("No site configure found from: [", `"`+strings.Join(sitesWithoutConfigures, `", "`)+`"`, "]")
	}
	if len(sitesWithConfigures) > 0 {
		fmt.Println("Found site configure for: [", `"`+strings.Join(sitesWithConfigures, `", "`)+`"`, "]")
	}

	return
}

func ReadSiteConfigure(dirSite string) (*SiteConfigure, error) {
	target := path.Join(dirSite, FileSiteConfigure)
	_, err := os.Stat(target)
	if os.IsNotExist(err) {
		// The site does exist and
		return nil, nil
	}
	bts, err := ioutil.ReadFile(target)
	if err != nil {
		return nil, err
	}
	fmt.Println(FileSiteConfigure, ":", string(bts))
	conf := new(SiteConfigure)
	err = json.Unmarshal(bts, conf)
	return conf, err
}