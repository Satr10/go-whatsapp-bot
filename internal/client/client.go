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

func RandomMeme() (gambar []byte, memeInfo string, err error) {
	var memeData *model.OneCakPost
	resp, err := http.Get("https://1cak-api.vercel.app/api/random")
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error baca responya")
		return nil, "", err
	}

	err = json.Unmarshal(body, &memeData)
	if err != nil {
		return nil, "", err
	}
	fmt.Println(memeData)

	req, err := http.NewRequest("GET", memeData.Meme.ImageURL, nil)
	if err != nil {
		fmt.Println("Error membuat permintaan:", err)
		return
	}

	req.Header.Set("Referer", memeData.Meme.ImageURL)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:137.0) Gecko/20100101 Firefox/137.0")

	client := &http.Client{}
	respGambar, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer respGambar.Body.Close()
	// Check for successful response
	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	gambar, err = io.ReadAll(respGambar.Body)
	if err != nil {
		return nil, "", err
	}

	memeInfo = fmt.Sprintf("Author: %v\nLink: %v", memeData.Meme.Uploader, memeData.Meme.PostURL)

	return gambar, memeInfo, err
}

// TODO:TAMBAH MEME PAKAI API https://github.com/D3vd/Meme_Api
