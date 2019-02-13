// Code attribution: https://github.com/miguelmota/ethereum-development-with-go-book/blob/master/code/generate_wallet.go
package main

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/labstack/echo"
	plivo "github.com/plivo/plivo-go"
	"golang.org/x/crypto/sha3"
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
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	fmt.Println(hexutil.Encode(privateKeyBytes)[2:]) // fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	fmt.Println(hexutil.Encode(publicKeyBytes)[4:]) // 9a7df67f79246283fdc93af76d4f8cdd62c4886e8cd870944e817dd0b97934fdd7719d0810951e03418205868a5c1b40b192451367f28e0088dd75e15de40c05

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Println(address) // 0x96216849c49358B10257cb55b28eA603c874b05E

	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicKeyBytes[1:])
	fmt.Println(hexutil.Encode(hash.Sum(nil)[12:])) // 0x96216849c49358b10257cb55b28ea603c874b05e

	return c.String(http.StatusOK, "Here's your wallet deposit address: 0x50Ebe9ad50DCf1Be1A35570E29587fa9F6eCDB46")
}

func getWalletAddress(c echo.Context) error {
	return c.String(http.StatusOK, "Here's your wallet deposit address: 0x50Ebe9ad50DCf1Be1A35570E29587fa9F6eCDB46")
}

func main() {
	e := echo.New()
	port := os.Getenv("PORT")
	if port == "" {
		e.Logger.Fatal("$PORT must be set")
	}

	// Set up Echo, configure server side validation, and hook into middleware.
	e.Server.Addr = ":" + port

	// Taking in mock JSON and mapping it to our data structure.
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
	e.POST("/messages", replyToMessage)
	e.GET("/wallets", getWalletAddress)
	e.POST("/wallets", createWallet)
	e.GET("/users", getBalance)

	// Gracefully shut down the server on interrupt.
	e.Logger.Fatal(gracehttp.Serve(e.Server))
}
