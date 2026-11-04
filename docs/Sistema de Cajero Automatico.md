# Proyecto: Sistema de Cajero Automático

## Lista de requerimientos y seguimiento

### 1\. Diagrama de Contexto

- [ ] Construir el Diagrama de Contexto del sistema **Cajero Automático**.  
- [ ] Incluir los 3 terminadores:  
  - Cliente.  
  - Sucursal.  
  - Red Central.  
- [ ] Incorporar todos los flujos de entrada y salida identificados en la sección 6.1.  
- [ ] Verificar que todos los flujos tengan un nombre.

### 2\. DFD Nivel 0

- [ ] Construir el DFD Nivel 0\.  
- [ ] Incluir los 8 procesos de la sección 6.2.  
- [ ] Incorporar los 3 almacenes de datos:  
  - D1.  
  - D2.  
  - D3.  
- [ ] Representar todos los flujos que conectan procesos y almacenes.  
- [ ] Verificar el balanceo con el Diagrama de Contexto: cada flujo externo del Contexto debe aparecer en el Nivel 0\.

### 3\. Diccionario de Datos (DD)

- [ ] Completar el Diccionario de Datos.  
- [ ] Incluir para cada elemento:  
  - Nombre.  
  - Descripción.  
  - Composición.  
  - Valores permitidos.

### 4\. Diagrama de Casos de Uso

- [ ] Construir el Diagrama de Casos de Uso.  
- [ ] Incluir los 5 casos de uso de la sección 3\.  
- [ ] Incorporar los 3 actores.  
- [ ] Representar correctamente las relaciones `«include»` y `«extend»`.  
- [ ] Ubicar los actores secundarios a la derecha del sistema.

### 5\. Diagrama de Clases

- [ ] Construir el Diagrama de Clases completo.  
- [ ] Incluir las 8 clases de la sección 5\.  
- [ ] Incorporar todos los atributos con tipo y visibilidad.  
- [ ] Incorporar todos los métodos con tipo de retorno y visibilidad.  
- [ ] Representar la herencia de las siguientes clases desde `Operacion`:  
  - Consulta.  
  - Extraccion.  
  - Transferencia.  
- [ ] Representar las relaciones de composición.  
- [ ] Representar las relaciones de asociación.  
- [ ] Indicar multiplicidad en todos los extremos de las relaciones.  
- [ ] Subrayar los nombres de los OBJETOS en el diagrama.

### 6\. Diagrama de Secuencia: CU 1 \- Autenticar Usuario

- [ ] Construir el Diagrama de Secuencia del caso de uso 1\.  
- [ ] Incluir los 4 participantes:  
  - Cliente.  
  - `:PantallaCajero`.  
  - `:SucursalDB`.  
  - `:CentralDB`.  
- [ ] Representar todos los mensajes sincrónicos con sus parámetros.  
- [ ] Incorporar los mensajes de retorno con líneas discontinuas.  
- [ ] Ubicar correctamente los focos de control.

### 7\. Diagrama de Secuencia: Operación

- [ ] Elegir una operación:  
  - [ ] Consulta.  
  - [ ] Extracción.  
  - [ ] Transferencia.  
- [ ] Construir el Diagrama de Secuencia correspondiente.  
- [ ] Incluir todos los participantes necesarios.  
- [ ] Representar los mensajes sincrónicos con parámetros.  
- [ ] Incorporar los mensajes de retorno con líneas discontinuas.  
- [ ] Ubicar correctamente los focos de control.

### 8\. Diagrama de Estado: Clase Sesion

- [ ] Construir el Diagrama de Estado de la clase `Sesion`.  
- [ ] Modelar los siguientes estados:  
  - [ ] Iniciada.  
  - [ ] Autenticando (tarjeta).  
  - [ ] Validando (PIN).  
  - [ ] Operando.  
  - [ ] Cerrando.  
  - [ ] Cerrada.  
- [ ] Representar todas las transiciones entre estados.  
- [ ] Identificar los eventos disparadores de cada transición.

### 9\. Diagrama de Actividad del Sistema

- [ ] Construir el Diagrama de Actividad.  
- [ ] Representar el flujo completo del sistema.  
- [ ] Incluir las decisiones y alternativas del flujo.  
- [ ] Identificar correctamente el inicio y el final del proceso.

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
| 1 | Diagrama de Contexto | Pendiente |
| 2 | DFD Nivel 0 | Pendiente |
| 3 | Diccionario de Datos | Pendiente |
| 4 | Diagrama de Casos de Uso | Pendiente |
| 5 | Diagrama de Clases | Pendiente |
| 6 | Secuencia: Autenticar Usuario | Pendiente |
| 7 | Secuencia: Operación | Pendiente |
| 8 | Diagrama de Estado: Sesion | Pendiente |
| 9 | Diagrama de Actividad | Pendiente |
| 10 | Implementación en Capas | Pendiente |
| 11 | README.md | Pendiente |

**Estados posibles:** Pendiente / En progreso / Completado.

&nbsp;