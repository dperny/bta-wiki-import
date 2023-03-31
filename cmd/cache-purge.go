package cmd 

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"bufio"
    "fmt"
    "net/http"
    "net/url"
    "strings"
    "time"
)

var purgeCmd = &cobra.Command{
	Use:   "cache-purge",
	Short: "pull the list of current mechs and run a cache purge against them",
	Run: func(cmd *cobra.Command, args []string) {
		    // URL for the text file containing page titles
			fileURL := "https://bta3062.com/files/list_of_mechs.txt"
			// API URL for the MediaWiki API
			apiURL := "https://bta3062.com/api.php"
		
			// Fetching the text file from the URL
			response, err := http.Get(fileURL)
			if err != nil {
				fmt.Println("Failed to fetch the text file:", err)
				return
			}
			defer response.Body.Close()
		
			// Splitting the text file into lines and iterating over each line
			scanner := bufio.NewScanner(response.Body)
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if line == "" {
					continue
				}
		
				params := url.Values{}
				params.Set("action", "purge")
				params.Set("titles", line)
				params.Set("format", "json")
		
				req, err := http.NewRequest("POST", apiURL, strings.NewReader(params.Encode()))
				if err != nil {
					fmt.Println("Failed to create the POST request:", err)
					return
				}
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		
				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					fmt.Println("Failed to send the POST request:", err)
					return
				}
				defer resp.Body.Close()
		
				// Printing the response body
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					fmt.Println("Failed to read the response body:", err)
					return
				}
				fmt.Println(string(body))
		
				// Sleeping for 10 seconds
				time.Sleep(10 * time.Second)
			}
		
			if err := scanner.Err(); err != nil {
				fmt.Println("Failed to read the text file:", err)
				return
			}
	},
}