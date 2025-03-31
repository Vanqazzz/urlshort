# urlshort

***Very simple*** url shortener, written in go

## Installation
Clone repository:

```cmd
git clone https://github.com/Vanqazzz/urlshort && cd urlshort
```

Edit  the `docker-compose.yml` and set base url:

    command: ["/shortr", "--url", "YOUR_URL"]

Change "YOUR_URL", to the your URL where your app served.

Now you can run the app with docker.

## Docker
``

    docker-compose up --build

Now app should be avaible by visiting localhost:8080