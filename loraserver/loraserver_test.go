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

import (
	"testing"
	"time"
)

func TestLogin(t *testing.T) {
	_, err := New("platform.ceit.aut.ac.ir:50013")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGatewayFrameStream(t *testing.T) {
	l, err := New("platform.ceit.aut.ac.ir:50013")
	if err != nil {
		t.Fatal(err)
	}

	c, err := l.GatewayFrameStream("b827ebffff47d1a5")
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		for d := range c {
			t.Log(d)
		}
	}()

	time.Sleep(1 * time.Second)
}
