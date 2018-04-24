/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 25-04-2018
 * |
 * | File Name:     loraserver/loraserver_test.go
 * +===============================================
 */

package loraserver

import "testing"

func TestLogin(t *testing.T) {
	_, err := New("https://platform.ceit.aut.ac.ir:50013/api")
	if err != nil {
		t.Fatal(err)
	}
}
