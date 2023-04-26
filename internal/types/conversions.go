package types

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	// client "github.com/pingidentity/pingaccess-go-client"
)

func StringToTF(v string) basetypes.StringValue {
	if v == "" {
		return types.StringNull()
	} else {
		return types.StringValue(v)
	}
}

func StringToInt64Pointer(value basetypes.StringValue) *int64 {
	valueToString := value.ValueString()
	newVal, _ := strconv.ParseInt(valueToString, 10, 64)
	return &newVal
}

func Int64PointerToString(value int64) string {
	return strconv.FormatInt(value, 10)
}

func Int64ToString(value types.Int64) string {
	return strconv.FormatInt(value.ValueInt64(), 10)
}

// Get a types.Set from a slice of strings
func GetStringSet(values string) types.Set {
	setValues := make([]attr.Value, len(values))
	set, _ := types.SetValue(types.StringType, setValues)
	return set
}

// Get a types.Set from a slice of int32
func GetInt64Set(values []int32) types.Set {
	setValues := make([]attr.Value, len(values))
	for i := 0; i < len(values); i++ {
		setValues[i] = types.Int64Value(int64(values[i]))
	}
	set, _ := types.SetValue(types.Int64Type, setValues)
	return set
}

// Get a types.String from the given string pointer, handling if the pointer is nil
func StringTypeOrNil(str *string, useEmptyStringForNil bool) types.String {
	if str == nil {
		// If a plan was provided and is using an empty string, we should use that for a nil string in the response.
		// For PingAccess nil and empty string is equivalent, but to Terraform they are distinct. So we
		// just want to match whatever is in the plan when we get a nil string back.
		if useEmptyStringForNil {
			// Use empty string instead of null to match the plan when resetting string properties.
			// This is useful for computed values being reset to null.
			return types.StringValue("")
		} else {
			return types.StringNull()
		}
	}
	return types.StringValue(*str)
}

// Get a types.Bool from the given bool pointer, handling if the pointer is nil
func BoolTypeOrNil(b *bool) types.Bool {
	if b == nil {
		return types.BoolNull()
	}
	return types.BoolValue(*b)
}

// Get a types.Int64 from the given int32 pointer, handling if the pointer is nil
func Int64TypeOrNil(i *int32) types.Int64 {
	if i == nil {
		return types.Int64Null()
	}

	return types.Int64Value(int64(*i))
}

// Get a types.Float64 from the given float32 pointer, handling if the pointer is nil
func Float64TypeOrNil(f *float32) types.Float64 {
	if f == nil {
		return types.Float64Null()
	}

	return types.Float64Value(float64(*f))
}

// Get types.Map form slice of Strings
func GetStringMap(m *string) types.Map {
	setValues := make(map[string]attr.Value)
	set, _ := types.MapValue(types.StringType, setValues)
	return set
}
