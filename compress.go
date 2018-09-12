package glog

import (
	"archive/tar"
	"archive/zip"
	// "compress/bzip2"
	"compress/gzip"
	"io"
	"os"
)

func compressNone(src *os.File) {
	defer src.Close()
}

func compressZip(src *os.File) {
	// 压缩文件
	srcInfo, err := src.Stat()
	if err != nil {
		return
	}

	header, err := zip.FileInfoHeader(srcInfo)
	if err != nil {
		return
	}
	srcName := header.Name

	defer os.Remove(srcName) // del source file
	defer src.Close()

	f, err := os.Create(srcName + ".zip")
	if err != nil {
		return
	}
	defer f.Close()
	w := zip.NewWriter(f)
	defer w.Close()

	writer, err := w.CreateHeader(header)
	if err != nil {
		return
	}
	io.Copy(writer, src) // src to zip buffer file
}

func compressGzip(src *os.File) {
	srcName := src.Name()
	defer os.Remove(srcName) // del source file
	defer src.Close()

	f, err := os.Create(srcName + ".tar.gz")
	if err != nil {
		return
	}
	defer f.Close()
	// archive --> compress
	gw := gzip.NewWriter(f)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	finfo, err := os.Stat(srcName)
	if err != nil {
		return
	}
	h := &tar.Header{
		Name:    finfo.Name(),
		Size:    finfo.Size(),
		Mode:    int64(finfo.Mode()),
		ModTime: finfo.ModTime(),
	}
	err = tw.WriteHeader(h)
	if err != nil {
		return
	}
	io.Copy(tw, src)
}

// TODO: 未完成 -.-
func compressBzip2(src *os.File) {
	// srcName := src.Name()
	// defer os.Remove(srcName) // del source file
	defer src.Close()
}
