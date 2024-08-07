package get_items

import (
	"encoding/csv"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
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

func New(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.get-items.New"

		log := logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		ids, err := parsUrl(r.RequestURI, log)

		items, err := findItem(ids, log)
		if err != nil {
			log.Error("failed to find items", "id", ids)

			render.Status(r, http.StatusInternalServerError)

			render.JSON(w, r, models.Error("internal server error"))
		}
		render.JSON(w, r, Response{
			Items:    items,
			Response: models.OK(),
		})
	}
}

func findItem(ids []string, log *slog.Logger) ([]map[string]interface{}, error) {
	file, err := os.Open("ueba.csv")
	if err != nil {
		log.Error("error while opening file", err)
	}
	defer file.Close()

	csvReader := csv.NewReader(file)

	some, err := csvReader.ReadAll()
	if err != nil {
		log.Error("error while reading file", err)

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

func parsUrl(rowUrl string, log *slog.Logger) ([]string, error) {
	parsedUrl, err := url.Parse(rowUrl)
	if err != nil {
		log.Error("failed to parse url", rowUrl)

		return nil, err
	}
	str := strings.Trim(parsedUrl.RawQuery, "ids=")
	ids := strings.Split(str, "%")

	return ids, nil
}
