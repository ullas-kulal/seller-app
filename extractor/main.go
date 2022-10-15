package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
)

// UrlDetails represents a scraped URL and product details
type UrlDetails struct {
	Url     string          `json:"url"`
	Product *ProductDetails `json:"product"`
}

// ProductDetails represents a product details
type ProductDetails struct {
	Name         string `json:"name"`
	ImageURL     string `json:"imageURL"`
	Description  string `json:"description"`
	Price        string `json:"price"`
	TotalReviews string `json:"totalReviews"`
}

// webScraper scarpes the url and fetchs the product details
func webScraper(ctx *gin.Context) {
	c := colly.NewCollector()
	Pd := ProductDetails{}
	c.OnHTML("div[id=centerCol]", func(h *colly.HTMLElement) {
		Pd.Name = string(h.ChildText("span.a-size-large.product-title-word-break"))
		Pd.TotalReviews = string(h.ChildText("a[id=acrCustomerReviewLink] > span[id=acrCustomerReviewText]"))
		Pd.Price = string(h.ChildText("span.a-offscreen"))
	})
	c.OnHTML("div[id=imageBlock]", func(h *colly.HTMLElement) {
		Pd.ImageURL = string(h.ChildAttr("img.a-dynamic-image", "src"))
	})
	c.OnHTML("div[id=feature-bullets]", func(h *colly.HTMLElement) {
		Pd.Description = string(h.ChildText("span.a-list-item"))
	})
	url := ctx.Query("url")
	fmt.Println(url)
	if url == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, nil)
		return
	}
	c.Visit(url)
	data := UrlDetails{Url: url, Product: &Pd}
	jsonReq, err := json.Marshal(data)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
	}
	response, err := http.Post("http://api-service:3001/products", "application/json", bytes.NewBuffer(jsonReq))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
	}
	fmt.Println(string(responseData))
	ctx.JSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": string(responseData),
	})
}

func main() {
	r := gin.Default()
	r.POST("/products", webScraper)
	r.Run("0:3001")
}
