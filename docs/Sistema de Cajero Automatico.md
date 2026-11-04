# Proyecto: Sistema de Cajero Automático

## Lista de requerimientos y seguimiento

### 1\. Diagrama de Contexto

- [x] Construir el Diagrama de Contexto del sistema **Cajero Automático**.  
- [x] Incluir los 3 terminadores:  
  - Cliente.  
  - Sucursal.  
  - Red Central.  
- [x] Incorporar todos los flujos de entrada y salida identificados en la sección 6.1.  
- [x] Verificar que todos los flujos tengan un nombre.

### 2\. DFD Nivel 0

- [x] Construir el DFD Nivel 0\.  
- [x] Incluir los 8 procesos de la sección 6.2.  
- [x] Incorporar los 3 almacenes de datos:  
  - D1.  
  - D2.  
  - D3.  
- [x] Representar todos los flujos que conectan procesos y almacenes.  
- [x] Verificar el balanceo con el Diagrama de Contexto: cada flujo externo del Contexto debe aparecer en el Nivel 0\.

### 3\. Diccionario de Datos (DD)

- [x] Completar el Diccionario de Datos.  
- [x] Incluir para cada elemento:  
  - Nombre.  
  - Descripción.  
  - Composición.  
  - Valores permitidos.

### 4\. Diagrama de Casos de Uso

- [x] Construir el Diagrama de Casos de Uso.  
- [x] Incluir los 5 casos de uso de la sección 3\.  
- [x] Incorporar los 3 actores.  
- [x] Representar correctamente las relaciones `«include»` y `«extend»`.  
- [x] Ubicar los actores secundarios a la derecha del sistema.

### 5\. Diagrama de Clases

- [x] Construir el Diagrama de Clases completo.  
- [x] Incluir las 8 clases de la sección 5\.  
- [x] Incorporar todos los atributos con tipo y visibilidad.  
- [x] Incorporar todos los métodos con tipo de retorno y visibilidad.  
- [x] Representar la herencia de las siguientes clases desde `Operacion`:  
  - Consulta.  
  - Extraccion.  
  - Transferencia.  
- [x] Representar las relaciones de composición.  
- [x] Representar las relaciones de asociación.  
- [x] Indicar multiplicidad en todos los extremos de las relaciones.  
- [x] Subrayar los nombres de los OBJETOS en el diagrama.

### 6\. Diagrama de Secuencia: CU 1 \- Autenticar Usuario

- [x] Construir el Diagrama de Secuencia del caso de uso 1\.  
- [x] Incluir los 4 participantes:  
  - Cliente.  
  - `:PantallaCajero`.  
  - `:SucursalDB`.  
  - `:CentralDB`.  
- [x] Representar todos los mensajes sincrónicos con sus parámetros.  
- [x] Incorporar los mensajes de retorno con líneas discontinuas.  
- [x] Ubicar correctamente los focos de control.

### 7\. Diagrama de Secuencia: Operación

- [x] Elegir una operación:  
  - [ ] Consulta.  
  - [x] Extracción.  
  - [ ] Transferencia.  
- [x] Construir el Diagrama de Secuencia correspondiente.  
- [x] Incluir todos los participantes necesarios.  
- [x] Representar los mensajes sincrónicos con parámetros.  
- [x] Incorporar los mensajes de retorno con líneas discontinuas.  
- [x] Ubicar correctamente los focos de control.

### 8\. Diagrama de Estado: Clase Sesion

- [x] Construir el Diagrama de Estado de la clase `Sesion`.  
- [x] Modelar los siguientes estados:  
  - [x] Iniciada.  
  - [x] Autenticando (tarjeta).  
  - [x] Validando (PIN).  
  - [x] Operando.  
  - [x] Cerrando.  
  - [x] Cerrada.  
- [x] Representar todas las transiciones entre estados.  
- [x] Identificar los eventos disparadores de cada transición.

### 9\. Diagrama de Actividad del Sistema

- [x] Construir el Diagrama de Actividad.  
- [x] Representar el flujo completo del sistema.  
- [x] Incluir las decisiones y alternativas del flujo.  
- [x] Identificar correctamente el inicio y el final del proceso.

### 10\. Implementación: Solución en Capas

#### 10.1. Capa de Datos

- [ ] Implementar la clase `SucursalDB`.  
- [ ] Implementar la clase `CentralDB`.  
- [ ] Incorporar métodos que simulen la comunicación con los sistemas externos.  
- [ ] Utilizar valores hardcodeados o aleatorios para simular las respuestas.

#### 10.2. Capa de Lógica

- [ ] Implementar la clase abstracta `Operacion`.  
- [ ] Implementar la clase concreta `Consulta`.  
- [ ] Implementar la clase concreta `Extraccion`.  
- [ ] Implementar la clase concreta `Transferencia`.  
- [ ] Implementar el método `ejecutar()` con `override` en las tres operaciones concretas.  
- [ ] Implementar la clase `Cliente`.  
- [ ] Implementar la clase `Tarjeta`.  
- [ ] Implementar la clase `Cuenta`.  
- [ ] Implementar la clase `Sesion`.  
- [ ] Incorporar atributos privados utilizando guion bajo (`_`).  
- [ ] Incorporar propiedades públicas (`get/set`).  
- [ ] Implementar los constructores correspondientes.  
- [ ] Implementar los métodos necesarios en cada clase.

#### 10.3. Capa de Presentación

- [ ] Crear el formulario principal `FormCajero`.  
- [ ] Permitir ingresar el número de tarjeta.  
- [ ] Permitir ingresar el PIN.  
- [ ] Simular el proceso de autenticación.  
- [ ] Mostrar el menú principal del cajero.  
- [ ] Implementar al menos una de las tres operaciones:  
  - [ ] Consulta.  
  - [ ] Extracción.  
  - [ ] Transferencia.  
- [ ] Mostrar el resultado de la operación en pantalla.

### 11\. Documentación: README.md

- [ ] Crear el archivo `README.md`.  
- [ ] Incorporar una descripción general del sistema.  
- [ ] Identificar a los integrantes del equipo.  
- [ ] Explicar cómo ejecutar el proyecto.  
- [ ] Incluir los requisitos necesarios para ejecutar la aplicación.

---

## Seguimiento general

| N.º | Requerimiento | Estado |
| :---- | :---- | :---- |
| 1 | Diagrama de Contexto | Completado |
| 2 | DFD Nivel 0 | Completado |
| 3 | Diccionario de Datos | Completado |
| 4 | Diagrama de Casos de Uso | Completado |
| 5 | Diagrama de Clases | Completado |
| 6 | Secuencia: Autenticar Usuario | Completado |
| 7 | Secuencia: Operación | Completado |
| 8 | Diagrama de Estado: Sesion | Completado |
| 9 | Diagrama de Actividad | Completado |
| 10 | Implementación en Capas | Pendiente |
| 11 | README.md | Pendiente |

**Estados posibles:** Pendiente / En progreso / Completado.

&nbsp;