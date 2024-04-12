package transport

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"online-lists/internal/helpers"
)

func SetupRouter(svc ListService) *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	r.GET("/ya_list", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": svc.GetYaList(),
		})
	})
	r.GET("/headers", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": svc.GetHeaders(),
		})
	})
	r.GET("/set_csv", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": svc.SetDefaultCsv(c.Query("filename")),
		})
	})
	r.GET("/list_csv", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": svc.ListCsv(),
		})
	})
	r.GET("/add", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": svc.Add(c.Query("header"), c.Query("value")),
		})
	})
	r.GET("/ya_file", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": svc.YAFile(c.Query("filename")),
		})
	})
	r.GET("/ya_upload", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": svc.YAUpload(c.Query("filename")),
		})
	})
	r.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, helpers.GetHomeTemplate())
	})

	return r
}
