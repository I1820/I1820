/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 22-02-2018
 * |
 * | File Name:     requests.go
 * +===============================================
 */

package main

type sendReq struct {
	Data          interface{} `json:"data" binding:"required"`
	ThingID       string      `json:"thing_id" binding:"required"`
	ApplicationID string      `json:"application_id" binding:"required"` // The ApplicationID can be retrieved using the API or from the web-interface, this is not the AppEUI!
	FPort         int         `json:"fport"`
	Confirmed     bool        `json:"confirmed"`
}
