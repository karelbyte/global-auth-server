# Utiliza una imagen base de Go para construir la aplicación
FROM golang:latest AS builder

# Establece el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copia los archivos go.mod y go.sum (si usas Go Modules)
COPY go.mod go.sum ./

# Descarga las dependencias
RUN go mod download -x

# Copia el resto del código fuente de la aplicación
COPY . .
COPY .env.docker .env
COPY templates ./templates

# Construye la aplicación Go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o global-auth .

# --- Imagen final y ligera para la ejecución ---
FROM alpine:latest
RUN apk --no-cache add ca-certificates

# Establece el directorio de trabajo
WORKDIR /app

# Copia el ejecutable construido desde la etapa 'builder'
COPY --from=builder /app/global-auth .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/.env .env
COPY --from=builder /app/certificates ./certificates

# Expone cualquier puerto que tu aplicación escuche (ejemplo: 8080)
EXPOSE 4000
RUN chmod +x /app/global-auth
# Comando para ejecutar la aplicación cuando el contenedor se inicie
CMD ["./global-auth"]