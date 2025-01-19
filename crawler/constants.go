package main

const (
	targetTestUrl = "https://html_test_server:8080"
	targetRealUrl = "https://dblp.uni-trier.de"

	dblpIndex      = "/pers"
	dblpIndexCount = 3

	userAgent = "Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Mobile Safari/537.36"
	maxDepth  = 2
)

type UrlPath int

const (
	UnkownPage UrlPath = iota
	IndexPage
	AuthorPage
)
