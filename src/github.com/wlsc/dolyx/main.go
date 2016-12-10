package main

import (
	"github.com/gin-gonic/gin"
	"os"
	"strconv"
	"strings"
	"net/http"
	"log"
	"os/exec"
	"regexp"
)

const DEFAULT_HOST string = "localhost"
const DEFAULT_PORT int64 = 8080

const CMD_DELIMITER string = "="

type Command struct {
	Name  string                `form:"name" json:"name" binding:"required"`
	Value string                `form:"value" json:"value" binding:"required"`
}

type Request struct {
	Type    string        `form:"type" json:"type" binding:"required"`
	Command Command        `form:"command" json:"command" binding:"required"`
}

func main() {

	paramsLength := len(os.Args)

	host := DEFAULT_HOST
	port := DEFAULT_PORT

	if paramsLength > 1 {
		mappings := getMappings(os.Args[1:])

		if mappings["host"] == "" {
			host = DEFAULT_HOST
		} else {
			host = mappings["host"]
		}

		if mappings["port"] == "" {
			port = DEFAULT_PORT
		} else {
			port, _ = strconv.ParseInt(mappings["port"], 0, 0)
		}
	}

	regex_multiple_spaces := regexp.MustCompile(`[\s\p{Zs}]{2,}`)

	log.Println("Preparing server on " + host + ":" + strconv.FormatInt(port, 10))

	router := gin.Default()

	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static")
	router.GET("/", func(c *gin.Context) {

		images := getImages()
		imagesHeader := strings.Split(regex_multiple_spaces.ReplaceAllString(images[0], ";"), ";")
		imagesContent := [][]string{}

		for _, line := range images[1:] {

			var splitted = strings.Split(regex_multiple_spaces.ReplaceAllString(line, ";"), ";")

			if len(splitted) > 1 {
				imagesContent = append(imagesContent, splitted)
			}
		}

		c.HTML(http.StatusOK, "index.tpl", gin.H{
			"images_list_header" : imagesHeader,
			"images_list" : imagesContent,
		})
	})

	router.POST("/control", func(c *gin.Context) {

		var request Request

		if c.BindJSON(&request) == nil {
			switch request.Type {
			case "images":
				handleImagesCommand(request.Command, c)
				break
			default:
				log.Println("default case fired")
				showError(c)
			}
		} else {
			showError(c)
		}
	})

	// listen on address:port
	router.Run(host + ":" + strconv.FormatInt(port, 10))
}

/**
 *	Handling command "images"
 */
func handleImagesCommand(command Command, c *gin.Context) {

	log.Println("Handling image command...")

	var commandName = command.Name
	var commandValue = command.Value

	switch commandName {
	case "remove":
		if removeImage(commandValue) {
			c.JSON(http.StatusOK, gin.H{"status": 0})
		} else {
			sendRemoveFailed(c)
		}
		break
	case "removeall":
		if removeAllImages() {
			c.JSON(http.StatusOK, gin.H{"status": 0})
		} else {
			sendRemoveFailed(c)
		}
		break
	}
}

/**
 *	Tells user unable remove an image
 */
func sendRemoveFailed(c *gin.Context) {
	c.JSON(http.StatusConflict, gin.H{"error": "Cannot remove an image"})
}

/**
 *	Tells user about general error
 */
func showError(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{"error": "Cannot parse command from client"})
}

/**
 *	Retrieves all available docker images
 */
func getImages() [] string {

	out, err := exec.Command("docker", "images").Output()

	if err != nil {
		log.Fatal(err)
		return []string{}
	}

	return strings.Split(string(out[:]), "\n")
}

/**
 *	Removes a given images
 */
func removeImage(id string) bool {

	out, err := exec.Command("docker", "rmi", "-f", id).Output()

	if err != nil {
		return false
	}

	log.Println(out)

	return true
}

/**
 *	Stops all running containers and removes all images
 */
func removeAllImages() bool {

	// stopping all containers
	out, err := exec.Command("docker", "ps", "-a", "-q").Output()

	if err != nil {
		log.Println("Cannot list running containers")
		return false
	}

	var runningContainers = strings.Split(string(out), "\n")

	for _, id := range runningContainers {

		if id == "" {
			continue
		}

		_, err := exec.Command("docker", "stop", id).Output()

		if err != nil {
			log.Println("Cannot stop container", id)
			return false
		}
	}

	// removing all images
	out, err = exec.Command("docker", "images", "-q").Output()

	if err != nil {
		log.Println("Cannot list images (IDs)")
		return false
	}

	var imageIds = strings.Split(string(out), "\n")

	for _, id := range imageIds {

		if id == "" {
			continue
		}

		_, err := exec.Command("docker", "rmi", "-f", id).Output()

		if err != nil {
			log.Println("Cannot remove image", id)
			return false
		}
	}

	return true
}

/**
 *	Returns mappings key -> value from arguments
 */
func getMappings(args []string) map[string]string {
	mappings := map[string]string{}
	argsLength := len(args)

	for i := 0; i < argsLength; i++ {
		items := strings.Split(args[i], CMD_DELIMITER)

		// no pair
		if len(items) < 2 {
			continue
		}

		mappings[items[0]] = items[1]
	}

	return mappings
}