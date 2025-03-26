# API for Threads Knock Off Web App

Documenting the basic initial steps to get the application started and going.

### Creating / Initializing a Go Application
```
go mod init github.com/ManoVikram/Threads-Knock-Off-API
```

### Install Packages
Install the below package to build RESTful APIs
```
go get -u github.com/gin-gonic/gin
```

Install the below package to load and handle environment variables
```
go get github.com/joho/godotenv
```

Install the below package to handl CORS (Cross Origin Resource Sharing)
```
go get github.com/gin-contrib/cors
```

Install the below package to handle UUID values from the DB
```
go get github.com/google/uuid
```