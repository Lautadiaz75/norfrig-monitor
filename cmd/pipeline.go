package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	gh "github.com/Lautadiaz75/norfrig-monitor/internal/github"
	"github.com/spf13/cobra"
)

var pipelineCmd = &cobra.Command{
	Use:   "pipeline",
	Short: "Muestra el estado del último run de GitHub Actions",
	RunE:  runPipeline,
}

func init() {
	rootCmd.AddCommand(pipelineCmd)
}

func runPipeline(cmd *cobra.Command, args []string) error {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return fmt.Errorf("variable GITHUB_TOKEN no configurada\nCreala en: github.com/settings/tokens")
	}

	client := gh.NewClient(token)
	run, err := client.GetLatestRun("Lautadiaz75", "norfrig-bi-v2")
	if err != nil {
		return fmt.Errorf("error consultando GitHub: %w", err)
	}

	printPipelineStatus(run)
	return nil
}

func printPipelineStatus(run gh.WorkflowRun) {
	duracion := run.UpdatedAt.Sub(run.CreatedAt).Round(time.Second)
	hace := time.Since(run.CreatedAt).Round(time.Minute)

	estado := formatConclusion(run.Conclusion, run.Status)
	commitShort := run.HeadCommit.ID
	if len(commitShort) > 7 {
		commitShort = commitShort[:7]
	}
	commitMsg := run.HeadCommit.Message
	if idx := strings.Index(commitMsg, "\n"); idx != -1 {
		commitMsg = commitMsg[:idx]
	}

	fmt.Println("PIPELINE — norfrig-bi-v2")
	fmt.Println("─────────────────────────────────────────────────────")
	fmt.Printf("Estado:   %s\n", estado)
	fmt.Printf("Iniciado: hace %s\n", formatDuracion(hace))
	fmt.Printf("Duración: %s\n", formatDuracion(duracion))
	fmt.Printf("Commit:   %s · %s\n", commitShort, commitMsg)
}

func formatConclusion(conclusion, status string) string {
	if status == "in_progress" {
		return "en progreso..."
	}
	switch conclusion {
	case "success":
		return "OK"
	case "failure":
		return "FALLO"
	case "cancelled":
		return "cancelado"
	default:
		return conclusion
	}
}

func formatDuracion(d time.Duration) string {
	if d >= time.Hour {
		h := int(d.Hours())
		m := int(d.Minutes()) % 60
		return fmt.Sprintf("%dh %dm", h, m)
	}
	if d >= time.Minute {
		m := int(d.Minutes())
		s := int(d.Seconds()) % 60
		return fmt.Sprintf("%dm %ds", m, s)
	}
	return fmt.Sprintf("%ds", int(d.Seconds()))
}
