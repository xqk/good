package http

import (
	"fmt"
	"github.com/xqk/good/pkg/ilog"
	"github.com/xqk/good/pkg/util/ifile"
	"io/ioutil"
	"os"
)

// GetFileName ...
func GetFileName(cacheKey string, cacheDir string) string {
	return cacheDir + string(os.PathSeparator) + cacheKey
}

// WriteConfigToFile ...
func WriteConfigToFile(cacheKey string, cacheDir string, content string) {
	if err := ifile.MkdirIfNecessary(cacheDir); err != nil {
		ilog.Errorf("[ERROR]:faild to MkdirIfNecessary config ,value:%s ,err:%s \n", string(content), err.Error())
		return
	}
	fileName := GetFileName(cacheKey, cacheDir)
	err := ioutil.WriteFile(fileName, []byte(content), 0666)
	if err != nil {
		ilog.Errorf("[ERROR]:faild to write config  cache:%s ,value:%s ,err:%s \n", fileName, string(content), err.Error())
	}
}

// ReadConfigFromFile ...
func ReadConfigFromFile(cacheKey string, cacheDir string) (string, error) {
	fileName := GetFileName(cacheKey, cacheDir)
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", fmt.Errorf("failed to read config cache file:%s,err:%s! ", fileName, err.Error())
	}
	return string(b), nil
}
