package main

import (
	"bufio"
	"fmt"
	"image/color"
	"log"
	"os"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg/draw"
)

func main() {
	xys, err := readData("tviq.txt")
	if err != nil {
		log.Fatalf("Could not read data: %v", err)
	}
	for _, xy := range xys {
		fmt.Println(xy.x, xy.y)
	}

	err = plotData("tviq.png", xys)
	if err != nil {
		log.Fatalf("Could not plot data: %v", err)
	}
}

type xy struct{ x, y float64 }

func readData(path string) ([]xy, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var xys []xy
	s := bufio.NewScanner(f)

	for s.Scan() {
		var x, y float64
		_, err = fmt.Sscanf(s.Text(), "%f\t%f", &x, &y)
		if err != nil {
			log.Printf("Discard data point %q: %v\n", s.Text(), err)
			continue
		}
		xys = append(xys, xy{x, y})
	}

	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %v", err)
	}
	return xys, nil
}

func plotData(path string, xys []xy) error {
	p, err := plot.New()
	if err != nil {
		return fmt.Errorf("create plot: %v", err)
	}

	sp, err := scatterPlot(xys)
	if err != nil {
		return fmt.Errorf("scatter plot: %v", err)
	}
	p.Add(sp)
	p.Add(plotter.NewGrid())

	err = p.Save(256, 256, path)
	if err != nil {
		return fmt.Errorf("save plot: %v", err)
	}

	return nil
}

func scatterPlot(xys []xy) (plot.Plotter, error) {
	pxys := make(plotter.XYs, len(xys))
	for i, xy := range xys {
		pxys[i].X = xy.x
		pxys[i].Y = xy.y
	}

	s, err := plotter.NewScatter(pxys)
	if err != nil {
		return nil, fmt.Errorf("bad data points: %v", err)
	}
	s.GlyphStyle.Shape = draw.CrossGlyph{}
	s.Color = color.RGBA{R: 255, A: 255}

	return s, nil
}
