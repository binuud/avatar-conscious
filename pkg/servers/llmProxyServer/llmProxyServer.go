package llmProxyServer

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/binuud/conscious/pkg/utils"
	"github.com/sirupsen/logrus"
)

var (
	LLM_SAFE_PATHS = []string{"/engines/v1/models", "/engines/v1/chat/completions", "/engines/v1/embeddings",
		"/v1/models", "/v1/chat/completions", "/v1/embeddings", "/v1/images/generations"}
)

func isPathAllowed(needle string) bool {

	// iterate using the for loop
	for _, item := range LLM_SAFE_PATHS {
		if item == needle {
			return true
		}
	}
	return false
}

// GetLlmProxyHandler()
// sanitizes the request path against list of safe paths
// and proxies the llm request, by sending same to LLM_HOST
func GetLlmProxyHandler() *http.HandlerFunc {

	LLM_HOST := utils.GetEnv("LLM_HOST", "http://localhost:11434")
	logrus.Infof("Handling LLM Proxy for : %s", LLM_HOST)

	targetURLSanitized, err := url.Parse(LLM_HOST)
	if err != nil {
		log.Fatalf("Error parsing target URL: %v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURLSanitized)

	coreHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		logrus.Infof("Handling llm proxy for path %s", r.URL.RequestURI())
		// Step A: Set CORS headers globally for all matching traffic
		w.Header().Set("Access-Control-Allow-Origin", "*") // Restrict in production!
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		// Step B: Intercept browser OPTIONS requests early
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.WriteHeader(http.StatusAccepted) // 204 No Content
			return
		}

		// Step C: Validate the path against the whitelist
		// r.URL.Path has already had "/myproxy" stripped off at this stage
		if !isPathAllowed(r.URL.Path) {
			http.Error(w, "Forbidden: Path not allowed by proxy", http.StatusForbidden)
			return
		}

		// Step D: Validate the HTTP method before forwarding
		switch r.Method {
		case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete:
			proxy.ServeHTTP(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	return &coreHandler

}
