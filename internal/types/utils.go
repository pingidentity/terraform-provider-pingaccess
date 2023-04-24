package types

import (
	"context"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Return true if this types.String represents an empty (but non-null and non-unknown) string
func IsEmptyString(str types.String) bool {
	return !str.IsNull() && !str.IsUnknown() && str.ValueString() == ""
}
func IsNonEmptyMap(m types.Map) bool {
	return !m.IsNull() && !m.IsUnknown() && m.Elements() != nil
}

func IsNonEmptyList(l types.List) bool {
	return !l.IsNull() && !l.IsUnknown() && l.Elements() != nil
}

func IsNonEmptyObj(obj types.Object) bool {
	return !obj.IsNull() && !obj.IsUnknown() && obj.Attributes() != nil
}

// Return true if this types.String represents a non-empty, non-null, non-unknown string
func IsNonEmptyString(str types.String) bool {
	return !str.IsNull() && !str.IsUnknown() && str.ValueString() != ""
}

// Return true if this value represents a defined (non-null and non-unknown) value
func IsDefined(value attr.Value) bool {
	return !value.IsNull() && !value.IsUnknown()
}

// Check if an attribute slice contains a value
func Contains(slice []attr.Value, value attr.Value) bool {
	for _, element := range slice {
		if element.Equal(value) {
			return true
		}
	}
	return false
}

// Check if a string slice contains a value
func StringSliceContains(slice []string, value string) bool {
	for _, element := range slice {
		if element == value {
			return true
		}
	}
	return false
}

// Check if two slices representing sets are equal
func SetsEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	// Assuming there are no duplicate elements since the slices represent sets
	for _, aElement := range a {
		found := false
		for _, bElement := range b {
			if bElement == aElement {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func CamelCaseToUnderscores(s string) string {
	re, _ := regexp.Compile(`([A-Z])`)
	res := re.ReplaceAllStringFunc(s, func(m string) string {
		return strings.ToLower("_" + m[0:])
	})
	return res
}

func UnderscoresToCamelCase(s string) string {
	re, _ := regexp.Compile(`(_[A-Za-z])`)
	res := re.ReplaceAllStringFunc(s, func(m string) string {
		return strings.ToUpper(m[1:])
	})
	return res
}

// Converts the basetypes.MapValue to map[string]interface{} required for PingAccess Client
func MapValuesToClientMap(mv basetypes.MapValue, con context.Context) *map[string]interface{} {
	type StringMap map[string]string
	var value StringMap
	mv.ElementsAs(con, &value, false)
	converted := map[string]interface{}{}
	for k, v := range value {
		converted[k] = v
	}
	return &converted
}

// Converts the basetypes.MapValue to map[string]interface{} required for PingAccess Client
func ObjValuesToClientMap(obj types.Object) *map[string]interface{} {
	attrs := obj.Attributes()
	converted := map[string]interface{}{}
	for key, value := range attrs {
		strvalue, ok := value.(basetypes.StringValue)
		if ok {
			// make this nicer
			if strvalue.IsNull() || strvalue.IsUnknown() {
				continue
			} else {
				converted[UnderscoresToCamelCase(key)] = strvalue.ValueString()
				continue
			}
		}
		boolvalue, ok := value.(basetypes.BoolValue)
		if ok {
			converted[key] = boolvalue.ValueBool()
			continue
		}
		int64value, ok := value.(basetypes.Int64Value)
		if ok {
			converted[key] = int64value.ValueInt64()
			continue
		}
	}

	return &converted
}

// Converts the map[string]attr.Type to basetypes.ObjectValue required for Terraform
func MaptoObjValue(attributeTypes map[string]attr.Type, attributeValues map[string]attr.Value, diags diag.Diagnostics) basetypes.ObjectValue {
	newObj, err := types.ObjectValue(attributeTypes, attributeValues)
	if err != nil {
		diags.AddError("ERROR: ", "An error occured while converting ")
	}
	return newObj
}

func StringValueOrNull(value interface{}) types.String {
	if value == nil {
		return basetypes.NewStringNull()
	} else {
		return types.StringValue(value.(string))
	}
}
