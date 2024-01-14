package main

import (
	"Netpbm/images/pbmImages"
	"Netpbm/images/pgmImages"
	"Netpbm/images/ppmImages"
	"fmt"
)

func main() {

	pbmImage, err := pbmImages.ReadPBM("images/pbmImages/example.pbm")
	if err != nil {
		fmt.Println("Erreur lors de la lecture de l'image PBM:", err)
		return
	}

	fmt.Println("Taille de l'image PBM:", pbmImage.Size())
	fmt.Println("Valeur du pixel (0,0) dans l'image PBM:", pbmImage.At(0, 0))

	pbmImage.Invert()
	fmt.Println("Image PBM inversée:", pbmImage)

	err = pbmImage.Save("images/pbmImages/inverted_example.pbm")
	if err != nil {
		fmt.Println("Erreur lors de l'enregistrement de l'image PBM inversée:", err)
		return
	}

	pgmImage, err := pgmImages.ReadPGM("images/pgmImages/example.pgm")
	if err != nil {
		fmt.Println("Erreur lors de la lecture de l'image PGM:", err)
		return
	}

	fmt.Println("Taille de l'image PGM:", pgmImage.Size())
	fmt.Println("Valeur du pixel (0,0) dans l'image PGM:", pgmImage.At(0, 0))

	pgmImage.Invert()
	fmt.Println("Image PGM inversée:", pgmImage)

	err = pgmImage.Save("images/pgmImages/inverted_example.pgm")
	if err != nil {
		fmt.Println("Erreur lors de l'enregistrement de l'image PGM inversée:", err)
		return
	}

	ppmImage, err := ppmImages.ReadPPM("images/ppmImages/example.ppm")
	if err != nil {
		fmt.Println("Erreur lors de la lecture de l'image PPM:", err)
		return
	}

	fmt.Println("Taille de l'image PPM:", ppmImage.Size())
	fmt.Println("Valeur du pixel (0,0) dans l'image PPM:", ppmImage.At(0, 0))

	ppmImage.Invert()
	fmt.Println("Image PPM inversée:", ppmImage)

	err = ppmImage.Save("images/ppmImages/inverted_example.ppm")
	if err != nil {
		fmt.Println("Erreur lors de l'enregistrement de l'image PPM inversée:", err)
		return
	}
}
