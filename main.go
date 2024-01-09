package main

import (
	"fmt"
	"log"

	"github.com/<username>/Netpbm/pbm"
	"github.com/<username>/Netpbm/pgm"
	"github.com/<username>/Netpbm/ppm"
)

func main() {
	pbmFilePath := "images/mon_image.pbm"

	pbmImage, err := pbm.ReadPBM(pbmFilePath)
	if err != nil {
		log.Fatal("Erreur lors de la lecture de l'image PBM :", err)
	}

	pbmImage.Invert()

	err = pbmImage.Save("images/mon_image_inversee.pbm")
	if err != nil {
		log.Fatal("Erreur lors de l'enregistrement de l'image PBM :", err)
	}

	pgmFilePath := "images/mon_image.pgm"

	pgmImage, err := pgm.ReadPGM(pgmFilePath)
	if err != nil {
		log.Fatal("Erreur lors de la lecture de l'image PGM :", err)
	}

	pgmImage.Invert()

	err = pgmImage.Save("images/mon_image_inversee.pgm")
	if err != nil {
		log.Fatal("Erreur lors de l'enregistrement de l'image PGM :", err)
	}

	ppmFilePath := "images/mon_image.ppm"

	ppmImage, err := ppm.ReadPPM(ppmFilePath)
	if err != nil {
		log.Fatal("Erreur lors de la lecture de l'image PPM :", err)
	}

	ppmImage.Invert()

	err = ppmImage.Save("images/mon_image_inversee.ppm")
	if err != nil {
		log.Fatal("Erreur lors de l'enregistrement de l'image PPM :", err)
	}

	fmt.Println("Opérations terminées avec succès.")
}
