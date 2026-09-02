package presentacion

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/nicos92/dis_cajero/internal/dominio/entidad"
	"github.com/nicos92/dis_cajero/internal/servicio"
)

type Consola struct {
	Cajero *servicio.CajeroService
	Bufio  *bufio.Reader
}

func NuevaConsola(cajero *servicio.CajeroService) *Consola {
	return &Consola{
		Cajero: cajero,
		Bufio:  bufio.NewReader(os.Stdin),
	}
}

func (consola *Consola) leerTexto(mensaje string) (string, error) {
	fmt.Print(mensaje)
	entrada, err := consola.Bufio.ReadString('\n')

	if err != nil {
		return "", errors.New("Error al leer la entrada de texto")
	}

	return strings.TrimSpace(entrada), nil
}

func (consola *Consola) menuLogin() (bool, error) {
	fmt.Println("Bienvenido a Banco Central")
	nroTarjeta, err := consola.leerTexto("Ingrese el número de la tarjeta")

	if err != nil {
		return false, err
	}

	nroCuenta, err := consola.leerTexto("Ingrese el número de la cuenta")

	if err != nil {
		return false, err
	}

	nroPin, err := consola.leerTexto("Ingrese el número de PIN")

	if err != nil {
		return false, err
	}

	tarjeta := entidad.Tarjeta{Numero: nroTarjeta}
	cuenta := entidad.Cuenta{Numero: nroCuenta, Pin: nroPin}

	sesion, err := consola.Cajero.AuntenticarUsuario(tarjeta, cuenta)

	if err != nil {
		fmt.Printf("Error de autenticacion: %s\n", err.Error())
		return false, err
	}
	if !sesion.Activa {
		fmt.Println("La sesión no está activa")
		return false, err
	}
	fmt.Println("Inicio de sesión exitoso")

	return true, nil
}

func (consola *Consola) menuOperaciones() (bool, error) {
	fmt.Printf("Bienvenido\nIngresaste con la cuenta %s\n", consola.Cajero.Sesion.CuentaCliente.Numero)
	fmt.Println("\n--- Menú Principal ---")
	fmt.Println("1. Consultar Saldo")
	fmt.Println("2. Extraer Efectivo")
	fmt.Println("3. Transferir Dinero")
	fmt.Println("4. Salir / Cerrar Sesión")

	opcion, err := consola.leerTexto("Ingrese una opción")

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
			fmt.Printf("Consulta: %v\nSu saldo es: %v$\n", estado, saldo)

		}
	case "2":
		{
			monto, err := consola.leerTexto("Ingrese el monto que sea retirar")
			if err != nil {
				fmt.Printf("Fallo: %s\n", err.Error())
				return false, err
			}
			nuevoMonto, err := strconv.ParseFloat(monto, 64)
			if err != nil {
				fmt.Printf("Fallo al convertir a número: %s\n", err.Error())
				return false, err

			}
			estado, id, err := consola.Cajero.ExtraerEfectivo(nuevoMonto)
			if err != nil {
				fmt.Printf("Fallo: %s\n", err.Error())
				return false, err

			}

			fmt.Printf("Extración número %d: %v\n", id, estado)

		}
	case "3":
		{
			monto, err := consola.leerTexto("Ingrese el monto que sea transferir")
			if err != nil {
				fmt.Printf("Fallo: %s\n", err.Error())
				return false, err
			}
			nuevoMonto, err := strconv.ParseFloat(monto, 64)
			if err != nil {
				fmt.Printf("Fallo al convertir a número: %s\n", err.Error())
				return false, err

			}

			nroCuentaDestino, err := consola.leerTexto("Ingrese el número de cuenta al que desea transferir")
			if err != nil {
				fmt.Printf("Fallo: %s\n", err.Error())
				return false, err
			}
			destino := entidad.Cuenta{Numero: nroCuentaDestino}
			estado, id, err := consola.Cajero.RealizarTransferencia(destino, nuevoMonto)
			if err != nil {
				fmt.Printf("Fallo: %s\n", err.Error())
				return false, err

			}

			fmt.Printf("Transferencia número %d: %v\n", id, estado)
		}
	case "4":
		{
			consola.Cajero.CerrarSesion()
			fmt.Println("Sesión cerrada. Puede retirar su tarjeta.")
		}
	default:
		fmt.Println("Esta opción no es valida.")
	}
	return true, nil
}

func (consola *Consola) Iniciar() {
	for {
		if !consola.Cajero.Sesion.Activa {
			consola.menuLogin()
		} else {
			consola.menuOperaciones()
		}
	}
}
