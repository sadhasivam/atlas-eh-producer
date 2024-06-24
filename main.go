package main

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/sadhasivam/atlas-eh-producer/schema"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
}

func main() {

	token, err := GetAuthToken()
	if err != nil {
		panic(err)
	}
	messages := CreateRFIDMessage()
	err = PostRFIDScans(token.AccessToken, messages)
	if err != nil {
		panic(err)
	}
}

func PostRFIDScans(token string, payload schema.AtlasPayload[schema.RfidHubScan]) error {
	ingestionUrl := viper.GetString("ATLAS_INGESTION_URL")

	// Marshal the payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %w", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", ingestionUrl, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %w", err)
	}

	// Set the authorization header
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("RFID scans posted successfully \n " + string(body))
	return nil
}

func CreateRFIDMessage() schema.AtlasPayload[schema.RfidHubScan] {
	rfidHubScan := schema.RfidHubScan{
		TagHexEpc:      "9222CE9CF38BFC80463E7023",
		Timestamp:      342343242,
		TrackingNumber: 134134134,
		IsFedexTag:     true,
		Center:         "RAJU",
		ExtraData:      map[string]interface{}{},
	}
	rfidMessage := schema.AtlasMessage[schema.RfidHubScan]{
		IsJson:      true,
		JsonMessage: rfidHubScan,
	}

	return schema.AtlasPayload[schema.RfidHubScan]{
		ExchangeId:      generateUniqueID(),
		PublishTimeInMs: getCurrentTimeInMs(),
		ReceiveTimeInMs: getCurrentTimeInMs() - 100,
		Messages:        []schema.AtlasMessage[schema.RfidHubScan]{rfidMessage},
	}
}

func GetAuthToken() (*schema.OktaToken, error) {
	method := "POST"
	oktaUrl := viper.GetString("OKTA_URL")
	clientId := viper.GetString("CLIENT_ID")
	clientSecret := viper.GetString("CLIENT_SECRET")

	payload := bytes.NewBufferString("grant_type=client_credentials&scope=Custom_Scope")

	client := &http.Client{}
	req, err := http.NewRequest(method, oktaUrl, payload)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Basic "+encodeBase64(clientId, clientSecret))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		return nil, err

	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		var errorResp schema.OktaError
		if err := json.Unmarshal(body, &errorResp); err != nil {
			return nil, fmt.Errorf("error unmarshaling error response: %w", err)
		}
		return nil, &errorResp
	} else {
		var oktaToken schema.OktaToken
		if err := json.Unmarshal(body, &oktaToken); err != nil {
			return nil, fmt.Errorf("error unmarshaling success response: %w", err)
		}
		return &oktaToken, nil
	}
}

func encodeBase64(clientId, clientSecret string) string {
	concatenated := clientId + ":" + clientSecret
	encoded := base64.StdEncoding.EncodeToString([]byte(concatenated))
	return encoded
}

func generateUniqueID() int {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return rand.Intn(10000)
}

func getCurrentTimeInMs() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
