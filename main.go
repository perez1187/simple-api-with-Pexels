// https://www.youtube.com/watch?v=P-wPx6jmFnM&list=PL5dTjWUk_cPYztKD7WxVFluHvpBNM28N9&index=16

package main

import (
	// to call api

	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

// two const variables for api
const (
	PhotoApi = "https://api.pexels.com/v1"
	VideoApi = "https://api.pexels.com/videos"
)

type Client struct {
	Token          string
	hc             http.Client
	RemainingTimes int32 // how many free downloading from pixels api
}

// function creating Client
func NewClient(token string) *Client { // return a point to Client
	c := http.Client{}
	return &Client{Token: token, hc: c}
}

// so struct is something that we define on only what we have in struct we can edit / use (?)
// for example nodejs automaticallt understand json, but GO, not
// additionally, more work but more control over aplication
type SearchResult struct {
	Page         int32   `json:"page"` //json value comes from api, and value like Page we use in our program
	PerPage      int32   `json:"per_page"`
	TotalResults int32   `json:"total_Results`
	NextPage     string  `json:"next_page`
	Photos       []Photo `json:"photos` // slice of photos
}

type Photo struct {
	Id              int32       `json:"id"`
	Width           int32       `json:"width"`
	Height          int32       `json:"height"`
	Url             string      `json:"url"`
	Photographer    string      `json:"photographer"`
	PhotographerUrl string      `json:"photographer_url"`
	Src             PhotoSource `json:"src"` // another struct
}

type CuratedResult struct {
	Page     int32   `json:"page"`
	PerPage  int32   `json:"per_page"`
	NextPage string  `json:"next_page"`
	Photos   []Photo `json:"photos"`
}

// in api you can choose a size of picture
type PhotoSource struct {
	Original  string `json:"original"`
	Large     string `json:"large"`
	Large2x   string `json:"large2x"`
	Medium    string `json:"medium"`
	Small     string `json:small"`
	Potrait   string `json:"portrait"`
	Square    string `json:"square"`
	Landscape string `json:"landscape"`
	Tiny      string `json:"tiny"`
}

type VideoSearchResult struct {
	Page         int32   `json:"page"`
	PerPage      int32   `json:"per_page"`
	TotalResults int32   `json:"total_results"`
	NextPage     string  `json:"next_page"`
	Videos       []Video `json:"videos"`
}

type Video struct {
	Id            int32           `json:"id"`
	Width         int32           `json:"width"`
	Height        int32           `json:"height"`
	Url           string          `json:"url"`
	Image         string          `json:"image"`
	FullRes       interface{}     `json:"full_res"`
	Duration      float64         `json:"duration"`
	VideoFiles    []VideoFiles    `json:"video_files"`
	VideoPictures []VideoPictures `json:"video_pictures"`
}

type PopularVideos struct {
	Page         int32   `json:"page"`
	PerPage      int32   `json:"per_page"`
	TotalResults int32   `json:"total_results"`
	Url          string  `json:"url"`
	Videos       []Video `json:"videos"`
}

type VideoFiles struct {
	Id       int32  `json:"id"`
	Quality  string `json:"quality"`
	FileType string `json:"file_type"`
	Width    int32  `json:""width`
	Height   int32  `json:"height"`
	Link     string `json:"link"`
}

type VideoPictures struct {
	Id      int32  `json:"id"`
	Picture string `json:"picture"`
	Nr      int32  `json:"nr"`
}

// when you see something like this c *Client it is a struct method, and after that is a method name: SearchPhotos
// and a way to acces is client.name (c.SearchPhotos)
func (c *Client) SearchPhotos(query string, perPage, page int) (*SearchResult, error) { // it return resutl or err
	url := fmt.Sprintf(PhotoApi+"/search?query=%s&per_page=%d&page=%d", query, perPage, page)
	resp, err := c.requestDoWithAuth("GET", url) // we make request here, and capture to resp
	defer resp.Body.Close()                      // we close this function,when the function is finish

	data, err := ioutil.ReadAll(resp.Body) // we use ioutil package to read the body

	if err != nil { // if error is not nil, show err
		return nil, err
	}
	var result SearchResult // SearchResult is a struct

	// we capture result
	// Unmarshal save data to result
	err = json.Unmarshal(data, &result) // https://pkg.go.dev/encoding/json#Unmarshal
	return &result, err                 // so we return SearchResult "object"
}
func (c *Client) requestDoWithAuth(method, url string) (*http.Response, error) { // response http.Response or error

	// we capture NewRequest to req
	// https://pkg.go.dev/net/http
	req, err := http.NewRequest(method, url, nil) //http that we imported, and NewRequest method from http

	if err != nil {
		return nil, err
	}

	// we setup Header
	req.Header.Add("Authorization", c.Token) // c give us acces to Token in struct, because c is the client
	resp, err := c.hc.Do(req)                // hc is http.Client (from struct), we call Do method

	if err != nil {
		return resp, err
	}
	times, err := strconv.Atoi(resp.Header.Get("X-Ratelimit-Remaining"))
	if err != nil {
		return resp, nil
	} else {
		c.RemainingTimes = int32(times) // we set remainigTimes in Struct field to int32
	}
	return resp, nil
}

func (c *Client) CuratedPhotos(perPage, page int) (*CuratedResult, error) {
	url := fmt.Sprintf(PhotoApi+"/curated?per_page=%d&page=%d", perPage, page)
	resp, err := c.requestDoWithAuth("GET", url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result CuratedResult
	err = json.Unmarshal(data, &result)
	return &result, err
}

func (c *Client) GetPhoto(id int32) (*Photo, error) { // return Photo struct
	url := fmt.Sprintf(PhotoApi+"/photos/%d", id)
	resp, err := c.requestDoWithAuth("GET", url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	var result Photo // we create result type Photo (so it is struct)
	err = json.Unmarshal(data, &result)
	return &result, err
}

func (c *Client) GetRandomPhoto() (*Photo, error) {

	rand.Seed(time.Now().Unix())
	randNum := rand.Intn(1001)
	result, err := c.CuratedPhotos(1, randNum)
	if err == nil && len(result.Photos) == 1 {
		return &result.Photos[0], nil
	}
	return nil, err
}

func (c *Client) SearchVideo(query string, perPage, page int) (*VideoSearchResult, error) {
	url := fmt.Sprintf(VideoApi+"/search?query=%s&per_page=%d&page=%d", query, perPage, page)
	resp, err := c.requestDoWithAuth("GET", url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result VideoSearchResult
	err = json.Unmarshal(data, &result)
	return &result, err
}

func (c *Client) GetRemainingRequestsInThisMonth() int32 {
	return c.RemainingTimes
}

func (c *Client) PopularVideo(perPage, page int) (*PopularVideos, error) {
	url := fmt.Sprintf(VideoApi+"/popular?per_page=%d&page=%d", perPage, page)
	resp, err := c.requestDoWithAuth("GET", url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result PopularVideos
	err = json.Unmarshal(data, &result)
	return &result, err

}

func (c *Client) GetRandomVideo() (*Video, error) {
	rand.Seed(time.Now().Unix())
	randNum := rand.Intn(1001)
	result, err := c.PopularVideo(1, randNum)
	if err == nil && len(result.Videos) == 1 {
		return &result.Videos[0], nil
	}
	return nil, err
}

func main() {

	// we manually set token
	os.Setenv("PexelsToken", "VKSbcM6xRE1punB1SWPXueRe6Osb6xAkYarp5G5Wy0QexaroT1drmQKL") // create account on pexels

	// we get token from env
	var TOKEN = os.Getenv("PexelsToken")

	//we create client to work with pexels api
	var c = NewClient(TOKEN)

	//now you can use the client to call  function like SearchPhotos
	//we have to ensure this function will return result or an error

	// c.SearchPhotos means that this is not just only funtion, but struct method
	// for Client more speciflicly
	// you can call diffrent functions , like searchVideos
	result, err := c.SearchPhotos("waves", 15, 1) //in brackects, query string, perPage, page int

	//now we can handle the error
	if err != nil {
		fmt.Errorf("Search error %v", err)
	}

	// if there are no results
	if result.Page == 0 {
		fmt.Errorf("search results is wrong")
	}

	//if everythink is ok, we print result
	fmt.Println(result)

}
