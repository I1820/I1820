/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 11-08-2018
 * |
 * | File Name:     main_downlink.go
 * +===============================================
 */

package actions

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type sendReq struct {
	Data          interface{} `json:"data" binding:"required"`
	ThingID       string      `json:"thing_id" binding:"required"`
	ApplicationID string      `json:"application_id" binding:"required"` // The ApplicationID can be retrieved using the API or from the web-interface, this is not the AppEUI!
	FPort         int         `json:"fport"`
	Confirmed     bool        `json:"confirmed"`

	SegmentSize  int   `json:"ss"`
	RepeatNumber int   `json:"rn"`
	Sleep        int64 `json:"sleep"` // Sleep interval between sends in seconds
}

// SendHandler handles downlink send request from applications
// This function is mapped to the path
// POST /send
func SendHandler(c echo.Context) error {
	var req sendReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	/*
		p, err := pm.GetThingProject(r.ThingID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		encoder := encoder.New(fmt.Sprintf("http://%s:%s", Config.Encoder.Host, p.Runner.Port))

		raw, err := encoder.Encode(r.Data, r.ThingID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		b, err := json.Marshal(lora.TxMessage{
			Reference: "abcd1234",
			FPort:     r.FPort,
			Data:      raw,
			Confirmed: r.Confirmed,
		})
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if err := cli.Publish(&client.PublishOptions{
			QoS:       mqtt.QoS0,
			TopicName: []byte(fmt.Sprintf("application/%s/node/%s/tx", r.ApplicationID, r.ThingID)),
			Message:   b,
		}); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		lan, err := json.Marshal(struct {
			Data []byte
		}{
			Data: b,
		})
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		go func() {

			req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/devices/%s/push", Config.LanServer.URL, r.ThingID), bytes.NewReader(lan))
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("Authorization", "aabbccddee11223344")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}

			if resp.StatusCode != 200 {
				c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("Invalid lan server response"))
				return
			}
		}()

		c.JSON(http.StatusOK, raw)
	*/
	return c.JSON(http.StatusNotImplemented, "downlink is under construction")
}

func sendRawHandler(c echo.Context) error {
	/*
		c.Header("Content-Type", "application/json")

		var r sendReq
		if err := c.BindJSON(&r); err != nil {
			return
		}

		b64, ok := r.Data.(string)
		if !ok {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("Invalid byte stream"))
			return
		}
		raw, err := base64.StdEncoding.DecodeString(b64)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		buffer := bytes.NewBuffer(raw)

		for raw := buffer.Next(r.SegmentSize); len(raw) != 0; raw = buffer.Next(r.SegmentSize) {
			log.Infof("Segment %v", raw)

			b, err := json.Marshal(lora.TxMessage{
				Reference: "abcd1234",
				FPort:     r.FPort,
				Data:      raw,
				Confirmed: r.Confirmed,
			})
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}

			for i := 0; i < r.RepeatNumber; i++ {
				if err := cli.Publish(&client.PublishOptions{
					QoS:       mqtt.QoS0,
					TopicName: []byte(fmt.Sprintf("application/%s/node/%s/tx", r.ApplicationID, r.ThingID)),
					Message:   b,
				}); err != nil {
					c.AbortWithError(http.StatusInternalServerError, err)
					return
				}
				log.Infof("MQTT Packet %s [%d]", b, i)

				time.Sleep(time.Duration(r.Sleep) * time.Second)
			}
		}

		c.JSON(http.StatusOK, raw)
	*/
	return c.JSON(http.StatusNotImplemented, "downlink is under construction")
}
