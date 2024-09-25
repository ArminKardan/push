package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"push/lib"
	"strings"
	"time"
)

// Struct for the SSH key payload
type SSHKeyPayload struct {
	Title string `json:"title"`
	Key   string `json:"key"`
}

// Struct for the repository creation payload
type RepoPayload struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Private     bool   `json:"private,omitempty"`
}

// GitHub API URL for adding SSH keys and creating repos
const (
	githubAPIURL = "https://api.github.com/user/keys"
	repoAPIURL   = "https://api.github.com/user/repos"
)

type GitHubUser struct {
	Login string `json:"login"`
}

func getGitHubUsername(token string) (string, error) {
	url := "https://api.github.com/user"

	// Create a new request with the token
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	// Set the Authorization header with the token
	req.Header.Set("Authorization", "token "+token)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check for a successful response
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get user info: %s", resp.Status)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Parse the JSON response
	var user GitHubUser
	if err := json.Unmarshal(body, &user); err != nil {
		return "", err
	}

	return user.Login, nil
}

// Add SSH key to GitHub account
func addSSHKey(githubToken, title, sshKey string) error {
	payload := SSHKeyPayload{
		Title: title,
		Key:   sshKey,
	}

	// Convert the payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", githubAPIURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("Authorization", "token "+githubToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Content-Type", "application/json")

	// Make HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body for error details
	responseBody, _ := io.ReadAll(resp.Body)

	// Check if the request was successful
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to add SSH key, status code: %d, response: %s", resp.StatusCode, string(responseBody))
	}

	fmt.Println("SSH key added successfully!")
	return nil
}

// Create a new repository on GitHub
func createRepository(githubToken, name, description string, isPrivate bool) error {
	payload := RepoPayload{
		Name:        name,
		Description: description,
		Private:     isPrivate,
	}

	// Convert the payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal repository payload: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", repoAPIURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create request for repository: %v", err)
	}

	// Set headers
	req.Header.Set("Authorization", "token "+githubToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Content-Type", "application/json")

	// Make HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request to create repository: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body for error details
	responseBody, _ := io.ReadAll(resp.Body)

	// Check if the request was successful
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to create repository, status code: %d, response: %s", resp.StatusCode, string(responseBody))
	}

	fmt.Println("Repository created successfully!")
	return nil
}

// Get the default SSH key location for the current OS
func getSSHKeyPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting user home directory:", err)
		os.Exit(1)
	}

	return filepath.Join(homeDir, ".ssh", "id_rsa.pub")
}

// Check if the SSH key exists, and generate one if it doesn't
func getOrGenerateSSHKey() (string, error) {
	sshKeyPath := getSSHKeyPath()

	// Check if the SSH key exists
	if _, err := os.Stat(sshKeyPath); os.IsNotExist(err) {
		fmt.Println("SSH key not found, generating a new one...")

		// Generate SSH key using ssh-keygen
		err := generateSSHKey()
		if err != nil {
			return "", fmt.Errorf("failed to generate SSH key: %v", err)
		}
	}

	// Read the SSH public key file
	sshKey, err := os.ReadFile(sshKeyPath)
	if err != nil {
		return "", fmt.Errorf("failed to read SSH key: %v", err)
	}

	return string(sshKey), nil
}

// Generate SSH key using ssh-keygen
func generateSSHKey() error {
	sshDir := filepath.Dir(getSSHKeyPath())

	// Create .ssh directory if it doesn't exist
	if _, err := os.Stat(sshDir); os.IsNotExist(err) {
		err := os.MkdirAll(sshDir, 0700)
		if err != nil {
			return fmt.Errorf("failed to create .ssh directory: %v", err)
		}
	}

	// Run ssh-keygen command to generate the key
	cmd := exec.Command("ssh-keygen", "-t", "rsa", "-b", "4096", "-f", filepath.Join(sshDir, "id_rsa"), "-N", "")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("ssh-keygen failed: %v", err)
	}

	return nil
}

func UpdateGitignore(gitignorePath string, linesToAdd []string) error {
	// Check if the .gitignore file exists
	if _, err := os.Stat(gitignorePath); os.IsNotExist(err) {
		// Create the .gitignore file if it doesn't exist
		file, err := os.Create(gitignorePath)
		if err != nil {
			return fmt.Errorf("error creating .gitignore file: %w", err)
		}
		defer file.Close()

		// Write the lines to the new .gitignore file
		for _, line := range linesToAdd {
			file.WriteString(line + "\n")
		}

		fmt.Println(".gitignore file created with default content.")
	} else {
		// If it exists, check for each line
		file, err := os.Open(gitignorePath)
		if err != nil {
			return fmt.Errorf("error opening .gitignore file: %w", err)
		}
		defer file.Close()

		// Read existing lines into a map for easy checking
		existingLines := make(map[string]struct{})
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			existingLines[scanner.Text()] = struct{}{}
		}

		// Check for missing lines and add them
		needsUpdate := false
		for _, line := range linesToAdd {
			if _, exists := existingLines[line]; !exists {
				needsUpdate = true
				break
			}
		}

		if needsUpdate {
			// Append missing lines to the .gitignore file
			file, err := os.OpenFile(gitignorePath, os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				return fmt.Errorf("error opening .gitignore file for appending: %w", err)
			}
			defer file.Close()

			for _, line := range linesToAdd {
				if _, exists := existingLines[line]; !exists {
					file.WriteString(line + "\n")
					fmt.Printf("Added missing line: %s\n", line)
				}
			}
		} else {
			fmt.Println(".gitignore file is already up to date.")
		}
	}
	return nil
}

func RunCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)

	// Set the output to standard output and standard error
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running command: %w", err)
	}

	return nil
}

func GetUsernameFromRemote(repoPath string) (string, error) {
	cmd := exec.Command("git", "-C", repoPath, "remote", "get-url", "origin")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to get remote origin: %s", out.String())
	}

	remoteURL := strings.TrimSpace(out.String())
	// Parse the URL
	parsedURL, err := url.Parse(remoteURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %w", err)
	}

	// Extract the path and split by '/'
	parts := strings.Split(parsedURL.Path, "/")
	if len(parts) >= 3 {
		return parts[1], nil // The username is usually the second part of the path
	}

	return "", fmt.Errorf("invalid repository URL format")
}

func main() {

	binaryName := filepath.Base(os.Args[0])
	if binaryName != "push.exe" {
		fmt.Println("File name should be: push.exe but it is:", binaryName)
		time.Sleep(5 * time.Second)
		os.Exit(4)
	}
	githubToken := ""
	reader := bufio.NewReader(os.Stdin)
	token, err := lib.SecureRead()
	if err != nil {
		fmt.Println("Please enter your github development token:")
		input, err := reader.ReadString('\n')
		if err == nil {
			input = strings.TrimSpace(input)
			githubToken = input
			lib.SecureWrite(githubToken)
		} else {
			os.Exit(3)
		}
	} else {
		githubToken = token
	}

	fmt.Println("Token found:" + githubToken)

	time.Sleep(3 * time.Second)

	// Replace with your actual GitHub token

	// Check if the token is empty
	if githubToken == "" {
		fmt.Println("GitHub token is required!")
		os.Exit(1)
	}

	// Get or generate the SSH key
	sshKey, err := getOrGenerateSSHKey()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// Add SSH key to GitHub
	addSSHKey(githubToken, "rgsshkey", sshKey)
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	os.Exit(1)
	// }

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}

	// Create a new repository
	repoName := filepath.Base(currentDir)
	repoDescription := "Turing research group."
	isPrivate := false // Set to true for a private repository
	if repoName[0] == '_' {
		isPrivate = true
	}

	createRepository(githubToken, repoName, repoDescription, isPrivate)

	gitignorePath := ".gitignore"

	// Define the lines to check in the .gitignore file
	linesToAdd := []string{
		"/push.exe",
		"/publish.exe",
		"/node_modules/*",
		"/.next",
		"/packages",
		"/chrome",
		"/bin",
		"/*/bin",
		"/*/obj",
	}

	// Call the function to update the .gitignore file
	if err := UpdateGitignore(gitignorePath, linesToAdd); err != nil {
		fmt.Println("Error:", err)
	}

	repoDir := "./" // Change this to your desired path

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(repoDir, 0755); err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}

	// Change to the specified directory
	if err := os.Chdir(repoDir); err != nil {
		fmt.Println("Error changing directory:", err)
		return
	}

	gitUsername, _ := getGitHubUsername(githubToken)

	fmt.Println("github username:", gitUsername)

	repousername, _ := GetUsernameFromRemote("./")

	if strings.EqualFold(repousername, gitUsername) {
		os.RemoveAll("./.git")
	}

	RunCommand("git", "init")
	RunCommand("git", "remote", "remove", "origin")
	RunCommand("git", "pull", "--rebase", "origin", "main")
	RunCommand("git", "pull", "--rebase", "origin", "master")
	RunCommand("git", "reset", "origin/main")
	RunCommand("git", "reset", "origin/master")
	RunCommand("git", "remote", "add", "origin", "https://"+githubToken+"@github.com/"+gitUsername+"/"+repoName+".git")
	RunCommand("git", "add", ".")
	RunCommand("git", "commit", "-m", "Update")
	RunCommand("git", "push", "-u", "origin", "master", "--force")
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	os.Exit(1)
	// }
	fmt.Println("Done, exiting...")
	time.Sleep(4 * time.Second)
}
