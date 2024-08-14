package teleram

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"timer/model"
)

func Textt(region string) (model.Nomoztime, error) {
	if region == "" {
		region = "Toshkent"
	}

	url := fmt.Sprintf("https://islomapi.uz/api/present/day?region=%s", region)
	res, err := http.Get(url)
	if err != nil {
		return model.Nomoztime{}, fmt.Errorf("error fetching data from API: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return model.Nomoztime{}, fmt.Errorf("error reading response body: %v", err)
	}

	var nomoztime model.Nomoztime
	if err := json.Unmarshal(body, &nomoztime); err != nil {
		return model.Nomoztime{}, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	return nomoztime, nil
}
