package main

import (
	"log"
	"net/http"

	"gopkg.in/yaml.v2"
)

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url, ok := pathsToUrls[r.URL.Path]

		if ok {
			log.Printf("Redirecting %s to %s\n", r.URL.Path, url)
			http.Redirect(w, r, url, http.StatusMovedPermanently)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

func parseYAML(content []byte) ([]map[string]string, error) {
	m := []map[string]string{}

	err := yaml.Unmarshal(content, &m)

	if err != nil {
		return nil, err
	}

	return m, nil
}

func buildMap(yml []map[string]string) map[string]string {
	mapping := make(map[string]string)

	for _, m := range yml {
		mapping[m["path"]] = m["url"]
	}

	return mapping
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}
