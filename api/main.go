package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ProductDetails represents a product details
type ProductDetails struct {
	Name         string `json:"name"`
	ImageURL     string `json:"imageURL"`
	Description  string `json:"description"`
	Price        string `json:"price"`
	TotalReviews int    `json:"totalReviews"`
}

// UrlDetails represents a scraped URL and product details
type UrlDetails struct {
	ID        string          `json:"id"`
	Url       string          `json:"url"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	Product   *ProductDetails `json:"product"`
}

// SaveProductDetails saves a product details into	MongoDB
func SaveProductDetails(ctx context.Context, product UrlDetails, db *mongo.Database) (string, error) {
	collection := db.Collection("products")
	res, err := collection.InsertOne(ctx, product)
	// (primitive.ObjectID).Hex() is used to convert object id to string
	return res.InsertedID.(primitive.ObjectID).Hex(), err
}
func GetProductDetails(ctx context.Context, ID string, db *mongo.Database) (*UrlDetails, error) {
	productDetails := UrlDetails{}
	objectId, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, err
	}
	err = db.Collection("products").FindOne(ctx, bson.D{{"_id", objectId}}).Decode(&productDetails)
	return &productDetails, err
}
func main() {
	// URI represents mongodb Connstring
	URI := os.Getenv("MONGODB_CONNSTRING")
	client, err := mongo.NewClient(options.Client().ApplyURI(URI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	fmt.Println("Connect to MongoDB server")
	db := client.Database("sellerapp")
	r := gin.Default()
	// POST request to insert productsdetails into database
	r.POST("/products", func(c *gin.Context) {
		product := UrlDetails{}
		c.ShouldBind(&product)
		ID, err := SaveProductDetails(context.TODO(), product, db)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err)
			return
		}
		product.ID = ID
		c.AbortWithStatusJSON(http.StatusCreated, product)
		return
	})
	// GET request to fetch productsdetails from database
	r.GET("/products/:id", func(c *gin.Context) {
		ID := c.Param("id")
		productDetails, err := GetProductDetails(context.TODO(), ID, db)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err)
			return
		}
		c.AbortWithStatusJSON(http.StatusOK, *productDetails)
		return
	})
	r.Run("0:3001")
}
