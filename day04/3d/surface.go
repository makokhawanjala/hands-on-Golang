// Professional version with error handling and configurability
package main

import (
	"fmt"
	"math"
	"os"
)

// Config holds rendering parameters
type Config struct {
	Width, Height int
	Cells         int
	XYRange       float64
	ZScale        float64
	ViewAngle     float64
}

// DefaultConfig returns sensible defaults
func DefaultConfig() *Config {
	return &Config{
		Width:     600,
		Height:    320,
		Cells:     100,
		XYRange:   30.0,
		ZScale:    128,
		ViewAngle: math.Pi / 6, // 30 degrees
	}
}

// Surface renders a 3D function as SVG
type Surface struct {
	config   *Config
	sin, cos float64 // Pre-calculated trig values
}

// NewSurface creates a new surface renderer
func NewSurface(cfg *Config) *Surface {
	if cfg == nil {
		cfg = DefaultConfig()
	}
	return &Surface{
		config: cfg,
		sin:    math.Sin(cfg.ViewAngle),
		cos:    math.Cos(cfg.ViewAngle),
	}
}

// Render generates the SVG output
func (s *Surface) Render(f func(x, y float64) float64) error {
	// Validate configuration
	if s.config.Cells <= 0 {
		return fmt.Errorf("cells must be positive, got %d", s.config.Cells)
	}

	// Write SVG header
	fmt.Printf("<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: #dc2626; fill: #fbbf24; stroke-width: 0.4; fill-opacity: 0.75' "+
		"width='%d' height='%d'>", s.config.Width, s.config.Height)

	// Generate mesh
	for i := 0; i < s.config.Cells; i++ {
		for j := 0; j < s.config.Cells; j++ {
			ax, ay, err := s.corner(i+1, j, f)
			if err != nil {
				continue // Skip invalid cells
			}
			bx, by, _ := s.corner(i, j, f)
			cx, cy, _ := s.corner(i, j+1, f)
			dx, dy, _ := s.corner(i+1, j+1, f)

			fmt.Printf("<polygon points='%g,%g %g,%g %g,%g %g,%g'/>\n",
				ax, ay, bx, by, cx, cy, dx, dy)
		}
	}

	fmt.Println("</svg>")
	return nil
}

// corner calculates screen coordinates for a grid cell corner
func (s *Surface) corner(i, j int, f func(x, y float64) float64) (float64, float64, error) {
	// Convert grid to math coordinates
	cells := float64(s.config.Cells)
	x := s.config.XYRange * (float64(i)/cells - 0.5)
	y := s.config.XYRange * (float64(j)/cells - 0.5)

	// Calculate surface height
	z := f(x, y)

	// Check for invalid values (NaN, Inf)
	if math.IsNaN(z) || math.IsInf(z, 0) {
		return 0, 0, fmt.Errorf("invalid z value at (%g, %g)", x, y)
	}

	// Calculate scaling factors
	xyscale := float64(s.config.Width) / 2 / s.config.XYRange

	// Project to 2D screen coordinates
	sx := float64(s.config.Width)/2 + (x-y)*s.cos*xyscale
	sy := float64(s.config.Height)/2 + (x+y)*s.sin*xyscale - z*s.config.ZScale

	return sx, sy, nil
}

// SincFunction is the classic sin(r)/r ripple function
func SincFunction(x, y float64) float64 {
	r := math.Hypot(x, y)
	if r == 0 {
		return 1.0 // Handle division by zero
	}
	return math.Sin(r) / r
}

// Main entry point
func main() {
	surface := NewSurface(nil) // Use default config

	if err := surface.Render(SincFunction); err != nil {
		fmt.Fprintf(os.Stderr, "Error rendering surface: %v\n", err)
		os.Exit(1)
	}
}
