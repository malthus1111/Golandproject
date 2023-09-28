# Utilisez l'image officielle de Golang.
FROM golang:1.21.1

# Définir le répertoire de travail dans le conteneur.
WORKDIR /app

# Copiez le code source dans le conteneur.
COPY src/ .

# Construisez l'application.
RUN go build -o /app/main .

# Commande pour exécuter l'application.
CMD ["/app/main"]
