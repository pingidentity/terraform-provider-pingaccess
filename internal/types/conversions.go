package types

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	// client "github.com/pingidentity/pingaccess-go-client"
)

// Get a types.Set from a slice of strings
func GetStringSet(values []string) types.Set {
	setValues := make([]attr.Value, len(values))
	for i := 0; i < len(values); i++ {
		setValues[i] = types.StringValue(values[i])
	}
	set, _ := types.SetValue(types.StringType, setValues)
	return set
}
