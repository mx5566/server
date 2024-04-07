package excelt

import (
	"fmt"
	"github.com/xiaonanln/goworld/engine/gwioutil"
	"testing"
)

func TestGetFileList(t *testing.T) {
	var listpath = "."
	//listpath, _ = os.Getwd()
	_, _ = fmt.Scanf("%s", &listpath)
	var filter gwioutil.FileFilter
	_ = filter.GetFileList(listpath, ".xslx")
	ListFileFunc(filter.ListFile)
}

func TestCombineKeys(t *testing.T) {
	t.Log(CombineKeys(1, "key1", 100, "key2"))

	//t.Log(reflect.TypeOf(nil).Name())
}
