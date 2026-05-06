package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Client hace llamadas HTTP a la API de Cloud Run de Norfrig.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// OrdenItem representa una fila de mart_compra_precalculadas.
type OrdenItem struct {
	Proveedor      string  `json:"proveedor"`
	SKU            string  `json:"sku"`
	NombreProducto string  `json:"nombre_producto"`
	StockActual    float64 `json:"stock_actual"`
	AvgDiarioProy  float64 `json:"avg_diario_proy"`
	DiasStockRest  float64 `json:"dias_stock_restante"`
	UnidadesAPedir float64 `json:"unidades_a_pedir"`
	Semaforo       string  `json:"semaforo"`
}

// GetOrden llama a GET /generar_orden?proveedor=<proveedor> y retorna los items.
func (c *Client) GetOrden(proveedor string) ([]OrdenItem, error) {
	url := fmt.Sprintf("%s/generar_orden?proveedor=%s", c.baseURL, proveedor)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("no se pudo conectar: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("la API respondió con estado %d", resp.StatusCode)
	}

	var items []OrdenItem
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, fmt.Errorf("error parseando respuesta JSON: %w", err)
	}

	return items, nil
}
