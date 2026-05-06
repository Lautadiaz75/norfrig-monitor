package cmd

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/Lautadiaz75/norfrig-monitor/internal/api"
	gh "github.com/Lautadiaz75/norfrig-monitor/internal/github"
	"github.com/spf13/cobra"
)

// Checker es cualquier cosa que pueda verificar su estado y devolver un resumen.
// PipelineChecker y StockChecker implementan esta interfaz sin declararlo explícitamente.
type Checker interface {
	Nombre() string
	Verificar() (string, error)
}

type checkResultado struct {
	nombre string
	texto  string
	err    error
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Estado general del sistema: pipeline + stock en paralelo",
	RunE:  runCheck,
}

func init() {
	rootCmd.AddCommand(checkCmd)
}

func runCheck(cmd *cobra.Command, args []string) error {
	checkers := []Checker{
		&PipelineChecker{token: os.Getenv("GITHUB_TOKEN")},
		&StockChecker{apiURL: os.Getenv("NORFRIG_API_URL")},
	}

	// El canal recibe un resultado por cada goroutine.
	// El tamaño del buffer (len(checkers)) evita que las goroutines se bloqueen.
	resultados := make(chan checkResultado, len(checkers))

	var wg sync.WaitGroup
	for _, c := range checkers {
		wg.Add(1)
		go func(c Checker) {
			defer wg.Done()
			texto, err := c.Verificar()
			resultados <- checkResultado{nombre: c.Nombre(), texto: texto, err: err}
		}(c)
	}

	// Cerramos el canal cuando todas las goroutines terminaron.
	// Esto habilita el range de abajo.
	go func() {
		wg.Wait()
		close(resultados)
	}()

	fmt.Println("NORFRIG BI — Estado General")
	fmt.Println("──────────────────────────────────────────")
	for r := range resultados {
		if r.err != nil {
			fmt.Printf("  %-12s ERROR: %v\n", r.nombre, r.err)
		} else {
			fmt.Printf("  %-12s %s\n", r.nombre, r.texto)
		}
	}
	return nil
}

// ─── PipelineChecker ──────────────────────────────────────────────────────────

type PipelineChecker struct {
	token string
}

func (p *PipelineChecker) Nombre() string { return "Pipeline" }

func (p *PipelineChecker) Verificar() (string, error) {
	if p.token == "" {
		return "", fmt.Errorf("GITHUB_TOKEN no configurado")
	}
	client := gh.NewClient(p.token)
	run, err := client.GetLatestRun("Lautadiaz75", "norfrig-bi-v2")
	if err != nil {
		return "", err
	}
	duracion := run.UpdatedAt.Sub(run.CreatedAt).Round(time.Second)
	hace := time.Since(run.CreatedAt).Round(time.Minute)
	estado := formatConclusion(run.Conclusion, run.Status)
	return fmt.Sprintf("%s · hace %s · duración %s",
		estado, formatDuracion(hace), formatDuracion(duracion)), nil
}

// ─── StockChecker ─────────────────────────────────────────────────────────────

type StockChecker struct {
	apiURL string
}

func (s *StockChecker) Nombre() string { return "Stock" }

func (s *StockChecker) Verificar() (string, error) {
	if s.apiURL == "" {
		return "", fmt.Errorf("NORFRIG_API_URL no configurado")
	}
	client := api.NewClient(s.apiURL)
	items, err := client.GetOrden("TODOS")
	if err != nil {
		return "", err
	}

	criticos, reorden := 0, 0
	for _, item := range items {
		switch item.Semaforo {
		case "CRITICO":
			criticos++
		case "REORDEN":
			reorden++
		}
	}
	return fmt.Sprintf("%d SKUs · %d CRITICOS · %d REORDEN",
		len(items), criticos, reorden), nil
}
