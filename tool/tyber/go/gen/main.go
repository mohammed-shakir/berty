package main

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"berty.tech/berty/v2/tool/tyber/go/parser"
	"berty.tech/weshnet/pkg/tyber"
)

func main() {
	fmt.Print("// generated by berty.tech/berty/v2/tool/tyber/gen\n\n")
	fmt.Printf("export %s\n", tsType(tyber.StatusType("")))
	fmt.Printf("export %s\n", tsType(&tyber.Detail{}))
	fmt.Printf("export %s\n", tsType(&tyber.Step{}))
	fmt.Printf("export %s\n", tsType(&parser.AppStep{}))
	fmt.Printf("export %s\n", tsType(&parser.CreateStepEvent{}))
	fmt.Printf("export %s\n", tsType(&parser.SubTarget{}))
	fmt.Printf("export %s\n", tsType(&parser.CreateTraceEvent{}))
	fmt.Printf("export %s\n", tsType(&parser.UpdateTraceEvent{}))
}

func primaryType(str string) string {
	if str == "Bool" {
		return "boolean"
	}
	if str == "String" {
		return "string"
	}
	if str == "TimeTime" {
		return "string"
	}
	return str
}

func finalTypeName(t reflect.Type) string {
	name := t.String()
	str := strings.Title(strings.Replace(strings.Join(strings.Split(name, "."), ""), "*", "", -1))
	if strings.HasPrefix(str, "[]") {
		return primaryType(strings.TrimPrefix(str, "[]")) + "[]"
	}
	return primaryType(str)
}

func endTypeName(t reflect.Type) string {
	name := t.String()
	if parts := strings.Split(name, "."); len(parts) > 1 {
		name = parts[len(parts)-1]
	}
	return name
}

func tsFields(level int, reflectElem reflect.Value) []string {
	fields := []string{}
	for i := 0; i < reflectElem.NumField(); i++ {
		member := reflectElem.Field(i)
		memberField := reflectElem.Type().Field(i)
		memberTypeName := finalTypeName(memberField.Type)
		if memberField.Name == endTypeName(memberField.Type) && memberField.Type.Kind() == reflect.Struct {
			fields = append(fields, tsFields(level, member)...)
		} else {
			memberName := memberField.Name
			re := regexp.MustCompile(`json:".+"`)
			if jsonTag := string(re.Find([]byte(memberField.Tag))); jsonTag != "" {
				memberName = jsonTag[len(`json:"`) : len(jsonTag)-len(`"`)]
			}
			prefix := ""
			for j := 0; j < level; j++ {
				prefix += "  "
			}
			fields = append(fields, fmt.Sprintf("%s%s: %s\n", prefix, memberName, memberTypeName))
		}
	}
	return fields
}

func tsType(goType interface{}) string {
	reflectValue := reflect.ValueOf(goType)
	if reflectValue.Kind() == reflect.String {
		return fmt.Sprintf("type %s = string\n", finalTypeName(reflectValue.Type()))
	}
	// str := fmt.Sprintf("Value: %s\n", reflectValue.Type())
	reflectElem := reflectValue.Elem()
	str := fmt.Sprintf("interface %s {\n", finalTypeName(reflectElem.Type()))
	str += strings.Join(tsFields(1, reflectElem), "")
	str += "}\n"
	return str
}
