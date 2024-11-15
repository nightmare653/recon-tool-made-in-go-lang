package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow requests from any origin
	},
}

func main() {
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/run", handleRunTask)

	log.Println("Server started at http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func handleRunTask(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		return
	}
	defer conn.Close()

	var msg struct {
		Task     string `json:"task"`
		Domain   string `json:"domain"`
	}

	_, msgBytes, err := conn.ReadMessage()
	if err != nil {
		log.Println("Failed to read message:", err)
		conn.WriteMessage(websocket.TextMessage, []byte("Error: Failed to read message"))
		return
	}

	err = json.Unmarshal(msgBytes, &msg)
	if err != nil {
		log.Println("Failed to unmarshal message:", err)
		conn.WriteMessage(websocket.TextMessage, []byte("Error: Invalid message format"))
		return
	}

	task := strings.TrimSpace(msg.Task) // Trim spaces from task
	domain := msg.Domain

	log.Printf("Received task: %s for domain: %s", task, domain)

	if domain == "" {
		conn.WriteMessage(websocket.TextMessage, []byte("Error: Domain is required"))
		return
	}

	// Create directory for domain
	dirPath := createDomainDirectory(domain)
	if dirPath == "" {
		conn.WriteMessage(websocket.TextMessage, []byte("Error: Could not create directory"))
		return
	}

	command, args := getCommand(task, domain)
	if command == "" {
		conn.WriteMessage(websocket.TextMessage, []byte("Error: Unknown task"))
		return
	}

	runCommand(conn, command, args, dirPath, task)
}

func createDomainDirectory(domain string) string {
	dirPath := filepath.Join(".", "results", domain)
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		log.Printf("Error creating directory %s: %v", dirPath, err)
		return ""
	}
	log.Printf("Created directory: %s", dirPath)
	return dirPath
}

func getCommand(task, domain string) (string, []string) {
	switch task {
	case "subfinder":
		return "subfinder", []string{"-d", domain, "-silent"}
	case "live_domains":
		return "aquatone", []string{"-out", "aquatone", "-scan-timeout", "500"}
	case "waybackurls":
		return "waybackurls", []string{domain}
	case "extract_ips":
		return "dig", []string{"+short", domain}
	case "whois":
		return "whois", []string{domain}
	case "waf_detection":
		return "wafw00f", []string{domain}
	case "arjun":
		return "arjun", []string{"-u", fmt.Sprintf("https://%s", domain), "--rate-limit", "10", "-t", "5"}
	case "katana":
		return "katana", []string{"-u", fmt.Sprintf("https://%s", domain), "-silent"}
	case "hakrawler":
		return "hakrawler", []string{"-url", fmt.Sprintf("https://%s", domain)}
	case "dirb":
		return "dirb", []string{fmt.Sprintf("http://%s", domain)}
	case "dalfox":
		return "dalfox", []string{"url", fmt.Sprintf("https://%s", domain)}
	case "nuclei":
		return "nuclei", []string{"-u", fmt.Sprintf("https://%s", domain), "-t", "cves/", "-silent"}
	case "amass":
		return "amass", []string{"enum", "-d", domain}
	case "bucket_finder":
		return "bucket_finder", []string{"-d", domain}
	case "cloudflair":
		return "CloudFlair", []string{domain}
	case "commix":
		return "commix", []string{"-u", fmt.Sprintf("https://%s", domain)}
	case "dirsearch":
		return "dirsearch", []string{"-u", fmt.Sprintf("http://%s", domain)}
	case "dnsenum":
		return "dnsenum", []string{domain}
	case "dnsrecon":
		return "dnsrecon", []string{"-d", domain}
	case "dotdotpwn":
		return "dotdotpwn", []string{"-m", "http", "-u", fmt.Sprintf("http://%s", domain)}
	case "fierce":
		return "fierce", []string{"-domain", domain}
	case "gobuster":
		return "gobuster", []string{"dir", "-u", fmt.Sprintf("http://%s", domain)}
	case "joomscan":
		return "joomscan", []string{"-u", fmt.Sprintf("https://%s", domain)}
	case "knockpy":
		return "knockpy", []string{"-d", domain}
	case "masscan":
		return "masscan", []string{"-p", "80,443", domain}
	case "massdns":
		return "massdns", []string{"-r", "/path/to/resolvers.txt", "-t", "A", domain}
	case "nikto":
		return "nikto", []string{"-h", fmt.Sprintf("http://%s", domain)}
	case "nmap":
		return "nmap", []string{"-sV", domain}
	case "recon-ng":
		return "recon-ng", []string{"-m", "recon/domains-hosts/brute_hosts", domain}
	case "s3recon":
		return "s3recon", []string{domain}
	case "sqlmap":
		return "sqlmap", []string{"-u", fmt.Sprintf("https://%s", domain), "--batch"}
	case "sublist3r":
		return "sublist3r", []string{"-d", domain}
	case "teh_s3_bucketeers":
		return "teh_s3_bucketeers", []string{domain}
	case "thc-hydra":
		return "hydra", []string{"-l", domain, "-P", "/usr/share/wordlists/rockyou.txt", "ssh://localhost"}
	case "theHarvester":
		return "theHarvester", []string{"-d", domain, "-b", "google"}
	case "virtual-host-discovery":
		return "virtual-host-discovery", []string{"-t", domain}
	case "wfuzz":
		return "wfuzz", []string{"-u", fmt.Sprintf("http://%s/FUZZ", domain)}
	case "whatweb":
		return "whatweb", []string{domain}
	case "wpscan":
		return "wpscan", []string{"--url", fmt.Sprintf("https://%s", domain)}
	case "xsstrike":
		return "XSStrike", []string{"-u", fmt.Sprintf("https://%s", domain)}
	default:
		log.Printf("Unknown task received: %s", task) // Log unknown tasks
		return "", nil
	}
}

func runCommand(conn *websocket.Conn, command string, args []string, dirPath, task string) {
	cmd := exec.Command(command, args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Error creating stdout pipe: %v", err)))
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Error creating stderr pipe: %v", err)))
		return
	}

	outputFilePath := filepath.Join(dirPath, fmt.Sprintf("%s_output.txt", task))
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Error creating output file: %v", err)))
		return
	}
	defer outputFile.Close()

	if err := cmd.Start(); err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Error starting command: %v", err)))
		return
	}

	go streamOutput(conn, stdout, outputFile)
	go streamOutput(conn, stderr, outputFile)

	if err := cmd.Wait(); err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Command execution failed: %v", err)))
	} else {
		conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Task %s completed. Results saved in %s", task, outputFilePath)))
	}
}

func streamOutput(conn *websocket.Conn, pipe io.ReadCloser, outputFile *os.File) {
	defer pipe.Close()

	reader := bufio.NewReader(pipe)
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()
		outputFile.WriteString(line + "\n")
		conn.WriteMessage(websocket.TextMessage, []byte(line))
	}
	if err := scanner.Err(); err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Error reading output: %v", err)))
	}
}
