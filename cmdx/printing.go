package cmdx

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/go-openapi/jsonpointer"
	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
)

type (
	TableHeader interface {
		Header() []string
	}
	Table interface {
		TableHeader
		Table() [][]string
		Interface() interface{}
		Len() int
	}

	format string
)

const (
	FormatQuiet       format = "quiet"
	FormatTable       format = "table"
	FormatJSON        format = "json"
	FormatJSONPath    format = "jsonpath"
	FormatJSONPointer format = "jsonpointer"
	FormatJSONPretty  format = "json-pretty"
	FormatYAML        format = "yaml"
	FormatDefault     format = "default"

	FlagFormat = "format"

	None = "<none>"
)

func filterJSONPointer(cmd *cobra.Command, data any) any {
	f, err := cmd.Flags().GetString(FlagFormat)
	// unexpected error
	Must(err, "flag access error: %s", err)
	_, jsonptr, found := strings.Cut(f, "=")
	if !found {
		_, _ = fmt.Fprintf(os.Stderr,
			"Format %s is missing a JSON pointer, e.g., --%s=%s=<jsonpointer>. The path syntax is described at https://datatracker.ietf.org/doc/html/draft-ietf-appsawg-json-pointer-07.",
			f, FlagFormat, f)
		os.Exit(1)
	}
	ptr, err := jsonpointer.New(jsonptr)
	Must(err, "invalid JSON pointer: %s", err)

	result, _, err := ptr.Get(data)
	Must(err, "failed to apply JSON pointer: %s", err)

	return result
}

func PrintTable(cmd *cobra.Command, table Table) {
	f := getFormat(cmd)

	switch f {
	case FormatQuiet:
		if table.Len() == 0 {
			fmt.Fprintln(cmd.OutOrStdout())
		}

		if idAble, ok := table.(interface{ IDs() []string }); ok {
			for _, row := range idAble.IDs() {
				fmt.Fprintln(cmd.OutOrStdout(), row)
			}
			break
		}

		for _, row := range table.Table() {
			fmt.Fprintln(cmd.OutOrStdout(), row[0])
		}
	case FormatJSON:
		printJSON(cmd.OutOrStdout(), table.Interface(), false, "")
	case FormatJSONPretty:
		printJSON(cmd.OutOrStdout(), table.Interface(), true, "")
	case FormatJSONPath:
		printJSON(cmd.OutOrStdout(), table.Interface(), true, getPath(cmd))
	case FormatJSONPointer:
		printJSON(cmd.OutOrStdout(), filterJSONPointer(cmd, table.Interface()), true, "")
	case FormatYAML:
		printYAML(cmd.OutOrStdout(), table.Interface())
	default:
		w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 8, 1, '\t', 0)

		for _, h := range table.Header() {
			fmt.Fprintf(w, "%s\t", h)
		}
		fmt.Fprintln(w)

		for _, row := range table.Table() {
			fmt.Fprintln(w, strings.Join(row, "\t")+"\t")
		}

		_ = w.Flush()
	}
}

func getQuiet(cmd *cobra.Command) bool {
	q, err := cmd.Flags().GetBool(FlagQuiet)
	// ignore the error here as we use this function also when the flag might not be registered
	if err != nil {
		return false
	}
	return q
}

func getFormat(cmd *cobra.Command) format {
	q := getQuiet(cmd)

	if q {
		return FormatQuiet
	}

	f, err := cmd.Flags().GetString(FlagFormat)
	// unexpected error
	Must(err, "flag access error: %s", err)

	switch {
	case f == string(FormatTable):
		return FormatTable
	case f == string(FormatJSON):
		return FormatJSON
	case strings.HasPrefix(f, string(FormatJSONPath)):
		return FormatJSONPath
	case strings.HasPrefix(f, string(FormatJSONPointer)):
		return FormatJSONPointer
	case f == string(FormatJSONPretty):
		return FormatJSONPretty
	case f == string(FormatYAML):
		return FormatYAML
	default:
		return FormatDefault
	}
}

func getPath(cmd *cobra.Command) string {
	f, err := cmd.Flags().GetString(FlagFormat)
	// unexpected error
	Must(err, "flag access error: %s", err)
	_, path, found := strings.Cut(f, "=")
	if !found {
		_, _ = fmt.Fprintf(os.Stderr,
			"Format %s is missing a path, e.g., --%s=%s=<path>. The path syntax is described at https://github.com/tidwall/gjson/blob/master/SYNTAX.md",
			f, FlagFormat, f)
		os.Exit(1)
	}

	return path
}

func printJSON(w io.Writer, v interface{}, pretty bool, path string) {
	if path != "" {
		temp, err := json.Marshal(v)
		Must(err, "Error encoding JSON: %s", err)
		v = gjson.GetBytes(temp, path).Value()
	}

	e := json.NewEncoder(w)
	if pretty {
		e.SetIndent("", "  ")
	}
	err := e.Encode(v)
	// unexpected error
	Must(err, "Error encoding JSON: %s", err)
}

func printYAML(w io.Writer, v interface{}) {
	j, err := json.Marshal(v)
	Must(err, "Error encoding JSON: %s", err)
	e, err := yaml.JSONToYAML(j)
	Must(err, "Error encoding YAML: %s", err)
	_, _ = w.Write(e)
}
