# CVAT Service

Serviço CVAT (Computer Vision Annotation Tool) configurado para rodar em containers OCI gerenciados pelo systemd.

## Arquitetura

O serviço é composto por 6 containers orquestrados:

- **cvat-db**: PostgreSQL 15 (banco de dados)
- **cvat-redis**: Redis 7.2 (cache e filas)
- **cvat-server**: Backend Django
- **cvat-ui**: Frontend React
- **cvat-worker-low**: Worker assíncrono (baixa prioridade)
- **cvat-worker-default**: Worker assíncrono (prioridade padrão)

Todos os containers compartilham a rede `cvat-network` e são gerenciados como uma única unidade pelo systemd.

## Habilitando o Serviço

Em um nó (ex: `whiterun`), adicione ao `default.nix`:

```nix
{
  services.cvat.enable = true;
}
```

## Uso

### Iniciando o CVAT

Para subir todos os containers:

```bash
sudo systemctl start cvat.service
```

Isso iniciará automaticamente todos os 6 containers na ordem correta.

### Inicializando o banco (primeira vez)

Após o primeiro start, execute as migrations e crie o superuser:

```bash
sudo systemctl start cvat-init.service
```

Isso criará o usuário padrão:
- **Username**: `admin`
- **Password**: `admin`
- **Email**: `admin@localhost`

### Parando o CVAT

Para parar todos os containers:

```bash
sudo systemctl stop cvat.service
```

Isso derrubará automaticamente todos os 6 containers.

### Status

Para verificar o status de todos os serviços:

```bash
sudo systemctl status cvat.service
sudo systemctl status podman-cvat-*
```

### Acessando a Interface

O CVAT estará disponível via ts-proxy no endereço:

```
https://cvat.stargazer-shark.ts.net
```

(ou o domínio configurado em `services.ts-proxy.network-domain`)

### Usando o cvat-manage

Um script wrapper `cvat-manage` está disponível para executar comandos Django manage.py:

```bash
# Como root ou com sudo
cvat-manage help

# Criar um novo superuser
cvat-manage createsuperuser

# Executar migrations
cvat-manage migrate

# Abrir shell Django
cvat-manage shell

# Ver todos os comandos disponíveis
cvat-manage help
```

O script automaticamente:
- Executa como usuário `cvat` se necessário
- Verifica se o container está rodando
- Executa o comando dentro do container `cvat-server`

## Diretórios

Todos os dados são armazenados em `/var/lib/cvat`:

```
/var/lib/cvat/
├── data/          # Dados de anotação (projetos, tasks, imagens)
├── logs/          # Logs do servidor
├── keys/          # Chaves SSH/API
├── postgres/      # Dados do PostgreSQL
└── redis/         # Dados do Redis
```

## Recursos

O serviço roda dentro do slice `cvat.slice` com os seguintes limites:

- **CPU**: 1 vCPU (100% de quota)
- **Memória**: 2GB máximo (1.8GB soft limit)

## Ports

A porta é alocada automaticamente pelo sistema de port allocation baseado no hash MD5 da chave "cvat".

Para descobrir a porta alocada:

```bash
# Via nix repl
nix repl
nix-repl> :lf .
nix-repl> nixosConfigurations.whiterun.config.networking.ports.cvat.port
```

## Troubleshooting

### Containers não iniciam

Verifique os logs dos containers individuais:

```bash
journalctl -u podman-cvat-db.service
journalctl -u podman-cvat-redis.service
journalctl -u podman-cvat-server.service
```

### Erro de permissão nos volumes

Verifique que os diretórios têm o ownership correto:

```bash
sudo chown -R cvat:cvat /var/lib/cvat
```

### Rede não existe

Recrie a rede manualmente:

```bash
sudo systemctl restart podman-network-cvat-network.service
```

### Reset completo

Para resetar completamente o CVAT (CUIDADO: apaga todos os dados):

```bash
sudo systemctl stop cvat.service
sudo rm -rf /var/lib/cvat/*
sudo systemctl start cvat.service
sudo systemctl start cvat-init.service
```

## Segurança

- O serviço **NÃO** é exposto à internet
- Acesso apenas via Tailscale (rede privada)
- TLS gerenciado pelo ts-proxy
- Credenciais padrão devem ser alteradas após primeira instalação

## Customização

### Alterar versão do CVAT

```nix
{
  services.cvat = {
    enable = true;
    serverImage = "cvat/server:v2.21.0";
    uiImage = "cvat/ui:v2.21.0";
  };
}
```

### Alterar diretório de dados

```nix
{
  services.cvat = {
    enable = true;
    dataDir = "/mnt/storage/cvat";
  };
}
```

## Referências

- [CVAT Documentation](https://docs.cvat.ai/)
- [CVAT GitHub](https://github.com/cvat-ai/cvat)
