/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 22-02-2018
 * |
 * | File Name:     data.go
 * +===============================================
 */

package main

// Data downlink
type Data struct {
	Data      string `json: "data"`
	ProjectID string `json: "project_id"`
}
