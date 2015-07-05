package yahooboss

import (
  "github.com/garyburd/go-oauth/oauth"
  "net/url"
  "net/http"
  "io/ioutil"
  "log"
  "fmt"
  "encoding/json"
)

type BossSearch struct {
    Token string
    Secret string
    SearchType string
}

type BossResultRow struct {
  url string
  title string
}

type BossResult struct {
    start int
    end int
    count int
    totalresults int
    results []BossResultRow
}

type BossResponse struct {
  responsecode int
  bossresult BossResult
}



func (bs *BossSearch) Search(text string) BossResponse {
  token := bs.Token
  secret := bs.Secret
  search_type := bs.SearchType
  boss_url := fmt.Sprintf("%s%s","https://yboss.yahooapis.com/ysearch/", search_type)

  cred := oauth.Credentials{}
  cred.Token = token
  cred.Secret = secret

  get_params := url.Values{}
  get_params.Set("q", fmt.Sprintf("\"%s\"", text))
  get_params.Add("format", "json")
  get_params.Add("title", "lyrics")
  client := oauth.Client{}
  client.SignatureMethod = 0
  client.Credentials = cred
  client.SignForm( nil, "GET", boss_url, get_params )
  res, err := http.Get(fmt.Sprintf("%s?%s", boss_url, get_params.Encode()))
  if err != nil {
		log.Fatal(err)
	}

  defer res.Body.Close()

  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    panic(err)
  }
  var dat map[string]interface{}

  if err := json.Unmarshal(body, &dat); err != nil {
    panic(err)
  }

  return parseSearchResult(dat, bs.SearchType)
}

/* Todo: Add err */
func parseSearchResult(json map[string]interface{}, search_type string) BossResponse {
  boss_response := json["bossresponse"].(map[string]interface{})
  search_type_results := boss_response[search_type].(map[string]interface{})
  results := search_type_results["results"].([]interface{})
  var br BossResponse
  var bresult BossResult
  var brows []BossResultRow
  for _, value := range results {
    map_value := value.(map[string]interface{})
    url := map_value["url"].(string)
    title := map_value["title"].(string)

    result := BossResultRow{url, title}
    brows = append(brows, result)
  }
  bresult.start = 0
  bresult.end = 50
  bresult.totalresults = 147
  bresult.results = brows

  br.bossresult = bresult
  br.responsecode = 200
  fmt.Println(br)
  return br
}
