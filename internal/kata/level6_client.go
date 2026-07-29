package kata

import (
	"encoding/json"
	"net/http"
)

// УРОВЕНЬ 6 — HTTP-клиент через httptest.NewServer.
//
// NameClient ходит за названием товара в чужой сервис.
//
// Контракт FetchName:
//   - сервер ответил 200 → название и nil
//   - сервер ответил 404 → ErrNotFound
//   - сервер ответил 5xx → ошибка (любая, но не nil)
//
// Задача: поднять фейковый сервер через httptest.NewServer, отдавать с него
// нужные статусы и проверить, что клиент реагирует по контракту.
// Настоящая сеть не нужна: httptest.NewServer слушает на localhost
// и отдаёт свой адрес в поле URL.
type NameClient struct {
	BaseURL string
	Client  *http.Client
}

func (c NameClient) FetchName(sku string) (string, error) {
	resp, err := c.Client.Get(c.BaseURL + "/name?sku=" + sku)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var body struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", err
	}
	return body.Name, nil
}
