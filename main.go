package main

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/amimof/huego"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

func main() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("cycler")

	token := viper.GetString("token")
	bridgeIP := viper.GetString("bridgeIP")
	if token == "" {
		fmt.Printf("missing token\n")
		return
	}
	bridge := huego.New(bridgeIP, token)

	// Server
	e := echo.New()
	e.GET("/status", func(c echo.Context) error {
		return c.JSON(http.StatusOK, jsonOK(c, "UP"))
	})
	e.POST("/cycle/:deviceID", cycleHandler(bridge))
	_ = e.Start(":8080")
}

type errorResponse struct {
	code    int
	message string
}

func jsonErr(e echo.Context, err error) error {
	return e.JSON(http.StatusBadRequest, errorResponse{
		code:    http.StatusBadRequest,
		message: err.Error(),
	})
}

func jsonOK(e echo.Context, msg string) error {
	return e.JSON(http.StatusOK, errorResponse{
		code:    http.StatusOK,
		message: msg,
	})
}

type cycleCache struct {
	sync.Mutex
	cache map[int]time.Time
}

func (c *cycleCache) canCycle(lightID int) (bool, time.Time) {
	c.Lock()
	defer c.Unlock()
	now := time.Now()
	if last, ok := c.cache[lightID]; ok {
		if last.Add(5 * time.Minute).Before(now) {
			return false, last
		}
	}
	// Cache miss, so init
	c.cache[lightID] = now
	return true, now
}

func cycleHandler(bridge *huego.Bridge) echo.HandlerFunc {
	lastCycle := &cycleCache{}

	return func(e echo.Context) error {
		deviceID := e.Param("deviceID")
		id, err := strconv.Atoi(deviceID)
		if err != nil {
			return jsonErr(e, err)
		}
		light, err := bridge.GetLight(id)
		if err != nil {
			return jsonErr(e, err)
		}
		cycle, last := lastCycle.canCycle(id)
		if cycle {
			// Cycle
			go func() {
				_ = light.Off()
				time.Sleep(10 * time.Second)
				_ = light.On()
			}()
			return jsonOK(e, fmt.Sprintf("cycling light %d at %s", id, last.Format(time.RFC3339)))
		}
		return jsonErr(e, fmt.Errorf("skippging light cycle %d, last: %s", id, last.Format(time.RFC3339)))
	}
}
