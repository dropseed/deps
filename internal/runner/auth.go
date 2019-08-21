package runner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/dropseed/deps/internal/output"
)

const appBaseURL = "https://3.dependencies.io"

type authorizer struct {
	token string
}

func newAuthorizer() (*authorizer, error) {
	// TODO will later be able to also check a DEPS_KEY for
	// users not connected to hosted service

	token := os.Getenv("DEPS_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("DEPS_TOKEN must be set. Ask your team admin or log in to %s to get your token.", appBaseURL)
	}
	return &authorizer{
		token: token,
	}, nil
}

func (auth *authorizer) validate() error {
	client := &http.Client{}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/usage/", appBaseURL), nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "token "+auth.token)
	req.Header.Add("User-Agent", "deps")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("API call failed with response %d\n\n%s", resp.StatusCode, string(body))
	}

	output.Debug("Authorized with %s", appBaseURL)

	return nil
}

func (auth *authorizer) incrementUsage(quantity int) error {
	inputJSON, err := json.Marshal(map[string]int{
		"quantity": quantity,
	})
	if err != nil {
		panic(err)
	}

	client := &http.Client{}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/usage/", appBaseURL), bytes.NewBuffer(inputJSON))
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "token "+auth.token)
	req.Header.Add("User-Agent", "deps")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("API call failed with response %d\n\n%s", resp.StatusCode, string(body))
	}

	return nil
}
