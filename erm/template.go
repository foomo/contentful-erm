package erm

import (
	"strings"
	"text/template"
)

func getFuncMap() template.FuncMap {
	return template.FuncMap{
		"firstCap":                 strings.Title,
		"fieldIsBasic":             fieldIsBasic,
		"fieldIsComplex":           fieldIsComplex,
		"fieldIsAsset":             fieldIsAsset,
		"fieldIsBoolean":           fieldIsBoolean,
		"fieldIsDate":              fieldIsDate,
		"fieldIsInteger":           fieldIsInteger,
		"fieldIsJSON":              fieldIsJSON,
		"fieldIsLink":              fieldIsLink,
		"fieldIsLocation":          fieldIsLocation,
		"fieldIsMultipleReference": fieldIsMultipleReference,
		"fieldIsMultipleAsset":     fieldIsMultipleAsset,
		"fieldIsNumber":            fieldIsNumber,
		"fieldIsReference":         fieldIsReference,
		"fieldIsRichText":          fieldIsRichText,
		"fieldIsSymbol":            fieldIsSymbol,
		"fieldIsSymbolList":        fieldIsSymbolList,
		"fieldIsText":              fieldIsText,
		"mapFieldType":             mapFieldType,
		"mapFieldTypeLiteral":      mapFieldTypeLiteral,
		"onlyLetters":              latinLettersOnly,
	}
}

// mapFieldType takes a ContentTypeField from the space model definition
// and returns a string that matches the type of the map[string] for the VO
func mapFieldType(contentTypeName string, field ContentTypeField) string {
	switch field.Type {
	case fieldTypeArray: // It's either a text list or a multiple reference
		switch field.Items.Type {
		case fieldItemsTypeSymbol:
			return "[]string"
		case fieldItemsTypeLink:
			return "[]ContentTypeSys"
		default:
			return ""
		}
	case fieldTypeBoolean:
		return "bool"
	case fieldTypeDate:
		return "string"
	case fieldTypeInteger:
		return "float64"
	case fieldTypeLink: // A single reference
		return "ContentTypeSys"
	case fieldTypeLocation:
		return "ContentTypeFieldLocation"
	case fieldTypeNumber: // Floating point
		return "float64"
	case fieldTypeJSON: // JSON field
		return "interface{}"
	case fieldTypeRichText:
		return "interface{}"
	case fieldTypeSymbol: // It's a text field
		return "string"
	case fieldTypeText: // It's a text field
		return "string"
	default:
		return ""
	}
}

// mapFieldTypeLiteral takes a ContentTypeField from the space model definition
// and returns an empty literal that matches the type of the map[string] for the VO
func mapFieldTypeLiteral(contentTypeName string, field ContentTypeField) string {
	switch field.Type {
	case fieldTypeBoolean:
		return "false"
	case fieldTypeDate, fieldTypeSymbol, fieldTypeText:
		return `""`
	case fieldTypeInteger, fieldTypeNumber:
		return "0"
	case fieldTypeArray, fieldTypeLink, fieldTypeLocation, fieldTypeJSON, fieldTypeRichText:
		return "nil"
	default:
		return ""
	}
}

func fieldIsAsset(field ContentTypeField) bool {
	return field.Type == fieldTypeLink && field.LinkType == fieldLinkTypeAsset
}

func fieldIsBoolean(field ContentTypeField) bool {
	return field.Type == fieldTypeBoolean
}

func fieldIsDate(field ContentTypeField) bool {
	return field.Type == fieldTypeDate
}

func fieldIsInteger(field ContentTypeField) bool {
	return field.Type == fieldTypeInteger
}

func fieldIsJSON(field ContentTypeField) bool {
	return field.Type == fieldTypeJSON
}

func fieldIsLink(field ContentTypeField) bool {
	return field.Type == fieldTypeLink
}

func fieldIsLocation(field ContentTypeField) bool {
	return field.Type == fieldTypeLocation
}

func fieldIsMultipleAsset(field ContentTypeField) bool {
	return field.Type == fieldTypeArray && field.Items.Type == fieldItemsTypeLink && field.Items.LinkType == fieldLinkTypeAsset
}
func fieldIsMultipleReference(field ContentTypeField) bool {
	return field.Type == fieldTypeArray && field.Items.Type == fieldItemsTypeLink && field.Items.LinkType == fieldLinkTypeEntry
}

func fieldIsNumber(field ContentTypeField) bool {
	return field.Type == fieldTypeNumber
}

func fieldIsReference(field ContentTypeField) bool {
	return field.Type == fieldTypeLink && field.LinkType == fieldLinkTypeEntry
}

func fieldIsRichText(field ContentTypeField) bool {
	return field.Type == fieldTypeRichText
}

func fieldIsSymbol(field ContentTypeField) bool {
	return field.Type == fieldTypeSymbol
}

func fieldIsSymbolList(field ContentTypeField) bool {
	return field.Type == fieldTypeArray && field.Items.Type == fieldItemsTypeSymbol
}

func fieldIsText(field ContentTypeField) bool {
	return field.Type == fieldTypeText
}

func fieldIsBasic(field ContentTypeField) bool {
	return fieldIsSymbolList(field) || fieldIsBoolean(field) || fieldIsInteger(field) || fieldIsNumber(field) || fieldIsSymbol(field) || fieldIsText(field) || fieldIsDate(field)
}
func fieldIsComplex(field ContentTypeField) bool {
	return field.Type == fieldTypeJSON || field.Type == fieldTypeLocation || field.Type == fieldTypeRichText
}
