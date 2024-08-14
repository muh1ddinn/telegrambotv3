package model

type Nomoztime struct {
	Region     string     `json:"region"`
	Date       string     `json:"date"`
	Weekday    string     `json:"weekday"`
	HijriDate  HijriDate  `json:"hijri_date"`
	DailyTimee DailyTimee `json:"times"`
}

type DailyTimee struct {
	Tong_saharlik string `json:"tong_saharlik"`
	Quyosh        string `json:"quyosh"`
	Peshin        string `json:"peshin"`
	Asr           string `json:"asr"`
	Shom_iftor    string `json:"shom_iftor"`
	Hufton        string `json:"hufton"`
}

type HijriDate struct {
	Month string `json:"month"`
	Day   int    `json:"day"`
}
