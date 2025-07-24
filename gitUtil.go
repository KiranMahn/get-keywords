package main

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func getDocumentCreationDate(path string) (time.Time, error) {

	// Ensure the path is absolute
	absPath, err := getAbsolutePath(path)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Ensure we are in a Git repository
	inValidRepo, err := inGitRepo(absPath)
	if !inValidRepo {
		return time.Time{}, fmt.Errorf("not a git repository (or any of the parent directories): %w", err)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Debug: Print the output of the git blame command
	output, err := runGitBlame(ctx, absPath)
	if err != nil {
		return time.Time{}, err
	}

	// Parse the output to find the last commit date
	lastCommitTime, err := getAuthorTime(output)
	if err != nil {
		return time.Time{}, err
	}

	return lastCommitTime, nil
}

func getLastModified(path string) (time.Time, error) {

	// Ensure the path is absolute
	absPath, err := getAbsolutePath(path)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Ensure we are in a Git repository
	inValidRepo, err := inGitRepo(absPath)
	if !inValidRepo {
		return time.Time{}, fmt.Errorf("not a git repository (or any of the parent directories): %w", err)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// run the git blame command
	dateCreated, err := runGitCommand(ctx, absPath)
	if err != nil {
		return time.Time{}, err
	}

	return dateCreated, nil
}

func runGitCommand(ctx context.Context, absPath string) (time.Time, error) {
	cmd := exec.Command("git", "log", "-1", "--format=%cd", "--date=local", absPath)
	cmd.Dir = filepath.Dir(absPath)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return time.Time{}, fmt.Errorf("git command timed out")
	}
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to run git cmd: %v, %s", err, stderr.String())
	}

	output := strings.TrimSpace(out.String())
	layout := "Mon Jan 2 15:04:05 2006"

	parsedTime, err := time.Parse(layout, output)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse time: %v, %s", err, stderr.String())
	}

	return parsedTime, nil
}

// to run this must clone repo first
func getAbsolutePath(path string) (string, error) {
	// Ensure the path is absolute

	absPath, err := filepath.Abs(path)
	if err != nil {
		return "error", fmt.Errorf("failed to get absolute path: %w", err)
	}
	return absPath, nil
}

func inGitRepo(absPath string) (bool, error) {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	cmd.Dir = filepath.Dir(absPath)
	if err := cmd.Run(); err != nil {
		return false, fmt.Errorf("not a git repository (or any of the parent directories): %w", err)
	}
	return true, nil
}

// get revision and author last modified for the 2 most recent lines
func runGitBlame(ctx context.Context, absPath string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "blame", "-p", "-L1,2", absPath)
	cmd.Dir = filepath.Dir(absPath)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return "NA", fmt.Errorf("git blame command timed out")
	}
	if err != nil {
		return "NA", fmt.Errorf("failed to run git blame: %v, %s", err, stderr.String())
	}

	output := out.String()
	return output, nil
}

func getAuthorTime(output string) (time.Time, error) {
	// Parse the output to find the last commit date
	var lastCommitTime time.Time
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "author-time ") {
			parts := strings.Fields(line)
			if len(parts) == 2 {
				commitTimeInt, err := strconv.ParseInt(parts[1], 10, 64)
				if err != nil {
					return time.Time{}, fmt.Errorf("failed to parse commit time: %w", err)
				}
				commitTime := time.Unix(commitTimeInt, 0)
				if commitTime.After(lastCommitTime) {
					lastCommitTime = commitTime
					return lastCommitTime, nil
				}
			}
		}
	}

	return time.Time{}, fmt.Errorf("no commit time found")
}

func getParentRepo(url string) string {
	parts := strings.Split(url, "/")
	parent := parts[1]
	return parent
}
