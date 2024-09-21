package log_test

import (
	"os"
	"strings"
	"testing"

	"github.com/base-server/go/common/config"
	"github.com/base-server/go/common/log"
	"github.com/common-library/go/file"
)

func Test(t *testing.T) {
	path, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	configFile := path + "/../config/config.yaml"

	if err := config.Read(configFile); err != nil {
		t.Fatal(err)
	}

	kind := "gRPC"
	fileName := config.Get(kind+".log.file.name").(string) + "." + config.Get(kind+".log.file.extensionName").(string)

	log.Initialize(kind)
	defer file.Remove(fileName)

	content := "test"
	log.Log.Info(content)
	log.Log.Flush()

	if data, err := file.Read(fileName); err != nil {
		t.Fatal(err)
	} else if strings.Contains(data, `"msg":"`+content+`"`) == false {
		t.Fatal("invalid :", data, ",", content)
	}
}
