package yahooboss

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/garyburd/go-oauth/oauth"
)

// BossSearch stores a Yahoo Boss Search, and implements the search functions
type BossSearch struct {
	Token      string
	Secret     string
	SearchType string
}

// BossResultRow contains a single row result
type BossResultRow struct {
	url   string
	title string
}

// BossResult contains basic data about the search
// and a list of all results rows
type BossResult struct {
	start        int
	end          int
	count        int
	totalresults int
	results      []BossResultRow
}

// BossResponse contains the most raw data about the search
type BossResponse struct {
	responsecode int
	bossresult   BossResult
}

func (bs *BossSearch) signQuery(url string, values *url.Values) {
	cred := oauth.Credentials{}
	cred.Token = bs.Token
	cred.Secret = bs.Secret

	client := oauth.Client{}
	client.SignatureMethod = oauth.HMACSHA1
	client.Credentials = cred

	client.SignForm(nil, "GET", url, *values)
}

func (bs *BossSearch) buildQuery(text string) string {
	searchType := bs.SearchType
	bossURL := fmt.Sprintf("%s%s", "https://yboss.yahooapis.com/ysearch/", searchType)

	getParams := url.Values{}
	getParams.Set("q", fmt.Sprintf("\"%s\"", text))
	getParams.Add("format", "json")
	getParams.Add("title", "lyrics")

	bs.signQuery(bossURL, &getParams)
	return fmt.Sprintf("%s?%s", bossURL, getParams.Encode())
}

// Search does a search using Yahoo Boss for "text"
func (bs *BossSearch) Search(text string) (BossResponse, error) {
	url := bs.buildQuery(text)

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	return parseSearchResult(body, bs.SearchType)
}

/* Todo: Add err */
func parseSearchResult(body []byte, searchType string) (BossResponse, error) {
	var dat map[string]interface{}

	if err := json.Unmarshal(body, &dat); err != nil {
		return BossResponse{}, errors.New("No JSON provided")
	}
	bossResponse := dat["bossresponse"].(map[string]interface{})
	searchTypeResults := bossResponse[searchType].(map[string]interface{})
	results := searchTypeResults["results"].([]interface{})
	var br BossResponse
	var bresult BossResult
	var brows []BossResultRow
	for _, value := range results {
		mapValue := value.(map[string]interface{})
		url := mapValue["url"].(string)
		title := mapValue["title"].(string)

		result := BossResultRow{url, title}
		brows = append(brows, result)
	}
	bresult.start = 0
	bresult.end = 50
	bresult.totalresults = 147
	bresult.results = brows

	br.bossresult = bresult
	br.responsecode = 200
	return br, nil
}
