package glog

import (
	"archive/zip"
	"archive/tar"
	//"compress/bzip2"
	"compress/gzip"
	"io"
	"os"
)

func compressNone(src *os.File) {
	defer src.Close()
}

func compressZip(src *os.File) {
	defer src.Close()
	f, err := os.Create(src.Name()+".zip")
	if err != nil {
		return
	}
	defer f.Close()
	w := zip.NewWriter(f)
	defer w.Close()

	fi, err := w.Create(src.Name())
	if err != nil {
		return
	}
	io.Copy(fi, src) // src to zip buffer file
}

func compressGzip(src *os.File) {
	defer src.Close()
	f, err := os.Create(src.Name()+".tar.gz")
	if err != nil {
		return
	}
	defer f.Close()
	// archive --> compress
	gw := gzip.NewWriter(f)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	finfo, err := os.Stat(src.Name())
	if err != nil {
		return
	}
	h := &tar.Header {
		Name: finfo.Name(),
		Size: finfo.Size(),
		Mode: int64(finfo.Mode()),
		ModTime: finfo.ModTime(),
	}
	err = tw.WriteHeader(h)
	if err !=nil {
		return
	}
	io.Copy(tw, src)
}

func compressBzip2(src *os.File) {
	defer src.Close()
}
