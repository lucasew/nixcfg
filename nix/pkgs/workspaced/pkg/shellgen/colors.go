package shellgen

import (
	"fmt"
	"strings"
	"workspaced/pkg/config"
)

// GenerateColors generates a shell script snippet to apply terminal colors.
//
// It reads the `desktop.palette.base16` configuration and emits raw ANSI escape codes
// directly into the shell script. This approach allows setting terminal colors
// instantaneously upon shell startup without spawning external processes (like `sed` or `cat`),
// significantly reducing shell startup latency.
func GenerateColors() (string, error) {
	cfg, err := config.Load()
	if err != nil {
		return "", err
	}

	var desktop map[string]interface{}
	if err := cfg.UnmarshalKey("desktop", &desktop); err != nil {
		return "", err
	}

	palette, ok := desktop["palette"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("palette not found")
	}

	base16, ok := palette["base16"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("base16 palette not found")
	}

	get := func(key string) string {
		val, _ := base16[key].(string)
		return val
	}

	// Base16 standard colors mapping
	colors := []string{
		get("base00"), get("base08"), get("base0B"), get("base0A"),
		get("base0D"), get("base0E"), get("base0C"), get("base05"),
		get("base03"), get("base08"), get("base0B"), get("base0A"),
		get("base0D"), get("base0E"), get("base0C"), get("base07"),
	}

	var codes strings.Builder
	for i, color := range colors {
		if color != "" {
			// OSC 4: Set color palette entry
			codes.WriteString(fmt.Sprintf("\\033]4;%d;#%s\\033\\\\", i, color))
		}
	}

	fg := get("base05")
	bg := get("base00")
	if fg != "" {
		// OSC 10: Set default foreground color
		codes.WriteString(fmt.Sprintf("\\033]10;#%s\\033\\\\", fg))
		// OSC 12: Set cursor color
		codes.WriteString(fmt.Sprintf("\\033]12;#%s\\033\\\\", fg))
	}
	if bg != "" {
		// OSC 11: Set default background color
		codes.WriteString(fmt.Sprintf("\\033]11;#%s\\033\\\\", bg))
	}

	var output strings.Builder
	output.WriteString("# Apply terminal colors (inline, no external calls)\n")
	output.WriteString("if [[ $- == *i* ]]; then\n")
	output.WriteString("\tprintf '")
	output.WriteString(codes.String())
	output.WriteString("'\n")
	output.WriteString("fi\n")

	return output.String(), nil
}
