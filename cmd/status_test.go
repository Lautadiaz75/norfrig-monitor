package cmd

import (
	"testing"

	"github.com/Lautadiaz75/norfrig-monitor/internal/api"
)

func TestFiltrarUrgentes_soloDevuelveCriticoYReorden(t *testing.T) {
	items := []api.OrdenItem{
		{SKU: "A", Semaforo: "CRITICO"},
		{SKU: "B", Semaforo: "OK"},
		{SKU: "C", Semaforo: "REORDEN"},
		{SKU: "D", Semaforo: "AGOTADO SIN VENTAS"},
		{SKU: "E", Semaforo: "SIN MOVIMIENTO"},
	}

	resultado := filtrarUrgentes(items)

	if len(resultado) != 2 {
		t.Fatalf("esperaba 2 items urgentes, got %d", len(resultado))
	}
	for _, item := range resultado {
		if item.Semaforo != "CRITICO" && item.Semaforo != "REORDEN" {
			t.Errorf("semaforo %q no deberia estar en urgentes", item.Semaforo)
		}
	}
}

func TestFiltrarUrgentes_sinUrgentesDevuelveVacio(t *testing.T) {
	items := []api.OrdenItem{
		{SKU: "A", Semaforo: "OK"},
		{SKU: "B", Semaforo: "SIN MOVIMIENTO"},
	}

	resultado := filtrarUrgentes(items)

	if len(resultado) != 0 {
		t.Errorf("esperaba slice vacío, got %d items", len(resultado))
	}
}

func TestFiltrarUrgentes_sliceVacioDeVacioDevuelveVacio(t *testing.T) {
	resultado := filtrarUrgentes([]api.OrdenItem{})

	if len(resultado) != 0 {
		t.Errorf("esperaba slice vacío, got %d items", len(resultado))
	}
}

func TestSortPorPrioridad_criticoAntesQueReorden(t *testing.T) {
	items := []api.OrdenItem{
		{SKU: "A", Semaforo: "REORDEN", UnidadesAPedir: 10},
		{SKU: "B", Semaforo: "OK", UnidadesAPedir: 0},
		{SKU: "C", Semaforo: "CRITICO", UnidadesAPedir: 5},
	}

	sortPorPrioridad(items)

	if items[0].Semaforo != "CRITICO" {
		t.Errorf("posicion 0: esperaba CRITICO, got %q", items[0].Semaforo)
	}
	if items[1].Semaforo != "REORDEN" {
		t.Errorf("posicion 1: esperaba REORDEN, got %q", items[1].Semaforo)
	}
}

func TestSortPorPrioridad_criticosOrdenadosPorUnidadesDesc(t *testing.T) {
	items := []api.OrdenItem{
		{SKU: "A", Semaforo: "CRITICO", UnidadesAPedir: 10},
		{SKU: "B", Semaforo: "CRITICO", UnidadesAPedir: 100},
		{SKU: "C", Semaforo: "CRITICO", UnidadesAPedir: 50},
	}

	sortPorPrioridad(items)

	if items[0].UnidadesAPedir != 100 {
		t.Errorf("primer CRITICO deberia tener 100 unidades, got %.0f", items[0].UnidadesAPedir)
	}
	if items[1].UnidadesAPedir != 50 {
		t.Errorf("segundo CRITICO deberia tener 50 unidades, got %.0f", items[1].UnidadesAPedir)
	}
	if items[2].UnidadesAPedir != 10 {
		t.Errorf("tercer CRITICO deberia tener 10 unidades, got %.0f", items[2].UnidadesAPedir)
	}
}
