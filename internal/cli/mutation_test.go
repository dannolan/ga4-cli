package cli

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestMutationCommandsRequireApplyFlag(t *testing.T) {
	root := NewRootCommand()
	paths := []string{
		"ga4 data audience-exports create",
		"ga4 admin accounts delete",
		"ga4 admin accounts patch",
		"ga4 admin properties create",
		"ga4 admin properties update-data-retention-settings",
		"ga4 admin property-resources conversion-events create",
		"ga4 admin property-resources custom-dimensions archive",
		"ga4 admin property-resources custom-metrics archive",
		"ga4 admin property-resources data-streams delete",
		"ga4 admin property-resources measurement-protocol-secrets patch",
		"ga4 admin property-resources firebase-links delete",
		"ga4 admin property-resources google-ads-links patch",
		"ga4 admin property-resources key-events create",
	}
	for _, path := range paths {
		cmd := findCommandByPath(root, path)
		if cmd == nil {
			t.Fatalf("missing command %s", path)
		}
		if cmd.Flags().Lookup("apply") == nil {
			t.Fatalf("mutation command %s does not expose --apply", path)
		}
	}
}

func findCommandByPath(root interface{ Commands() []*cobra.Command }, path string) *cobra.Command {
	for _, child := range root.Commands() {
		if child.CommandPath() == path {
			return child
		}
		if found := findCommandByPath(child, path); found != nil {
			return found
		}
	}
	return nil
}
