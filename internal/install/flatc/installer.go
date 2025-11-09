package flatc

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/machanirobotics/buffman/internal/install/utilities"
)

// flatc specific variables (moved from previous example)
var (
	flatcVersion      = "25.2.10"
	flatcDownloadArch = map[string]string{
		"amd64": "x86_64",
		"arm64": "aarch64",
	}
	flatcDownloadOS = map[string]string{
		"darwin":  "macos",
		"linux":   "linux",
		"windows": "windows",
	}
	flatcDownloadURL = map[string]string{
		"linux-x86_64":    "https://github.com/google/flatbuffers/releases/download/v" + flatcVersion + "/Linux.flatc.binary.g++-13.zip",
		"linux-aarch64":   "https://github.com/google/flatbuffers/releases/download/v" + flatcVersion + "/Linux.flatc.binary.g++-13.zip",
		"macos-x86_64":    "https://github.com/google/flatbuffers/releases/download/v" + flatcVersion + "/MacIntel.flatc.binary.zip",
		"macos-aarch64":   "https://github.com/google/flatbuffers/releases/download/v" + flatcVersion + "/Mac.flatc.binary.zip",
		"windows-x86_64":  "https://github.com/google/flatbuffers/releases/download/v" + flatcVersion + "/Windows.flatc.binary.zip",
		"windows-aarch64": "https://github.com/google/flatbuffers/releases/download/v" + flatcVersion + "/Windows.flatc.binary.zip",
	}
)

// Installer is the interface for buffer installers.
type FlatBuffersInstaller interface {
	Install() error
	Validate() error // To check if the installed tool works correctly (e.g., flatc --version)
	Uninstall() error
	Exists() bool                // To check if the tool is already present
	GetVersion() (string, error) // To retrieve the installed tool's version
}

// FlatbuffersInstaller struct implements the Installer interface for flatc.
type FlatbuffersInstaller struct{}

// NewFlatbuffersInstaller creates and returns a new FlatbuffersInstaller.
func NewFlatbuffersInstaller() FlatBuffersInstaller {
	return &FlatbuffersInstaller{}
}

// Install method for FlatbuffersInstaller.
func (f *FlatbuffersInstaller) Install() error {
	if f.Exists() {
		fmt.Println("flatc is already installed. Skipping installation.")
		return nil
	}

	osType := runtime.GOOS
	arch := runtime.GOARCH

	mappedOS, okOS := flatcDownloadOS[osType]
	mappedArch, okArch := flatcDownloadArch[arch]

	if !okOS || !okArch {
		return fmt.Errorf("unsupported operating system (%s) or architecture (%s) for flatc", osType, arch)
	}

	downloadKey := fmt.Sprintf("%s-%s", mappedOS, mappedArch)
	url, ok := flatcDownloadURL[downloadKey]
	if !ok {
		return fmt.Errorf("no download URL found for flatc %s-%s", mappedOS, mappedArch)
	}

	fmt.Printf("Downloading flatc from: %s\n", url)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download flatc: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download flatc, status code: %d", resp.StatusCode)
	}

	tempDir, err := os.MkdirTemp("", "flatc-install")
	if err != nil {
		return fmt.Errorf("failed to create temporary directory for flatc: %w", err)
	}
	defer os.RemoveAll(tempDir) // Clean up temp directory

	zipFilePath := filepath.Join(tempDir, "flatc.zip")
	out, err := os.Create(zipFilePath)
	if err != nil {
		return fmt.Errorf("failed to create flatc zip file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save downloaded flatc file: %w", err)
	}
	out.Close() // Close the file before opening for unzipping

	fmt.Printf("Extracting %s...\n", zipFilePath)
	err = utilities.Unzip(zipFilePath, tempDir)
	if err != nil {
		return fmt.Errorf("failed to unzip flatc: %w", err)
	}

	var installDir string
	if runtime.GOOS == "windows" {
		installDir = filepath.Join(os.Getenv("APPDATA"), "flatc", "bin")
	} else {
		installDir = "/usr/local/bin"
		if _, err := os.Stat(installDir); os.IsPermission(err) {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("could not determine home directory: %w", err)
			}
			installDir = filepath.Join(homeDir, ".local", "bin")
			if _, err := os.Stat(installDir); os.IsNotExist(err) {
				if err := os.MkdirAll(installDir, 0755); err != nil {
					return fmt.Errorf("failed to create directory %s: %w", installDir, err)
				}
			}
		}
	}

	fmt.Printf("Installing flatc to: %s\n", installDir)

	flatcExecutableName := "flatc"
	if runtime.GOOS == "windows" {
		flatcExecutableName += ".exe"
	}

	extractedFlatcPath := ""
	filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == flatcExecutableName {
			extractedFlatcPath = path
			return nil // Found it, stop walking
		}
		return nil
	})

	if extractedFlatcPath == "" {
		return fmt.Errorf("flatc executable not found in the extracted archive")
	}

	destinationPath := filepath.Join(installDir, flatcExecutableName)

	if _, err := os.Stat(installDir); os.IsNotExist(err) {
		if err := os.MkdirAll(installDir, 0755); err != nil {
			return fmt.Errorf("failed to create installation directory %s: %w", installDir, err)
		}
	}

	if runtime.GOOS != "windows" && strings.HasPrefix(installDir, "/usr/local/bin") {
		fmt.Println("Admin/sudo privileges might be required to install flatc to /usr/local/bin.")
		fmt.Println("Please enter your password if prompted.")
		cmd := exec.Command("sudo", "mv", extractedFlatcPath, destinationPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to move flatc with sudo: %w", err)
		}
	} else {
		err = os.Rename(extractedFlatcPath, destinationPath)
		if err != nil {
			return fmt.Errorf("failed to move flatc executable: %w", err)
		}
	}

	if runtime.GOOS != "windows" {
		err = os.Chmod(destinationPath, 0755)
		if err != nil {
			return fmt.Errorf("failed to set executable permissions for flatc: %w", err)
		}
	}

	if runtime.GOOS == "windows" {
		fmt.Printf("Please ensure '%s' is added to your system's PATH environment variable for flatc.\n", installDir)
	}

	return nil
}

// Validate method for FlatbuffersInstaller.
func (f *FlatbuffersInstaller) Validate() error {
	cmd := exec.Command("flatc", "--version")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("flatc validation failed: %w", err)
	}
	fmt.Printf("flatc validation successful. Version: %s\n", strings.TrimSpace(string(output)))
	return nil
}

// Uninstall method for FlatbuffersInstaller.
func (f *FlatbuffersInstaller) Uninstall() error {
	flatcExecutableName := "flatc"
	if runtime.GOOS == "windows" {
		flatcExecutableName += ".exe"
	}

	// Try common installation paths
	possiblePaths := []string{}
	if runtime.GOOS == "windows" {
		appDataPath := os.Getenv("APPDATA")
		if appDataPath != "" {
			possiblePaths = append(possiblePaths, filepath.Join(appDataPath, "flatc", "bin", flatcExecutableName))
		}
	} else {
		possiblePaths = append(possiblePaths, "/usr/local/bin/"+flatcExecutableName)
		homeDir, err := os.UserHomeDir()
		if err == nil {
			possiblePaths = append(possiblePaths, filepath.Join(homeDir, ".local", "bin", flatcExecutableName))
		}
	}

	foundPath := ""
	for _, p := range possiblePaths {
		if _, err := os.Stat(p); err == nil {
			foundPath = p
			break
		}
	}

	if foundPath == "" {
		return fmt.Errorf("flatc not found in common installation paths for uninstallation")
	}

	fmt.Printf("Uninstalling flatc from: %s\n", foundPath)

	if runtime.GOOS != "windows" && strings.HasPrefix(foundPath, "/usr/local/bin") {
		fmt.Println("Admin/sudo privileges might be required to uninstall flatc from /usr/local/bin.")
		fmt.Println("Please enter your password if prompted.")
		cmd := exec.Command("sudo", "rm", foundPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to uninstall flatc with sudo: %w", err)
		}
	} else {
		err := os.Remove(foundPath)
		if err != nil {
			return fmt.Errorf("failed to uninstall flatc: %w", err)
		}
	}

	fmt.Println("flatc uninstalled successfully.")
	return nil
}

// Exists method for FlatbuffersInstaller.
func (f *FlatbuffersInstaller) Exists() bool {
	_, err := exec.LookPath("flatc")
	return err == nil
}

// GetVersion method for FlatbuffersInstaller.
func (f *FlatbuffersInstaller) GetVersion() (string, error) {
	if !f.Exists() {
		return "", fmt.Errorf("flatc is not installed")
	}
	cmd := exec.Command("flatc", "--version")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get flatc version: %w", err)
	}
	// The output usually looks like "flatc version 1.12.0". We need to parse this.
	versionStr := strings.TrimSpace(string(output))
	parts := strings.Split(versionStr, " ")
	if len(parts) >= 3 && parts[0] == "flatc" && parts[1] == "version" {
		return parts[2], nil
	}
	return versionStr, nil // Fallback if format is unexpected
}
