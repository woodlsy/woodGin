package helper

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func Proxy(c *gin.Context, targetURL string) {
	req, err := http.NewRequest(c.Request.Method, targetURL, c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	for key, values := range c.Request.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	c.Status(resp.StatusCode)
	if resp.Body != nil {
		_, err = io.Copy(c.Writer, resp.Body)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}
}
