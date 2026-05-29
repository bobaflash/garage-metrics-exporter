package internal

import (
	"encoding/json"

	"io"
	"log"
	"log/slog"
	"net/http"
	"net/url"
)

type GarageCollector struct {
	client  *http.Client
	baseUrl string
	token   string
}

func NewGarageCollector(baseUrl, token string) *GarageCollector {
	return &GarageCollector{
		token:   token,
		baseUrl: baseUrl,
		client:  &http.Client{},
	}
}

func (r *GarageCollector) request(url string, queryParameter map[string]string) ([]byte, error) {

	req, err := http.NewRequest("GET", url, nil)

	req.Header.Add("Authorization", r.token)

	q := req.URL.Query()
	if len(queryParameter) > 0 {
		for k, v := range queryParameter {
			q.Add(k, v)
		}
	}

	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
	}
	return []byte(body), nil

}

func (r *GarageCollector) GetBucketIds() ([]string, error) {

	requestUrl, err := url.JoinPath(r.baseUrl, "/v2/ListBuckets")
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	body, err := r.request(requestUrl, nil)
	if err != nil {
		slog.Error(err.Error())
	}

	var result []string

	var jsonData []struct {
		Id string `json:"id"`
	}

	if err := json.Unmarshal(body, &jsonData); err != nil {
		return nil, err
	}

	for _, v := range jsonData {
		result = append(result, v.Id)
	}

	return result, nil
}

type bucketInfo struct {
	GlobalAliases []string `json:"globalAliases"`
	Bytes         int      `json:"bytes"`
}

func (r *GarageCollector) GetBucketIinfo(bucketId string, metrics *Metrics) error {

	requestUrl, err := url.JoinPath(r.baseUrl, "/v2/GetBucketInfo")

	if err != nil {
		slog.Error(err.Error())
		return err
	}

	body, err := r.request(requestUrl, map[string]string{"id": bucketId})

	if err != nil {
		slog.Error(err.Error())
	}

	var result bucketInfo

	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}

	metrics.Bytes.WithLabelValues(result.GlobalAliases...).Set(float64(result.Bytes))

	return nil
}
