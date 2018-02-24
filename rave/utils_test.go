// Test utility functions

package rave

import "testing"

// TestCheckRequiredParametersFail : Test Check required parameters function's failure
func TestCheckRequiredParametersFail(*testing.T) {
	params := map[string]interface{}{"first_name": "fred", "last_name": "quimby"}

	checkRequiredParameters(params, []string{"address"})
}

// TestCheckRequiredParametersSuccess : Test Check required parameters function's success
func TestCheckRequiredParametersSuccess(*testing.T) {
	params := map[string]interface{}{"first_name": "fred", "last_name": "quimby"}

	checkRequiredParameters(params, []string{"first_name", "last_name"})
}
