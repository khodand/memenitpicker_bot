package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"go.uber.org/zap"

	"github.com/khodand/memenitpicker_bot/cmd/config"
	"github.com/khodand/memenitpicker_bot/internal/meme"
	"github.com/khodand/memenitpicker_bot/internal/store"
	"github.com/khodand/memenitpicker_bot/pkg/postgres"
)

const prefixPath = "cmd/importer/ChatExport_2024-07-20/"

func main() {
	cfg := config.Config{
		Postgres: postgres.Config{
			Database:    "meme_bot",
			Username:    "meme_bot",
			Password:    "meme_bot",
			HostPrimary: "localhost",
			Port:        "5432",
			SSLMode:     "disable",
		},
	}
	database, err := postgres.NewPgxPool(cfg.Postgres)
	if err != nil {
		log.Fatal("failed to init database", zap.Error(err))
	}
	memesDB := store.NewMemes(database)
	memeService := meme.New(memesDB)

	file, err := os.Open(prefixPath + "result.json")
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		log.Printf("Failed to read file: %s", err)
		return
	}

	var in tInput
	err = json.Unmarshal(byteValue, &in)
	if err != nil {
		log.Printf("Failed to parse JSON: %s", err)
		return
	}

	output := strings.Builder{}
	output.WriteString(htmlStart)
	processed := make(map[int]string)
	const chatID = -4228130256
	ctx := context.Background()
	for _, msg := range in.Messages {
		if msg.Photo == "" {
			continue
		}

		imageFile, err := os.Open(prefixPath + msg.Photo)
		if err != nil {
			log.Printf("Failed to open file: %s", err)
			return
		}

		dID, aID, pID, err := memeService.InsertReturnALl(ctx, msg.ID, chatID, imageFile)
		if err != nil {
			log.Printf("Failed to insert file: %s", err)
			return
		}
		if err = imageFile.Close(); err != nil {
			log.Printf("Failed to close file: %s", err)
			return
		}

		if (dID == aID) && (aID == pID) && (aID == msg.ID) {
			processed[msg.ID] = msg.Photo
			continue
		}

		if (dID == aID) && (aID == pID) && (aID != msg.ID) {
			processed[msg.ID] = msg.Photo
			if !comparePhoto(processed[dID], msg.Photo) {
				log.Printf("there is duplicate message id by ALL: %d != %d", dID, msg.ID)
				draw(&output, "ALL", processed[dID], msg.Photo)
			}
			continue
		}

		prefix := ""
		if aID != msg.ID && !comparePhoto(processed[aID], msg.Photo) {
			prefix += "A"
		}
		if dID != msg.ID && !comparePhoto(processed[dID], msg.Photo) {
			prefix += "D"
		}
		if pID != msg.ID && !comparePhoto(processed[pID], msg.Photo) {
			prefix += "P"
		}
		switch {
		case prefix == "AD" && !comparePhoto(processed[aID], processed[dID]):
			draw(&output, "A", processed[aID], msg.Photo)
			draw(&output, "D", processed[dID], msg.Photo)
		case prefix == "AP" && !comparePhoto(processed[aID], processed[pID]):
			draw(&output, "A", processed[aID], msg.Photo)
			draw(&output, "P", processed[pID], msg.Photo)
		case prefix == "DP" && !comparePhoto(processed[dID], processed[pID]):
			draw(&output, "D", processed[dID], msg.Photo)
			draw(&output, "P", processed[pID], msg.Photo)
		case prefix == "A":

		default:
			draw(&output, prefix, max(processed[aID], processed[dID], processed[pID]), msg.Photo)
		}

		log.Printf("duplicate %s", prefix)
		processed[msg.ID] = msg.Photo
	}

	output.WriteString(`</body></html>`)
	err = os.WriteFile("index.html", []byte(output.String()), 0o644)
	if err != nil {
		log.Printf("Output failed to write: %s\n", err)
		return
	}

	log.Println("Output successfully written to import.txt")
}

func draw(b *strings.Builder, prefix, path1, path2 string) {
	path1 = prefixPath + path1
	path2 = prefixPath + path2
	htmlSnippet := fmt.Sprintf(`
<div class="photo-pair">
	<span class="prefix">%s:</span>
	<img src="%s">
	<img src="%s">
</div>
`+"\n", prefix, path1, path2)
	b.WriteString(htmlSnippet)
}

func comparePhoto(path1, path2 string) bool {
	return getId(path1) == getId(path2)
}

func getId(photoPath string) string {
	parts := strings.Split(photoPath, "/")
	lastPart := parts[len(parts)-1]
	idParts := strings.Split(lastPart, "@")
	return idParts[0]
}

const htmlStart = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Photo Comparison</title>
    <style>
    body { 
        font-family: Arial, sans-serif; 
        margin: 20px; 
    }
    .photo-pair { 
        display: flex; /* Keeps images in line */
        align-items: center; /* Center images vertically */
        justify-content: center; /* Center images horizontally */
        margin-bottom: 20px;
        border: 1px solid #ccc;
        padding: 10px;
        background-color: #f9f9f9;
    }
    .photo-pair img { 
        max-width: 45%; /* Limits width to 45% of the container */
        max-height: 70vh; /* Limits height to 70% of the viewport height */
        height: auto; /* Keeps the aspect ratio intact */
        margin-right: 5%; /* Space between images */
        object-fit: contain; /* Ensures the image is resized to maintain its aspect ratio while fitting within the element’s box */
    }
    .photo-pair img:last-child {
        margin-right: 0; /* No margin on the right for the last image */
    }
    .prefix { 
        font-weight: bold; 
        margin-right: 10px; /* Space between the prefix and the images */
    }
</style>
</head>
<body>
<h1>Photo Comparison</h1>
`
