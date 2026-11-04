package presentacion

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/nicos92/dis_cajero/internal/dominio/entidad"
	"github.com/nicos92/dis_cajero/internal/presentacion/inputs"
	"github.com/nicos92/dis_cajero/internal/presentacion/ports"
)

// Códigos de color ANSI
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	White  = "\033[37m"

	// Estilos
	Bold = "\033[1m"
)

type consolaa struct {
	Cajero ports.Cajero
	Bufio  *bufio.Reader
}

func NuevaConsola(cajero ports.Cajero) *consolaa {
	return &consolaa{
		Cajero: cajero,
		Bufio:  bufio.NewReader(os.Stdin),
	}
}

func (consola *consolaa) leerTexto(mensaje string) (string, error) {
	fmt.Print(mensaje)
	entrada, err := consola.Bufio.ReadString('\n')

	if err != nil {
		return "", errors.New("error al leer la entrada de texto")
	}

	return strings.TrimSpace(entrada), nil
}

func (consola *consolaa) menuLogin() (bool, error) {
	inputs.CallClear()
	fmt.Println(Green + Bold + "Bienvenido al Banco Central" + Reset)

	nroTarjeta, err := consola.leerTexto("Ingrese el número de la tarjeta: ")
	if err != nil {
		return false, err
	}

	tarjeta := entidad.Tarjeta{Numero: nroTarjeta}

	if err := consola.Cajero.ValidarTarjeta(tarjeta); err != nil {
		return false, err
	}

	fmt.Println(Green + "Tarjeta válida" + Reset)

	nroPin, err := consola.leerTexto("Ingrese el número de PIN: ")
	if err != nil {
		return false, err
	}

	if _, err := consola.Cajero.AutenticarUsuario(tarjeta, nroPin); err != nil {
		return false, err
	}

	cuentas, err := consola.Cajero.ObtenerCuentasAsociadas()
	if err != nil {
		return false, err
	}

	if err := consola.seleccionarCuenta(cuentas); err != nil {
		return false, err
	}

	fmt.Println(Green + "Inicio de sesión exitoso" + Reset)
	time.Sleep(1000 * time.Millisecond)
	inputs.CallClear()

	return true, nil
}

func (consola *consolaa) seleccionarCuenta(cuentas []entidad.Cuenta) error {
	fmt.Println(Bold + Yellow + "\nCuentas asociadas a la tarjeta" + Reset)

	for indice, cuenta := range cuentas {
		estado := "activa"
		if !cuenta.Activa {
			estado = "inactiva"
		}

		fmt.Printf("%d. Cuenta %s (%s)\n", indice+1, cuenta.Numero, estado)
	}

	opcion, err := consola.leerTexto("Seleccione una cuenta: ")
	if err != nil {
		return err
	}

	indice, err := strconv.Atoi(opcion)
	if err != nil || indice < 1 || indice > len(cuentas) {
		return errors.New("opción de cuenta inválida")
	}

	_, err = consola.Cajero.SeleccionarCuenta(cuentas[indice-1].Numero)

	return err
}

func (consola *consolaa) menuOperaciones() (bool, error) {
	sesion := consola.Cajero.SesionActual()

	fmt.Printf("%sBienvenido\nIngresaste con la cuenta %s%s\n", Green, sesion.CuentaCliente.Numero, Reset)
	fmt.Println(Bold + Yellow + "\nMenú Principal" + Reset)
	fmt.Println("1. Consultar Saldo")
	fmt.Println("2. Extraer Efectivo")
	fmt.Println("3. Transferir Dinero")
	fmt.Println("4. Salir / Cerrar Sesión")

	opcion, err := consola.leerTexto("Ingrese una opción: ")

	if err != nil {
		return false, err
	}

	switch opcion {
	case "1":
		{
			saldo, estado, err := consola.Cajero.ConsultarSaldo()
			if err != nil {
				fmt.Printf("Error al consultar el saldo: %s\n", err.Error())
				return false, err
			}
			fmt.Printf("%sConsulta número: %v\nSu saldo es: %v$\n%s", Green, estado, saldo, Reset)

		}
	case "2":
		{
			monto, err := consola.leerTexto("Ingrese el monto que desea retirar: ")
			if err != nil {
				return false, err
			}
			nuevoMonto, err := strconv.ParseFloat(monto, 64)
			if err != nil {
				return false, errors.New("error al ingresar el monto")

			}
			estado, id, err := consola.Cajero.ExtraerEfectivo(nuevoMonto)
			if err != nil {
				return false, err

			}

			fmt.Printf("Extracción número %d: %v\n", id, estado)

		}
	case "3":
		{
			monto, err := consola.leerTexto("Ingrese el monto que desea transferir: ")
			if err != nil {
				fmt.Printf("Fallo: %s\n", err.Error())
				return false, err
			}
			nuevoMonto, err := strconv.ParseFloat(monto, 64)
			if err != nil {
				fmt.Printf("Fallo al convertir a número: %s\n", err.Error())
				return false, err

			}

			nroCuentaDestino, err := consola.leerTexto("Ingrese el número de cuenta al que desea transferir: ")
			if err != nil {
				fmt.Printf("Fallo: %s\n", err.Error())
				return false, err
			}
			destino := entidad.Cuenta{Numero: nroCuentaDestino}
			estado, id, err := consola.Cajero.RealizarTransferencia(destino, nuevoMonto)
			if err != nil {
				return false, err

			}

			fmt.Printf("Transferencia número %d: %v\n", id, estado)
		}
	case "4":
		{
			if _, err := consola.Cajero.CerrarSesion(); err != nil {
				return false, err
			}
			fmt.Println("Sesión cerrada. Puede retirar su tarjeta.")
		}
	default:
		fmt.Println("Esta opción no es valida.")
	}
	return true, nil
}

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

		if !consola.Cajero.SesionActual().Activa {
			if _, err := consola.menuLogin(); err != nil {
				fmt.Printf("%sError: %v\n", Red, err)
				time.Sleep(2 * time.Second)
			}
		} else {
			if _, err := consola.menuOperaciones(); err != nil {
				fmt.Printf("%sError: %v%s\n", Red, err, Reset)
				time.Sleep(1 * time.Second)
			}
		}
	}
}
