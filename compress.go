package glog

import (
	// "archive/tar"
	// "archive/zip"
	// "compress/gzip"
	"fmt"
	// "io"
	// "os"
	"time"
)

// the -logcompress flag is of compress method(zip/gzip/bzip2(todo)/none)
type compress func([]string)

var logCompress compress = compressNone

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

var uncompresses chan<- string

// timely detection of uncompressed logs and archiving
func detectUncompressed() {
	var ticker *time.Ticker
	var uncompress chan string
	// init
	// normal update
	uncompresses = uncompress
	select {
	case <-ticker.C:
	case <-uncompress:
	default:
	}
}

func compressNone(src []string) {
}

func compressZip(src []string) {
	// 压缩文件
	// srcInfo, err := src.Stat()
	// if err != nil {
	// 	return
	// }

	// header, err := zip.FileInfoHeader(srcInfo)
	// if err != nil {
	// 	return
	// }
	// srcName := header.Name

	// defer os.Remove(srcName) // del source file
	// defer src.Close()

	// f, err := os.Create(srcName + ".zip")
	// if err != nil {
	// 	return
	// }
	// defer f.Close()
	// w := zip.NewWriter(f)
	// defer w.Close()

	// writer, err := w.CreateHeader(header)
	// if err != nil {
	// 	return
	// }
	// io.Copy(writer, src) // src to zip buffer file
}

func compressGzip(src []string) {
	// srcName := src.Name()
	// defer os.Remove(srcName) // del source file
	// defer src.Close()

	// f, err := os.Create(srcName + ".tar.gz")
	// if err != nil {
	// 	return
	// }
	// defer f.Close()
	// // archive --> compress
	// gw := gzip.NewWriter(f)
	// defer gw.Close()
	// tw := tar.NewWriter(gw)
	// defer tw.Close()

	// finfo, err := os.Stat(srcName)
	// if err != nil {
	// 	return
	// }
	// h := &tar.Header{
	// 	Name:    finfo.Name(),
	// 	Size:    finfo.Size(),
	// 	Mode:    int64(finfo.Mode()),
	// 	ModTime: finfo.ModTime(),
	// }
	// err = tw.WriteHeader(h)
	// if err != nil {
	// 	return
	// }
	// io.Copy(tw, src)
}

// TODO: 未完成 -.-
func compressBzip2(src []string) {
}
