package cmd

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)


var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "This command will get the desired Gopher",
	Long:  `This get command will call the GitHub repository to retrieve the desired Gopher.`,
	Run: func(cmd *cobra.Command, args []string) {
		gopherName := "dr-who.png"

		if len(args) >= 1 && args[0] != "" {
			gopherName = args[0]
		}

		URL := fmt.Sprintf("https://github.com/scraly/gophers/raw/main/%s.png", gopherName)

		fmt.Printf("Trying to get '%s' Gopher...\n", gopherName)

		response, err := http.Get(URL)
		if err != nil {
			log.Fatal(err)
		}
		defer response.Body.Close()

		if response.StatusCode == http.StatusOK {
			imageDir := "images"
			if err := os.MkdirAll(imageDir, 0755); err != nil {
				log.Fatal(err)
			}

			entries, err := os.ReadDir(imageDir)
			if err != nil {
				log.Fatal(err)
			}

			for _, entry := range entries {
				fileName := entry.Name()

				if fileName == gopherName+".png" {
					gopherName = gopherName + randSeq(6)
				}
			}

			outPath := filepath.Join(imageDir, gopherName+".png")
			out, err := os.Create(outPath)
			if err != nil {
				log.Fatal(err)
			}
			defer out.Close()

			_, err = io.Copy(out, response.Body)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Saved as '%s'!\n", out.Name())
		} else {
			fmt.Printf("Error: '%s' does not exist!\n", gopherName)
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	rand.Seed(time.Now().UnixNano())
}