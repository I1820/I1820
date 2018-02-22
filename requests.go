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
	Data      string `json:"data" binding:"required"`
	ProjectID string `json:"project_id"`
	ThingID   string `json:"thing_id" binding:"required"`
}
