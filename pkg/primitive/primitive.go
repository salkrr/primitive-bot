package primitive

import (
	"math/rand"
	"runtime"
	"time"

	"github.com/fogleman/primitive/primitive"
	"github.com/nfnt/resize"
)

type Shape int

const (
	ShapeAny Shape = iota
	ShapeTriangle
	ShapeRectangle
	ShapeEllipse
	ShapeCircle
	ShapeRotatedRectangle
	ShapeBezier
	ShapeRotatedEllipse
	ShapePolygon
)

type Config struct {
	workers    int
	OutputSize int
	Shape      Shape
	Iterations int
	Repeat     int
	Alpha      int
	Extension  Extension
}

func NewConfig() Config {
	return Config{
		workers:    runtime.NumCPU(),
		OutputSize: 1280,
		Shape:      ShapeAny,
		Iterations: 200,
		Repeat:     1,
		Alpha:      128,
		Extension:  JPG,
	}
}

type Extension string

const (
	PNG = "png"
	JPG = "jpg"
	SVG = "svg"
	GIF = "gif"
)

func (c Config) Create(inputPath, outputPath string) error {
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
	model := primitive.NewModel(input, bg, c.OutputSize, c.workers)
	for i := 0; i < c.Iterations; i++ {
		// find optimal shape and add it to the model
		model.Step(primitive.ShapeType(c.Shape), c.Alpha, c.Repeat)
	}

	// write output image
	switch c.Extension {
	case PNG:
		err = primitive.SavePNG(outputPath, model.Context.Image())
		if err != nil {
			return err
		}
	case JPG:
		err = primitive.SaveJPG(outputPath, model.Context.Image(), 95)
		if err != nil {
			return err
		}
	case SVG:
		err = primitive.SaveFile(outputPath, model.SVG())
		if err != nil {
			return err
		}
	case GIF:
		frames := model.Frames(0.001)
		err = primitive.SaveGIFImageMagick(outputPath, frames, 50, 250)
		if err != nil {
			return err
		}
	}

	return nil
}
