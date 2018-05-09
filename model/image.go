package model

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"github.com/nfnt/resize"
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

const (
	IMAGE_PNG  = "png"
	IMAGE_JPEG = "jpeg"
	IMAGE_GIF  = "gif"
)

type Image struct {
	Hash      string    `db:"hash char(32) pk"`
	Name      string    `db:"name varchar(256)"`
	Width     int       `db:"width int"`
	Height    int       `db:"height int"`
	Format    string    `db:"format varchar(64)"`
	Size      int64     `db:"size int"`
	CreatedAt time.Time `db:"created_at datetime"`
	UpdatedAt time.Time `db:"updated_at datetime"`
}

func (image *Image) TableName() string {
	return "image"
}

func (image *Image) PK() string {
	return "hash"
}

func (image *Image) Pathfile() string {
	return conf.image_dir + image.Name
}

func (image *Image) FillWithMultipart(src mp.File, header *mp.FileHeader) error {
	image.Size = header.Size
	image.Name = header.Filename
	var conf goImage.Config
	var err error
	switch path.Ext(header.Filename) {
	case ".png":
		conf, err = png.DecodeConfig(src)
		image.Format = IMAGE_PNG
	case ".gif":
		conf, err = gif.DecodeConfig(src)
		image.Format = IMAGE_GIF
	default:
		conf, err = jpeg.DecodeConfig(src)
		image.Format = IMAGE_JPEG
	}
	if err != nil {
		return err
	}
	image.Width = conf.Width
	image.Height = conf.Height
	src.Seek(0, 0)
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

func (image *Image) FileExisted() (exists bool) {
	exists = true
	if _, err := os.Stat(image.Pathfile()); os.IsNotExist(err) {
		exists = false
	}

	return exists
}

func (image *Image) MimeType() string {
	var contentType string
	switch image.Format {
	case IMAGE_PNG:
		contentType = "image/png"
	case IMAGE_GIF:
		contentType = "image/gif"
	default:
		contentType = "image/jpeg"
	}
	return contentType
}

func (image *Image) Resize(w io.Writer, width, height uint, interp resize.InterpolationFunction) error {
	var err error
	var rImage goImage.Image
	f, err := os.Open(image.Pathfile())
	if err != nil {
		return err
	}
	defer f.Close()
	switch image.Format {
	case IMAGE_PNG:
		rImage, err = png.Decode(f)
	case IMAGE_GIF:
		rImage, err = gif.Decode(f)
	default:
		rImage, err = jpeg.Decode(f)
	}
	if err != nil {
		return err
	}
	result := resize.Resize(width, height, rImage, interp)
	switch image.Format {
	case IMAGE_PNG:
		return png.Encode(w, result)
	case IMAGE_GIF:
		return gif.Encode(w, result, nil)
	default:
		return jpeg.Encode(w, result, nil)
	}
	return nil
}

func NewImage() *Image {
	return new(Image)
}

func NewImageRepo() (imageRepo *Repo, err error) {
	return CreateRepo(new(Image))
}
