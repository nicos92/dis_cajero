# Code Review: dis_cajero - Simulador de Cajero Automático

**Fecha:** 2026-09-02
**Revisor:** Staff Engineer (Go 1.22+, Clean Architecture)
**Archivos analizados:** 15 archivos `.go`, `go.mod`

---

## 1. [CRÍTICO] Bugs, Recursividad y Riesgos

### BUG P0-1: Loop infinito sin salida limpia

**Archivo:** `internal/presentacion/consola.go:182-190`

```go
func (consola *consolaa) Iniciar() {
    for { // ← Nunca se rompe. No hay Ctrl+C handling ni retorno.
        if !consola.Cajero.Sesion.Activa {
            consola.menuLogin()  // ← Error retornado se ignora silenciosamente
        } else {
            consola.menuOperaciones() // ← Idem
        }
    }
}
```

**Problema:** El usuario nunca puede salir de forma graceful. Si `menuLogin()` retorna un error (EOF en stdin por ejemplo), el loop se repite infinitamente sin mostrar nada al usuario.

**Fix:**

```go
func (consola *consolaa) Iniciar() {
    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
    defer stop()

    for {
        select {
        case <-ctx.Done():
            fmt.Println("\nSesión terminada.")
            return
        default:
        }

        if !consola.Cajero.Sesion.Activa {
            if _, err := consola.menuLogin(); err != nil {
                fmt.Printf("Error: %v\n", err)
                time.Sleep(1 * time.Second)
            }
        } else {
            if _, err := consola.menuOperaciones(); err != nil {
                fmt.Printf("Error: %v\n", err)
                time.Sleep(1 * time.Second)
            }
        }
    }
}
```

---

### BUG P0-2: Error nil retornado con sesión fallida

**Archivos:** `internal/servicio/cajero_service.go:30-32` + `internal/presentacion/consola.go:84-86`

En el servicio:

```go
if !validacionTarjeta || !validacionPin {
    return *sesion, err  // ← err es NIL si ValidarPin no falló técnicamente
}
```

En la consola:

```go
if !sesion.Activa {
    fmt.Println("La sesión no está activa")
    return false, err  // ← err es nil aquí. El error se silencia completamente.
}
```

**Problema:** Cuando la tarjeta o PIN son inválidos pero no hubo error técnico, se retorna `(sesion{Activa:false}, nil)`. El caller no puede distinguir entre "no hubo error" y "autenticación fallida".

**Fix en servicio:**

```go
if !validacionTarjeta || !validacionPin {
    return entidad.Sesion{}, fmt.Errorf("credenciales inválidas: tarjeta o PIN incorrectos")
}
```

---

### BUG P0-3: Desincronización de saldo entre service y storage

**Archivos:** `internal/servicio/cajero_service.go:74` y `internal/servicio/cajero_service.go:96`

```go
// En ExtraerEfectivo (linea 74):
cajeroService.Sesion.CuentaCliente.Saldo -= monto  // ← Modifica un snapshot local

// En RealizarTransferencia (linea 96):
cajeroService.Sesion.CuentaCliente.Saldo -= monto  // ← Mismo problema
```

**Problema:** El storage (`RedCentralAlmacenamiento`) ya debitó el saldo en su mapa interno. El service modifica su copia local en `Sesion.CuentaCliente`. Si luego `ConsultarSaldo()` lee del storage, el saldo mostrado no coincide con el de la sesión. El servicio no debería cachear estado mutable perteneciente al storage.

**Fix:** No restar saldo en el service. El servicio debe re-leer del storage cuando necesite el saldo actualizado, o la sesión no debería contener datos mutables del dominio.

---

### BUG P0-4: Transaction ID colisionable

**Archivo:** `internal/almacenamiento/red_central_almacenamiento.go:88-89`

```go
id := len(redCentral.Transacciones)
redCentral.Transacciones[strconv.Itoa(id)] = id
```

**Problema:** Si se eliminan transacciones en el futuro o hay acceso concurrente, `len()` puede reutilizar IDs ya asignados. Un mapa `map[string]int` para tracking de IDs es innecesario y frágil.

**Fix:** Usar un `atomic.Int64` como counter y una lista real de transacciones:

```go
type RedCentralAlmacenamiento struct {
    Cuentas       map[string]entidad.Cuenta
    Transacciones []entidad.Transaccion
    nextID        atomic.Int64
}
```

---

### RIESGO P1-1: `go.mod` declara versión inexistente

**Archivo:** `go.mod:3`

```
go 1.27.0
```

Go 1.27 no existe (la versión estable más reciente es 1.22/1.23). Esto puede causar fallos de compilación o comportamiento indefinido.

**Fix:** Cambiar a `go 1.22.0`.

---

### RIESGO P1-2: Typos en nombres públicos

| Ubicación | Actual | Correcto |
| --- | --- | --- |
| `cajero_service.go:16` | `AuntenticarUsuario` | `AutenticarUsuario` |
| `sucursar_almacenamiento.go` (archivo) | `sucursar_almacenamiento.go` | `sucursal_almacenamiento.go` |
| `helpers/IngresoNumEntero.go` (archivo) | `IngresoNumEntero.go` | `ingreso_entero.go` |
| `helpers/CallClear.go` (archivo) | `CallClear.go` | `call_clear.go` |

---

### RIESGO P1-3: Mapas inicializados como nil en el futuro

Los constructores `NuevoSucursalAlmacenamiento` y `NuevoRedCentralAlmacenamiento` sí usan `make()` correctamente. Sin embargo, si alguien crea la struct directamente (`RedCentralAlmacenamiento{}`) obtendrá un panic por nil map. Se recomienda usarconstructores privados o validaciones.

---

## 2. [ARQUITECTURA] Clean Architecture

### ISSUE A1: `consola.go` depende de struct concreta, no de interfaz

**Archivo:** `internal/presentacion/consola.go:32`

```go
type consolaa struct {
    Cajero *servicio.CajeroService  // ← Acoplamiento directo al concreto
}
```

**Problema:** No puedes testear la capa de presentación sin levantar todo el servicio con sus dependencias. Viola el Dependency Inversion Principle.

**Fix:** Extraer una interfaz:

```go
// internal/presentacion/ports.go (nuevo archivo)
package presentacion

import (
    "context"
    "github.com/nicos92/dis_cajero/internal/dominio/entidad"
)

type Cajero interface {
    AutenticarUsuario(ctx context.Context, tarjeta entidad.Tarjeta, cuenta entidad.Cuenta) (entidad.Sesion, error)
    ConsultarSaldo(ctx context.Context) (float64, entidad.EstadoTransaccion, error)
    ExtraerEfectivo(ctx context.Context, monto float64) (entidad.EstadoTransaccion, int, error)
    RealizarTransferencia(ctx context.Context, cuentaDestino entidad.Cuenta, monto float64) (entidad.EstadoTransaccion, int, error)
    CerrarSesion(ctx context.Context) (bool, error)
}
```

Y en `consola.go`:

```go
type consolaa struct {
    Cajero Cajero  // ← Ahora es una interfaz
}
```

---

### ISSUE A2: `Sesion` es mutable state acoplado al servicio

**Archivo:** `internal/servicio/cajero_service.go:13`

```go
type CajeroService struct {
    Sesion entidad.Sesion  // ← Estado mutable acoplado al servicio
}
```

**Problema:** El servicio no debería ser stateful. La sesión debería pasarse como parámetro o vivir en un contexto/manager separado. Esto rompe la posibilidad de testing paralelo y viola el principio de que los usecases son stateless.

**Fix:** Pasar `Sesion` como parámetro o inyectar un `SessionManager`:

```go
type CajeroService struct {
    ClienteRedCentral ports.ClienteRedCentral
    ClienteSucursal   ports.ClienteSucursal
}

func (s *CajeroService) ConsultarSaldo(ctx context.Context, sesion entidad.Sesion) (float64, entidad.EstadoTransaccion, error) {
    // ...
}
```

---

### ISSUE A3: `Sesion` contiene `Cuenta` completa (incluyendo PIN y Saldo)

**Archivo:** `internal/dominio/entidad/sesion.go:4`

```go
type Sesion struct {
    CuentaCliente Cuenta  // ← Snapshot completo, incluyendo Pin y Saldo
}
```

**Problema:** La sesión expone el PIN y el saldo como datos estáticos. Si el storage actualiza el saldo, la sesión queda obsoleta. Además, el PIN no debería residir en ningún objeto en memoria más de lo estrictamente necesario.

**Fix:** `Sesion` debería contener solo un identificador:

```go
type Sesion struct {
    CuentaNumero string  // Solo el ID, no la entidad completa
    Activa       bool
}
```

---

### ISSUE A4: Dead code en `presentacion/helpers`

`LeerEntero()` y `LeerTexto()` en `helpers/` nunca son llamados desde `consola.go`, que reimplementa `leerTexto` internamente. Esto es código muerto que genera confusión.

**Recomendación:** Eliminar `helpers/` o integrar sus funciones en `consola.go`.

---

### Lo que está bien

| Aspecto | Evaluación |
| --- | --- |
| `dominio/entidad` no importa infrastructure | Correcto |
| `servicio/ports` define interfaces | Correcto |
| `almacenamiento` implementa los ports | Correcto |
| Flujo `main → servicio → ports → almacenamiento` | Correcto |
| Entidades simples y sin lógica | Correcto |

---

## 3. [MEJORAS GO IDIOMÁTICO]

### IDIOM-1: Errores sin envolver con `%w`

Múltiples archivos usan `errors.New()` puro sin wrapping, lo que impide usar `errors.Is()` y `errors.As()` en los callers.

```go
// red_central_almacenamiento.go:49
return false, errors.New("La cuenta no existe")

// cajero_service.go:45
return 0, entidad.EstadoCancelado, errors.New("La cuenta sesion no esta activa")
```

**Fix:** Definir errores sentinelas y envolver con `%w`:

```go
// internal/dominio/entidad/errors.go (nuevo archivo)
package entidad

import "errors"

var (
    ErrSesionInactiva      = errors.New("sesión no activa")
    ErrSaldoInsuficiente   = errors.New("saldo insuficiente")
    ErrCredencialesInvalidas = errors.New("credenciales inválidas")
    ErrCuentaNoExiste      = errors.New("cuenta no existe")
    ErrTarjetaInvalida     = errors.New("tarjeta inválida")
    ErrTarjetaVencida      = errors.New("tarjeta vencida")
)
```

Uso:

```go
return false, fmt.Errorf("cuenta %s no existe: %w", cuenta.Numero, ErrCuentaNoExiste)
```

---

### IDIOM-2: Sin `context.Context` en ningún lado

Ninguna función del servicio o ports recibe `context.Context`. Esto impide:

- Timeouts en operaciones de storage (cuando se conecte a SQLite real)
- Cancelación graceful de operaciones
- Propagación de tracing/logging estructurado

**Fix de signatures:**

```go
// En ports/cliente_red_central.go
ValidarPin(ctx context.Context, cuenta entidad.Cuenta) (bool, error)
ConsultarSaldo(ctx context.Context, cuenta entidad.Cuenta) (float64, error)

// En servicio/cajero_service.go
AutenticarUsuario(ctx context.Context, tarjeta entidad.Tarjeta, cuenta entidad.Cuenta) (entidad.Sesion, error)
```

---

### IDIOM-3: Manejo inconsistente de punteros vs valores

En `cajero_service.go:18`:

```go
sesion := new(entidad.Sesion)  // ← Crea puntero
sesion.Activa = false          // ← Modifica via puntero
// ...
cajeroService.Sesion = *sesion // ← Copia por valor al campo
```

**Recomendación:** Para entidades pequeñas (`Sesion`, `Cuenta`), pasar por valor es correcto y más eficiente. Pero el patrón actual de crear un puntero, modificarlo, y desreferenciarlo es innecesariamente complejo:

```go
// Más simple y directo:
sesion := entidad.Sesion{Activa: false}
```

---

### IDIOM-4: Convenciones de nombres en Go

| Ubicación actual | Convención Go | Razón |
| --- | --- | --- |
| `sucursar_almacenamiento.go` | `sucursal_almacenamiento.go` | Typo |
| `helpers/IngresoNumEntero.go` | `helpers/ingreso_entero.go` | Go usa snake_case en archivos |
| `helpers/CallClear.go` | `helpers/call_clear.go` | Go usa snake_case en archivos |
| `AuntenticarUsuario` | `AutenticarUsuario` | Typo |
| Paquete `helpers` | Paquete `input` o `ui` | "helpers" es un anti-pattern de nombre |

---

### IDIOM-5: Logging con `fmt.Println` en lugar de `log/slog`

```go
// Actual (múltiples archivos)
fmt.Printf("Error al consultar el saldo: %s\n", err.Error())
fmt.Println("Sesión cerrada. Puede retirar su tarjeta.")
```

**Fix con `log/slog` (stdlib desde Go 1.21):**

```go
import "log/slog"

slog.Error("consulta de saldo fallida",
    "error", err,
    "cuenta", sesion.CuentaCliente.Numero,
)
slog.Info("sesión cerrada", "cuenta", cuenta.Numero)
```

Ventajas: output estructurado (JSON), nivel de log configurable, integración con herramientas de observabilidad.

---

### IDIOM-6: Mensajes de error con mayúscula inicial

En Go, los mensajes de error no empiezan con mayúscula (convención `go vet` / `golint`):

```go
// Mal
errors.New("La cuenta no existe")
errors.New("Saldo insuficiente en la cuenta.")

// Correcto
errors.New("cuenta no existe")
errors.New("saldo insuficiente")
```

---

### IDIOM-7: Uso de `var` implícito en `TipoTransaccion`

**Archivo:** `internal/dominio/entidad/transaccion.go:26-28`

```go
const (
    Consulta      = "Consulta"     // ← Tipo inferido como string, no TipoTransaccion
    Extraccion    = "Extraccion"
    Transferencia = "Transferencia"
)
```

**Fix:** Tipar correctamente las constantes:

```go
const (
    Consulta      TipoTransaccion = "Consulta"
    Extraccion    TipoTransaccion = "Extraccion"
    Transferencia TipoTransaccion = "Transferencia"
)
```

---

## 4. [NUEVAS HERRAMIENTAS] Stack Moderno para CLI

### T1: Middleware pattern para logging, recovery y métricas

**Librería necesaria:** Ninguna extra. Se implementa con闭包 y interfaces de Go.

```go
// internal/servicio/middleware.go
package servicio

import (
    "context"
    "log/slog"
    "time"

    "github.com/nicos92/dis_cajero/internal/dominio/entidad"
    "github.com/nicos92/dis_cajero/internal/servicio/ports"
)

type CajeroService struct {
    ClienteRedCentral ports.ClienteRedCentral
    ClienteSucursal   ports.ClienteSucursal
}

type CajeroMiddleware func(next *CajeroService) *CajeroService

// LoggingMiddleware registra cada operación con duración y resultado
func LoggingMiddleware(logger *slog.Logger) CajeroMiddleware {
    return func(next *CajeroService) *CajeroService {
        // Implementar wrappers que loguean antes/después de cada método
        // usando un wrapper struct que delega al next
        return next
    }
}

// RecoveryMiddleware captura panics y los convierte en errores
func RecoveryMiddleware(logger *slog.Logger) CajeroMiddleware {
    return func(next *CajeroService) *CajeroService {
        return next
    }
}
```

Ejemplo de wrapper concreto:

```go
type loggingRedCentral struct {
    inner  ports.ClienteRedCentral
    logger *slog.Logger
}

func (l *loggingRedCentral) ConsultarSaldo(ctx context.Context, cuenta entidad.Cuenta) (float64, error) {
    start := time.Now()
    saldo, err := l.inner.ConsultarSaldo(ctx, cuenta)
    l.logger.Info("consultar_saldo",
        "cuenta", cuenta.Numero,
        "saldo", saldo,
        "duration_ms", time.Since(start).Milliseconds(),
        "error", err,
    )
    return saldo, err
}
```

Encadenamiento en `main.go`:

```go
var cajero ports.CajeroInterface = servicio.NewCajeroService(sucursal, redCentral)
cajero = servicio.ApplyMiddleware(cajero,
    servicio.LoggingMiddleware(logger),
    servicio.RecoveryMiddleware(logger),
)
```

---

### T2: CLI Framework — Cobra + Viper

**Instalación:**

```bash
go get github.com/spf13/cobra
go get github.com/spf13/viper
```

**Por qué:** `flag` de stdlib es funcional pero limitado. Cobra ofrece subcomandos, help automático, validación de args, y shell completion. Viper maneja config desde archivos `.yaml`, variables de entorno y flags de forma unificada.

**Ejemplo de migración:**

```go
// cmd/cajero/main.go
package main

import (
    "context"
    "log/slog"
    "os"
    "os/signal"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

func main() {
    var cfg struct {
        DBPath   string `mapstructure:"db_path"`
        LogLevel string `mapstructure:"log_level"`
    }

    rootCmd := &cobra.Command{
        Use:   "cajero",
        Short: "Simulador de cajero automático",
        RunE: func(cmd *cobra.Command, args []string) error {
            ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt)
            defer stop()

            logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
                Level: parseLogLevel(cfg.LogLevel),
            }))

            // Setup de dependencias...
            _ = ctx
            _ = logger
            return consola.Iniciar(ctx)
        },
    }

    rootCmd.Flags().StringVar(&cfg.DBPath, "db", "cajero.db", "Ruta a la base de datos SQLite")
    rootCmd.Flags().StringVar(&cfg.LogLevel, "log-level", "info", "Nivel de log (debug, info, warn, error)")

    viper.AutomaticEnv()
    viper.BindPFlag("db_path", rootCmd.Flags().Lookup("db"))
    viper.BindPFlag("log_level", rootCmd.Flags().Lookup("log-level"))

    if err := rootCmd.Execute(); err != nil {
        os.Exit(1)
    }
}
```

---

### T3: Validación de DTOs — `go-playground/validator`

**Instalación:**

```bash
go get github.com/go-playground/validator/v10
```

**Por qué:** Actualmente no hay validación de inputs. El usuario puede ingresar strings vacíos, PINs de 1 dígito, o montos negativos. `validator` permite declarar reglas con tags en structs.

**Ejemplo:**

```go
// internal/dominio/dto/login.go (nuevo archivo)
package dto

import "github.com/go-playground/validator/v10"

type LoginRequest struct {
    NumeroTarjeta string `validate:"required,min=6,max=20"`
    NumeroCuenta  string `validate:"required,min=4,max=20"`
    Pin           string `validate:"required,min=4,max=8,numeric"`
}

type ExtraccionRequest struct {
    Monto float64 `validate:"required,gt=0,lte=100000"`
}

type TransferenciaRequest struct {
    CuentaDestino string  `validate:"required,min=4,max=20"`
    Monto         float64 `validate:"required,gt=0,lte=500000"`
}

var validate = validator.New()

func (r LoginRequest) Validate() error {
    return validate.Struct(r)
}
```

Uso en el servicio:

```go
func (s *CajeroService) AutenticarUsuario(req dto.LoginRequest) (entidad.Sesion, error) {
    if err := req.Validate(); err != nil {
        return entidad.Sesion{}, fmt.Errorf("datos de entrada inválidos: %w", err)
    }
    // ...
}
```

---

### T4: TUI interactiva — BubbleTea + Lipgloss

**Instalación:**

```bash
go get github.com/charmbracelet/bubbletea
go get github.com/charmbracelet/lipgloss
```

**Por qué:** Si la consola es interactiva (menús, inputs, formateo visual), BubbleTea提供 un state machine elegante con soporte de teclado, rendering diferencial, y estilos con Lipgloss. Reemplaza los `fmt.Println` + `bufio` actuales.

**Ejemplo de migración:**

```go
package presentacion

import (
    "fmt"
    "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

var (
    estiloExito  = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true)
    estiloError  = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true)
    estiloTitulo = lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Bold(true)
)

type estadoMenu int

const (
    estadoLogin estadoMenu = iota
    estadoMenuPrincipal
    estadoConsultarSaldo
    estadoExtraer
    estadoTransferir
    estadoSalir
)

type model struct {
    sesion       entidad.Sesion
    estado       estadoMenu
    input        string
    resultado    string
    cajero       ports.Cajero
}

func (m model) Init() tea.Cmd {
    return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "q", "ctrl+c":
            return m, tea.Quit
        case "enter":
            return m.handleEnter()
        default:
            m.input += msg.String()
            return m, nil
        }
    }
    return m, nil
}

func (m model) View() string {
    switch m.estado {
    case estadoLogin:
        return m.viewLogin()
    case estadoMenuPrincipal:
        return m.viewMenuPrincipal()
    default:
        return ""
    }
}
```

**Nota:** La migración completa a BubbleTea es un esfuerzo significativo. Se recomienda hacerlo incrementalmente, manteniendo la versión de `fmt` actual como fallback.

---

### T5: Testing — Table-driven tests + testify

**Instalación:**

```bash
go get github.com/stretchr/testify
go get github.com/stretchr/testify/mock  # Para mocks
```

**Por qué:** Actualmente hay **0 archivos de test** en todo el proyecto. La interfaz `ports` facilita mocking. Table-driven tests son la convención estándar en Go.

**Ejemplo completo:**

```go
// internal/servicio/cajero_service_test.go
package servicio_test

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"

    "github.com/nicos92/dis_cajero/internal/dominio/entidad"
    "github.com/nicos92/dis_cajero/internal/servicio"
)

// Mock de ClienteSucursal
type mockSucursal struct {
    validarFn func(entidad.Tarjeta) (bool, error)
}

func (m *mockSucursal) ValidarTarjeta(t entidad.Tarjeta) (bool, error) {
    return m.validarFn(t)
}

// Mock de ClienteRedCentral
type mockRedCentral struct {
    validarPinFn    func(entidad.Cuenta) (bool, error)
    consultarSaldoFn func(entidad.Cuenta) (float64, error)
}

func (m *mockRedCentral) ValidarPin(c entidad.Cuenta) (bool, error) {
    return m.validarPinFn(c)
}

func (m *mockRedCentral) ConsultarSaldo(c entidad.Cuenta) (float64, error) {
    return m.consultarSaldoFn(c)
}

func TestConsultarSaldo(t *testing.T) {
    tests := []struct {
        name       string
        sesion     entidad.Sesion
        saldo      float64
        saldoErr   error
        wantSaldo  float64
        wantErr    error
    }{
        {
            name:    "sesion inactiva retorna error",
            sesion:  entidad.Sesion{Activa: false},
            wantErr: entidad.ErrSesionInactiva,
        },
        {
            name:      "sesion activa retorna saldo",
            sesion:    entidad.Sesion{Activa: true, CuentaCliente: entidad.Cuenta{Numero: "1111"}},
            saldo:     234.234,
            wantSaldo: 234.234,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            rc := &mockRedCentral{
                consultarSaldoFn: func(c entidad.Cuenta) (float64, error) {
                    return tt.saldo, tt.saldoErr
                },
            }

            svc := &servicio.CajeroService{
                ClienteRedCentral: rc,
                Sesion:            tt.sesion,
            }

            saldo, estado, err := svc.ConsultarSaldo(context.Background())

            if tt.wantErr != nil {
                assert.ErrorIs(t, err, tt.wantErr)
                assert.Equal(t, entidad.EstadoCancelado, estado)
            } else {
                require.NoError(t, err)
                assert.Equal(t, tt.wantSaldo, saldo)
                assert.Equal(t, entidad.EstadoConfirmado, estado)
            }
        })
    }
}
```

---

## Resumen de Prioridades

| Prioridad | Cantidad | Items |
| ----------- | ---------- | ------- |
| **P0 (Crítico)** | 4 | Loop sin salida, error nil, desincronización saldo, ID colisionable |
| **P1 (Alta)** | 4 | go.mod inválido, typos, dead code, convenciones de nombre |
| **P2 (Media)** | 3 | context.Context, slog, wrapping de errores |
| **P3 (Nice to have)** | 5 | Cobra/Viper, BubbleTea, validator, testify, middleware |

---

## Checklist: 5 Tareas para Implementar Hoy

- [ ] **1.** Fix error nil retornado con sesión fallida en `AuntenticarUsuario` (línea 30-32) + envolver errores con `%w` en `cajero_service.go` y `consola.go`
- [ ] **2.** Agregar signal handling (SIGINT) en `Iniciar()` + manejar errores de `menuLogin`/`menuOperaciones` en `consola.go`
- [ ] **3.** Corregir `go.mod` a `go 1.22.0`, renombrar `AuntenticarUsuario` → `AutenticarUsuario`, renombrar archivo `sucursar` → `sucursal`
- [ ] **4.** Extraer interfaz `Cajero` en `presentacion/` para desacoplar `consola.go` del struct concreto `CajeroService`
- [ ] **5.** Eliminar dead code en `helpers/` o integrarlo. Renombrar archivos a snake_case
