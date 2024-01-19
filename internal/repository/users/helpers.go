package users

import (
	"encoding/json"
	"io"
	"net/http"
)

func (r *userRepository) fetchURL(url string, target any) error {
	response, err := http.Get(url)
	if err != nil {
		// r.logger.Error("making get request", "error", err.Error())
		return err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		// r.logger.Error("reading response body", "error", err.Error())
		return err
	}

	if err := json.Unmarshal(body, target); err != nil {
		// r.logger.Error("unmarshaling json", "error", err.Error())
		return err
	}

	return nil
}
