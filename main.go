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
		"logc starting. (logd: %s, org: %s, merge: %v)\n",
		conf.LogdAddr, conf.OrgName, conf.MergeData,
	)
	logd, err := logdpkg.New(conf.LogdAddr, conf.MergeData)
	if err != nil {
		log.Fatal(err)
	}
	createOrgFiles(logd, conf.OrgName, conf.Files)
	listenOrgFiles(logd, conf.OrgName, conf.Files)
}

func createOrgFiles(logd *logdpkg.Logd, orgName string, filesAry []File) {
	var filesMapping []map[string]interface{}
	for _, file := range filesAry {
		filesMapping = append(filesMapping, map[string]interface{}{
			`name`: file.Name, `mapping`: file.Mapping,
		})
	}
	if err := logd.Create(orgName, filesMapping); err != nil {
		log.Fatal(err)
	}
}

func listenOrgFiles(logd *logdpkg.Logd, orgName string, filesAry []File) {
	wg := sync.WaitGroup{}
	for _, info := range filesAry {
		if file := filespkg.New(orgName, info.Name, info.Path, logd); file != nil {
			wg.Add(1)
			go func() {
				defer wg.Done()
				file.Listen()
			}()
		}
	}
	wg.Wait()
}
