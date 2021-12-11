# Mi solución (parte 2)

En esta segunda parte el objetivo es incluir un límite de tiempo (configurable por el usuario) para completar la prueba.

## Nuevo *flag* para establecer el límite de tiempo

Definimos un nuevo *flag* llamado `limite` para establecer el tiempo límite para realizar la prueba:

```go
func main() {
    csvFile := flag.String("csv", "problems.csv", "Fichero de problemas en formato 'enunciado,respuesta' (CSV)")
    limite := flag.Int("limite", 30, "Tiempo limite para completar la prueba en segundos")
    flag.Parse()
```

## Temporizador

El paquete [`time`](https://pkg.go.dev/time) incluye el tipo `NewTimer`, que permite definr un temporizador.

El temporizador envía un mensaje en su *channel* una vez pasado el tiempo definido.

Para no hacer perder tiempo al usuario mientras todavía estamos realizando el *setup* de la prueba, definimos el *timer* con el valor especificado en el *flag* `limite`, en segundos, después de inicializar la puntuación a cero:

```go
var puntuacion int = 0

temporizador := time.NewTimer(time.Duration(*limite) * time.Second)

for i := range problemas {
    ...
```

## Canales

El problema de los *channel* es que bloquean la ejecución del programa hasta que reciban un mensaje (como se indica, por ejemplo, en el [Go Tour: Channels](https://go.dev/tour/concurrency/2)).

Como la ejecución del programa también se detiene mientras esperamos la respuesta a la pregunta planteada al usuario, modificamos la función `getUserAnswer()`:

- modificamos los *signature* de la función para indicar que ahora no devuelve un *string*
- enviamos la respuesta del usuario al canal `canalRespuestaUsuario`

Como el canal `canalRespuestaUsuario` tiene que estar disponible en la función `main` y en `getUserAnswer()`, lo definimos antes de la función `main`:

```go
var canalRespuestaUsuario = make(chan string)

func main() {
    ...
```

Y la función  getUserAnswer()`:

```go
func getUserAnswer() {

    r := bufio.NewReader(os.Stdin)
    respuestaUsuario, err := r.ReadString('\n')
    if err != nil {
        exit("Error al leer la respuesta.\n[ERROR] %s\n", err, ERROR_NO_ENTIENDO_RESPUESTA)
    }
    canalRespuestaUsuario <- strings.Trim(respuestaUsuario, "\n ")
    return
}
```

## `select`

El temporizador se ejecuta *en segundo plano* hasta que pasa el tiempo especificado y entonces, se escribe en el canal asociado al temporaizador; en mi caso, `temporizador.C`.

Para que la llamada a la función `getUserAnswer()` no bloquee la ejecución del programa, la *convertimos* en una *goroutine* llamándola precedida de `go`:

```go
for i := range problemas {
    fmt.Printf("Pregunta %d:\t %s = ", i+1, problemas[i].enunciado)
    go getUserAnswer()
    ...
```

De esta forma, programa continúa.

En este punto, tenemos que esperar a recibir un mensaje a través de alguno de los canales de las *goroutines* que se ejecutan en paralelo; el *timer* y la obtención de la respuesta por parte del usuario.

Usamos `select` para esperar a la recepción de múltiples operaciones de comunicación, como dice el [Go Tour: Select](https://go.dev/tour/concurrency/5).

Si recibimos respuesta por parte del usuario, comprobamos si es correcta y modificamos la puntuación como corresponda; si se ha alcanzado el tiempo límite de la prueba, mostramos un mensaje.

Una vez finalizado el tiempo de la prueba, deberíamos acabar. Si usamos `break`, salimos del bloque `select`, pero seguimos dentro del `for`, por lo que se seguirían mostrando el resto de las preguntas pendientes... Y eso no es lo que queremos.

Si se ha acabado el tiempo, queremos acabar el programa, pero debemos mostrar el número de preguntas acertadas antes de salir... También queremos mostrar el resultado si no se ha acabado el tiempo pero se han completado todas las respuestas de la prueba...

Mi solución, ha sido mover el mensaje que muestra la puntuación a una función, que se llama en los dos casos; la función muestra la puntuación y termina el program (con éxito):

```go
func printResult(puntuacion int, total int) {
    fmt.Printf("Ha respondido %d respuestas correctas de %d preguntas.\n", puntuacion, total)
    os.Exit(PRUEBA_FINALIZADA)
}

func main() {
    ...
    for i := range problemas {
        fmt.Printf("Pregunta %d:\t %s = ", i+1, problemas[i].enunciado)
        go getUserAnswer()

        select {
        case msg := <-canalRespuestaUsuario:
            if problemas[i].respuesta == msg {
                puntuacion++
            }
        case <-temporizador.C:
            fmt.Println("\n¡Campana y se acabó!")
            printResult(puntuacion, len(problemas))
        }
    }
    printResult(puntuacion, len(problemas))
}
```

## Cierre del fichero

Revisando el código final he visto que el fichero CSV desde el que se cargan las preguntas del test se dejaba abierto :(

Para hacer que Go cierre el fichero al salir de la función (en este caso, `main()`, uso `defer`:

```go
...
var fHandle *os.File = openFile(*csvFile)
defer fHandle.Close()

r := csv.NewReader(fHandle)
...
```
