package controllers

import (
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"

	"github.com/astaxie/beego"
	"github.com/nfnt/resize"
)

// Photo struct represents a pohto
type Photo struct {
	Origin string
	Big    string
	Small  string
	Title  string
}

// Link struct represents a link
type Link struct {
	Url   string
	Title string
}

// MainController main controller
type MainController struct {
	beego.Controller
}

func traverse(rootPath string) (photos []Photo, links []Link) {
	walkFunc := func(itemPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if info.Name() != "__s__" && info.Name() != "__b__" && itemPath != rootPath {
				encodedPath := base64.StdEncoding.EncodeToString([]byte(info.Name()))
				encodedPath = strings.Replace(encodedPath, "/", ":slash:", -1)
				links = append(links, Link{
					Url:   "/p/" + encodedPath,
					Title: info.Name(),
				})
			}
			if itemPath == rootPath {
				return nil
			}
			return filepath.SkipDir
		}
		ext := strings.ToLower(filepath.Ext(itemPath))
		if ext == ".png" || ext == ".jpg" || ext == ".jpeg" {
			encodedPath := base64.StdEncoding.EncodeToString([]byte(info.Name()))
			encodedPath = strings.Replace(encodedPath, "/", ":slash:", -1)
			photos = append(photos, Photo{
				Origin: "/i/" + encodedPath,
				Small:  "/s/" + encodedPath,
				Big:    "/b/" + encodedPath,
				Title:  info.Name(),
			})
		}
		return nil
	}
	filepath.Walk(rootPath, walkFunc)
	return
}

// Get return main page
func (c *MainController) Get() {
	rootPath := beego.AppConfig.String("docroot")
	fmt.Println("root directory", rootPath)
	photos, links := traverse(rootPath)
	c.Data["Photos"] = photos
	c.Data["Title"] = "Our Photos"
	c.Data["Links"] = links
	c.TplName = "index.tpl"
}

// GetImage return image
func (c *MainController) GetImage() {
	path := c.Ctx.Input.Param(":path")
	path = strings.Replace(path, ":slash:", "/", -1)
	rawPath, err := base64.StdEncoding.DecodeString(path)
	if err != nil {
		return
	}
	http.ServeFile(c.Ctx.ResponseWriter, c.Ctx.Request, string(rawPath))
}

func isFileExists(path string) (bool, error) {
	stat, err := os.Stat(path)
	if err == nil {
		if stat.Mode()&os.ModeType == 0 {
			return true, nil
		}
		return false, errors.New(path + " exists but is not regular file")
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func saveImage(img *image.Image, savePath string) (err error) {
	dir := filepath.Dir(savePath)
	if b, e := isFileExists(dir); !b || e != nil {
		os.MkdirAll(dir, 0755)
	}
	var file *os.File
	if f, err := os.OpenFile(savePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644); err != nil {
		log.Fatal(savePath, err)
	} else {
		file = f
	}
	defer file.Close()
	ext := filepath.Ext(savePath)
	switch strings.ToLower(ext) {
	case ".png":
		err = png.Encode(file, *img)
	default:
		err = jpeg.Encode(file, *img, &jpeg.Options{100})
	}

	if err != nil {
		log.Println(savePath, err)
		return err
	}
	return nil
}

func scaleImage(filepath string, fileSavePath string, width int, height int) error {
	reader, err := os.Open(filepath)
	if err != nil {
		log.Println(filepath, err)
		return err
	}
	defer reader.Close()
	m, _, err := image.Decode(reader)
	if err != nil {
		log.Println(filepath, err)
		return err
	}
	// scale it
	if width == 75 && height == 75 {
		im := resize.Resize(uint(width), uint(height), m, resize.Bilinear)
		return saveImage(&im, fileSavePath)
	}

	bounds := m.Bounds()

	if bounds.Size().X > width || bounds.Size().Y > height {
		// scale it
		var im image.Image
		if bounds.Size().X > width &&
			bounds.Size().Y*width/bounds.Size().X < height {
			im = resize.Resize(uint(width), 0, m, resize.Bilinear)
		} else {
			im = resize.Resize(0, uint(height), m, resize.Bilinear)
		}
		return saveImage(&im, fileSavePath)
	}

	return saveImage(&m, fileSavePath)
}

// GetSmallImage return image
func (c *MainController) GetSmallImage() {
	path := c.Ctx.Input.Param(":path")
	path = strings.Replace(path, ":slash:", "/", -1)
	rawPath, err := base64.StdEncoding.DecodeString(path)
	if err != nil {
		return
	}
	docroot := beego.AppConfig.String("docroot")
	imgPath := fmt.Sprintf("%s%c%s", docroot, os.PathSeparator, string(rawPath))
	idx := strings.LastIndexByte(imgPath, os.PathSeparator)
	smallImgPath := fmt.Sprintf("%s%c__s__%c%s", imgPath[:idx], os.PathSeparator, os.PathSeparator, imgPath[idx+1:])
	if b, e := isFileExists(smallImgPath); !b || e != nil {
		scaleImage(imgPath, smallImgPath, 75, 75)
	}
	http.ServeFile(c.Ctx.ResponseWriter, c.Ctx.Request, smallImgPath)
}

// GetBigImage return image
func (c *MainController) GetBigImage() {
	path := c.Ctx.Input.Param(":path")
	path = strings.Replace(path, ":slash:", "/", -1)
	rawPath, err := base64.StdEncoding.DecodeString(path)
	if err != nil {
		return
	}
	docroot := beego.AppConfig.String("docroot")
	imgPath := fmt.Sprintf("%s%c%s", docroot, os.PathSeparator, string(rawPath))
	idx := strings.LastIndexByte(imgPath, os.PathSeparator)
	bigImgPath := fmt.Sprintf("%s%c__b__%c%s", imgPath[:idx], os.PathSeparator, os.PathSeparator, imgPath[idx+1:])
	if b, e := isFileExists(bigImgPath); !b || e != nil {
		scaleImage(imgPath, bigImgPath, 1024, 1024)
	}
	http.ServeFile(c.Ctx.ResponseWriter, c.Ctx.Request, bigImgPath)
}

// GetPage return a specified page
func (c *MainController) GetPage() {
	path := c.Ctx.Input.Param(":path")
	path = strings.Replace(path, ":slash:", "/", -1)
	rawPath, err := base64.StdEncoding.DecodeString(path)
	if err != nil {
		fmt.Println("can't decode path", err)
		return
	}
	currentPath := string(rawPath)
	docroot := beego.AppConfig.String("docroot")
	rootPath := fmt.Sprintf("%s%c%s", docroot, os.PathSeparator, currentPath)

	var photos []Photo
	var links []Link
	walkFunc := func(itemPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if info.Name() != "__s__" && info.Name() != "__b__" && itemPath != rootPath {
				encodedPath := base64.StdEncoding.EncodeToString([]byte(itemPath[len(docroot):]))
				encodedPath = strings.Replace(encodedPath, "/", ":slash:", -1)
				links = append(links, Link{
					Url:   "/p/" + encodedPath,
					Title: info.Name(),
				})
			}
			if itemPath == rootPath {
				return nil
			}
			return filepath.SkipDir
		}

		ext := strings.ToLower(filepath.Ext(itemPath))
		if ext == ".png" || ext == ".jpg" || ext == ".jpeg" {
			encodedPath := base64.StdEncoding.EncodeToString([]byte(itemPath[len(docroot):]))
			encodedPath = strings.Replace(encodedPath, "/", ":slash:", -1)
			photos = append(photos, Photo{
				Origin: "/i/" + encodedPath,
				Small:  "/s/" + encodedPath,
				Big:    "/b/" + encodedPath,
				Title:  info.Name(),
			})
		}
		return nil
	}

	filepath.Walk(rootPath, walkFunc)

	idx := strings.LastIndexByte(currentPath, os.PathSeparator)
	if idx >= 0 {
		encodedPath := base64.StdEncoding.EncodeToString([]byte(currentPath[:idx]))
		encodedPath = strings.Replace(encodedPath, "/", ":slash:", -1)
		if encodedPath == "" {
			links = append(links, Link{
				Url:   "/",
				Title: "返回上一级目录",
			})
		} else {
			links = append(links, Link{
				Url:   "/p/" + encodedPath,
				Title: "返回上一级目录",
			})
		}
	} else {
		links = append(links, Link{
			Url:   "/",
			Title: "返回上一级目录",
		})
	}

	c.Data["Photos"] = photos
	c.Data["Title"] = currentPath
	c.Data["Links"] = links
	c.TplName = "index.tpl"
}
