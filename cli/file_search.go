package cli

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	helpers "github.com/irevenko/what-anime-cli/helpers"
	types "github.com/irevenko/what-anime-cli/types"
	"github.com/muesli/termenv"
)

const (
	fileSearchURL = "https://trace.moe/api/search"
)

// SearchByImageFile is for finding the anime scene by existing image file
func SearchByImageFile(imagePath string) {
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		if err != nil {
			log.Fatal("Invalid file path")
		}
	}

	termenv.HideCursor()
	defer termenv.ShowCursor()

	s := spinner.New(spinner.CharSets[33], 100*time.Millisecond)
	s.Prefix = "🔎 Searching for the anime: "
	s.FinalMSG = color.GreenString("✔️  Found!\n")

	go catchInterrupt(s)

	s.Start()

	imageFile, err := os.Open(imagePath)
	helpers.HandleError(err)

	reader := bufio.NewReader(imageFile)
	content, err := ioutil.ReadAll(reader)
	helpers.HandleError(err)

	encodedImage := base64.StdEncoding.EncodeToString(content)

	reqBody, err := json.Marshal(map[string]string{"image": encodedImage})
	helpers.HandleError(err)

	resp, err := http.Post(fileSearchURL, "application/json", bytes.NewBuffer(reqBody))
	helpers.HandleError(err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	helpers.HandleError(err)

	var animeResp types.Response
	json.Unmarshal(body, &animeResp)

	s.Stop()

	fmt.Println("🌸 Title Native: " + animeResp.Docs[0].TitleNative)
	fmt.Println("🗻 Title Romaji: " + animeResp.Docs[0].TitleRomanji)
	fmt.Println("🗽 Title English: " + animeResp.Docs[0].TitleEnglish)
	fmt.Print("📊 Similarity: ")
	helpers.PrintAnimeSimilarity(strconv.FormatFloat(animeResp.Docs[0].Similarity, 'f', 6, 64))
	fmt.Println("📺 Episode Number: " + color.MagentaString(strconv.Itoa(animeResp.Docs[0].Episode)))
	fmt.Print("⌚ Scene At: ")
	helpers.PrintSceneAt(animeResp.Docs[0].At)
	fmt.Println("📅 Year & Season: " + color.CyanString(animeResp.Docs[0].Season))
	fmt.Print("🍓 Is Adult: ")
	helpers.PrintIsAdult(animeResp.Docs[0].IsAdult)
	//fmt.Println(string(body))
}
