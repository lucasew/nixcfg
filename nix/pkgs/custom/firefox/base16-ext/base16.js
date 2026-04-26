
/**
 * Initializes and applies the base16 color palette as CSS custom properties
 * on the document body. This allows web extensions and customized pages to
 * inherit the system-wide colors dynamically.
 *
 * The %COLORS% placeholder is replaced at build-time with the actual JSON
 * payload defined in the Nix derivation.
 */
(function() {
    'use strict';
    console.time("base16")
    const colors = %COLORS%
    Object.keys(colors.colors).forEach(k => {
        document.body.style.setProperty(`--${k}`, `#${colors.colors[k]}`)
    })
    console.timeEnd("base16")
    console.log("base16 kicked in")
})();
