# WebSocket Task Runner

A simple web application to execute various reconnaissance and penetration testing tasks using WebSocket. This tool allows users to select tasks, input a domain, and view results in a terminal-like output interface.

---

## Features
- **Task Selection**: Choose from a wide range of reconnaissance, scanning, and exploitation tools.
- **Domain Input**: Input a domain for task execution.
- **Real-time Output**: Results are streamed live and displayed in a terminal-style output box.
- **WebSocket Integration**: Ensures seamless communication between client and server.

---

## Prerequisites
- Go installed.
- A WebSocket server running on `ws://localhost:8080`.

---

## Installation

1. Clone the repository:

   git clone https://github.com/your-username/websocket-task-runner.git

2. Navigate to the project directory:

  cd websocket-task-runner

3. Run your WebSocket server:

  go run server.go

Usage
Open index.html in your browser.

Select a task from the dropdown.

Enter the target domain.

Click Run to start the task.

View the results in the terminal-style output box.

Supported Tasks

Subdomain Enumeration: subfinder, amass, sublist3r, etc.

URL Discovery: waybackurls, hakrawler, katana, etc.

Web Scanning: dirsearch, gobuster, nikto, etc.

Vulnerability Scanning: nmap, nuclei, sqlmap, etc.

CMS Scanning: wpscan, joomscan, etc.

Cloud and S3 Tools: bucket_finder, s3recon, etc.

Technologies Used

HTML/CSS: User interface.

JavaScript: WebSocket communication.

Go: WebSocket server for task execution.
