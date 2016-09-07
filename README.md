# RamenBot
RamenBot utilizes the [discordgo](https://github.com/bwmarrin/discordgo) library, a free and open source library. RamenBot requires Go 1.4 or higher.

## Changes to original bot
- Added sound-bytes from other forks and some clips that I added
- Support for tracking time spent 'playing' games in Discord
- Added basic sqlite database support (Didn't want to deal with Redis...)
- Added Markov chain type chatbot interactions
- Reorganized code into multiple files

## Usage
RamenBot Bot has two components, a bot client that handles a plethora of cool features, and a web server that implements OAuth2 and stats. Once added to your server, RamenBot can be summoned by running `!!*root_command* *function* *arguments*`.  Full documentation on root commands and their functions to come!


### Running the Bot

**First install the bot:**
```
go get github.com/t11230/ramenbot/cmd/bot
go install github.com/t11230/ramenbot/cmd/bot
```
 **Then run the setup script:**

 ```
./setup
 ```
 You will be prompted to enter your Bot Account Token, location to run MongoDB, and to specify which modules you want enabled.

 **Then run the following command:**

```
bot
```

### Running the Web Server
First install the webserver: `go install github.com/t11230/ramenbot`, then run `make static`, finally run:

```
./airhornweb -r "localhost:6379" -i MY_APPLICATION_ID -s 'MY_APPLICATION_SECRET"
```

Note, the webserver requires a redis instance to track statistics

## Thanks
Thanks to the awesome (one might describe them as smart... loyal... appreciative...) [iopred](https://github.com/iopred) and [bwmarrin](https://github.com/bwmarrin/discordgo) for helping code review the initial release.
