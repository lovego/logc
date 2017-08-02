package main

import (
	"log"
	"sync"

	filespkg "github.com/lovego/logc/files"
	logdpkg "github.com/lovego/logc/logd"
)

func main() {
	conf := getConfig()
	log.Printf(
		"logc starting. (logd: %s, merge: %v)\n",
		conf.LogdAddr, conf.MergeData,
	)
	logd, err := logdpkg.New(conf.LogdAddr, conf.MergeData)
	if err != nil {
		log.Fatal(err)
	}
	createOrgFiles(logd, conf.Files)
	listenOrgFiles(logd, conf.Files)
}

func createOrgFiles(logd *logdpkg.Logd, filesAry []File) {
	filesMappings := make(map[string][]map[string]interface{})
	for _, file := range filesAry {
		orgName := file.OrgName
		if filesMappings[orgName] == nil {
			filesMappings[orgName] = []map[string]interface{}{}
		}
		filesMappings[orgName] = append(filesMappings[orgName], map[string]interface{}{
			`name`: file.Name, `mapping`: file.Mapping,
		})
	}
	for orgName, filesMapping := range filesMappings {
		if err := logd.Create(orgName, filesMapping); err != nil {
			log.Fatal(err)
		}
	}
}

func listenOrgFiles(logd *logdpkg.Logd, filesAry []File) {
	wg := sync.WaitGroup{}
	for _, info := range filesAry {
		if file := filespkg.New(info.OrgName, info.Name, info.Path, logd); file != nil {
			wg.Add(1)
			go func() {
				defer wg.Done()
				file.Listen()
			}()
		}
	}
	wg.Wait()
}
