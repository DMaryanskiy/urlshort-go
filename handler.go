package urlshortgo

import (
	"encoding/json"
	"net/http"

	"gopkg.in/yaml.v3"
)

type URLMap struct {
	Path string `yaml:"path" json:"path"`
	URL  string `yaml:"url" json:"url"`
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if url, ok := pathsToUrls[r.URL.Path]; !ok {
			fallback.ServeHTTP(w, r)
		} else {
			http.Redirect(w, r, url, http.StatusMovedPermanently)
		}
	})
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yaml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYAML(yaml)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}

func parseYAML(yamlRaw []byte) ([]URLMap, error) {
	parsedData := []URLMap{}

	err := yaml.Unmarshal(yamlRaw, &parsedData)
	if err != nil {
		return nil, err
	}

	return parsedData, nil
}

func buildMap(parsedYaml []URLMap) map[string]string {
	pathMap := make(map[string]string, len(parsedYaml))
	for _, item := range parsedYaml {
		pathMap[item.Path] = item.URL
	}

	return pathMap
}

func JSONHandler(json []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedJSON, err := parseJSON(json)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedJSON)
	return MapHandler(pathMap, fallback), nil
}

func parseJSON(jsonRaw []byte) ([]URLMap, error) {
	parsedData := []URLMap{}

	err := json.Unmarshal(jsonRaw, &parsedData)
	if err != nil {
		return nil, err
	}

	return parsedData, nil
}
