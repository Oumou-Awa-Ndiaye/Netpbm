package netpbm

import (
	"os"
	"reflect"
	"testing"
)

const pgmTestFileP2 = "./imagesTest/pgmImages/imageP2.pgm"

var testPGMDataP2 = [][]uint8{
	{100, 150, 200},
	{50, 75, 100},
}

const pgmTestFileP5 = "./imagesTest/pgmImages/imageP5.pgm"

var testPGMDataP5 = [][]uint8{
	{50, 53, 10, 51, 32, 50, 10, 50, 53, 53, 10, 49, 48, 48, 32, 49, 53, 48, 32, 50, 48, 48, 10, 53, 48, 32, 55, 53, 32, 49, 48, 48, 10},
	{53, 32, 55, 53, 32, 49, 48, 48, 10},
}

func TestReadPGM(t *testing.T) {
	/*I Create a temporary test file pour faciliter mon raisonnement*/
	testFilename := "./imagesTest/pgmImages/imageP1.pgm"
	file, err := os.Create(testFilename)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(testFilename)

	/*Ici j'écris les données P2 dans le fichier*/
	file.WriteString("P2\n")
	file.WriteString("3 2\n")
	file.WriteString("255\n")
	file.WriteString("100 150 200\n")
	file.WriteString("50 75 100\n")
	file.Close()

	pgm, err := ReadPGM(testFilename)
	if err != nil {
		t.Fatal(err)
	}

	if pgm.magicNumber != "P2" {
		t.Error("Wrong magic number")
	}

	if pgm.width != 3 || pgm.height != 2 {
		t.Error("Wrong size")
	}

	if pgm.max != 255 {
		t.Error("Wrong max value")
	}

	expectedData := [][]uint8{{100, 150, 200}, {50, 75, 100}}
	if !reflect.DeepEqual(pgm.data, expectedData) {
		t.Error("Wrong data for P2")
	}

}

func TestSizePGM(t *testing.T) {
	pgm := &PGM{data: [][]uint8{{100, 150, 200}, {50, 75, 100}}, width: 3, height: 2, magicNumber: "P2", max: 255}
	w, h := pgm.Size()
	if w != 3 || h != 2 {
		t.Error("Wrong size")
	}
}

func TestAtPGM(t *testing.T) {
	pgm := &PGM{data: [][]uint8{{100, 150, 200}, {50, 75, 100}}, width: 3, height: 2, magicNumber: "P2", max: 255}
	if pgm.At(0, 1) != 150 {
		t.Error("Wrong value")
	}
}

func TestInvertPGM(t *testing.T) {
	pgm := &PGM{data: [][]uint8{{100, 150, 200}, {50, 75, 100}}, width: 3, height: 2, magicNumber: "P2", max: 255}
	pgm.Invert()

	expectedData := [][]uint8{{155, 105, 55}, {205, 180, 155}}
	if !reflect.DeepEqual(pgm.data, expectedData) {
		t.Error("Wrong data after inverting") /*Cette ligne compare les deux matrices.
		Si les matrices sont différentes, la fonction renvoie true, ce qui signifie que les données après inversion ne correspondent pas aux données attendues*/
	}
}

func TestFlipPGM(t *testing.T) {
	pgm := &PGM{data: [][]uint8{{100, 150, 200}, {50, 75, 100}}, width: 3, height: 2, magicNumber: "P2", max: 255}
	pgm.Flip()

	expectedData := [][]uint8{{200, 150, 100}, {100, 75, 50}}
	if !reflect.DeepEqual(pgm.data, expectedData) {
		t.Error("Wrong data after flipping")
	}
}

func TestFlopPGM(t *testing.T) {
	pgm := &PGM{data: [][]uint8{{100, 150, 200}, {50, 75, 100}}, width: 3, height: 2, magicNumber: "P2", max: 255}
	pgm.Flop()

	expectedData := [][]uint8{{50, 75, 100}, {200, 150, 100}}
	if !reflect.DeepEqual(pgm.data, expectedData) {
		t.Error("Wrong data after flopping")
	}
}

func TestRotate90CWPGM(t *testing.T) {
	pgm := &PGM{data: [][]uint8{{100, 150, 200}, {50, 75, 100}}, width: 3, height: 2, magicNumber: "P2", max: 255}
	pgm.Rotate90CW()

	expectedData := [][]uint8{{50, 100}, {75, 150}, {100, 200}}
	if !reflect.DeepEqual(pgm.data, expectedData) {
		t.Error("Wrong data after rotating 90 degrees clockwise")
	}
}

func TestToPBM(t *testing.T) {
	pgm := &PGM{data: [][]uint8{{100, 150, 200}, {50, 75, 100}}, width: 3, height: 2, magicNumber: "P2", max: 255}
	pbm := pgm.ToPBM()

	expectedData := [][]bool{{false, false, true}, {false, true, true}}
	if !reflect.DeepEqual(pbm.data, expectedData) {
		t.Error("Wrong data after converting to PBM")
	}
}
