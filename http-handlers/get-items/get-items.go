package get_items

import (
	"encoding/csv"
	"github.com/go-chi/render"
	"log"
	"net/http"
	"net/url"
	"os"
	"simple_http/models"
	"strings"
)

type Response struct {
	Items []map[string]interface{} `json:"items"`
	models.Response
}

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ids, err := parsUrl(r.RequestURI)

		items, err := findItem(ids)
		if err != nil {
			log.Printf("failed to find items: %v", ids)

			render.Status(r, http.StatusBadRequest)

			render.JSON(w, r, models.Error("bad request"))
		}
		render.JSON(w, r, Response{
			Items:    items,
			Response: models.OK(),
		})
	}
}

func findItem(ids []string) ([]map[string]interface{}, error) {
	file, err := os.Open("ueba.csv")
	if err != nil {
		log.Printf("error while opening file: %v", err)
	}
	defer file.Close()

	csvReader := csv.NewReader(file)

	some, err := csvReader.ReadAll()
	if err != nil {
		log.Printf("error while reading file: %v", err)

		return nil, err
	}

	headers := some[0]

	var data []map[string]interface{}

	for _, id := range ids {
		for _, val := range some {
			if val[1] == id {
				m := make(map[string]interface{})
				for i, val := range val {
					m[headers[i]] = val
				}
				data = append(data, m)
			}
		}
	}
	return data, nil
}

func parsUrl(rowUrl string) ([]string, error) {
	parsedUrl, err := url.Parse(rowUrl)
	if err != nil {
		log.Printf("failed to parse url: %v", err)

		return nil, err
	}
	str := strings.Trim(parsedUrl.RawQuery, "id=")
	ids := strings.Split(str, "%")

	return ids, nil
}
