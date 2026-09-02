package helpers

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func LeerEntero(mensaje string) (int, error) {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println(mensaje)

	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return 0, err
		}
		return 0, fmt.Errorf("lectura de entrada finalizada")
	}

	entrada := strings.TrimSpace(scanner.Text())

	numero, err := strconv.Atoi(entrada)
	if err != nil {
		// return 0, fmt.Errorf("ingreso no válido: %w", err)
		return 0, fmt.Errorf("Ingreso no válido")
	}
	return numero, nil
}

func LeerTexto(mensaje string, buf bufio.Reader) (string, error) {
	fmt.Println(mensaje)
	entrada, err := buf.ReadString('\n')

	if err != nil {
		return "", errors.New("Error al leer la entrada de texto.")
	}

	return strings.TrimSpace(entrada), nil
}
