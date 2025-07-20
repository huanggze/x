package flagx

import (
	"github.com/spf13/cobra"

	"github.com/huanggze/x/cmdx"
)

// MustGetBool returns a bool flag or fatals if an error occurs.
// Deprecated: just handle the error properly, this breaks command testing
func MustGetBool(cmd *cobra.Command, name string) bool {
	ok, err := cmd.Flags().GetBool(name)
	if err != nil {
		cmdx.Fatalf(err.Error())
	}
	return ok
}

// MustGetInt returns a int flag or fatals if an error occurs.
// Deprecated: just handle the error properly, this breaks command testing
func MustGetInt(cmd *cobra.Command, name string) int {
	ss, err := cmd.Flags().GetInt(name)
	if err != nil {
		cmdx.Fatalf(err.Error())
	}
	return ss
}
