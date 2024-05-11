package controllers

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/redvant/lenslocked/context"
	"github.com/redvant/lenslocked/errors"
	"github.com/redvant/lenslocked/models"
)

type Galleries struct {
	Templates struct {
		New           Template
		Show          Template
		Edit          Template
		Index         Template
		ShowPublished Template
	}
	GalleryService *models.GalleryService
}

func (g Galleries) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Title string
	}
	data.Title = r.FormValue("title")
	g.Templates.New.Execute(w, r, data)
}

func (g Galleries) Create(w http.ResponseWriter, r *http.Request) {
	var data struct {
		UserID int
		Title  string
	}
	data.UserID = context.User(r.Context()).ID
	data.Title = r.FormValue("title")

	gallery, err := g.GalleryService.Create(data.Title, data.UserID)
	if err != nil {
		g.Templates.New.Execute(w, r, data, err)
		return
	}
	editPath := fmt.Sprintf("/galleries/%d/edit", gallery.ID)
	http.Redirect(w, r, editPath, http.StatusFound)
}

func (g Galleries) ShowPublished(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r, galleryMustBePublished)
	if err != nil {
		return
	}
	galleryData, err := g.galleryData(w, gallery)
	if err != nil {
		return
	}
	g.Templates.ShowPublished.Execute(w, r, galleryData)
}

func (g Galleries) Show(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r, userMustOwnGallery)
	if err != nil {
		return
	}
	galleryData, err := g.galleryData(w, gallery)
	if err != nil {
		return
	}
	g.Templates.Show.Execute(w, r, galleryData)
}

func (g Galleries) Edit(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r, userMustOwnGallery)
	if err != nil {
		return
	}
	var data struct {
		ID        int
		Title     string
		Published bool
	}
	data.ID = gallery.ID
	data.Title = gallery.Title
	data.Published = gallery.Published
	g.Templates.Edit.Execute(w, r, data)
}

func (g Galleries) Update(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r, userMustOwnGallery)
	if err != nil {
		return
	}
	gallery.Title = r.FormValue("title")
	err = g.GalleryService.Update(gallery)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	editPath := fmt.Sprintf("/galleries/%d/edit", gallery.ID)
	http.Redirect(w, r, editPath, http.StatusFound)
}

func (g Galleries) Index(w http.ResponseWriter, r *http.Request) {
	type Gallery struct {
		ID        int
		Title     string
		Published bool
	}
	var data struct {
		Galleries []Gallery
	}
	user := context.User(r.Context())
	galleries, err := g.GalleryService.ByUserID(user.ID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}
	for _, gallery := range galleries {
		data.Galleries = append(data.Galleries, Gallery{
			ID:        gallery.ID,
			Title:     gallery.Title,
			Published: gallery.Published,
		})
	}
	g.Templates.Index.Execute(w, r, data)
}

func (g Galleries) Delete(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r, userMustOwnGallery)
	if err != nil {
		return
	}
	err = g.GalleryService.Delete(gallery.ID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/galleries/", http.StatusFound)
}

func (g Galleries) Publish(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r, userMustOwnGallery)
	if err != nil {
		return
	}
	err = g.GalleryService.Publish(gallery.ID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/galleries/", http.StatusFound)
}

func (g Galleries) Unpublish(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r, userMustOwnGallery)
	if err != nil {
		return
	}
	err = g.GalleryService.Unpublish(gallery.ID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/galleries/", http.StatusFound)
}

func (g Galleries) Image(w http.ResponseWriter, r *http.Request) {
	filename := r.PathValue("filename")
	gallery, err := g.galleryByID(w, r, userMustOwnGallery)
	if err != nil {
		return
	}
	image, err := g.GalleryService.Image(gallery.ID, filename)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "Image not found", http.StatusNotFound)
			return
		}
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	http.ServeFile(w, r, image.Path)
}

func (g Galleries) PublishedImage(w http.ResponseWriter, r *http.Request) {
	filename := r.PathValue("filename")
	gallery, err := g.galleryByID(w, r, galleryMustBePublished)
	if err != nil {
		return
	}
	image, err := g.GalleryService.Image(gallery.ID, filename)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "Image not found", http.StatusNotFound)
			return
		}
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	http.ServeFile(w, r, image.Path)
}

type galleryOpt func(http.ResponseWriter, *http.Request, *models.Gallery) error

func (g Galleries) galleryByID(w http.ResponseWriter, r *http.Request, opts ...galleryOpt) (*models.Gallery, error) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid gallery ID", http.StatusBadRequest)
		return nil, err
	}
	gallery, err := g.GalleryService.ByID(id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "Gallery not found", http.StatusNotFound)
			return nil, err
		}
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return nil, err
	}
	for _, opt := range opts {
		err = opt(w, r, gallery)
		if err != nil {
			return nil, err
		}
	}
	return gallery, nil
}

func userMustOwnGallery(w http.ResponseWriter, r *http.Request, gallery *models.Gallery) error {
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "You are not authorized to access this gallery", http.StatusForbidden)
		return fmt.Errorf("user does not have access to this gallery")
	}
	return nil
}

func galleryMustBePublished(w http.ResponseWriter, r *http.Request, gallery *models.Gallery) error {
	if !gallery.Published {
		http.Error(w, "This gallery is not public", http.StatusForbidden)
		return fmt.Errorf("the gallery is not published")
	}
	return nil
}

type Image struct {
	GalleryID       int
	Filename        string
	FilenameEscaped string
}
type GalleryData struct {
	ID     int
	Title  string
	Images []Image
}

func (g Galleries) galleryData(w http.ResponseWriter, gallery *models.Gallery) (GalleryData, error) {
	var data GalleryData
	images, err := g.GalleryService.Images(gallery.ID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return data, err
	}
	data.ID = gallery.ID
	data.Title = gallery.Title
	for _, image := range images {
		data.Images = append(data.Images, Image{
			GalleryID:       image.GalleryID,
			Filename:        image.Filename,
			FilenameEscaped: url.PathEscape(image.Filename),
		})
	}
	return data, nil
}
