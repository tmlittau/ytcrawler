package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

const missingClientSecretsMessage = `
Please configure OAuth 2.0
`

type ChannelStats struct {
	ChannelID   string
	ChannelName string
	SubCount    int
	ViewCount   int
	VideoCount  int
}

func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(cacheFile, tok)
	}
	return config.Client(ctx, tok)
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("youtube-go-quickstart.json")), err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func handleError(err error, message string) {
	if message == "" {
		message = "Error making API call"
	}
	if err != nil {
		log.Fatalf(message+": %v", err.Error())
	}
}

func addSubsToList(response *youtube.SubscriptionListResponse, SubList []ChannelStats) []ChannelStats {

	for i := 0; i < len(response.Items); i++ {
		SubList = append(SubList, ChannelStats{ChannelID: response.Items[i].Snippet.ResourceId.ChannelId})
	}

	return SubList
}

func subListByChannelID(service *youtube.Service, part []string, channelID string) []ChannelStats {
	call := service.Subscriptions.List(part)
	call = call.ChannelId(channelID)
	response, err := call.Do()
	handleError(err, "")

	nTotal := response.PageInfo.TotalResults
	npPage := response.PageInfo.ResultsPerPage
	nPages := int(math.Ceil(float64(nTotal) / float64(npPage)))

	var SubList []ChannelStats

	SubList = addSubsToList(response, SubList)

	nextPageToken := response.NextPageToken

	for i := 0; i < nPages-2; i++ {
		call = call.PageToken(nextPageToken)
		response, err = call.Do()
		handleError(err, "")

		SubList = addSubsToList(response, SubList)
		nextPageToken = response.NextPageToken
	}

	call = call.PageToken(nextPageToken)
	response, err = call.Do()
	handleError(err, "")

	SubList = addSubsToList(response, SubList)

	return SubList
}

func getChannelStats(service *youtube.Service, part []string, container ChannelStats) ChannelStats {
	call := service.Channels.List(part)
	call = call.Id(container.ChannelID)

	response, err := call.Do()
	handleError(err, "")

	container.ChannelName = response.Items[0].Snippet.Title
	container.SubCount = int(response.Items[0].Statistics.SubscriberCount)
	container.ViewCount = int(response.Items[0].Statistics.ViewCount)
	container.VideoCount = int(response.Items[0].Statistics.VideoCount)

	return container
}

func main() {
	ctx := context.Background()

	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/youtube-go-quickstart.json
	config, err := google.ConfigFromJSON(b, youtube.YoutubeReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(ctx, config)
	service, err := youtube.New(client)

	handleError(err, "Error creating YouTube client")

	buf, _ := ioutil.ReadFile("channelID.txt")
	chID := string(buf)

	SubList := subListByChannelID(service, []string{"snippet", "contentDetails"}, chID)

	for i := 0; i < len(SubList); i++ {
		SubList[i] = getChannelStats(service, []string{"snippet", "contentDetails", "statistics"}, SubList[i])
	}

	fmt.Println(SubList)
}
