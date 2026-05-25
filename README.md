# App Sara — Backend

API REST para App Sara. Autenticación JWT, roles (user/veterinarian/admin), gestión de usuarios.

## Stack

- **Go 1.26** + **Fiber v2** (HTTP)
- **PostgreSQL** + **GORM** (ORM)
- **JWT** (access 15min + refresh 7d)
- **bcrypt** (contraseñas)
- **Rate limiting**, **validator**, **CORS-ready**

## Requisitos

- Go 1.26+
- PostgreSQL (o Supabase)
- Docker (opcional)

## Configuración

```bash
cp .env.example .env
# Editar .env con tus credenciales
```

| Variable | Descripción |
|---|---|
| `DATABASE_URL` | DSN de PostgreSQL |
| `JWT_SECRET` | Clave para firmar access tokens |
| `JWT_REFRESH_SECRET` | Clave para firmar refresh tokens |
| `PORT` | Puerto del servidor (default: 3000) |
| `ADMIN_SECRET` | Secreto para registrar admins |

## Desarrollo

```bash
# Descargar dependencias
go mod download

# Iniciar servidor (hot reload con air)
go run ./cmd/server
```

Servidor en `http://localhost:3000`.

## Docker

```bash
# Construir imagen
docker build -t app-sara-backend .

# Ejecutar
docker run -p 3000:3000 --env-file .env app-sara-backend
```

## API

Base: `/api/v1`

### Auth (público, rate limit: 10 req/min)

| Método | Ruta | Descripción |
|---|---|---|
| POST | `/auth/register` | Registro |
| POST | `/auth/login` | Login |
| POST | `/auth/refresh` | Refresh token |

### Usuarios (rate limit: 30 req/min)

| Método | Ruta | Auth | Admin | Descripción |
|---|---|---|---|---|
| GET | `/users/me` | ✅ | ❌ | Perfil propio |
| PUT | `/users/me` | ✅ | ❌ | Actualizar perfil |
| GET | `/users` | ✅ | ✅ | Listar usuarios |
| GET | `/users/:id` | ✅ | ✅ | Obtener usuario |
| DELETE | `/users/:id` | ✅ | ✅ | Eliminar (soft delete) |

### Roles

- `user` — usuario estándar
- `veterinarian` — veterinario
- `admin` — administrador (requiere `admin_secret` en registro)

## Estructura

```
cmd/server/          → Entry point
internal/
  config/            → Config desde env vars
  database/          → Conexión PostgreSQL + AutoMigrate
  domain/            → Modelos y DTOs
  auth/              → Registro, login, refresh
  user/              → CRUD de usuarios
  middleware/        → Auth JWT, roles, rate limit, validación
  router/            → Definición de rutas
```

## Endpoints clave

### Register

```json
POST /api/v1/auth/register
{
  "name": "Juan",
  "email": "juan@email.com",
  "password": "123456",
  "role": "user",
  "admin_secret": ""  // requerido si role = "admin"
}
```

### Login

```json
POST /api/v1/auth/login
{
  "email": "juan@email.com",
  "password": "123456"
}
```

### Respuesta auth

```json
{
  "user": { ... },
  "access_token": "eyJ...",
  "refresh_token": "eyJ..."
}
```

### Headers para endpoints autenticados

```
Authorization: Bearer <access_token>
```

## Postman

Colección incluida: `app_sara.postman_collection.json`
