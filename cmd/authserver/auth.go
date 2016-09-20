package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/antage/eventsource"
	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/handlers"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var (
	// Permission Constants
	ADMINISTRATOR    = 0x00000008
	MANAGE_GUILD     = 0x00000020
	READ_MESSAGES    = 0x00000400
	SEND_MESSAGES    = 0x00000800
	CONNECT          = 0x00100000
	SPEAK            = 0x00200000
	CHANGE_NICKNAME  = 0x04000000
	MANAGE_NICKNAMES = 0x08000000
	MANAGE_ROLES     = 0x10000000

	// Oauth2 Config
	oauthConf *oauth2.Config

	// Used for storing session information in a cookie
	store *sessions.CookieStore

	// Used for pushing live stat updates to the client
	es eventsource.EventSource

	// Base URL of the discord API
	apiBaseUrl = "https://discordapp.com/api"
)

// Represents a JSON struct of stats that are updated every second and pushed to the client
type CountUpdate struct {
	Total          string `json:"total"`
	UniqueUsers    string `json:"unique_users"`
	UniqueGuilds   string `json:"unique_guilds"`
	UniqueChannels string `json:"unique_channels"`
	SecretCount    string `json:"secret_count"`
}

func (c *CountUpdate) ToJSON() []byte {
	data, _ := json.Marshal(c)
	return data
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// Return a random character sequence of n length
func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// Returns the current session or aborts the request
func getSessionOrAbort(w http.ResponseWriter, r *http.Request) *sessions.Session {
	session, err := store.Get(r, "session")

	if session == nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to get session")
		http.Error(w, "Invalid or corrupted session", http.StatusInternalServerError)
		return nil
	}

	return session
}

// Redirects to the oauth2
func handleLogin(w http.ResponseWriter, r *http.Request) {
	session := getSessionOrAbort(w, r)
	if session == nil {
		return
	}

	// Create a random state
	session.Values["state"] = randSeq(32)
	session.Save(r, w)

	// OR the permissions we want
	perms := READ_MESSAGES | SEND_MESSAGES | CONNECT | SPEAK | CHANGE_NICKNAME | MANAGE_NICKNAMES | MANAGE_ROLES

	// Return a redirect to the ouath provider
	url := oauthConf.AuthCodeURL(session.Values["state"].(string), oauth2.AccessTypeOnline)
	http.Redirect(w, r, url+fmt.Sprintf("&permissions=%v", perms), http.StatusTemporaryRedirect)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	session := getSessionOrAbort(w, r)
	if session == nil {
		return
	}

	// Check the state string is correct
	state := r.FormValue("state")
	if state != session.Values["state"] {
		log.WithFields(log.Fields{
			"expected": session.Values["state"],
			"received": state,
		}).Error("Invalid OAuth state")
		http.Redirect(w, r, "/?key_to_success=0", http.StatusTemporaryRedirect)
		return
	}

	errorMsg := r.FormValue("error")
	if errorMsg != "" {
		log.WithFields(log.Fields{
			"error": errorMsg,
		}).Error("Received OAuth error from provider")
		http.Redirect(w, r, "/?key_to_success=0", http.StatusTemporaryRedirect)
		return
	}

	token, err := oauthConf.Exchange(oauth2.NoContext, r.FormValue("code"))
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"token": token,
		}).Error("Failed to exchange token with provider")
		http.Redirect(w, r, "/?key_to_success=0", http.StatusTemporaryRedirect)
		return
	}

	body, _ := json.Marshal(map[interface{}]interface{}{})
	req, err := http.NewRequest("GET", apiBaseUrl+"/users/@me", bytes.NewBuffer(body))
	if err != nil {
		log.WithFields(log.Fields{
			"body":  body,
			"req":   req,
			"error": err,
		}).Error("Failed to create @me request")
		http.Error(w, "Failed to retrieve user profile", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Authorization", token.Type()+" "+token.AccessToken)
	client := &http.Client{Timeout: (20 * time.Second)}
	resp, err := client.Do(req)
	if err != nil {
		log.WithFields(log.Fields{
			"error":  err,
			"client": client,
			"resp":   resp,
		}).Error("Failed to request @me data")
		http.Error(w, "Failed to retrieve user profile", http.StatusInternalServerError)
		return
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"body":  resp.Body,
		}).Error("Failed to read data from HTTP response")
		http.Error(w, "Failed to retrieve user profile", http.StatusInternalServerError)
		return
	}

	user := discordgo.User{}
	err = json.Unmarshal(respBody, &user)
	if err != nil {
		log.WithFields(log.Fields{
			"data":  respBody,
			"error": err,
		}).Error("Failed to parse JSON payload from HTTP response")
		http.Error(w, "Failed to retrieve user profile", http.StatusInternalServerError)
		return
	}

	// Finally write some information to the session store
	session.Values["token"] = token.AccessToken
	session.Values["username"] = user.Username
	session.Values["tag"] = user.Discriminator
	delete(session.Values, "state")
	session.Save(r, w)

	// And redirect the user back to the dashboard
	http.Redirect(w, r, "/?key_to_success=1", http.StatusTemporaryRedirect)
}

func handleMe(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")

	body, err := json.Marshal(map[string]interface{}{
		"username": session.Values["username"],
		"tag":      session.Values["tag"],
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func server() {
	server := http.NewServeMux()
	server.Handle("/", http.FileServer(http.Dir("static/dist")))
	server.HandleFunc("/me", handleMe)
	server.HandleFunc("/login", handleLogin)
	server.HandleFunc("/callback", handleCallback)

	port := os.Getenv("PORT")
	if port == "" {
		port = "14000"
	}

	log.WithFields(log.Fields{
		"port": port,
	}).Info("Starting HTTP Server")

	// If the requests log doesn't exist, make it
	if _, err := os.Stat("requests.log"); os.IsNotExist(err) {
		ioutil.WriteFile("requests.log", []byte{}, 0600)
	}

	// Open the log file in append mode
	logFile, err := os.OpenFile("requests.log", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to open requests log file")
		return
	}
	defer logFile.Close()

	// Actually start the server
	loggedRouter := handlers.LoggingHandler(logFile, server)
	http.ListenAndServe(":"+port, loggedRouter)
}

func main() {
	var (
		ClientID     = flag.String("i", "", "OAuth2 Client ID")
		ClientSecret = flag.String("s", "", "OAtuh2 Client Secret")
	)
	flag.Parse()

	// Create a cookie store
	store = sessions.NewCookieStore([]byte(*ClientSecret))

	// Setup the OAuth2 Configuration
	endpoint := oauth2.Endpoint{
		AuthURL:  apiBaseUrl + "/oauth2/authorize",
		TokenURL: apiBaseUrl + "/oauth2/token",
	}

	oauthConf = &oauth2.Config{
		ClientID:     *ClientID,
		ClientSecret: *ClientSecret,
		Scopes:       []string{"bot", "identify"},
		Endpoint:     endpoint,
		RedirectURL:  "https://airhornbot.com/callback",
	}

	server()
}
