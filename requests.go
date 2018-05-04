/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 07-02-2018
 * |
 * | File Name:     requests.go
 * +===============================================
 */

package main

// project request payload
type projectReq struct {
	Name string `json:"name" binding:"required"`
	// TODO adds docker constraints and envs
}

// thing request payload
type thingReq struct {
	Name string `json:"name" binding:"required"`
}
