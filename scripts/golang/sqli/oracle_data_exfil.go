package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil" // Added for dumping requests
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Config holds the command-line parameters and parsed request data
type Config struct {
	RequestFile     string
	WordlistFile    string
	OutputFile      string
	PlaceholderChar string // For the character from wordlist
	PlaceholderPos  string // For the character position
	SuccessCode     int
	Threads         int
	Verbose         bool
	TimeoutSeconds  int
	ExfilLength     int // New: Length of the string to exfiltrate
	// Parsed from request file
	BaseMethod string
	BaseHost   string
	BasePath   string
	BaseScheme string
	RawHeaders string
	RawBody    string
}

var httpClient *http.Client

const maxInitialAttempts = 3 // Number of attempts before a major pause and final try

// sleepWithContext sleeps for the given duration or until the context is done.
// Returns true if the sleep completed, false if context was done.
func sleepWithContext(ctx context.Context, duration time.Duration) bool {
	select {
	case <-time.After(duration):
		return true // Sleep completed
	case <-ctx.Done():
		return false // Context was cancelled during sleep
	}
}

func main() {
	config := Config{}

	flag.StringVar(&config.RequestFile, "r", "", "Path to the raw HTTP request file (required)")
	flag.StringVar(&config.WordlistFile, "w", "", "Path to the wordlist file (e.g., a-zA-Z0-9) (required)")
	flag.StringVar(&config.OutputFile, "o", "", "Output file to append successful characters (optional)")
	flag.StringVar(&config.PlaceholderChar, "pc", "§PAYLOAD§", "Placeholder for the character from wordlist")
	flag.StringVar(&config.PlaceholderPos, "pp", "§CHAR_POS§", "Placeholder for the character position (1-indexed)")
	flag.IntVar(&config.SuccessCode, "s", 0, "HTTP status code that indicates success (required)")
	flag.IntVar(&config.Threads, "threads", 1, "Number of concurrent goroutines per character position")
	flag.BoolVar(&config.Verbose, "v", false, "Verbose output")
	flag.IntVar(&config.TimeoutSeconds, "timeout", 10, "HTTP request timeout in seconds")
	flag.IntVar(&config.ExfilLength, "length", 0, "Total number of characters to exfiltrate (required for exfil mode)")
	flag.Parse()

	fmt.Printf("%s\n\n", time.Now().Format("2006-01-02 15:04:05"))

	if config.RequestFile == "" || config.WordlistFile == "" || config.SuccessCode == 0 || config.ExfilLength <= 0 {
		fmt.Println("Usage: go_exfiltrator -r <req_file> -w <wordlist> -s <success_code> -length <num_chars> [options]")
		flag.PrintDefaults()
		os.Exit(1)
	}
	if config.Threads < 1 {
		log.Fatal("--threads must be at least 1.")
	}

	var outputFileHandle *os.File
	if config.OutputFile != "" {
		var err error
		outputFileHandle, err = os.OpenFile(config.OutputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Failed to open output file '%s': %v", config.OutputFile, err)
		}
		defer outputFileHandle.Close()
		if config.Verbose {
			log.Printf("Logging successful characters to: %s", config.OutputFile)
		}
	}

	httpClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
			DisableKeepAlives: true,
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	err := parseRequestTemplate(&config)
	if err != nil {
		log.Fatalf("Error parsing request file '%s': %v", config.RequestFile, err)
	}
	if config.Verbose {
		log.Printf("Base request parsed. Target: %s://%s%s.", config.BaseScheme, config.BaseHost, config.BasePath)
		log.Printf("Character Placeholder: '%s', Position Placeholder: '%s'", config.PlaceholderChar, config.PlaceholderPos)
		log.Printf("Success condition: HTTP Status Code %d", config.SuccessCode)
		log.Printf("Exfiltrating %d characters.", config.ExfilLength)
		log.Printf("Using %d threads per character position.", config.Threads)
		log.Printf("Retry logic: Up to %d initial attempts, then 10s pause, then 1 final attempt before script exit on persistent failure for a character.", maxInitialAttempts)
	}

	var exfiltratedString strings.Builder

	for charPos := 1; charPos <= config.ExfilLength; charPos++ {
		if config.Verbose {
			log.Printf("\nAttempting to find character at position %d...", charPos)
		}

		payloadsChan := make(chan string, config.Threads*2)
		successCharChan := make(chan string, 1)
		var wg sync.WaitGroup
		posCtx, posCancel := context.WithCancel(context.Background())

		for i := 0; i < config.Threads; i++ {
			wg.Add(1)
			go worker(posCtx, &config, charPos, payloadsChan, successCharChan, &wg)
		}

		wordlistFile, err := os.Open(config.WordlistFile)
		if err != nil {
			posCancel()
			log.Fatalf("Error re-opening wordlist file for position %d: %v", charPos, err)
		}

		go func() {
			defer close(payloadsChan)
			defer wordlistFile.Close()

			scanner := bufio.NewScanner(wordlistFile)
			for scanner.Scan() {
				select {
				case <-posCtx.Done():
					return
				case payloadsChan <- scanner.Text():
					// Payload sent
				}
			}
			if err := scanner.Err(); err != nil {
				select {
				case <-posCtx.Done():
				default:
					log.Printf("Error reading wordlist for position %d: %v", charPos, err)
				}
			}
		}()

		foundCharForPos := ""
		overallCharPosTimeoutDuration := time.Duration(config.TimeoutSeconds*(maxInitialAttempts+1)*5+15) * time.Second // Heuristic

		select {
		case foundChar := <-successCharChan:
			foundCharForPos = foundChar
		case <-time.After(overallCharPosTimeoutDuration):
			if config.Verbose {
				log.Printf("Overall timeout waiting for character at position %d.", charPos)
			}
		}

		posCancel()
		wg.Wait()
		close(successCharChan)

		if foundCharForPos != "" {
			fmt.Printf("%s", foundCharForPos)
			exfiltratedString.WriteString(foundCharForPos)
			if outputFileHandle != nil {
				if _, err := outputFileHandle.WriteString(foundCharForPos); err != nil {
					log.Printf("Error writing char '%s' to output file: %v", foundCharForPos, err)
				}
			}
		} else {
			log.Printf("\nFailed to determine character at position %d after all attempts or timeout.", charPos)
			log.Printf("Current exfiltrated string: %s", exfiltratedString.String())
			if outputFileHandle != nil {
				outputFileHandle.WriteString("\n")
			}
			os.Exit(1)
		}
	}

	finalResult := exfiltratedString.String()
	fmt.Println("\n\nExfiltration complete.")
	log.Printf("Final exfiltrated string: %s", finalResult)
	if outputFileHandle != nil && finalResult != "" {
		outputFileHandle.WriteString("\n")
	}
	if config.Verbose {
		log.Println("Intruder run finished.")
	}
}

func parseRequestTemplate(config *Config) error {
	fileContent, err := os.ReadFile(config.RequestFile)
	if err != nil {
		return err
	}
	contentStr := string(fileContent)
	contentStr = strings.ReplaceAll(contentStr, "\r\n", "\n")
	contentStr = strings.TrimRight(contentStr, "\n")
	lines := strings.Split(contentStr, "\n")
	if len(lines) == 0 {
		return fmt.Errorf("request file is empty")
	}
	firstLineParts := strings.SplitN(lines[0], " ", 3)
	if len(firstLineParts) < 2 {
		return fmt.Errorf("invalid request line: %s", lines[0])
	}
	config.BaseMethod = strings.TrimSpace(firstLineParts[0])
	config.BasePath = strings.TrimSpace(firstLineParts[1])
	headerEndIndex := -1
	for i := 1; i < len(lines); i++ {
		if lines[i] == "" {
			headerEndIndex = i
			break
		}
	}
	if headerEndIndex == -1 {
		headerEndIndex = len(lines)
	}
	config.RawHeaders = strings.Join(lines[1:headerEndIndex], "\n")
	tempReqReader := bufio.NewReader(strings.NewReader(config.RawHeaders + "\n\n"))
	parsedHeadersReq, err := http.ReadRequest(tempReqReader)
	if err == nil && parsedHeadersReq.Host != "" {
		config.BaseHost = parsedHeadersReq.Host
	} else {
		re := hostHeaderRegex
		match := re.FindStringSubmatch(config.RawHeaders)
		if len(match) > 1 {
			config.BaseHost = strings.TrimSpace(match[1])
		}
	}
	if config.BaseHost == "" {
		return fmt.Errorf("could not parse Host header")
	}
	config.BaseScheme = "https"
	if headerEndIndex < len(lines) && headerEndIndex+1 <= len(lines) {
		config.RawBody = strings.Join(lines[headerEndIndex+1:], "\n")
	} else {
		config.RawBody = ""
	}
	return nil
}

func worker(ctx context.Context, config *Config, charPos int, payloadsChan <-chan string, successCharChan chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	charPosStr := strconv.Itoa(charPos)

	for payloadChar := range payloadsChan {
		// ** ADDED CHECK TO SKIP PAYLOAD IF IT'S A PLACEHOLDER ITSELF **
		if payloadChar == config.PlaceholderChar || payloadChar == config.PlaceholderPos {
			if config.Verbose {
				log.Printf("Pos %d: Skipping payloadChar '%s' because it matches a defined placeholder string.", charPos, payloadChar)
			}
			continue // Skip this iteration, get next payloadChar
		}

		select {
		case <-ctx.Done():
			return
		default:
		}

		var finalAttemptAfterPauseMade = false
		var requestSuccessfulForPayloadChar = false

		for currentAttempt := 1; ; currentAttempt++ {
			select {
			case <-ctx.Done():
				return
			default:
			}

			tempPath := strings.ReplaceAll(config.BasePath, config.PlaceholderPos, charPosStr)
			currentPath := strings.ReplaceAll(tempPath, config.PlaceholderChar, payloadChar)
			tempHeadersStr := strings.ReplaceAll(config.RawHeaders, config.PlaceholderPos, charPosStr)
			currentHeadersStr := strings.ReplaceAll(tempHeadersStr, config.PlaceholderChar, payloadChar)
			tempBodyStr := strings.ReplaceAll(config.RawBody, config.PlaceholderPos, charPosStr)
			currentBodyStr := strings.ReplaceAll(tempBodyStr, config.PlaceholderChar, payloadChar)
			currentBodyReader := strings.NewReader(currentBodyStr)
			targetURL := fmt.Sprintf("%s://%s%s", config.BaseScheme, config.BaseHost, currentPath)

			attemptCtx, attemptCancel := context.WithTimeout(ctx, time.Duration(config.TimeoutSeconds)*time.Second)

			req, err := http.NewRequestWithContext(attemptCtx, config.BaseMethod, targetURL, currentBodyReader)
			if err != nil {
				if config.Verbose {
					log.Printf("Pos %d, Char '%s': Error creating request (attempt %d): %v", charPos, payloadChar, currentAttempt, err)
				}
				attemptCancel()
				if currentAttempt < maxInitialAttempts {
					if !sleepWithContext(ctx, 1*time.Second) {
						return
					}
					continue
				}
				if !finalAttemptAfterPauseMade {
					log.Printf("Pos %d, Char '%s' (req creation err): Failed %d initial attempts. Pausing 10s.", charPos, payloadChar, maxInitialAttempts)
					if !sleepWithContext(ctx, 10*time.Second) {
						return
					}
					finalAttemptAfterPauseMade = true
					log.Printf("Pos %d, Char '%s': Resuming final attempt (req creation err).", charPos, payloadChar)
					continue
				}
				log.Fatalf("CRITICAL: Pos %d, Char '%s' (req creation err): Failed final attempt. Exiting. Last error: %v", charPos, payloadChar, err)
			}

			headerScanner := bufio.NewScanner(strings.NewReader(currentHeadersStr))
			for headerScanner.Scan() {
				line := headerScanner.Text()
				if line == "" {
					continue
				}
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					hName, hVal := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
					if strings.ToLower(hName) != "host" && strings.ToLower(hName) != "content-length" {
						req.Header.Add(hName, hVal)
					}
				} else if config.Verbose && strings.TrimSpace(line) != "" {
					log.Printf("Warning: Pos %d, Char '%s': Malformed header line (attempt %d): \"%s\"", charPos, payloadChar, currentAttempt, line)
				}
			}

			if config.Verbose {
				if currentAttempt == 1 && !finalAttemptAfterPauseMade {
					log.Printf("Pos %d: Attempting char '%s' (URL: %s)", charPos, payloadChar, targetURL)
				} else {
					log.Printf("Pos %d: Retrying char '%s' (Attempt: %d of %d / Final: %t, URL: %s)",
						charPos, payloadChar, currentAttempt, maxInitialAttempts, finalAttemptAfterPauseMade, targetURL)
				}
				dump, dumpErr := httputil.DumpRequestOut(req, true)
				if dumpErr != nil {
					log.Printf("Error dumping request: %v", dumpErr)
				} else {
					log.Printf("---\n%s\n+++", string(dump))
				}
			}

			resp, err := httpClient.Do(req)
			if err != nil {
				if err == context.Canceled || err == context.DeadlineExceeded {
					if ctx.Err() == context.Canceled {
						attemptCancel()
						return
					}
				}

				if config.Verbose {
					log.Printf("Pos %d, Char '%s': Attempt %d HTTP request failed with status N/A: %v", charPos, payloadChar, currentAttempt, err)
				}
				attemptCancel()
				if currentAttempt < maxInitialAttempts {
					if !sleepWithContext(ctx, 1*time.Second) {
						return
					}
					continue
				}
				if !finalAttemptAfterPauseMade {
					log.Printf("Pos %d, Char '%s': Failed %d initial attempts (HTTP error). Pausing this worker 10s.", charPos, payloadChar, maxInitialAttempts)
					if !sleepWithContext(ctx, 10*time.Second) {
						return
					}
					finalAttemptAfterPauseMade = true
					log.Printf("Pos %d, Char '%s': Resuming. Final attempt (HTTP error).", charPos, payloadChar)
					continue
				}
				log.Fatalf("CRITICAL: Pos %d, Char '%s': Failed final attempt (HTTP error). Exiting. Last error: %v", charPos, payloadChar, err)
			}

			if config.Verbose {
				log.Printf("Pos %d, Char '%s': Attempt %d received status %d", charPos, payloadChar, currentAttempt, resp.StatusCode)
			}

			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
			attemptCancel()

			if resp.StatusCode == config.SuccessCode {
				select {
				case successCharChan <- payloadChar:
					if config.Verbose {
						log.Printf(">>> SUCCESS: Pos %d, Char '%s' confirmed with status %d <<<", charPos, payloadChar, resp.StatusCode)
					}
				case <-ctx.Done():
					if config.Verbose {
						log.Printf("Pos %d, Char '%s': Context cancelled while trying to report success.", charPos, payloadChar)
					}
				}
				requestSuccessfulForPayloadChar = true
			}
			break
		}

		if requestSuccessfulForPayloadChar {
			return
		}
	}
}

var hostHeaderRegex = regexp.MustCompile(`(?im)^Host:\s*(.*)$`)

func init() {
	log.SetFlags(0)
}
