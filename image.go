package main

import (
	"image/jpeg"
	"os"
	"path/filepath"
	"strconv"

	"github.com/nfnt/resize"
	"gopkg.in/h2non/bimg.v0"
)

func MakeThumbnailVIPS(path string, name string, w int, h int) string {

	newname := strconv.Itoa(int(w)) + "_" + name
	fo := filepath.Join(path, name)
	fn := filepath.Join(path, newname)

	if _, err := os.Stat(fn); err == nil {
		Trace.Printf("thumbnail already exists: %s", fn)
		return newname
	} else {
		Info.Printf("creating thumbnail: %s", fn)
	}

	options := bimg.Options{
		Width:   w,
		Height:  h,
		Crop:    true,
		Quality: 99,
	}

	buffer, err := bimg.Read(fo)
	if err != nil {
		Error.Fatal(err)
	}

	newImage, err := bimg.NewImage(buffer).Process(options)
	if err != nil {
		Error.Fatal(err)
	}

	bimg.Write(fn, newImage)
	return newname

}

func MakeThumbnail(path string, name string, w uint, h uint) string {

	newname := strconv.Itoa(int(w)) + "_" + name
	fo := filepath.Join(path, name)
	fn := filepath.Join(path, newname)

	if _, err := os.Stat(fn); err == nil {
		Trace.Printf("thumbnail already exists: %s", fn)
		return newname
	} else {
		Info.Printf("creating thumbnail: %s", fn)
	}

	file, err := os.Open(fo)
	if err != nil {
		Error.Fatal(err)
	}

	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		Error.Fatal(err)
	}
	file.Close()

	// resize to width [wxh] using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Resize(w, h, img, resize.MitchellNetravali)

	out, err := os.Create(fn)
	if err != nil {
		Error.Fatal(err)
	}
	defer out.Close()

	// write new image to file
	jpeg.Encode(out, m, nil)
	return newname

}
