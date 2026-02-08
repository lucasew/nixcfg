# Base16 GTK Theme

Tema GTK autocontido baseado em [FlatColor](https://github.com/jasperro/FlatColor) que usa a paleta do workspaced.

## Como funciona

- Tema completo já incluído no dotfiles (gtk-2.0, gtk-3.0, gtk-3.20)
- Arquivos `colors2` e `colors3` são gerados automaticamente via templates
- Tema recarrega automaticamente quando a paleta muda

## Setup (fazer uma vez)

### 1. Aplicar o tema:
```bash
gsettings set org.gnome.desktop.interface gtk-theme base16
```

### 2. Criar tema dummy para reload automático:
```bash
ln -Ts ~/.themes/base16 ~/.themes/dummy
```

## Uso

Sempre que você mudar a paleta no `settings.toml` e rodar `workspaced sync`, o tema GTK será atualizado automaticamente!

Apps GTK que já estão rodando vão recarregar o tema sem precisar reiniciar.
