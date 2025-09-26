package meilisearch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	ms "github.com/meilisearch/meilisearch-go"
)

type Client interface {
	Connect(host string, port int, apiKey string) error
	GetMeili() RawClient // dışarıdan "Index(...)" bekleniyor
	Add(index string, model interface{}) error
	Search(index string, query string, page int64) *SearchResponse
	GetByDate(index string, query string, startTime time.Time, endTime time.Time, page int64) *SearchResponse
	Delete(index string, id string) error

	Log(index Index, doc map[string]interface{})
	LogUserActivity(message string, userId uint)
	LogUserErrLog(message string, userId uint)
	LogSystemErrLog(message string)
}

type Index string

const (
	UserActivityLog Index = "user_activity_log"
	UserErrLog      Index = "user_err_log"
	SystemErrLog    Index = "system_err_log"
)

// ------- REST MODELLERİ (SDK bağımsız) -------
type SearchResponse struct {
	Hits               []json.RawMessage `json:"hits"`
	Offset             int64             `json:"offset,omitempty"`
	Limit              int64             `json:"limit,omitempty"`
	EstimatedTotalHits int64             `json:"estimatedTotalHits,omitempty"`
	TotalHits          *int64            `json:"totalHits,omitempty"`
	Page               *int64            `json:"page,omitempty"`
	TotalPages         *int64            `json:"totalPages,omitempty"`
}

type indexesQueryResponse struct {
	Results []struct {
		UID string `json:"uid"`
	} `json:"results"`
}

// --------------------------------------------
type client struct {
	baseURL string
	apiKey  string
	http    *http.Client
}

func New() Client { return &client{} }

// timeout’lu HTTP client
func makeHTTPClient() *http.Client {
	tr := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: 3 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   3 * time.Second,
		ResponseHeaderTimeout: 5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConns:          20,
		IdleConnTimeout:       90 * time.Second,
	}
	return &http.Client{
		Timeout:   6 * time.Second,
		Transport: tr,
	}
}

func (c *client) Connect(host string, port int, apiKey string) error {
	c.baseURL = fmt.Sprintf("http://%s:%d", host, port)
	c.apiKey = apiKey
	c.http = makeHTTPClient()

	// 1) Preflight /version
	if err := c.headOrGet("/version"); err != nil {
		return fmt.Errorf("[MEILI] preflight failed: %w", err)
	}

	// 2) Indexleri hazırla
	if err := c.ensureIndexes(); err != nil {
		return err
	}

	return nil
}

// -------------------------------------------------------------------
// RAW CLIENT SHIM — repository'ndeki raw.Index(...).Search(...) için
// -------------------------------------------------------------------

type RawClient interface {
	Index(name string) *IndexShim
}

type rawRoot struct{ c *client }

func (r *rawRoot) Index(name string) *IndexShim { return &IndexShim{c: r.c, name: name} }

type IndexShim struct {
	c    *client
	name string
}

// Search shim: repository ms.SearchRequest/ms.SearchResponse bekliyor
func (s *IndexShim) Search(q string, req *ms.SearchRequest) (*ms.SearchResponse, error) {
	payload := map[string]interface{}{
		"q": q,
	}
	if req != nil {
		// Sık kullanılan alanları geçiriyoruz; gerekirse genişletilebilir
		if req.Filter != "" {
			payload["filter"] = req.Filter
		}
		if req.Page != 0 {
			payload["page"] = req.Page
		}
		if req.Limit != 0 {
			payload["limit"] = req.Limit
		}
		if req.Offset != 0 {
			payload["offset"] = req.Offset
		}
		if len(req.Sort) > 0 {
			payload["sort"] = req.Sort
		}
		if len(req.AttributesToRetrieve) > 0 {
			payload["attributesToRetrieve"] = req.AttributesToRetrieve
		}
		if len(req.AttributesToSearchOn) > 0 {
			payload["attributesToSearchOn"] = req.AttributesToSearchOn
		}
	}
	body, _ := json.Marshal(payload)

	b, err := s.c.doJSON("POST", fmt.Sprintf("/indexes/%s/search", s.name), body)
	if err != nil {
		return nil, err
	}
	var res ms.SearchResponse
	if err := json.Unmarshal(b, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *client) GetMeili() RawClient {
	return &rawRoot{c: c}
}

// ------------- Public işlemler (adapter API) -------------
func (c *client) Add(index string, documents interface{}) error {
	body, _ := json.Marshal(documents)
	_, err := c.doJSON("POST", fmt.Sprintf("/indexes/%s/documents", index), body)
	return err
}

func (c *client) Search(index string, query string, page int64) *SearchResponse {
	payload := map[string]interface{}{
		"q":    query,
		"page": page,
	}
	body, _ := json.Marshal(payload)

	var res SearchResponse
	if b, err := c.doJSON("POST", fmt.Sprintf("/indexes/%s/search", index), body); err == nil {
		_ = json.Unmarshal(b, &res)
	} else {
		log.Println("[MEILI] Search error:", err)
	}
	return &res
}

func (c *client) GetByDate(index string, query string, startTime time.Time, endTime time.Time, page int64) *SearchResponse {
	filter := fmt.Sprintf("created_at_ts >= %d AND created_at_ts < %d", startTime.Unix(), endTime.Unix())
	payload := map[string]interface{}{
		"q":      query,
		"filter": filter,
		"page":   page,
	}
	body, _ := json.Marshal(payload)

	var res SearchResponse
	if b, err := c.doJSON("POST", fmt.Sprintf("/indexes/%s/search", index), body); err == nil {
		_ = json.Unmarshal(b, &res)
	} else {
		log.Println("[MEILI] GetByDate error:", err)
	}
	return &res
}

func (c *client) Delete(index string, id string) error {
	_, err := c.doJSON("DELETE", fmt.Sprintf("/indexes/%s/documents/%s", index, id), nil)
	return err
}

// ------------- İç yardımcılar -------------
func (c *client) ensureIndexes() error {
	// mevcut indexleri çek
	b, err := c.doJSON("GET", "/indexes", nil)
	if err != nil {
		return fmt.Errorf("[MEILI] get indexes failed: %w", err)
	}
	var list indexesQueryResponse
	if err := json.Unmarshal(b, &list); err != nil {
		return fmt.Errorf("[MEILI] parse indexes failed: %w", err)
	}
	have := map[Index]bool{}
	for _, it := range list.Results {
		have[Index(it.UID)] = true
	}

	// yoksa oluştur + filterable attrs
	filterables := []string{"message", "created_at", "created_at_ts"}
	for _, idx := range []Index{UserActivityLog, UserErrLog, SystemErrLog} {
		if have[idx] {
			log.Printf("[MEILI] index %s already exists", idx)
		} else {
			log.Printf("[MEILI] create index %s", idx)

			// create index → taskUid dönüyor
			create := map[string]string{"uid": string(idx), "primaryKey": "id"}
			taskBody, err := c.doJSON("POST", "/indexes", mustJSON(create))
			if err != nil {
				return fmt.Errorf("[MEILI] create index %s failed: %w", idx, err)
			}
			tuid, _ := parseTaskUID(taskBody)
			if tuid != 0 {
				if err := c.waitTask(tuid, 12*time.Second); err != nil {
					return fmt.Errorf("[MEILI] create index %s task failed: %w", idx, err)
				}
			}
		}

		// update filterable attributes → **PUT** kullan
		taskBody, err := c.doJSON(
			"PUT",
			fmt.Sprintf("/indexes/%s/settings/filterable-attributes", idx),
			mustJSON(filterables),
		)
		if err != nil {
			return fmt.Errorf("[MEILI] update filterables %s failed: %w", idx, err)
		}
		tuid, _ := parseTaskUID(taskBody)
		if tuid != 0 {
			if err := c.waitTask(tuid, 12*time.Second); err != nil {
				return fmt.Errorf("[MEILI] update filterables %s task failed: %w", idx, err)
			}
		}
	}
	return nil
}

func (c *client) headOrGet(path string) error {
	req, _ := http.NewRequest("GET", c.baseURL+path, nil)
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("status %s", resp.Status)
	}
	return nil
}

func (c *client) doJSON(method, path string, body []byte) ([]byte, error) {
	req, _ := http.NewRequest(method, c.baseURL+path, bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	out, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("meili %s %s → %s: %s", method, path, resp.Status, string(out))
	}
	return out, nil
}

// ---------- Log convenience ----------
func (c *client) Log(index Index, doc map[string]interface{}) {
	if err := c.Add(string(index), doc); err != nil {
		log.Println("[MEILI] Log error:", err)
	}
}

// küçük JSON helper
func mustJSON(v any) []byte {
	b, _ := json.Marshal(v)
	return b
}

// task body'den taskUid'i çek
func parseTaskUID(b []byte) (int64, error) {
	var x struct {
		TaskUID int64 `json:"taskUid"`
	}
	if err := json.Unmarshal(b, &x); err != nil {
		return 0, err
	}
	return x.TaskUID, nil
}

// /tasks/{id} tamamlanana kadar bekle
func (c *client) waitTask(taskID int64, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		if time.Now().After(deadline) {
			return fmt.Errorf("task %d timed out after %s", taskID, timeout)
		}
		b, err := c.doJSON("GET", fmt.Sprintf("/tasks/%d", taskID), nil)
		if err != nil {
			return err
		}
		var t struct {
			Status string `json:"status"`
			Error  any    `json:"error"`
		}
		if err := json.Unmarshal(b, &t); err != nil {
			return err
		}
		switch t.Status {
		case "succeeded":
			return nil
		case "failed", "canceled":
			return fmt.Errorf("task %d %s: %v", taskID, t.Status, t.Error)
		case "enqueued", "processing":
			time.Sleep(250 * time.Millisecond)
		default:
			// bilinmeyen durum: kısa bekleyip tekrar dene
			time.Sleep(250 * time.Millisecond)
		}
	}
}
