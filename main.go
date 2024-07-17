package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	sqlite "_/internal/db"
	chgkdb "_/pkg/chgkdb"
)

func main() {
	chdbPath := flag.String("chdb", "baza", "Path to the ЧГК database dir")
	sqlitePath := flag.String("sqlite", "chgk.db", "Path to the SQLite database file")
	cvsCheckout := flag.Bool("cvs", false, "Checkout the latest ЧГК database from CVS")
	flag.Parse()

	if *cvsCheckout {
		if err := checkoutFromCVS(); err != nil { // Simplified error check
			log.Fatalf("Failed to checkout from CVS: %v", err)
		}
	}

	questions, err := chgkdb.LoadQuestions(*chdbPath)
	if err != nil {
		log.Fatalf("Failed to load ЧГК database: %v", err)
	}

	db, err := sqlite.InitializeDatabase(*sqlitePath)
	if err != nil {
		log.Fatalf("Failed to initialize SQLite database: %v", err)
	}
	defer db.Close()

	if err = sqlite.InsertQuestions(db, questions); err != nil { // Simplified error check
		log.Fatalf("Failed to insert questions into SQLite: %v", err)
	}

	fmt.Println("Database export successful")
}

func cvsLogin() error {
	cmd := exec.Command("cvs", "-d", ":pserver:anonymous@bilbo.dynip.com:/home/cvsroot", "login")

	// Create a buffer to capture stdout and stderr
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	// Provide the password via stdin
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdin pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	// Write the password to stdin
	if _, err := stdin.Write([]byte("anonymous\n")); err != nil {
		return fmt.Errorf("failed to write to stdin: %w", err)
	}
	stdin.Close()

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("command failed: %w\nOutput: %s", err, out.String())
	}

	fmt.Println(out.String())
	return nil
}

func cvsCheckout() error {
	err := cvsLogin()
	if err != nil {
		return err
	}
	err = cvsCheckout()
	return err
}

func checkoutFromCVS() error {
	cmds := []struct {
		cmd  string
		args []string
	}{
		{"sh", []string{"-c", "echo anonymous | cvs -d :pserver:anonymous@bilbo.dynip.com:/home/cvsroot login"}},
		{"cvs", []string{"-d", ":pserver:anonymous@bilbo.dynip.com:/home/cvsroot", "checkout", "baza"}},
	}

	for _, c := range cmds {
		cmd := exec.Command(c.cmd, c.args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}
