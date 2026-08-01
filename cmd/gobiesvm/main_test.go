package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestFixtureOutputs(t *testing.T) {
	tests := map[string]string{
		"arithmetic.rb":         "10\n23\n46\n23\n11.500000000000\n",
		"array.rb":              "[7, 2, 3, 4]\n",
		"block_times.rb":        strings.Repeat("Hello\n", 10),
		"block_times_params.rb": "0\n1\n2\n0\n1\n2\n0\n1\n2\n",
		"factorial.rb":          "2432902008176640000\n",
		"p.rb":                  "'Hello, '\n'world'\n",
		"puts.rb":               "kerker\n12345\n",
		"strconcat.rb":          "kerker2herher\n",
	}

	for fixture, expected := range tests {
		t.Run(fixture, func(t *testing.T) {
			stdout, stderr, err := runFixture(t, fixture)
			if err != nil {
				t.Fatalf("run %s: %v\nstderr:\n%s", fixture, err, stderr)
			}
			if stdout != expected {
				t.Fatalf("unexpected output\nwant:\n%s\ngot:\n%s", expected, stdout)
			}
		})
	}
}

func TestNestedEachFixture(t *testing.T) {
	var expected strings.Builder
	for _, letter := range []string{"a", "b", "c", "d", "e", "f"} {
		expected.WriteString(letter + "\n1\n2\n3\n4\n")
	}

	stdout, stderr, err := runFixture(t, "each.rb")
	if err != nil {
		t.Fatalf("run each.rb: %v\nstderr:\n%s", err, stderr)
	}
	if stdout != expected.String() {
		t.Fatalf("unexpected output\nwant:\n%s\ngot:\n%s", expected.String(), stdout)
	}
}

func TestReadlinesFixture(t *testing.T) {
	stdout, stderr, err := runFixture(t, "readlines.rb")
	if err != nil {
		t.Fatalf("run readlines.rb: %v\nstderr:\n%s", err, stderr)
	}
	for _, expected := range []string{"['a = 10\n',", "'p a.to_f / 4.to_f # 11.5\n', '']"} {
		if !strings.Contains(stdout, expected) {
			t.Fatalf("output does not contain %q:\n%s", expected, stdout)
		}
	}
}

func TestHashFixtureIgnoresMapIterationOrder(t *testing.T) {
	stdout, stderr, err := runFixture(t, "hash.rb")
	if err != nil {
		t.Fatalf("run hash.rb: %v\nstderr:\n%s", err, stderr)
	}
	lines := strings.Split(strings.TrimSpace(stdout), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected two hashes, got %d lines:\n%s", len(lines), stdout)
	}
	for _, key := range []string{"cat", "dog", "cow", "tiger", "lion", "wolf", "fish"} {
		if !strings.Contains(lines[0], key+"=>0") {
			t.Errorf("first hash is missing %s=>0: %s", key, lines[0])
		}
		expected := key + "=>0"
		if key == "dog" {
			expected = "dog=>2"
		}
		if !strings.Contains(lines[1], expected) {
			t.Errorf("second hash is missing %s: %s", expected, lines[1])
		}
	}
}

func TestConcurrentFixtures(t *testing.T) {
	t.Run("safety", func(t *testing.T) {
		stdout, stderr, err := runFixture(t, "safety.rb")
		if err != nil {
			t.Fatalf("run safety.rb: %v\nstderr:\n%s", err, stderr)
		}
		if stdout != "10000\n" {
			t.Fatalf("expected 10000 elements, got %q", stdout)
		}
	})

	t.Run("pi", func(t *testing.T) {
		if raceEnabled {
			t.Skip("the safety fixture covers STM concurrency under the race detector")
		}
		stdout, stderr, err := runFixture(t, "pi.rb")
		if err != nil {
			t.Fatalf("run pi.rb: %v\nstderr:\n%s", err, stderr)
		}
		value, err := strconv.ParseFloat(strings.TrimSpace(stdout), 64)
		if err != nil {
			t.Fatalf("parse pi output %q: %v", stdout, err)
		}
		if value < 2.5 || value > 3.8 {
			t.Fatalf("pi estimate is out of range: %v", value)
		}
	})
}

func TestRunReportsInputErrors(t *testing.T) {
	t.Run("missing argument", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		err := run(nil, &stdout, &stderr)
		if err == nil || !strings.Contains(err.Error(), "expected one .rb file") {
			t.Fatalf("expected argument error, got %v", err)
		}
		if !strings.Contains(stderr.String(), "Usage of gobiesvm") {
			t.Fatalf("expected usage text, got %q", stderr.String())
		}
	})

	t.Run("unreadable source", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		err := run([]string{filepath.Join(t.TempDir(), "missing.rb")}, &stdout, &stderr)
		if err == nil || !strings.Contains(err.Error(), "read") {
			t.Fatalf("expected read error, got %v", err)
		}
	})

	t.Run("invalid syntax", func(t *testing.T) {
		filename := writeScript(t, "@")
		var stdout, stderr bytes.Buffer
		err := run([]string{filename}, &stdout, &stderr)
		if err == nil || !strings.Contains(err.Error(), "parse") {
			t.Fatalf("expected parse error, got %v", err)
		}
	})

	t.Run("IO.readlines failure", func(t *testing.T) {
		missing := filepath.Join(t.TempDir(), "missing.txt")
		filename := writeScript(t, fmt.Sprintf("lines = IO.readlines('%s')\n", missing))
		var stdout, stderr bytes.Buffer
		err := run([]string{filename}, &stdout, &stderr)
		if err == nil || !strings.Contains(err.Error(), "IO.readlines") {
			t.Fatalf("expected IO.readlines error, got %v", err)
		}
	})
}

func TestCLIFlagsRemainAvailable(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{name: "ast", args: []string{"-ast", "arithmetic.rb"}, expected: "NODE_ROOT"},
		{name: "bytecode", args: []string{"-bytecode", "arithmetic.rb"}, expected: "BC_PUTOBJ"},
		{name: "log", args: []string{"-log", "arithmetic.rb"}, expected: "Transaction Logging Status"},
		{name: "bench", args: []string{"-bench=1", "arithmetic.rb"}, expected: "Average Time in 1 iterations"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stdout, stderr, err := runCommand(t, test.args...)
			if err != nil {
				t.Fatalf("run %v: %v\nstderr:\n%s", test.args, err, stderr)
			}
			if !strings.Contains(stdout, test.expected) {
				t.Fatalf("output does not contain %q:\n%s", test.expected, stdout)
			}
		})
	}
}

func TestGobiesVMHelperProcess(t *testing.T) {
	if os.Getenv("GOBIESVM_HELPER_PROCESS") != "1" {
		return
	}
	separator := -1
	for i, arg := range os.Args {
		if arg == "--" {
			separator = i
			break
		}
	}
	if separator == -1 {
		fmt.Fprintln(os.Stderr, "missing helper argument separator")
		os.Exit(2)
	}
	if err := run(os.Args[separator+1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(0)
}

func runFixture(t *testing.T, fixture string) (string, string, error) {
	t.Helper()
	return runCommand(t, fixture)
}

func runCommand(t *testing.T, args ...string) (string, string, error) {
	t.Helper()
	commandArgs := []string{"-test.run=^TestGobiesVMHelperProcess$", "--"}
	commandArgs = append(commandArgs, args...)
	cmd := exec.Command(os.Args[0], commandArgs...)
	cmd.Dir = filepath.Join("..", "..", "tests")
	cmd.Env = append(os.Environ(), "GOBIESVM_HELPER_PROCESS=1")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func writeScript(t *testing.T, source string) string {
	t.Helper()
	filename := filepath.Join(t.TempDir(), "script.rb")
	if err := os.WriteFile(filename, []byte(source), 0o600); err != nil {
		t.Fatalf("write script: %v", err)
	}
	return filename
}
