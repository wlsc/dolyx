package main

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const DEFAULT_HOST string = "localhost"
const DEFAULT_PORT int64 = 8080

type Command struct {
	Name  string `form:"name" json:"name" binding:"required"`
	Value string `form:"value" json:"value" binding:"required"`
}

type Request struct {
	Type    string  `form:"type" json:"type" binding:"required"`
	Command Command `form:"command" json:"command" binding:"required"`
}

type Image struct {
	Id      string
	Tag     string
	Created string
	Size    string
}

func main() {

	paramsLength := len(os.Args)

	host := DEFAULT_HOST
	port := DEFAULT_PORT

	if paramsLength > 1 {
		mappings := getProgramArguments(os.Args[1:])

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

	log.Println("Preparing server on " + host + ":" + strconv.FormatInt(port, 10))

	router := gin.Default()

	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tpl", gin.H{
			"images": getImages(),
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

	_ = router.Run(host + ":" + strconv.FormatInt(port, 10))
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
func getImages() []Image {

	cli, ctx := getClient()

	imagesRaw, err := cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		log.Fatal(err)
		return []Image{}
	}

	var images []Image
	for _, image := range imagesRaw {
		allImageIds := strings.ReplaceAll(image.ID, "sha256:", "")
		images = append(images,
			Image{
				Id:      allImageIds,
				Tag:     strings.Join(image.RepoTags, ","),
				Created: time.Unix(image.Created, 0).String(),
				Size:    ByteSize(image.Size)})
		log.Println(image)
	}

	return images
}

const (
	KILOBYTE = 1 << 10
	MEGABYTE = 1 << 20
	GIGABYTE = 1 << 30
)

func ByteSize(bytes int64) string {
	unit := ""
	value := float64(bytes)

	switch {
	case bytes >= GIGABYTE:
		unit = "G"
		value = value / GIGABYTE
	case bytes >= MEGABYTE:
		unit = "M"
		value = value / MEGABYTE
	case bytes >= KILOBYTE:
		unit = "K"
		value = value / KILOBYTE
	}

	result := strconv.FormatFloat(value, 'f', 1, 64)
	result = strings.TrimSuffix(result, ".0")
	return result + unit
}

/**
 *	Removes a given images
 */
func removeImage(id string) bool {

	cli, ctx := getClient()

	results, err := cli.ImageRemove(ctx, id, types.ImageRemoveOptions{true, true})
	if err != nil {
		return false
	}

	for _, result := range results {
		log.Println(result.Deleted)
		log.Println(result.Untagged)
	}

	return true
}

/**
 *	Stops all running containers and removes all images
 */
func removeAllImages() bool {

	cli, ctx := getClient()

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		log.Print("Killing container ", container.ID[:10])

		err := cli.ContainerKill(ctx, container.ID, "");
		if err != nil {
			log.Println(err.Error())
			return false
		}

		err = cli.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{})
		if err != nil {
			log.Println(err.Error())
			return false
		}
	}

	images, err := cli.ImageList(ctx, types.ImageListOptions{All: true})
	if err != nil {
		log.Println(err.Error())
		return false
	}

	for _, image := range images {
		log.Println("Removing image " + strings.Join(image.RepoTags, ","))
		_, err := cli.ImageRemove(ctx, image.ID, types.ImageRemoveOptions{Force: true, PruneChildren: true})
		if err != nil {
			log.Println(err.Error())
			return false
		}
	}

	log.Println("All images were removed!")

	return true
}

/**
 *	Returns mappings key -> value from arguments
 */
func getProgramArguments(args []string) map[string]string {
	mappings := map[string]string{}
	argsLength := len(args)

	for i := 0; i < argsLength; i++ {
		items := strings.Split(args[i], "=")

		// no pair
		if len(items) < 2 {
			continue
		}

		mappings[items[0]] = items[1]
	}

	return mappings
}

func getClient() (*client.Client, context.Context) {

	ctx := context.Background()
	cli, err := client.NewClientWithOpts()
	if err != nil {
		panic(err)
	}

	cli.NegotiateAPIVersion(ctx)

	return cli, ctx
}
