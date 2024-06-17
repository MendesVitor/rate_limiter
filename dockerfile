# Imagem base
FROM golang:latest

# Configuração do diretório de trabalho
WORKDIR /app

# Copiando o código fonte e o arquivo .env
COPY . .

# Compilando o aplicativo
RUN go build -o app .

# Comando padrão para executar o aplicativo
CMD ["./app"]
