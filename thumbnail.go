package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ovhlabs/functions/go-sdk/event"

	"github.com/disintegration/gift"
)

// Thumb est une "OVH function" qui va retourner une miniature
// de l'image qui lui est transmise spar son URL
// Formats d'entré: jpeg, png ou gif
// Format sortie: png de 100 pixels de large
func Thumb(event event.Event) (string, error) {
	// debug
	fmt.Println(event)

	// on récupère l'URL de l'image à traiter passer en paramétre
	// que l'on devrait donc trouver dans event.Params
	picURL, ok := event.Params["pic"]

	// si ce parametre est manquant ou si c'est une chaine vide
	// on retourne une erreur
	if !ok || strings.TrimSpace(picURL) == "" {
		return "", errors.New("parameter pic is missing")
	}

	// on pourrait vérifier que le paramètre picURL est bien une URL
	// mais le http.Get du dessous va le faire et générera un erreur
	// si ce n'est pas le cas

	// on va chercher l'image
	resp, err := http.Get(picURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// on instancie une nouvelle stucture Image
	imgIn, _, err := image.Decode(resp.Body)
	if err != nil {
		return "", err
	}

	// on instancie gift et on lui passe les modifications à faire
	g := gift.New(
		// on redimensione l'image pour que la largeur soit de 100px
		// en utilisant l'algorithme Lanczos
		gift.Resize(100, 0, gift.LanczosResampling),
	)

	// on crée un nouvelle image qui recevra l'image redimensionée
	imgOut := image.NewRGBA(g.Bounds(imgIn.Bounds()))

	// Draw va prendre l'image de départ (img), va appliquer les filtres
	// et copier le résultat dans
	g.Draw(imgOut, imgIn)

	// on encode l'image en png dans un buffer
	imgOutByte := bytes.NewBuffer([]byte{})
	err = png.Encode(imgOutByte, imgOut)
	if err != nil {
		return "", err
	}

	// on recupere le contenu du buffer dans un variable
	out, err := ioutil.ReadAll(imgOutByte)
	if err != nil {
		return "", err
	}

	// on retourne notre slice de bytes sous forme d'une string
	// puisque c'est ce qu'attend le gestionnaire
	return string(out), nil
}
