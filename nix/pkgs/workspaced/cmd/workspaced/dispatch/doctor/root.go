package doctor

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
	"text/tabwriter"
	"workspaced/pkg/driver"

	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "doctor",
	Short: "Check status of all registered drivers",
	Run: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")
		report := driver.Doctor(cmd.Context())

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "TYPE\tID\tDRIVER\tWEIGHT\tSTATUS\tMESSAGE")

		for _, iface := range report {
			// Use full interface name if verbose, otherwise friendly name
			typeName := iface.Name
			if !verbose {
				typeName = getFriendlyInterfaceName(iface.Name)
			}

			for _, d := range iface.Drivers {
				status := "❌ Unavailable"
				msg := ""
				if d.Available {
					if d.Selected {
						status = "🎯 Selected"
					} else {
						status = "✅ Available"
					}
					if d.Weight == 0 {
						msg = "Warning: implicit selection (weight 0). Consider setting explicit weight."
					}
				} else if d.Error != nil {
					if errors.Is(d.Error, driver.ErrIncompatible) {
						status = "❌ Incompatible"
						// Strip the "driver is incompatible: " prefix if present
						reason := d.Error.Error()
						reason = strings.TrimPrefix(reason, driver.ErrIncompatible.Error()+": ")
						msg = reason
					} else {
						msg = d.Error.Error()
					}
				}

				// Format ID based on verbose flag
				providerID := d.ID
				if verbose && d.ProviderType != nil {
					// Show full provider struct path
					providerID = getProviderTypeName(d.ProviderType)
				}

				// In verbose mode, show driver name with slug ID
				driverName := d.Name
				if verbose {
					driverName = fmt.Sprintf("%s (%s)", d.Name, d.ID)
				}

				fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\t%s\n", typeName, providerID, driverName, d.Weight, status, msg)
			}
		}
		w.Flush()
	},
}

func init() {
	Command.Flags().BoolP("verbose", "v", false, "Show full interface and driver names")
}

// getFriendlyInterfaceName extracts a user-friendly type name from the full interface path
func getFriendlyInterfaceName(fullPath string) string {
	// Extract the part after the last "/"
	// e.g., "workspaced/pkg/driver/audio.Driver" -> "audio.Driver"
	// or   "workspaced/pkg/driver/dialog.Chooser" -> "dialog.Chooser"
	parts := strings.Split(fullPath, "/")
	if len(parts) == 0 {
		return strings.ToLower(fullPath)
	}

	lastPart := parts[len(parts)-1]

	// Split on "." to get package and type
	// e.g., "audio.Driver" -> ["audio", "Driver"]
	// or   "dialog.Chooser" -> ["dialog", "Chooser"]
	dotParts := strings.Split(lastPart, ".")
	if len(dotParts) != 2 {
		return strings.ToLower(lastPart)
	}

	pkg := dotParts[0]
	typeName := dotParts[1]

	// If it's the main "driver" package (e.g., "driver.Driver"), use the parent package
	// e.g., "workspaced/pkg/driver/audio.Driver" -> "audio"
	if pkg == "driver" && len(parts) >= 2 {
		parentPkg := parts[len(parts)-2]
		return strings.ToLower(parentPkg)
	}

	// For typed interfaces like "dialog.Chooser", include both
	return strings.ToLower(pkg) + "." + strings.ToLower(typeName)
}

// getProviderTypeName returns the full path of a provider type
func getProviderTypeName(t any) string {
	rt, ok := t.(reflect.Type)
	if !ok {
		return fmt.Sprintf("%v", t)
	}

	// If it's a pointer, get the underlying type
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}

	// Get package path and name
	pkgPath := rt.PkgPath()
	name := rt.Name()

	if pkgPath != "" && name != "" {
		return pkgPath + "." + name
	}

	return rt.String()
}
