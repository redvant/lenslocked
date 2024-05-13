package models

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const DefaultImagesDir = "images"

type Image struct {
	GalleryID int
	Path      string
	Filename  string
}

type ImagesConfig struct {
	ImagesDir           string   `env:"IMAGES_DIR"`
	AllowedExtensions   []string `env:"IMAGES_ALLOWED_EXTENSIONS" envSeparator:","`
	AllowedContentTypes []string `env:"IMAGES_ALLOWED_TYPES" envSeparator:","`
}

type ImageService struct {
	// ImagesDir is used to tell the ImageService where to store
	// and locate images. If not set, the ImageService will default
	// to usign the DefaultImagesDir directory.
	ImagesDir string

	// AllowedExtensions is used to restrict the extensions that
	// a file to be uploaded can have. If not set, the ImageService
	// will default to: .png, .jpg, .jpeg, .gif.
	AllowedExtensions []string

	// AllowedContentTypes is used to restrict the content type that
	// a file to be uploaded can have. If not set, the ImageService
	// will default to: image/png, image/jpeg, image/gif.
	AllowedContentTypes []string
}

func NewImageService(config ImagesConfig) *ImageService {
	is := ImageService(config)
	return &is
}

func (is *ImageService) Images(galleryID int) ([]Image, error) {
	globPattern := filepath.Join(is.galleryDir(galleryID), "*")
	allFiles, err := filepath.Glob(globPattern)
	if err != nil {
		return nil, fmt.Errorf("retriving gallery-%d images: %w", galleryID, err)
	}
	var images []Image
	for _, file := range allFiles {
		if hasExtension(file, is.extensions()) {
			images = append(images, Image{
				GalleryID: galleryID,
				Path:      file,
				Filename:  filepath.Base(file),
			})
		}
	}
	return images, nil
}

func (is *ImageService) Image(galleryID int, filename string) (Image, error) {
	imagePath := filepath.Join(is.galleryDir(galleryID), filename)
	_, err := os.Stat(imagePath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return Image{}, ErrNotFound
		}
		return Image{}, fmt.Errorf("querying for image %q: %w", filename, err)
	}
	return Image{
		GalleryID: galleryID,
		Path:      imagePath,
		Filename:  filename,
	}, nil
}

func (is *ImageService) CreateImage(galleryID int, filename string, contents io.ReadSeeker) error {
	err := checkContentType(contents, is.imageContentTypes())
	if err != nil {
		return fmt.Errorf("creating image %v: %w", filename, err)
	}
	err = checkExtension(filename, is.extensions())
	if err != nil {
		return fmt.Errorf("creating image %v: %w", filename, err)
	}
	galleryDir := is.galleryDir(galleryID)
	err = os.MkdirAll(galleryDir, 0o755)
	if err != nil {
		return fmt.Errorf("creating gallery-%d images directory: %w", galleryID, err)
	}
	imagePath := filepath.Join(galleryDir, filename)
	dst, err := os.Create(imagePath)
	if err != nil {
		return fmt.Errorf("creating image file: %w", err)
	}
	defer dst.Close()
	_, err = io.Copy(dst, contents)
	if err != nil {
		return fmt.Errorf("copying contents to image: %w", err)
	}
	return nil
}

func (is *ImageService) DeleteImage(galleryID int, filename string) error {
	image, err := is.Image(galleryID, filename)
	if err != nil {
		return fmt.Errorf("deleting image: %w", err)
	}
	err = os.Remove(image.Path)
	if err != nil {
		return fmt.Errorf("deleting image: %w", err)
	}
	return nil
}

func (is *ImageService) DeleteImages(galleryID int) error {
	err := os.RemoveAll(is.galleryDir(galleryID))
	if err != nil {
		return fmt.Errorf("delete images: %w", err)
	}
	return nil
}

func (is *ImageService) extensions() []string {
	if len(is.AllowedExtensions) > 0 {
		return is.AllowedExtensions
	}
	return []string{".png", ".jpg", ".jpeg", ".gif"}
}

func (is *ImageService) GetAllowedContentTypesString() string {
	types := is.imageContentTypes()
	return strings.Join(types, ", ")
}

func (is *ImageService) imageContentTypes() []string {
	if len(is.AllowedContentTypes) > 0 {
		return is.AllowedContentTypes
	}
	return []string{"image/png", "image/jpeg", "image/gif"}
}

func (is *ImageService) galleryDir(id int) string {
	imagesDir := is.ImagesDir
	if imagesDir == "" {
		imagesDir = DefaultImagesDir
	}
	return filepath.Join(imagesDir, fmt.Sprintf("gallery-%d", id))
}

func hasExtension(file string, extensions []string) bool {
	for _, ext := range extensions {
		file = strings.ToLower(file)
		ext = strings.ToLower(ext)
		if filepath.Ext(file) == ext {
			return true
		}
	}
	return false
}
