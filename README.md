# socialbot

![Socialbot in action.](https://i.imgur.com/x9I1cnP.png)

## Installation
`git clone github.com/808hype/socialbot.git`

## Prerequisites

[Discord bot token](https://github.com/reactiflux/discord-irc/wiki/Creating-a-discord-bot-&-getting-a-token)

[Google API key](https://developers.google.com/youtube/registering_an_application)

The [YouTube channel ID](https://commentpicker.com/youtube-channel-id.php) can easily be obtained from the URL. If the user has a custom URL, click one of his/her videos, then copy the link address of the channel name and get the ID from there.

## Usage

Assuming you already have a Go development environment set up, these are the only steps you need to take.

First, enter your Discord bot token and Google API key into the config.json file.

Next, run the program. That's it.

`go run main.go`

Don't run the bot on your home computer 24/7, that's just ridiculous. Put it on Heroku if you want to run it 24/7.

At the moment, the only command is `!youtube <channel-id>`, which fetches basic information about the specified YouTube channel.

In the future, I may implement Twitter and Instagram.
