package shellgen

import (
	"fmt"
	"strings"
	"workspaced/pkg/config"
)

// GenerateColors generates terminal color codes inline
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

	colors := []string{
		get("base00"), get("base08"), get("base0B"), get("base0A"),
		get("base0D"), get("base0E"), get("base0C"), get("base05"),
		get("base03"), get("base08"), get("base0B"), get("base0A"),
		get("base0D"), get("base0E"), get("base0C"), get("base07"),
	}

	var codes strings.Builder
	for i, color := range colors {
		if color != "" {
			codes.WriteString(fmt.Sprintf("\\033]4;%d;#%s\\033\\\\", i, color))
		}
	}

	fg := get("base05")
	bg := get("base00")
	if fg != "" {
		codes.WriteString(fmt.Sprintf("\\033]10;#%s\\033\\\\", fg))
		codes.WriteString(fmt.Sprintf("\\033]12;#%s\\033\\\\", fg))
	}
	if bg != "" {
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
