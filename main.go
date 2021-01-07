package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"google.golang.org/api/option"

	"github.com/bwmarrin/discordgo"
	"github.com/tidwall/gjson"
	youtube "google.golang.org/api/youtube/v3"
)

type youtubeChannel struct {
	URL            string
	Name           string
	Description    string
	ProfilePicture string
	Subscribers    string
	LatestVideo    string
	LatestVideoURL string
}

func main() {
	token, _, err := readConfig()
	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			log.Fatalf("Config file not found.")
		} else {
			log.Fatalf("Failed to read config file: %v", err)
		}
	}
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Failed to create bot session: %v", err)
	}
	dg.AddHandler(messageCreate)
	// Connects to Discord via websocket to listen for a message.
	err = dg.Open()
	if err != nil {
		log.Fatalf("Failed to connect to Discord: %v", err)
	}
	// The channel blocks until an OS signal is sent. Copied from the discordgo examples.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	dg.Close()
}

func readConfig() (string, string, error) {
	// We know the file isn't large so using ioutil.ReadFile should be fine.
	data, err := ioutil.ReadFile("./config.json")
	if err != nil {
		return "", "", err
	}
	token := gjson.Get(string(data), "token").String()
	apiKey := gjson.Get(string(data), "apiKey").String()
	return token, apiKey, nil
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Checks if the bot is the one who sent the message.
	if m.Author.ID == s.State.User.ID {
		return
	}
	if strings.HasPrefix(m.Content, "!youtube") {
		channelID := strings.Split(m.Content, " ")[1]
		channel, err := getChannel(channelID)
		embed := &discordgo.MessageEmbed{}
		if err != nil {
			embed = &discordgo.MessageEmbed{
				Title: "Failed to retrieve channel!",
				Color: 16711680,
			}
		} else {
			embed = &discordgo.MessageEmbed{
				URL:         channel.URL,
				Title:       channel.Name,
				Description: channel.Description,
				Color:       65280,
				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL: channel.ProfilePicture,
				},
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:  "Subscriber Count",
						Value: channel.Subscribers,
					},
					{
						Name:  "Latest Video",
						Value: fmt.Sprintf("[%s](%s)", channel.LatestVideo, channel.LatestVideoURL),
					},
				},
			}
		}
		_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
		if err != nil {
			return
		}
	}
}

func getChannel(id string) (*youtubeChannel, error) {
	// I can't think of a better way to implement this so if someone has a fix let me know.
	_, apiKey, err := readConfig()
	if err != nil {
		return nil, err
	}
	service, err := youtube.NewService(context.Background(), option.WithAPIKey(apiKey))
	if err != nil {
		fmt.Printf("Failed to create service. This could be because of an unauthorized API key - %v", err)
		return nil, err
	}
	// Utilizes the YouTube Data API Go client library. Will probably clean this up eventually.
	channelCall := service.Channels.List("snippet,statistics,contentDetails").Id(id)
	channelResp, err := channelCall.Do()
	if err != nil {
		fmt.Printf("Failed to fetch channel - %v", err)
		return nil, err
	}
	channelTitle := channelResp.Items[0].Snippet.Title
	description := channelResp.Items[0].Snippet.Description
	profilePic := channelResp.Items[0].Snippet.Thumbnails.Medium.Url
	// Converts a uint64, which is how the API returns the subscriber count, to a string.
	subscribers := strconv.FormatUint(channelResp.Items[0].Statistics.SubscriberCount, 10)
	uploads := channelResp.Items[0].ContentDetails.RelatedPlaylists.Uploads
	uploadCall := service.PlaylistItems.List("snippet,contentDetails").PlaylistId(uploads)
	uploadResp, err := uploadCall.Do()
	if err != nil {
		fmt.Printf("Failed to fetch latest video - %v", err)
		return nil, err
	}
	videoID := uploadResp.Items[0].ContentDetails.VideoId
	videoTitle := uploadResp.Items[0].Snippet.Title
	ytChannel := &youtubeChannel{
		URL:            "https://youtube.com/channel/" + id,
		Name:           channelTitle,
		Description:    description,
		ProfilePicture: profilePic,
		Subscribers:    subscribers,
		LatestVideo:    videoTitle,
		LatestVideoURL: "https://youtube.com/watch?v=" + videoID,
	}
	return ytChannel, nil
} 
