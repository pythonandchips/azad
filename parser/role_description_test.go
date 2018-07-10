package parser

import (
	"testing"

	"github.com/pythonandchips/azad/expect"
	"github.com/stretchr/testify/assert"
)

func TestRoleDescriptioGroups(t *testing.T) {
	roleDescriptionGroups := roleDescriptionGroups{
		{name: "ruby"},
		{name: "java"},
		{name: "erlang"},
	}
	t.Run("RoleFor", func(t *testing.T) {
		t.Run("returns roleDescription group when role exists in group", func(t *testing.T) {
			roleDescriptionGroup, err := roleDescriptionGroups.RoleFor("ruby")
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			assert.Equal(t, roleDescriptionGroup.name, "ruby")
		})
		t.Run("error when role does not exist", func(t *testing.T) {
			_, err := roleDescriptionGroups.RoleFor("nginx")
			if err == nil {
				t.Fatalf("expected error but none returned")
			}
		})
	})
	t.Run("FilterMany", func(t *testing.T) {
		t.Run("returns all matching roleDescriptionGroup", func(t *testing.T) {
			matchingRoleDescriptionGroup, err := roleDescriptionGroups.FilterMany([]string{"ruby", "java"})
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			expect.EqualFatal(t, len(matchingRoleDescriptionGroup), 2)
		})
		t.Run("returns error of any name cannot be matched", func(t *testing.T) {
			_, err := roleDescriptionGroups.FilterMany([]string{"ruby", "elm"})
			if err == nil {
				t.Fatalf("expected error but none returned")
			}
		})
	})
	t.Run("Filter", func(t *testing.T) {
		t.Run("returns matching roleDescriptionGroup", func(t *testing.T) {
			matchingRoleDescriptionGroup, err := roleDescriptionGroups.Filter("ruby")
			if err != nil {
				t.Fatalf("Unexpected error: %s", err)
			}
			assert.Equal(t, matchingRoleDescriptionGroup.name, "ruby")
		})
		t.Run("returns error if roleDescriptionGroup not found", func(t *testing.T) {
			_, err := roleDescriptionGroups.Filter("elm")
			if err == nil {
				t.Fatalf("Expected error but none returned")
			}
		})
	})
	t.Run("ReplaceWith", func(t *testing.T) {
		t.Run("replaces the role with matching name", func(t *testing.T) {
			updatedRoleDescriptionGroup := roleDescriptionGroup{
				name: "java",
				files: []roleFile{
					{name: "main"},
				},
			}

			roleDescriptionGroups.ReplaceWith(updatedRoleDescriptionGroup)

			expect.EqualFatal(t, len(roleDescriptionGroups), 3)

			assert.Equal(t, len(updatedRoleDescriptionGroup.files), len(roleDescriptionGroups[1].files))
		})
	})
}
