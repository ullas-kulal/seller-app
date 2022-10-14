package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
)

type UrlDetails struct {
	Url     string          `json:"url"`
	Product *ProductDetails `json:"product"`
}

type ProductDetails struct {
	Name         string `json:"name"`
	ImageURL     string `json:"imageURL"`
	Description  string `json:"description"`
	Price        string `json:"price"`
	TotalReviews string `json:"totalReviews"`
}

func webScraper(ctx *gin.Context) {
	c := colly.NewCollector()
	Pd := ProductDetails{}
	c.OnHTML("div[id=centerCol]", func(h *colly.HTMLElement) {
		Pd.Name = string(h.ChildText("span.a-size-large.product-title-word-break"))
		Pd.TotalReviews = string(h.ChildText("a[id=acrCustomerReviewLink] > span[id=acrCustomerReviewText]"))
		// Pd.Price = string(h.ChildText( "span.a-offscreen"))
		Pd.Price = string(h.ChildText("span[id=priceblock_ourprice]"))

	})
	c.OnHTML("div[id=imageBlock]", func(h *colly.HTMLElement) {
		Pd.ImageURL = string(h.ChildAttr("img.a-dynamic-image", "src"))
	})
	c.OnHTML("div[id=feature-bullets]", func(h *colly.HTMLElement) {
		Pd.Description = string(h.ChildText("span.a-list-item"))
	})

	// url := "https://www.amazon.com/DualShock-Wireless-Controller-PlayStation-Black-4/dp/B01LWVX2RG/ref=d_bmx_dp_dbd5zd7n_sccl_2_2/142-9550004-6467604?pd_rd_w=u3Fxl&content-id=amzn1.sym.e8434dc0-bea5-41a2-9054-015496b4c898&pf_rd_p=e8434dc0-bea5-41a2-9054-015496b4c898&pf_rd_r=FWRDZKSFC6TR6J7MZZJT&pd_rd_wg=WXXI7&pd_rd_r=6da9f145-b5c6-4281-b172-652a1dcb96cf&pd_rd_i=B01LWVX2RG&psc=1"
	url := "https://www.amazon.com/PlayStation-4-Pro-1TB-Console/dp/B01LOP8EZC/"
	c.Visit(url)
	data := UrlDetails{}
	data.Url = url
	data.Product = &Pd
	ctx.JSONP(http.StatusOK, data)

}

func main() {
	r := gin.Default()
	r.GET("/scrape", webScraper)
	r.POST("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run("0:3001")
}
