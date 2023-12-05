package main

var (
	styleSheets = []string{
		"./css/SlidingText.css",
	}
)

type SlidingText struct {
	LeadIn string
	Words  []string
}

type Index struct {
	SlidingText *SlidingText
	StyleSheets []string
}
