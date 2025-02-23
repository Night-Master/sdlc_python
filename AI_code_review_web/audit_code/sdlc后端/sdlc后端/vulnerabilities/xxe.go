package vulnerabilities

import (
	"encoding/xml"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func XML_parse(c *gin.Context) {
	var xmlData struct {
		XMLName xml.Name `xml:"data"`
		Content string   `xml:",innerxml"`
	}

	if err := c.BindXML(&xmlData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": 0, "message": "Invalid XML input"})
		return
	}

	// This is vulnerable to XXE attacks
	c.JSON(http.StatusOK, gin.H{"status": 1, "message": "XML processed successfully", "data": xmlData.Content})
}
