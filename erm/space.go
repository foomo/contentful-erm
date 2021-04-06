package erm

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/foomo/contentful"
)

// spaceConf is the space config object passed to the template
type spaceConf struct {
	FuncMap      map[string]interface{}
	PackageName  string
	PackageDir   string
	Locales      []Locale
	ContentTypes []ContentType
	ContentType  ContentType
}

// GetLocales retrieves locale definition from Contentful
func getLocales(CMA *contentful.Contentful, spaceID string) (locales []Locale, err error) {

	col, err := CMA.Locales.List(spaceID).GetAll()
	if err != nil {
		return nil, fmt.Errorf("fetching locales: %w", err)
	}
	for _, item := range col.Items {
		var locale Locale
		byteArray, err := json.Marshal(item)
		if err != nil {
			return nil, fmt.Errorf("JSON encoding item: %w", err)
		}
		if err := json.Unmarshal(byteArray, &locale); err != nil {
			return nil, fmt.Errorf("JSON decoding item: %w", err)
		}
		locales = append(locales, locale)
	}
	return
}

// GetContentTypes retrieves content type definition from Contentful
func getContentTypes(CMA *contentful.Contentful, spaceID string) (contentTypes []ContentType, err error) {

	col := CMA.ContentTypes.List(spaceID)
	if _, err := col.GetAll(); err != nil {
		return nil, fmt.Errorf("fetching content types: %w", err)
	}
	for _, item := range col.Items {
		var contentType ContentType
		byteArray, err := json.Marshal(item)
		if err != nil {
			return nil, fmt.Errorf("JSON encoding item: %w", err)
		}
		if err := json.Unmarshal(byteArray, &contentType); err != nil {
			return nil, fmt.Errorf("JSON decoding item: %w", err)
		}
		contentTypes = append(contentTypes, contentType)
	}
	return
}

func getData(spaceID, cmaKey string, flagContentTypes []string) (finalContentTypes []ContentType, locales []Locale, err error) {
	// Get client
	CMA := contentful.NewCMA(cmaKey)
	CMA.Debug = false

	// Get space locales
	locales, err = getLocales(CMA, spaceID)
	if err != nil {
		return nil, nil, fmt.Errorf("getting locales: %w", err)
	}
	fmt.Println("Locales found:", locales)

	// Get content types
	contentTypes, err := getContentTypes(CMA, spaceID)
	if err != nil {
		return nil, nil, fmt.Errorf("getting content types: %w", err)
	}
	fmt.Println("Content types found:", len(contentTypes))

	finalContentTypes = []ContentType{}
	if len(flagContentTypes) == 0 {
		finalContentTypes = contentTypes
	} else {
		for _, ct := range contentTypes {
			if sliceIncludes(flagContentTypes, ct.Sys.ID) {
				finalContentTypes = append(finalContentTypes, ct)
			}
		}
	}
	var finalContentTypesString []string
	for _, finalContentType := range finalContentTypes {
		finalContentTypesString = append(finalContentTypesString, finalContentType.Name)
	}
	fmt.Println("Filtered Content types:", len(finalContentTypes), strings.Join(finalContentTypesString, ", "))
	return finalContentTypes, locales, nil
}

// GenerateAPI calls the generators
func GenerateAPI(dir, packageName, spaceID, cmaKey string, flagContentTypes []string) error {
	contentTypes, locales, err := getData(spaceID, cmaKey, flagContentTypes)
	if err != nil {
		return err
	}

	packageDir := filepath.Join(dir, packageName)
	if err := os.MkdirAll(packageDir, 0766); err != nil {
		return fmt.Errorf("creating package dir: %w", err)
	}
	funcMap := getFuncMap()
	conf := spaceConf{
		PackageDir:   packageDir,
		FuncMap:      funcMap,
		PackageName:  packageName,
		Locales:      locales,
		ContentTypes: contentTypes,
	}
	return generateCode(conf)
}
