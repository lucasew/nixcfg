(local lspconfig (require :lspconfig))
(local coq (require :coq))
(local icons (require :nvim-web-devicons))
(local lsp_signature (require :lsp_signature))

(local v (require :adapter.nvim.utils))

(lsp_signature.setup { :bind true })

(fn lsp [name options]
  (local opts (or options {}))
  (local coqed (coq.lsp_ensure_capabilities opts))
  (assert (. lspconfig name) (.. "The server " name " is not defined on lspconfig"))
  ((. (. lspconfig name) :setup) coqed))


(lsp :bashls)
(lsp :ccls { ;; C/C++
  :init_options {
    :cache 
      {:directory "/tmp/.ccls-cache"}
    :completion 
      {:detailedLabel false}
}})

(lsp :cmake)
(lsp :dockerls) ;; Dockerfile
(lsp :dotls) ;; GraphViz
(lsp :gopls) ;; Golang
(lsp :graphql) ;; GraphQL
(lsp :hls) ;; Haskell
(lsp :java_language_server { ;; Java
  :cmd [:java_language_server]
})
(lsp :rnix) ;; Nix
(lsp :rust_analyzer) ;; Rust
(lsp :terraformls) ;; Terraform
(lsp :texlab) ;; LaTeX
(lsp :tsserver) ;; TypeScript
(lsp :vimls) ;; VimL
(lsp :yamlls) ;; YAML
(lsp :zls) ;; Zig
(lsp :svelte) ;; Svelte

(v.cmd "nmap <C-p> :Telescope<CR>")
(v.cmd "nmap <C-.> :Telescope lsp_code_actions<CR>")

(print "hello")