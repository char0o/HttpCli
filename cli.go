package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const url = "http://127.0.0.1:8000/sales.json"
const username = "charo"
const password = "12345"

const (
	QUIT = iota
	CONNECT
)

type AuthClient struct {
	Username string
	Password string
}

func NewAuthClient(username, password string) *AuthClient {
	return &AuthClient{
		Username: username,
		Password: password,
	}
}
func (ac *AuthClient) Connect(path string) (*http.Response, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	req.SetBasicAuth(ac.Username, ac.Password)

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Authentication failed", resp.Status)
	}
	return resp, nil
}
func main() {
	fmt.Println("Menu\n----\n0.Quit\n1.Connect to server")
	choice := PromptInt(0, 1, "Make a hoice:")
	switch choice {
	case QUIT:
		os.Exit(0)
		break
	case CONNECT:
		cl := LoginScreen()
		ListFiles(cl)
		break
	}
}
func ListFiles(cl *AuthClient) {
	resp, err := cl.Connect("http://127.0.0.1:8000/files")

	if err != nil {
		fmt.Println("Error connecting: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	fmt.Println(string(body))
}
func LoginScreen() *AuthClient {
	for {
		var username, password string
		fmt.Print("Enter username: ")
		fmt.Scan(&username)
		fmt.Print("Enter password: ")
		fmt.Scan(&password)
		authClient := NewAuthClient(username, password)

		resp, err := authClient.Connect("http://127.0.0.1:8000/")

		if err != nil {
			fmt.Println("Error: ", err)
			continue
		}
		defer resp.Body.Close()
		return authClient
	}
}
func PromptInt(min int, max int, prompt string) int {
	var num int
	var err error
	for {
		fmt.Print(prompt)
		_, err = fmt.Scan(&num)
		if err != nil {
			fmt.Println("Error: ", err)
			continue
		}
		if num < min || num > max {
			fmt.Printf("Number must be between %d and %d\n", min, max)
			continue
		}
		break
	}
	return num
}
