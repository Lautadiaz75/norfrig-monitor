package cmd

import (
	"fmt"
	"os"

	"github.com/Lautadiaz75/norfrig-monitor/internal/api"
	"github.com/spf13/cobra"
)

var proveedor string

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Muestra el semáforo de stock por SKU",
	RunE:  runStatus,
}

func init() {
	statusCmd.Flags().StringVarP(&proveedor, "proveedor", "p", "TODOS", "Filtrar por proveedor")
	rootCmd.AddCommand(statusCmd)
}

func runStatus(cmd *cobra.Command, args []string) error {
	baseURL := os.Getenv("NORFRIG_API_URL")
	if baseURL == "" {
		return fmt.Errorf("variable NORFRIG_API_URL no configurada\nEjemplo: set NORFRIG_API_URL=https://tu-servicio.run.app")
	}

	fmt.Printf("Consultando API para proveedor: %s...\n\n", proveedor)

	client := api.NewClient(baseURL)
	items, err := client.GetOrden(proveedor)
	if err != nil {
		return fmt.Errorf("error consultando la API: %w", err)
	}

	if len(items) == 0 {
		fmt.Println("No hay SKUs para mostrar.")
		return nil
	}

	printTable(items)
	return nil
}

func printTable(items []api.OrdenItem) {
	fmt.Printf("%-20s %-12s %-35s %8s %8s\n",
		"SEMÁFORO", "SKU", "PRODUCTO", "STOCK", "A PEDIR")
	fmt.Println("─────────────────────────────────────────────────────────────────────────────────────")

	for _, item := range items {
		nombre := item.NombreProducto
		if len(nombre) > 35 {
			nombre = nombre[:32] + "..."
		}
		fmt.Printf("%-20s %-12s %-35s %8d %8d\n",
			item.Semaforo,
			item.SKU,
			nombre,
			item.StockActual,
			item.UnidadesAPedir,
		)
	}

	fmt.Printf("\nTotal SKUs: %d\n", len(items))
}
