local lspconfig = require'lspconfig'
local coq = require'coq'

for name, value in ipairs({
    arduino_language_server = {},
    bashls = {},
    ccls = {}, -- c/c++
    cmake = {},
    dockerls = {},
    dotls = {}, -- dot/graphviz
    emmet_ls = {},
    gopls = {}, -- golang
    graphql = {},
    hls = {}, -- haskell
    rnix = {}, -- nix
    rust_analyzer = {}, -- rust
    terraformls = {}, -- terraform
    texlab = {}, -- latex
    tsserver = {}, -- typescript
    vimls = {}, -- vimscript
    yamlls = {}, -- yaml
    zls = {}, -- zig
    svelte = {} -- svelte
}) do
    print("Setting up language server " .. name .. "...")
    local coqed = coq.lsp_ensure_capatibilities(value)
    lspconfig[name].setup(coqed)
end
