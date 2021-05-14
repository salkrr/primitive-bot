package main

import (
	"math/rand"
	"time"

	"github.com/fogleman/primitive/primitive"
	"github.com/nfnt/resize"
)

type extension int

const (
	png extension = iota
	jpg
	svg
	gif
)

func (ext extension) String() string {
	return [...]string{"png", "jpg", "svg", "gif"}[ext]
}

type config struct {
	workers    int
	outputSize int
	shape      primitive.ShapeType
	iterations int
	repeat     int
	alpha      int
	ext        extension
}

func (c config) create(inputPath, outputPath string) error {
	// seed random number generator
	rand.Seed(time.Now().UTC().UnixNano())

	// read input image
	input, err := primitive.LoadImage(inputPath)
	if err != nil {
		return err
	}

	// scale down input image if needed
	size := uint(256)
	input = resize.Thumbnail(size, size, input, resize.Bilinear)

	// determine background color
	bg := primitive.MakeColor(primitive.AverageImageColor(input))

	// run algorithm
	model := primitive.NewModel(input, bg, c.outputSize, c.workers)
	for i := 0; i < c.iterations; i++ {
		// find optimal shape and add it to the model
		model.Step(c.shape, c.alpha, c.repeat)
	}

	// write output image
	switch c.ext {
	case png:
		err = primitive.SavePNG(outputPath, model.Context.Image())
		if err != nil {
			return err
		}
	case jpg:
		err = primitive.SaveJPG(outputPath, model.Context.Image(), 95)
		if err != nil {
			return err
		}
	case svg:
		err = primitive.SaveFile(outputPath, model.SVG())
		if err != nil {
			return err
		}
	case gif:
		frames := model.Frames(0.001)
		err = primitive.SaveGIFImageMagick(outputPath, frames, 50, 250)
		if err != nil {
			return err
		}
	}

	return nil
}
