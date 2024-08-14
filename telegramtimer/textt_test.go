package teleram

import (
	"fmt"
	"testing"

	"timer/model"
	"timer/teleram"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestTextt(t *testing.T) {
	// Initialize the HTTP mock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Mock the API response
	mockResponse := `{
		"region": "Toshkent",
		"date": "2024-08-14",
		"weekday": "Chorshanba",
		"hijri_date": {
			"month": "safar",
			"day": 9
		},
		"times": {
			"tong_saharlik": "04:03",
			"quyosh": "05:31",
			"peshin": "12:28",
			"asr": "17:21",
			"shom_iftor": "19:26",
			"hufton": "20:48"
		}
	}`

	httpmock.RegisterResponder("GET", "https://islomapi.uz/api/present/day?region=Toshkent",
		httpmock.NewStringResponder(200, mockResponse))

	// Call the function
	nomoztime, err := teleram.Textt("Toshkent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Define the expected result
	expectedNomoztime := model.Nomoztime{
		Region:    "Toshkent",
		Date:      "2024-08-14",
		Weekday:   "Chorshanba",
		HijriDate: model.HijriDate{Month: "safar", Day: 9},
		DailyTimee: model.DailyTimee{
			Tong_saharlik: "04:03",
			Quyosh:        "05:31",
			Peshin:        "12:28",
			Asr:           "17:21",
			Shom_iftor:    "19:26",
			Hufton:        "20:48",
		},
	}

	// Compare the result with the expected result
	assert.Equal(t, expectedNomoztime, nomoztime)
}

func TestTextt_Error(t *testing.T) {
	// Initialize the HTTP mock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Simulate an error response
	httpmock.RegisterResponder("GET", "https://islomapi.uz/api/present/day?region=Toshkent",
		httpmock.NewErrorResponder(fmt.Errorf("simulated error")))

	// Call the function
	nomoztime, err := teleram.Textt("Toshkent")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	// Check that the returned value is the zero value for the Nomoztime struct
	assert.Equal(t, model.Nomoztime{}, nomoztime)
}
