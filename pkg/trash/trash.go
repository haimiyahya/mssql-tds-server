package trash

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// MoveToTrash moves a file/directory to the OS recycle bin/trash
func MoveToTrash(filePath string) error {
	// Check if file exists
	_, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("file not found: %w", err)
	}

	// Create a unique name for trash (with timestamp and path)
	trashName := filepath.Base(filePath) + "_" + time.Now().Format("20060102_150405")

	switch runtime.GOOS {
	case "windows":
		return moveWindowsToTrash(filePath)
	case "darwin":
		return moveMacToTrash(filePath)
	case "linux":
		return moveLinuxToTrash(filePath, trashName)
	default:
		// Fallback: move to ./trash directory
		return moveToManualTrash(filePath, trashName)
	}
}

// moveWindowsToTrash moves file to Windows Recycle Bin
func moveWindowsToTrash(filePath string) error {
	// Use PowerShell to move to Recycle Bin
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// PowerShell command to move file to Recycle Bin
	psCommand := fmt.Sprintf(
		"Add-Type -AssemblyName Microsoft.VisualBasic; [Microsoft.VisualBasic.FileIO.FileSystem]::DeleteFile('%s', 'OnlyErrorDialogs', 'SendToRecycleBin')",
		absPath,
	)

	cmd := exec.Command("powershell", "-Command", psCommand)
	cmd.Stdin = strings.NewReader("")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to move to recycle bin: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// moveMacToTrash moves file to macOS Trash
func moveMacToTrash(filePath string) error {
	// Use AppleScript to move file to Trash
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Escape path for AppleScript
	escapedPath := strings.ReplaceAll(absPath, "'", "'\\''")

	// AppleScript to move file to Trash
	script := fmt.Sprintf("tell application \"Finder\" to delete POSIX file \"%s\"", escapedPath)

	cmd := exec.Command("osascript", "-e", script)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to move to trash: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// moveLinuxToTrash moves file to Linux trash
func moveLinuxToTrash(filePath, trashName string) error {
	// Try to use gio (GNOME) first
	if gioTrash(filePath) == nil {
		return nil
	}

	// Try to use trash-cli
	if trashCliTrash(filePath) == nil {
		return nil
	}

	// Fallback to manual trash directory
	return moveToManualTrash(filePath, trashName)
}

// gioTrash uses gio command (GNOME) to move to trash
func gioTrash(filePath string) error {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Check if gio is available
	if _, err := exec.LookPath("gio"); err != nil {
		return fmt.Errorf("gio not available: %w", err)
	}

	// Use gio trash command
	cmd := exec.Command("gio", "trash", absPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("gio trash failed: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// trashCliTrash uses trash-cli command to move to trash
func trashCliTrash(filePath string) error {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Check if trash-cli is available
	if _, err := exec.LookPath("trash-put"); err != nil {
		return fmt.Errorf("trash-put not available: %w", err)
	}

	// Use trash-put command
	cmd := exec.Command("trash-put", absPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("trash-put failed: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// moveToManualTrash moves file to manual trash directory
func moveToManualTrash(filePath, trashName string) error {
	// Create ./trash directory if it doesn't exist
	trashDir := "./trash"
	if err := os.MkdirAll(trashDir, 0755); err != nil {
		return fmt.Errorf("failed to create trash directory: %w", err)
	}

	// Build destination path
	destPath := filepath.Join(trashDir, trashName)

	// Move file to trash directory
	if err := os.Rename(filePath, destPath); err != nil {
		// If rename fails (different filesystem), try copy and delete
		return copyAndDelete(filePath, destPath)
	}

	return nil
}

// copyAndDelete copies file to destination and deletes original
func copyAndDelete(src, dst string) error {
	// Open source file
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	// Create destination file
	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	// Copy file contents
	if _, err := dstFile.ReadFrom(srcFile); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	// Close destination file to flush
	if err := dstFile.Close(); err != nil {
		return fmt.Errorf("failed to close destination file: %w", err)
	}

	// Delete source file
	if err := os.Remove(src); err != nil {
		return fmt.Errorf("failed to delete source file: %w", err)
	}

	return nil
}

// RestoreFromTrash restores a file from trash
func RestoreFromTrash(trashPath, restorePath string) error {
	// Check if trash file exists
	if _, err := os.Stat(trashPath); err != nil {
		return fmt.Errorf("trash file not found: %w", err)
	}

	// Check if restore path exists
	if _, err := os.Stat(restorePath); err == nil {
		return fmt.Errorf("restore path already exists: %s", restorePath)
	}

	// Move file from trash to restore path
	if err := os.Rename(trashPath, restorePath); err != nil {
		// If rename fails (different filesystem), try copy and delete
		return copyAndDelete(trashPath, restorePath)
	}

	return nil
}

// ListTrash lists files in manual trash directory
func ListTrash() ([]os.FileInfo, error) {
	trashDir := "./trash"

	// Check if trash directory exists
	if _, err := os.Stat(trashDir); os.IsNotExist(err) {
		return []os.FileInfo{}, nil
	}

	// Read trash directory
	files, err := os.ReadDir(trashDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read trash directory: %w", err)
	}

	// Convert to []os.FileInfo
	result := make([]os.FileInfo, 0, len(files))
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			continue
		}
		result = append(result, info)
	}

	return result, nil
}

// EmptyTrash empties the manual trash directory
func EmptyTrash() error {
	trashDir := "./trash"

	// Check if trash directory exists
	if _, err := os.Stat(trashDir); os.IsNotExist(err) {
		return nil
	}

	// Remove all files in trash directory
	files, err := os.ReadDir(trashDir)
	if err != nil {
		return fmt.Errorf("failed to read trash directory: %w", err)
	}

	for _, file := range files {
		filePath := filepath.Join(trashDir, file.Name())
		if err := os.RemoveAll(filePath); err != nil {
			return fmt.Errorf("failed to delete %s: %w", file.Name(), err)
		}
	}

	return nil
}
