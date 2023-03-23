package types

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	// client "github.com/pingidentity/pingaccess-go-client"
)

// Get a types.Set from a slice of strings
func GetStringSet(values string) types.Set {
	setValues := make([]attr.Value, len(values))
	set, _ := types.SetValue(types.StringType, setValues)
	return set
}
