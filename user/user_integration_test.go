// +build integration

package user_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/jzs/ajournal/utils/testhelpers"
)

func TestEndpoints(t *testing.T) {

	// Setup server
	s, log, shutdown := testhelpers.InitTestServer()
	defer func() {
		err := shutdown()
		if err != nil {
			log.Flush()
			t.Fatalf("Expected no error, got %v", err)
		}
	}()

	errmsg := func(err error) {
		if err != nil {
			log.Flush()
			t.Fatalf("Expected no error, got %v", err)
		}
	}

	c := s.Client()

	usr := map[string]interface{}{
		"Username": "integration_user",
		"Password": "integration_pass",
	}
	in, _ := json.Marshal(usr)
	_, err := c.Post(fmt.Sprintf("%v%v", s.URL, "/users"), "application/json", bytes.NewReader(in))
	errmsg(err)

	_, err = c.Get(fmt.Sprintf("%v%v", s.URL, "/api/users/me"))
	errmsg(err)
}
