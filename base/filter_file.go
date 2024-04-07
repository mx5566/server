package base

import (
	"github.com/mx5566/logm"
	"os"
	"path/filepath"
	"strings"
)

// 查找指定的文件列表
type FileFilter struct {
	// file list in directory
	ListFile []string
	// 后缀 eg:.go .xlsx .txt ...
	Suffix string
}

func (this *FileFilter) Listfunc(path string, f os.FileInfo, err error) error {
	var strRet string
	/*strRet, _ = os.Getwd()

	if ostype == "windows" {
		strRet += "\\"
	} else if ostype == "linux" {
		strRet += "/"
	}*/

	if f == nil {
		return err
	}
	if f.IsDir() {
		return nil
	}

	strRet += path

	//用strings.HasSuffix(src, suffix)//判断src中是否包含 suffix结尾
	ok := strings.HasSuffix(strRet, this.Suffix)
	if ok {
		this.ListFile = append(this.ListFile, strRet) //将目录push到listfile []string中
	}

	return nil
}

func (this *FileFilter) GetFileList(path, suffix string) error {
	this.Suffix = suffix
	//var strRet string
	err := filepath.Walk(path, this.Listfunc)

	if err != nil {
		logm.FatalfE("filepath.Walk() returned %v\n", err)
		return err
	}

	return nil
}
