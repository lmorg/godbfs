package main

import (
	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/pathfs"
	"log"
	"strings"
)

type filesystem struct {
	pathfs.FileSystem
}

// SplitPath separates file name from path
func SplitPath(pathfile string) (path, file string) {
	split := strings.Split(pathfile, "/")

	switch {
	case len(split) == 1:
		file = split[0]
	case len(split) > 1:
		file = split[len(split)-1]
		path = strings.Join(split[:len(split)-1], "/")
	}

	return
}

func (fs *filesystem) GetAttr(pathfile string, context *fuse.Context) (*fuse.Attr, fuse.Status) {
	path, file := SplitPath(pathfile)

	// Get file metadata
	row := db.QueryRow(sqlGetMetaAttr, path, file)
	if row == nil {
		log.Println("sqlGetMetaAttr returned nothing")
		return nil, fuse.ENOENT
	}

	attr := new(fuse.Attr)
	err := row.Scan(&attr.Ino, &attr.Atime, &attr.Ctime, &attr.Mtime, &attr.Uid, &attr.Gid, &attr.Size, &attr.Mode)
	if err != nil {
		log.Println("Error scanning sqlGetMetaAttr:", err, pathfile)
		return nil, fuse.ENOENT
	}

	return attr, fuse.OK
}
