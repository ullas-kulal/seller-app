package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
)

// UrlDetails represents a scraped URL and product details
type UrlDetails struct {
	ID        string          `json:"id"`
	Url       string          `json:"url"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	Product   *ProductDetails `json:"product"`
}

// ProductDetails represents a product details
type ProductDetails struct {
	Name         string `json:"name"`
	ImageURL     string `json:"imageURL"`
	Description  string `json:"description"`
	Price        string `json:"price"`
	TotalReviews int    `json:"totalReviews"`
}

// webScraper scarpes the url and fetchs the product details
func webScraper(ctx *gin.Context) {
	c := colly.NewCollector()
	Pd := ProductDetails{}
	c.OnHTML("div[id=centerCol]", func(h *colly.HTMLElement) {
		Pd.Name = string(h.ChildText("span.a-size-large.product-title-word-break"))
		Reviews := string(h.ChildText("a[id=acrCustomerReviewLink] > span[id=acrCustomerReviewText]"))
		// Here we are using regex to fetch total review interger values from the string
		re := regexp.MustCompile("[0-9]+")
		num := re.FindAllString(Reviews, -1)
		totalReviewStr := strings.Join(num, "")
		TotalReviewsInt, err := strconv.Atoi(totalReviewStr)
		if err != nil {
			panic(err)
		}
		Pd.TotalReviews = TotalReviewsInt
		Pd.Price = string(h.ChildText("span.a-offscreen"))
	})
	c.OnHTML("div[id=imageBlock]", func(h *colly.HTMLElement) {
		Pd.ImageURL = string(h.ChildAttr("img.a-dynamic-image", "src"))
	})
	c.OnHTML("div[id=feature-bullets]", func(h *colly.HTMLElement) {
		Pd.Description = string(h.ChildText("span.a-list-item"))
	})
	url := ctx.Query("url")
	if url == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, nil)
		return
	}
	c.Visit(url)
	data := UrlDetails{Url: url, CreatedAt: time.Now(), UpdatedAt: time.Now(), Product: &Pd}
	jsonReq, err := json.Marshal(data)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
	}
	// Calls the API to Save the information in the mongodb
	response, err := http.Post("http://api-service:3001/products", "application/json", bytes.NewBuffer(jsonReq))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
	}
	respObj := UrlDetails{}
	json.Unmarshal(responseData, &respObj)
	ctx.JSON(http.StatusCreated, gin.H{
		"code":         http.StatusCreated,
		"product_info": respObj,
	})
}

func main() {
	r := gin.Default()
	r.POST("/products", webScraper)
	r.Run("0:3001")
}
