Rol y Contexto:
Actua como un Desarrollador Senior Backend especializado en Go (Golang) y en arquitectura basada en dominios (Domain-Driven Design). Necesito refactorizar el flujo de autenticacion y seleccion de cuentas en una aplicacion de consola que simula un cajero automatico (ATM).

Requerimientos Funcionales a Implementar:
1. Validacion de Tarjeta: Al iniciar el flujo, el usuario ingresa el numero de tarjeta. El sistema debe validar inmediatamente si la tarjeta se encuentra vencida.
2. Autenticacion: Si la tarjeta NO esta vencida, el sistema debe solicitar la contraseña asociada a la tarjeta.
3. Seleccion de Cuenta: Una vez autenticado exitosamente el usuario, el sistema debe listar y permitir elegir entre las cuentas bancarias asociadas y registradas para esa tarjeta/cliente.
4. Operaciones: Despues de seleccionar la cuenta, el flujo debe continuar permitiendo realizar las operaciones bancarias habituales ya existentes en la aplicacion.

Directrices Tecnicas:
- Lenguaje: Go (Golang).
- Interfaz: Aplicacion de consola (CLI).
- Arquitectura: Arquitectura basada en dominios (separacion clara entre capa de dominio, casos de uso/servicios de aplicacion, infraestructura y presentacion/consola).
- Manejo de Errores: Utilizar manejo de errores idiomatico en Go, tipados y con mensajes claros para la interfaz de consola.
- Clean Code: Respetar los principios SOLID, mantener la logica de negocio desacoplada de la entrada/salida por consola (inyeccion de dependencias).

Entregable Esperado:
- Actualizacion de los casos de uso / servicios de dominio necesarios para soportar este nuevo orden de validacion.
- Ajuste en la capa de presentacion de consola para guiar paso a paso al usuario segun los nuevos requerimientos.
- Codigo limpio, comentado donde sea estrictamente necesario y listo para integrar.
