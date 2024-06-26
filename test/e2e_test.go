package test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/stretchr/testify/assert"
	// Import your main package
)

var baileysHammerCmd *exec.Cmd

func getChomeDPContext() context.Context {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),    // disable headless mode
		chromedp.Flag("disable-gpu", false), // enable GPU acceleration (optional)
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := context.WithTimeout(allocCtx, 15*time.Second)
	defer cancel()

	return ctx
}

func TestMain(m *testing.M) {
	// set DATABASE_URL to ../tmp/test.db
	os.Setenv("DATABASE_URL", "../tmp/test.db")
	os.Setenv("PASS", "test_pass")

	// Start the "baileys-hammer" program
	baileysHammerCmd = exec.Command("../tmp/main")
	baileysHammerCmd.Stdout = os.Stdout
	baileysHammerCmd.Stderr = os.Stderr
	if err := baileysHammerCmd.Start(); err != nil {
		fmt.Printf("Failed to start baileys-hammer: %v\n", err)
		os.Exit(1)
	}

	ctx := getChomeDPContext()

	err := chromedp.Run(ctx,
		chromedp.Navigate(fmt.Sprintf("http://localhost:8080/finemaster/%s?f=false&fl=false&m=false&p=true&pf=false#players-manage", "test_pass")),
		chromedp.WaitVisible(`#name`, chromedp.ByID),
		chromedp.SendKeys(`#name`, "test player", chromedp.ByID),
		chromedp.Click("#add-player", chromedp.ByID),
	)

	if err != nil {
		fmt.Printf("Failed to add player: %v\n", err)
		os.Exit(1)
	}

	// Wait for the program to start up (adjust time as necessary)
	time.Sleep(5 * time.Second)

	// Run the tests
	code := m.Run()

	// Stop the "baileys-hammer" program
	if err := baileysHammerCmd.Process.Signal(os.Interrupt); err != nil {
		fmt.Printf("Failed to stop baileys-hammer: %v\n", err)
	}

	os.Exit(code)
}

func TestPageLoad(t *testing.T) {
	ctx := getChomeDPContext()

	var title string
	err := chromedp.Run(ctx,
		chromedp.Navigate("http://localhost:8080"),
		chromedp.Title(&title),
	)

	assert.NoError(t, err)
	assert.Equal(t, "🔨 Baileys Hammer 🔨", title)
}

func TestFormSubmission(t *testing.T) {
	ctx := getChomeDPContext()

	var result string
	err := chromedp.Run(ctx,
		chromedp.Navigate("http://localhost:8080/form"),
		chromedp.WaitVisible(`#select-fine-ts-control`, chromedp.ByID),
		chromedp.SendKeys(`#select-fine-ts-control`, "test fine", chromedp.ByID),

		chromedp.SendKeys(`#select-player`, "test player", chromedp.ByID),
		chromedp.Click(`#submit-button`, chromedp.ByID),
	)

	assert.NoError(t, err)
	assert.Contains(t, result, "Expected Result")
}
