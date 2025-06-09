package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"github.com/gin-gonic/gin"
)

// album represents data about a record album.
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// albums slice to seed record album data.
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	logDir := "mydata"
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.Mkdir(logDir, 0755) // Creates the directory with read/write/execute permissions for the owner, and read/execute for others
		if err != nil {
			log.Fatal(err)
		}
	}

	// Create a log file
	logFilePath := filepath.Join(logDir, "bgs.txt")
	logFile, err := os.Create(logFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	// Use a custom logger that writes to the file
	gin.DefaultWriter = logFile
	gin.DefaultErrorWriter = logFile

	router := gin.Default()
	router.GET("/albums", getAlbums)

	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)
	router.DELETE("/albums/:id", deleteAlbumByID)
	router.PUT("/albums/:id", updateAlbumByID)

	router.Run("0.0.0.0:8080")
}

// teste

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

// postAlbums adds an album from JSON received in the request body.
func postAlbums(c *gin.Context) {
	var newAlbum album

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	// Add the new album to the slice.
	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	// Loop through the list of albums, looking for
	// an album whose ID value matches the parameter.
	for _, a := range albums {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

func deleteAlbumByID(c *gin.Context) {
	id := c.Param("id")
	for i, a := range albums {
		if a.ID == id {
			albums = append(albums[:i], albums[i+1:]...)
			c.IndentedJSON((http.StatusOK), gin.H{"message": "album deleted"})
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

func updateAlbumByID(c *gin.Context) {
	var newAlbum album
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}
	id := c.Param("id")
	for i, a := range albums {
		if a.ID == id {
			albums[i] = newAlbum
			c.IndentedJSON(http.StatusAccepted, gin.H{"message": "album updated"})
			return
		}
	}
	//newAlbum.ID = id
	//albums = append(albums, newAlbum)
	//c.IndentedJSON(http.StatusCreated, gin.H{"message": "album created"})
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}
