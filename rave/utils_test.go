// Test utility functions

package rave

import (
	"reflect"
	"testing"
)

func assertEqual(t *testing.T, val1 interface{}, val2 interface{}) {
	if val1 != val2 {
		t.Fatalf(
			"'%s'(%s) is not Equal to '%s'(%s)",
			val1, reflect.TypeOf(val1), val2, reflect.TypeOf(val2),
		)
	}
}

// TestCheckRequiredParametersFail : Test Check required parameters function's failure
func TestCheckRequiredParametersFail(t *testing.T) {
	t.Parallel()

	params := map[string]interface{}{"first_name": "fred", "last_name": "quimby"}

	err := checkRequiredParameters(params, []string{"address"})

	assertEqual(t, err.Error(), "\"address\" is a required parameter for this method")
}

// TestCheckRequiredParametersSuccess : Test Check required parameters function's success
func TestCheckRequiredParametersSuccess(t *testing.T) {
	t.Parallel()

	params := map[string]interface{}{"first_name": "fred", "last_name": "quimby"}

	err := checkRequiredParameters(params, []string{"first_name", "last_name"})

	if err != nil {
		t.Fatal("Failed.")
	}
}
