package main

import (
	"Netpbm/pbm"
	"Netpbm/pgm"
	"Netpbm/ppm"
	"fmt"
)

func main() {
	pbmImage, err := pbm.ReadPBM("pbm/mon_image.pbm")
	if err != nil {
		fmt.Println("Erreur lors de la lecture de PBM :", err)
		return
	}

	pgmImage, err := pgm.ReadPGM("pgm/mon_image.pgm")
	if err != nil {
		fmt.Println("Erreur lors de la lecture de PGM :", err)
		return
	}

	ppmImage, err := ppm.ReadPPM("ppm/mon_image.ppm")
	if err != nil {
		fmt.Println("Erreur lors de la lecture de PPM :", err)
		return
	}

}
