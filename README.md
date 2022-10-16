# seller-app
# Preparation

You need followings:

- Docker installed

# Try this example

Clone this repository.

```console
$ git clone git@github.com:ullas-kulal/seller-app.git
$ cd seller-app/
```

Build and Run Docker container.

```console
$ docker-compose build
$ docker-compose up
```

To stop, in the console where `docker-compose` is running, hit `Ctrl + C` and wait.

### POST API to scrape an Amazon web page given its URL


Create – POST  http://localhost:3001/products

Querey Params – url = {Amazon url to scrape}

Example: Create – POST  http://localhost:3001/products?url=https%3A%2F%2Fwww.amazon.com%2FPlayStation-4-Pro-1TB-Console%2Fdp%2FB01LOP8EZC%2F

Response body:

    [
        {
          "code": 201,
          "product_info": {
            "id": "634ba0098f993eac4bf59cc2",
            "url": "https://www.amazon.com/PlayStation-4-Pro-1TB-Console/dp/B01LOP8EZC/",
            "created_at": "2022-10-16T06:09:13.640077Z",
            "updated_at": "2022-10-16T06:09:13.64007728Z",
            "product": {
              "name": "PlayStation 4 Pro 1TB Console",
              "imageURL": "https://m.media-amazon.com/images/W/IMAGERENDERING_521856-T1/images/I/6118ctEjpoL.__AC_SX300_SY300_QL70_ML2_.jpg",
              "description": "Make sure this fits by entering your model number. \n Heighten your experiences. Enrich your adventures. Let the super charged PS4 Pro lead the way   4K TV Gaming : PS4 Pro outputs gameplay to your 4K TV   More HD Power: Turn on Boost Mode to give PS4 games access to the increased power of PS4 Pro   HDR Technology : With an HDR TV, compatible PS4 games display an unbelievably vibrant and life like range of colors",
              "price": "",
              "totalReviews": 12409
            }
          }
        }
    ]

