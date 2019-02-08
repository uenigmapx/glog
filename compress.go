package glog

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	// "time"
)

// the -logcompress flag is of compress method(zip/gzip/bzip2(todo)/none)
type compress func([]string)

var logCompress compress = compressNone
var uncompresses chan<- []string

// String is part of the flag.Value interface. how to `String`
func (c *compress) String() string {
	return fmt.Sprintf("%d", *c)
}

// Get is part of the flag.Value interface.
func (c *compress) Set(value string) error {
	if value == "zip" {
		logCompress = compressZip
	} else if value == "gzip" {
		logCompress = compressGzip
	} else if value == "bzip2" {
		logCompress = compressBzip2
	} else if value == "none" || value == "" {
		logCompress = compressNone
	} else {
		return fmt.Errorf("logcompress flag has wrong")
	}
	return nil
}

type countPerCompress int32

var logCountPerCompress countPerCompress

// String is part of the flag.Value interface. how to `String`
func (c *countPerCompress) String() string {
	return fmt.Sprintf("%d", *c)
}

// Get is part of the flag.Value interface.
func (c *countPerCompress) Set(value string) error {
	// default
	if value == "" {
		*c = 0
		return nil
	}
	// special
	i, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return err
	}
	*c = countPerCompress(i)
	return nil
}

// timely detection of uncompressed logs and archiving
func detectUncompressed() {
	// 未压缩文件列表
	var uncompressList []string
	// 传过来的未压缩文件
	var uncompress chan []string
	var mu sync.Mutex
	// init
	// normal update
	uncompresses = uncompress
	for {
		select {
		case once := <-uncompress:
			mu.Lock()
			uncompressList = append(uncompressList, once...)
			if logCountPerCompress < 1 ||
				countPerCompress(len(uncompressList)) < logCountPerCompress {
				continue
			}
			// 开始压缩
			go logCompress(uncompressList)
			// 重置列表
			uncompressList = []string{}
			mu.Unlock()
		}
	}
}

func compressNone(src []string) {
}

func compressZip(src []string) {
	availName, err := compressName(src[0], "", ".zip")
	// 压缩失败
	if err != nil {
		return
	}
	f, err := os.Create(availName)
	if err != nil {
		return
	}
	defer f.Close()
	w := zip.NewWriter(f)
	defer w.Close()

	// 压缩文件
	for _, v := range src {
		srcInfo, err := os.Stat(v)
		if err != nil {
			continue
		}

		header, err := zip.FileInfoHeader(srcInfo)
		if err != nil {
			continue
		}

		writer, err := w.CreateHeader(header)
		if err != nil {
			continue
		}
		reader, err := os.Open(v)
		if err != nil {
			continue
		}
		io.Copy(writer, reader) // src to zip buffer file
		reader.Close()
		os.Remove(v) // del source file
	}
}

func compressGzip(src []string) {
	availName, err := compressName(src[0], "", ".tar.gz")
	// 压缩失败
	if err != nil {
		return
	}

	f, err := os.Create(availName)
	if err != nil {
		return
	}
	defer f.Close()
	// archive --> compress
	gw := gzip.NewWriter(f)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	for _, v := range src {
		finfo, err := os.Stat(v)
		if err != nil {
			continue
		}
		h := &tar.Header{
			Name:    finfo.Name(),
			Size:    finfo.Size(),
			Mode:    int64(finfo.Mode()),
			ModTime: finfo.ModTime(),
		}
		err = tw.WriteHeader(h)
		if err != nil {
			continue
		}
		reader, err := os.Open(v)
		if err != nil {
			continue
		}
		io.Copy(tw, reader)
		os.Remove(v) // del source file
	}
}

// TODO: 未完成 -.-
func compressBzip2(src []string) {
}

func compressName(src string, prefix string, suffix string) (availName string, err error) {
	logsys := strings.SplitN(src, ".", 6)
	tmpLogName := strings.Join(append(logsys[:4], logsys[5]), ".")
	var i int
	for i = 0; i < 8; i++ {
		availName = *logDir + string(os.PathSeparator) + prefix +
			fmt.Sprintf("%s-%d-", tmpLogName, i) + suffix
		_, err := os.Stat(availName)
		if os.IsExist(err) {
			continue
		}
		break
	}
	// 没有不存在的可用文件
	if i == 8 {
		return "", fmt.Errorf("No file name available")
	}
	return availName, nil
}
