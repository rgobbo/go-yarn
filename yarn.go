package go_yarn


import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"github.com/rgobbo/go-yarn/extractor"
	"github.com/rgobbo/fileutils"
	"strings"
)

//Yarn - Struture to hold yarn.json configuration
type Yarn struct {
	ID          string  `json:"_id"`
	Rev         string  `json:"_rev"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	DistTags    DistTag `json:"dist-tags"`
	Versions    map[string]Version
}

type DistTag struct {
	Latest string `json:"latest"`
	Old    string `json:"old"`
	Next   string `json:"next"`
}

type Version struct {
	Name    string `json:"name"`
	Ver     string `json:"version"`
	Distrib Dist   `json:"dist"`
}

type Dist struct {
	Shasum  string `json:"shasum"`
	Tarball string `json:"tarball"`
}

type YarnConf struct {
	Dependencies []YarnConfItem `json:"dependencies"`
}

type YarnConfItem struct {
	Lib      string `json:"lib"`
	Version  string `json:"version"`
	Resolved string `json:"resolved"`
}

//YarnInstall  - function to install libraries from yarn.conf
//Example : - YarnInstall("./conf/yarn.json","./static/vendor/")
func YarnInstall(confPath string, publishPath string) error {

	if !strings.HasSuffix(publishPath,"/") {
		publishPath = publishPath + "/"
	}

	urlRegistry := "https://registry.yarnpkg.com/"
	var yarnConf YarnConf
	var yarnConfNew YarnConf
	err := fileutils.LoadJson(confPath, &yarnConf)
	if err != nil {
		return err
	}

	for _, dep := range yarnConf.Dependencies {
		resp, err := http.Get(urlRegistry + dep.Lib)
		if err != nil {
			return err
		}

		var yarnResp Yarn
		err = json.NewDecoder(resp.Body).Decode(&yarnResp)
		if err != nil {
			return err
		}
		resp.Body.Close()

		version := ""
		if dep.Version != "" {
			version = dep.Version
		} else {
			version = yarnResp.DistTags.Latest
		}

		for key, ver := range yarnResp.Versions {
			if key == version {
				err := DownloadFile(ver.Distrib.Tarball, publishPath +dep.Lib+".tgz")
				if err != nil {
					return fmt.Errorf("error download dependencie %v -version %v , url : %v", dep.Lib, version, ver.Distrib.Tarball)
				}
				ext := extractor.NewTgz()
				err = ext.Extract(publishPath+dep.Lib+".tgz", publishPath + "tmp")
				if err != nil {
					return err
				}
				os.Remove(publishPath + dep.Lib + ".tgz")
				os.RemoveAll(publishPath + dep.Lib)
				os.Rename(publishPath + "tmp/package", publishPath +dep.Lib)
				err = os.RemoveAll(publishPath + "tmp")
				if err != nil {
					log.Println("Yarn : ", err)
				}
				dep.Version = version
				dep.Resolved = ver.Distrib.Tarball

				yarnConfNew.Dependencies = append(yarnConfNew.Dependencies, dep)

				break
			}
		}

	}

	err = fileutils.SaveJson(confPath, yarnConfNew)
	if err != nil {
		return err
	}

	return nil
}

//DownloadFile - function to downlaod and save file from a given url
func DownloadFile(url string, fileName string) error {

	output, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("error while creating %v - %v", fileName, err)
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error while downloading %v - %v", url, err)
	}
	defer response.Body.Close()

	_, err = io.Copy(output, response.Body)
	if err != nil {
		return fmt.Errorf("error while downloading %v - %v", url, err)
	}

	return nil
}
