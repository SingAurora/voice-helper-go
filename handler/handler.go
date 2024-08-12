package handler

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
	"voice-helper-go/api"
	"voice-helper-go/util"
)

var (
	fileDataMap = make(map[string][]byte)
	mapMutex    sync.Mutex
)

type SseMessage struct {
	Text   string `json:"text"`
	Voice  []byte `json:"voice"`
	Offset int    `json:"offset"`
}

func generateRandomString(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func Sts(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		util.RespFail(c, 500, err.Error())
		return
	}
	defer file.Close()

	out, err := os.Create("recording.wav")
	if err != nil {
		util.RespFail(c, 500, err.Error())
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		util.RespFail(c, 500, err.Error())
		return
	}

	// Read file content into memory
	fileContent, err := os.ReadFile("recording.wav")
	if err != nil {
		util.RespFail(c, 500, err.Error())
		return
	}

	// Generate a random string as the key
	key, err := generateRandomString(16)
	if err != nil {
		util.RespFail(c, 500, err.Error())
		return
	}

	// Store file content in the map
	mapMutex.Lock()
	fileDataMap[key] = fileContent
	mapMutex.Unlock()

	// Return the key to the client
	util.RespOK(c, key)
}

func Sse(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		c.String(http.StatusBadRequest, "key is required")
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	c.Stream(func(w io.Writer) bool {
		var err error
		mapMutex.Lock()
		data, exists := fileDataMap[key]
		if exists {
			delete(fileDataMap, key)
		}
		mapMutex.Unlock()

		if !exists {
			util.Sse(c, "error", "File not found")
			return false
		}
		stt, err := api.Stt(data)
		if err != nil {
			util.Sse(c, "error", err.Error())
			return false
		}
		ai, err := api.Ai(stt)
		if err != nil {
			util.Sse(c, "error", err.Error())
			return false
		}
		defer ai.Close()

		buffer := ""
		offset := 0
		start := time.Now()
		for {
			recv, err := ai.Recv()
			if err == io.EOF {
				if len(buffer) > 0 {
					util.Sse(c, "last", SseMessage{Text: buffer, Voice: []byte(""), Offset: offset})
					c.Writer.Flush()
				} else {
					util.Sse(c, "last", SseMessage{Text: "", Voice: []byte(""), Offset: offset})
					c.Writer.Flush()
				}
				return false
			}
			if err != nil {
				util.Sse(c, "error", err.Error())
				c.Writer.Flush()
				return false
			}

			if len(recv.Choices) > 0 {
				buffer += recv.Choices[0].Delta.Content
				if len(buffer) >= 20 {
					if strings.ContainsAny(buffer, "。？！.!?") {
						cleanBuffer := util.RemoveInvisibleChars(buffer)
						tts, err := api.Tts(cleanBuffer)
						if err != nil {
							util.Sse(c, "error", err.Error())
							return false
						}
						util.Sse(c, "message", SseMessage{Text: buffer, Voice: tts, Offset: offset})
						c.Writer.Flush()
						offset += 1
						buffer = ""
						end := time.Now()
						fmt.Println("消耗时间:", end.Sub(start).Seconds(), "s")
						start = end
					}
				}
			}
		}
		return false
	})
}
