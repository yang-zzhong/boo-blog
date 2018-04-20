package model

import (
	"bytes"
	"crypto/md5"
	// "github.com/djimenez/iconv-go"
	"encoding/hex"
	. "github.com/yang-zzhong/go-model"
	goImage "image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	mp "mime/multipart"
	"os"
	"path"
	"time"
)

type Image struct {
	Name    string `db:"name varchar(256)"`
	Width   int    `db:"width int"`
	Height  int    `db:"height int"`
	Format  string `db:"format varchar(64)"`
	Size    int64  `db:"size int"`
	Hash    string `db:"hash char(32) pk"`
	GroupId string `db:"group_id char(32) nil"`

	CreatedAt time.Time `db:"created_at datetime"`
	UpdatedAt time.Time `db:"updated_at datetime"`
	DeletedAt NullTime  `db:"deleted_at datetime nil"`
}

func (image *Image) TableName() string {
	return "image"
}

func (image *Image) PK() string {
	return "hash"
}

func (image *Image) Pathfile() string {
	return "./" + image.Name
}

func (image *Image) FillWithMultipart(src mp.File, header *mp.FileHeader) error {
	image.Size = header.Size
	image.Name = header.Filename
	var conf goImage.Config
	var err error
	switch path.Ext(header.Filename) {
	case ".png":
		conf, err = png.DecodeConfig(src)
		image.Format = ".png"
	case ".jpg":
	case ".jpeg":
		conf, err = jpeg.DecodeConfig(src)
		image.Format = ".jpeg"
	case ".gif":
		conf, err = gif.DecodeConfig(src)
		image.Format = ".gif"
	default:
		conf, err = jpeg.DecodeConfig(src)
		image.Format = ".jpeg"
	}
	if err != nil {
		return err
	}
	image.Width = conf.Width
	image.Height = conf.Height

	buf := new(bytes.Buffer)
	buf.ReadFrom(src)
	md5Sumb := md5.Sum(buf.Bytes())
	image.Hash = hex.EncodeToString(md5Sumb[:])

	return nil
}

/**
 * 判断文件是否存在
 */
func (image *Image) RecordExisted() (exists bool, err error) {
	repo, err := NewImageRepo()
	if err != nil {
		log.Fatal(err)
		return
	}
	repo.Where("hash", image.Hash)
	exists = repo.Count() > 0

	return
}

func (image *Image) FileExisted() (exists bool, err error) {
	exists = false
	return
}

func (image *Image) SaveFile(reader io.Reader) error {
	dist, err := os.Create(image.Pathfile())
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer dist.Close()
	io.Copy(dist, reader)

	return nil
}

func NewImage() *Image {
	return new(Image)
}

func NewImageRepo() (imageRepo *Repo, err error) {
	return CreateRepo(new(Image))
}
