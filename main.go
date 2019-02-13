package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/labstack/echo"
	plivo "github.com/plivo/plivo-go"
)

type message struct {
	ID          int    `json:"id"`
	Sender      string `json:"sender"`
	MessageTime string `json:"messageTime"`
	Text        string `json:"text"`
	Receiver    string `json:"receiver"`
}

// Handlers
func replyToMessage(c echo.Context) error {
	plivoAuthId := os.Getenv("PLIVO_AUTH_ID")
	plivoAuthToken := os.Getenv("PLIVO_AUTH_TOKEN")
	client, err := plivo.NewClient(plivoAuthId, plivoAuthToken, &plivo.ClientOptions{})
	if err != nil {
		panic(err)
	}

	response, err := client.Messages.Create(plivo.MessageCreateParams{
		Src:  "+14155696002",
		Dst:  "+17025305234",
		Text: "Toshitext just sent some crypto 😎!",
	})

	if err != nil {
		panic(err)
	}

	fmt.Println("Response: %#v\n", response)
	return c.String(http.StatusOK, "We have successfully sent the message!")
}

func receiveMessage(c echo.Context) error {
	return c.String(http.StatusOK, "The Toshitext service is listening to your commands.")
}

func getHelp(c echo.Context) error {
	return c.String(http.StatusOK, "Here's the list of available commands:...")
}

func renderHome(c echo.Context) error {
	return c.String(http.StatusOK, "Toshitext v1 - Send crypto with a text message.")
}

func getBalance(c echo.Context) error {
	return c.String(http.StatusOK, "Here's your wallet total balance: $100")
}

func createWallet(c echo.Context) error {
	return c.String(http.StatusOK, "Your Toshitext wallet has been created. Here's you wallet address and current balance.")
}

func getWalletAddress(c echo.Context) error {
	return c.String(http.StatusOK, "Here's your wallet deposit address: 0x50Ebe9ad50DCf1Be1A35570E29587fa9F6eCDB46")
}

func main() {
	e := echo.New()
	port := os.Getenv("process.env.PORT")
	if port == "" {
		e.Logger.Fatal("$PORT must be set")
	}

	// Set up Echo, configure server side validation, and hook into middleware.
	e.Server.Addr = ":" + port

	// Taking in JSON and mapping it to our data structure.
	jsn := `[{"ID":"id","Sender":"sender","MessageTime":"messageTime","Receiver":"receiver","Text":"text"}]`

	// MDR -> Message Detail Records
	details := []message{}
	fmt.Printf("Go data: %+v\n", details)

	err := json.Unmarshal([]byte(jsn), &details)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Go data: %+v\n", details)

	for i, v := range details {
		fmt.Println(i, v)
	}

	// Routes
	e.GET("/", renderHome)
	e.POST("/help", getHelp)
	e.POST("/messages", receiveMessage)
	e.GET("/wallets", getWalletAddress)
	e.POST("/wallets", createWallet)
	e.GET("/users", getBalance)

	// Gracefully shut down the server on interrupt.
	e.Logger.Fatal(gracehttp.Serve(e.Server))
}
