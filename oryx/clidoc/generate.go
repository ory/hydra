package clidoc

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/spf13/cobra"
)

// Generate generates markdown documentation for a cobra command and its children.
func Generate(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("command expects one argument which is the path to the output directory")
	}

	return generate(cmd, args[0])
}

func trimExt(s string) string {
	return strings.ReplaceAll(strings.TrimSuffix(s, filepath.Ext(s)), "_", "-")
}

func generate(cmd *cobra.Command, dir string) error {
	cmd.DisableAutoGenTag = true
	for _, c := range cmd.Commands() {
		if !c.IsAvailableCommand() || c.IsAdditionalHelpTopicCommand() {
			continue
		}
		if err := generate(c, dir); err != nil {
			return err
		}
	}

	basename := strings.Replace(cmd.CommandPath(), " ", "-", -1)
	if err := os.MkdirAll(filepath.Join(dir), 0750); err != nil {
		return err
	}

	filename := filepath.Join(dir, basename) + ".md"
	f, err := os.Create(filename) //#nosec:G304
	if err != nil {
		return err
	}
	defer (func() { _ = f.Close() })()

	if _, err := io.WriteString(f, fmt.Sprintf(`---
id: %s
title: %s
description: %s
---

<!--
This file is auto-generated.

To improve this file please make your change against the appropriate "./cmd/*.go" file.
-->
`,
		basename,
		cmd.CommandPath(),
		cmd.CommandPath(),
	)); err != nil {
		return err
	}

	var b bytes.Buffer
	if err := GenMarkdownCustom(cmd, &b, trimExt); err != nil {
		return err
	}

	_, err = f.WriteString(b.String())
	return err
}
