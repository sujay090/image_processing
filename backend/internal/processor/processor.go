package processor

import (
	"fmt"

	"github.com/h2non/bimg"
)

// TransformOptions represents all possible image transformations
type TransformOptions struct {
	Resize  *ResizeOpts `json:"resize,omitempty"`
	Crop    *CropOpts   `json:"crop,omitempty"`
	Rotate  *int        `json:"rotate,omitempty"`
	Format  *string     `json:"format,omitempty"`
	Flip    *bool       `json:"flip,omitempty"`
	Mirror  *bool       `json:"mirror,omitempty"`
	Filters *FilterOpts `json:"filters,omitempty"`
	Quality *int        `json:"quality,omitempty"`
}

type ResizeOpts struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type CropOpts struct {
	Width  int `json:"width"`
	Height int `json:"height"`
	X      int `json:"x"`
	Y      int `json:"y"`
}

type FilterOpts struct {
	Grayscale bool `json:"grayscale"`
	Sepia     bool `json:"sepia"`
}

// Process applies transformations to the image buffer and returns the processed buffer
func Process(buf []byte, opts TransformOptions) ([]byte, error) {
	img := bimg.NewImage(buf)
	var err error

	// Apply resize
	if opts.Resize != nil {
		img = bimg.NewImage(buf)
		buf, err = img.Resize(opts.Resize.Width, opts.Resize.Height)
		if err != nil {
			return nil, fmt.Errorf("resize failed: %w", err)
		}
	}

	// Apply crop
	if opts.Crop != nil {
		img = bimg.NewImage(buf)
		buf, err = img.Extract(opts.Crop.Y, opts.Crop.X, opts.Crop.Width, opts.Crop.Height)
		if err != nil {
			return nil, fmt.Errorf("crop failed: %w", err)
		}
	}

	// Apply rotation
	if opts.Rotate != nil {
		img = bimg.NewImage(buf)
		angle := bimg.Angle(*opts.Rotate)
		buf, err = img.Rotate(angle)
		if err != nil {
			return nil, fmt.Errorf("rotate failed: %w", err)
		}
	}

	// Apply flip (vertical)
	if opts.Flip != nil && *opts.Flip {
		img = bimg.NewImage(buf)
		buf, err = img.Flip()
		if err != nil {
			return nil, fmt.Errorf("flip failed: %w", err)
		}
	}

	// Apply mirror (horizontal flip)
	if opts.Mirror != nil && *opts.Mirror {
		img = bimg.NewImage(buf)
		buf, err = img.Flop()
		if err != nil {
			return nil, fmt.Errorf("mirror failed: %w", err)
		}
	}

	// Apply filters
	if opts.Filters != nil {
		if opts.Filters.Grayscale {
			img = bimg.NewImage(buf)
			buf, err = img.Process(bimg.Options{
				Interpretation: bimg.InterpretationBW,
			})
			if err != nil {
				return nil, fmt.Errorf("grayscale filter failed: %w", err)
			}
		}
	}

	// Apply format conversion
	if opts.Format != nil {
		imgType := parseImageType(*opts.Format)
		img = bimg.NewImage(buf)
		buf, err = img.Convert(imgType)
		if err != nil {
			return nil, fmt.Errorf("format conversion failed: %w", err)
		}
	}

	// Apply quality/compression
	if opts.Quality != nil {
		img = bimg.NewImage(buf)
		buf, err = img.Process(bimg.Options{
			Quality: *opts.Quality,
		})
		if err != nil {
			return nil, fmt.Errorf("compression failed: %w", err)
		}
	}

	return buf, nil
}

func parseImageType(format string) bimg.ImageType {
	switch format {
	case "jpeg", "jpg":
		return bimg.JPEG
	case "png":
		return bimg.PNG
	case "webp":
		return bimg.WEBP
	case "tiff":
		return bimg.TIFF
	case "gif":
		return bimg.GIF
	case "avif":
		return bimg.AVIF
	default:
		return bimg.JPEG
	}
}
