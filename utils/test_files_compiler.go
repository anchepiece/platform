// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package utils

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func goMod(t *testing.T, dir string, args ...string) {
	cmd := exec.Command("go", append([]string{"mod"}, args...)...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to %s: %s", strings.Join(args, " "), string(output))
	}
}

func CompileGo(t *testing.T, sourceCode, outputPath string) {
	dir, err := ioutil.TempDir(".", "")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	dir, err = filepath.Abs(dir)
	require.NoError(t, err)

	// Write out main.go given the source code.
	main := filepath.Join(dir, "main.go")
	err = ioutil.WriteFile(main, []byte(sourceCode), 0600)
	require.NoError(t, err)

	_, sourceFile, _, ok := runtime.Caller(0)
	require.True(t, ok)
	serverPath := filepath.Dir(filepath.Dir(sourceFile))

	out := &bytes.Buffer{}
	cmd := exec.Command("go", "build", "-mod=vendor", "-o", outputPath, main)
	cmd.Dir = serverPath
	cmd.Stdout = out
	cmd.Stderr = out
	err = cmd.Run()
	if err != nil {
		t.Log("Go compile errors:\n", out.String())
	}
	require.NoError(t, err, "failed to compile go")
}
