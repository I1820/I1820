/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 07-02-2018
 * |
 * | File Name:     pm/pm_test.go
 * +===============================================
 */

package pm

import "testing"

func TestGeThing(t *testing.T) {
	p := New("http://127.0.0.1:8080")

	thing, err := p.GetThing("parham")

	if err != nil {
		t.Fatalf("GetThing error: %s\n", err)
	}

	t.Logf("http://127.0.0.1:%s\n", thing.Project.Runner.Port)
}
