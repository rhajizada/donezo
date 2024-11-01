package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/rhajizada/donezo/internal/auth"
	"github.com/rhajizada/donezo/internal/config"
)

func main() {
	configPath := flag.String("config", "/etc/donezo/config.yaml", "Path to configuration file")
	baseUrl := flag.String("url", "localhost", "Server base URL, e.g. \"localhost\"")
	expiration := flag.Duration("expiration", 24*time.Hour, "Token duration")
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadServerConfig(*configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	token, err := auth.GenerateToken(cfg.JWT.Secret, *expiration)
	if err != nil {
		log.Fatalf("Failed to generate JWT token: %v", err)
	}
	message := fmt.Sprintf(`
cat <<EOF > ~/.config/donezo/config.yaml
baseURL: http://%s:%d
apiToken: %s
duration: 2s
EOF
  `, *baseUrl, cfg.Port, token)
	fmt.Print(message)
}
