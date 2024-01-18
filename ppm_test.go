package netpbm

import (
	"os"
	"reflect"
	"testing"
)

const (
	imageWidthPPM  = 3
	imageHeightPPM = 2
)

func TestReadPPM(t *testing.T) {
	/*Here if you look closely I used the same procedure as pbm and pgm*/
	testFilename := "/imagesTest/ppmImages/imageP3.ppm"
	file, err := os.Create(testFilename)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(testFilename)

	file.WriteString("P3\n")
	file.WriteString("3 2\n")
	file.WriteString("255\n")
	file.WriteString("255 0 0   0 255 0   0 0 255\n")
	file.WriteString("255 255 0 0 255 255 0 0 255\n")
	file.Close()

	ppm, err := ReadPPM(testFilename)
	if err != nil {
		t.Fatal(err)
	}

	if ppm.magicNumber != "P3" {
		t.Error("Wrong magic number")
	}

	if ppm.width != 3 || ppm.height != 2 {
		t.Error("Wrong size")
	}

	if ppm.max != 255 {
		t.Error("Wrong max value")
	}

	expectedData := [][]Pixel{
		{{255, 0, 0}, {0, 255, 0}, {0, 0, 255}},
		{{255, 255, 0}, {0, 255, 255}, {0, 0, 255}},
	}
	if !reflect.DeepEqual(ppm.data, expectedData) {
		t.Error("Wrong data for P3")
	}
}

func TestSizePPM(t *testing.T) {
	ppm := &PPM{
		data:        [][]Pixel{{{255, 0, 0}, {0, 255, 0}, {0, 0, 255}}, {{255, 255, 0}, {0, 255, 255}, {0, 0, 255}}},
		width:       3,
		height:      2,
		magicNumber: "P3",
		max:         255,
	}
	w, h := ppm.Size()
	if w != 3 || h != 2 {
		t.Error("Wrong size")
	}
}

func TestAtPPM(t *testing.T) {
	ppm := &PPM{
		data:        [][]Pixel{{{255, 0, 0}, {0, 255, 0}, {0, 0, 255}}, {{255, 255, 0}, {0, 255, 255}, {0, 0, 255}}},
		width:       3,
		height:      2,
		magicNumber: "P3",
		max:         255,
	}
	pixelValue := ppm.At(0, 1)

	if pixelValue.R != 0 || pixelValue.G != 255 || pixelValue.B != 0 {
		t.Error("Wrong value")
	}
}

func TestInvertPPM(t *testing.T) {
	ppm := &PPM{
		data:        [][]Pixel{{{255, 0, 0}, {0, 255, 0}, {0, 0, 255}}, {{255, 255, 0}, {0, 255, 255}, {0, 0, 255}}},
		width:       3,
		height:      2,
		magicNumber: "P3",
		max:         255,
	}
	ppm.Invert()

	expectedData := [][]Pixel{
		{{0, 255, 255}, {255, 0, 255}, {255, 255, 0}},
		{{0, 0, 255}, {255, 0, 0}, {255, 255, 0}},
	}
	if !reflect.DeepEqual(ppm.data, expectedData) {
		t.Error("Wrong data after inverting")
	}
}

func TestFlipPPM(t *testing.T) {
	ppm := &PPM{
		data:        [][]Pixel{{{255, 0, 0}, {0, 255, 0}, {0, 0, 255}}, {{255, 255, 0}, {0, 255, 255}, {0, 0, 255}}},
		width:       3,
		height:      2,
		magicNumber: "P3",
		max:         255,
	}
	ppm.Flip()

	expectedData := [][]Pixel{
		{{0, 0, 255}, {0, 255, 0}, {255, 0, 0}},
		{{0, 0, 255}, {255, 255, 0}, {255, 0, 0}},
	}
	if !reflect.DeepEqual(ppm.data, expectedData) {
		t.Error("Wrong data after flipping")
	}
}

func TestFlopPPM(t *testing.T) {
	ppm := &PPM{
		data:        [][]Pixel{{{255, 0, 0}, {0, 255, 0}, {0, 0, 255}}, {{255, 255, 0}, {0, 255, 255}, {0, 0, 255}}},
		width:       3,
		height:      2,
		magicNumber: "P3",
		max:         255,
	}
	ppm.Flop()

	expectedData := [][]Pixel{
		{{255, 0, 0}, {0, 255, 0}, {0, 0, 255}},
		{{255, 0, 0}, {255, 255, 0}, {0, 0, 255}},
	}
	if !reflect.DeepEqual(ppm.data, expectedData) {
		t.Error("Wrong data after flopping")
	}
}

func TestRotate90CWPPM(t *testing.T) {
	ppm := &PPM{
		data:        [][]Pixel{{{255, 0, 0}, {0, 255, 0}, {0, 0, 255}}, {{255, 255, 0}, {0, 255, 255}, {0, 0, 255}}},
		width:       3,
		height:      2,
		magicNumber: "P3",
		max:         255,
	}
	ppm.Rotate90CW()

	expectedData := [][]Pixel{
		{{0, 255, 0}, {0, 0, 255}},
		{{255, 255, 0}, {0, 255, 0}},
		{{0, 0, 255}, {255, 0, 0}},
	}
	if !reflect.DeepEqual(ppm.data, expectedData) {
		t.Error("Wrong data after rotating 90 degrees clockwise")
	}
}

func TestToPGMPPM(t *testing.T) {
	ppm := &PPM{
		data:        [][]Pixel{{{255, 0, 0}, {0, 255, 0}, {0, 0, 255}}, {{255, 255, 0}, {0, 255, 255}, {0, 0, 255}}},
		width:       3,
		height:      2,
		magicNumber: "P3",
		max:         255,
	}
	pgm := ppm.ToPGM()

	expectedData := [][]uint8{
		{85, 85, 85},
		{170, 170, 85},
		{85, 170, 255},
	}
	if !reflect.DeepEqual(pgm.data, expectedData) {
		t.Error("Wrong data after converting to PGM")
	}
}

func TestToPBMPPM(t *testing.T) {
	ppm := &PPM{
		data:        [][]Pixel{{{255, 0, 0}, {0, 255, 0}, {0, 0, 255}}, {{255, 255, 0}, {0, 255, 255}, {0, 0, 255}}},
		width:       3,
		height:      2,
		magicNumber: "P3",
		max:         255,
	}
	pbm := ppm.ToPBM()

	expectedData := [][]bool{
		{true, true, true},
		{true, true, true},
		{true, true, true},
	}
	if !reflect.DeepEqual(pbm.data, expectedData) {
		t.Error("Wrong data after converting to PBM")
	}
}

//Je n'ai pas pu finir le ppm_test.go
