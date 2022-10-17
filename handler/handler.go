package urlshortener

import (
	"net/http"

	"gopkg.in/yaml.v2"
)

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}

		fallback.ServeHTTP(w, r)
	}
}

func YAMLHandler(yamlIn []byte, fallback http.Handler) (http.HandlerFunc, error) {
	//parse yaml
	var pathUrls []pathURL
	err := yaml.Unmarshal(yamlIn, pathUrls)

	if err != nil {
		return nil, err
	}
	//convert yaml array to map
	pathToUrls := make(map[string]string)
	for _, pu := range pathUrls {
		pathToUrls[pu.Path] = pu.URL
	}

	//return map handler
	return MapHandler(pathToUrls, fallback), nil
}

type pathURL struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}
