package apiclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Satr10/go-whatsapp-bot/internal/model"
)

func GetQuote() (quote *model.Quote, err error) {
	rn := time.Now()
	resp, err := http.Get("https://qapi.vercel.app/api/random")
	if err != nil {
		fmt.Println("Error ngambil api nya Bang")
		return nil, err
	}
	defer resp.Body.Close()

	fmt.Println(time.Since(rn))

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error baca responya")
		return nil, err
	}

	err = json.Unmarshal(body, &quote)
	if err != nil {
		println("Error Unmarshal")
		return nil, err
	}

	return quote, err

}

func RandomFufufafa() (quote *model.FufufafaQuote, err error) {

	resp, err := http.Get("https://fufufafa-api.vercel.app/api/random")
	if err != nil {
		fmt.Println("Error ngambil api nya Bang")
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error baca responya")
		return nil, err
	}

	err = json.Unmarshal(body, &quote)
	if err != nil {
		println("Error Unmarshal")
		return nil, err
	}

	return quote, err

}
