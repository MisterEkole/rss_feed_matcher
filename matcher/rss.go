package matcher

import (
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"M1/search"
)

type (

	item struct{
		XMLName xml.Name `xml:"item"`
		PubDate string `xml:"pubDate"`
		Title string `xml:"title"`
		Description string `xml:"description"`
		Link string `xml:"link"`
		GUID string `xml:"guid"`
		
		GeoRSSPoint string `xml:"georss:point"`
	}

	image struct {
		XMLName xml.Name `xml:"image"`
		URL string `xml:"url"`
		Title string `xml:"title"`
		Link string `xml:"link"`

	}

	channel struct{
		XMLName xml.Name `xml:"channel"`
		Title string `xml:"title"`
		Description string `xml:"description"`
		Link string `xml:"link"`
		PubDate string `xml:"pubDate"`
		LastBuildDate string `xml:"lastBuildDate"`
		TTL string `xml:"ttl"`
		Language string `xml:"language"`
		ManagingEditor string `xml:"managingEditor"`
		Image image `xml:"image"`
		Item []item `xml:"item"`

	}

	rssDocument struct{

		XMLName xml.Name `xml:"rss"`
		Channel channel `xml:"channel"`

	}


)

type rssMatcher struct{}
//init func that registers matcher
func init(){
	var matcher rssMatcher
	search.Register("rss", matcher)
}

// look at document for specified search item

func (m rssMatcher) Search(feed *search.Feed,searchTerm string)([]*search.Result, error){
	var results []*search.Result
log.Printf("Search Feed Type: %s, Site: %s, for URI %s\n",feed.Type, feed.Name, feed.URI)

//receive data to seach
doc, err:= m.retrieve(feed)
if err!=nil{
	return nil, err
}

for _,channelItem:= range doc.Channel.Item{
	//check title of matched term

	matched, err:= regexp.MatchString(searchTerm, channelItem.Title)

	if err!= nil{
		return nil, err
	}
// if match found save the match
	if matched {
		results = append(results, &search.Result{
			Field: "Title",
			Content: channelItem.Title,
		})

	//check desctiption of searched term
	matched, err:= regexp.MatchString(searchTerm, channelItem.Description)

	if err!= nil{
		return nil, err
	}
	if matched{
		results= append(results, &search.Result{
			Field : "Title",
			Content: channelItem.Description,
		})
	}

	}
}
return results, nil

}

func (m rssMatcher) retrieve(feed*search.Feed)(*rssDocument, error){
	if feed.URI==""{
		return nil, errors.New("No rss uri provided")


	}
	resp, err:=http.Get(feed.URI)
	if err!=nil{
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode!=200{
		return nil, fmt.Errorf("Http response error %d\n", resp.StatusCode)

	}

	var doc rssDocument
	err = xml.NewDecoder(resp.Body).Decode(&doc)
	return &doc, err
}