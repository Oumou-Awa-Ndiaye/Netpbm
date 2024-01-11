package main

import (
	"Netpbm/pbm"
	"Netpbm/pgm"
	"Netpbm/ppm"
	"fmt"
	"log"
)

func main() {

	pbmFile, err := pbm.ReadPBM("pbm/image.pbm")
	if err != nil {
		log.Fatal(err)
	}

	pbmFile.Invert()
	pbmFile.Flip()

	err = pbmFile.Save("image_pbm_modifie.pbm")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Image PBM modifiée sauvegardée avec succès.")

	pgmFile, err := pgm.ReadPGM("pgm/image.pgm")
	if err != nil {
		log.Fatal(err)
	}

	pgmFile.Invert()
	pgmFile.Rotate90CW()

	err = pgmFile.Save("image_pgm_modifie.pgm")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Image PGM modifiée sauvegardée avec succès.")

	ppmFile, err := ppm.ReadPPM("ppm/image.ppm")
	if err != nil {
		log.Fatal(err)
	}

	ppmFile.Invert()
	ppmFile.Flip()

	err = ppmFile.Save("image_ppm_modifie.ppm")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Image PPM modifiée sauvegardée avec succès.")
}
