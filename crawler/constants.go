package main

const (
	base_url = "http://localhost:8080"
	real_url = "https://dblp.uni-trier.de"

	citations_article1      = "/citations_article1"
	citations_article2      = "/citations_article2"
	citations_page_1        = "/citations_page_1"
	citations_page_2        = "/citations_page_2"
	citations_page_fragment = "/citations_page_fragment"
	dblp_dc                 = "/dblp_dc"
	dblp_se                 = "/dblp_se"
	dblp_pers_index         = "/dblp_pers_index"

	dblp_index          = "/pers"
	dblp_index_count    = 3
	dblp_index_next_300 = "/pers?pos=301"
	dblp_author_1       = "/pid/80/2813"
	dblp_author_2       = "/pid/20/123"

	userAgent = "Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Mobile Safari/537.36"
	maxDepth  = 2
)

type UrlPath int

const (
	UnkownPage UrlPath = iota
	IndexPage
	AuthorPage
)
