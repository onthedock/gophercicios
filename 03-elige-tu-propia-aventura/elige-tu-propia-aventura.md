# Elige tu propia aventura

El objetivo del ejercicio es reproducir en formato web uno de los libros de "[Elige tu propia aventura](https://es.wikipedia.org/wiki/Elige_tu_propia_aventura)".

Las historias se proporcionan en formato JSON con la siguiente estructura:

> Al haber cambiado los nombres de los campos en la *struct* que contiene la historia *importada* desde el fichero JSON, el *matching* entre los campos cuyos nombres cambiaron `story --> paragraphs` y `arc --> chapter` fallaba y los campos en el *struct* aparecían vacíos. Esto hacía que no se insertara nigún valor en el *template* HTML, y que el texto de cada *página* de la historia se mostrara en blanco.
>
> Ha sido necesario modificar el fichero `gopher.json` y actualizar los nombres de los campos para que se muestre el valor de los campos en el navegador.

```json
{
  // Each story arc will have a unique key that represents
  // the name of that particular arc.
  "story-arc": {
    "title": "A title for that story arc. Think of it like a chapter title.",
    "story": [
      "A series of paragraphs, each represented as a string in a slice.",
      "This is a new paragraph in this particular story arc."
    ],
    // Options will be empty if it is the end of that
    // particular story arc. Otherwise it will have one or
    // more JSON objects that represent an "option" that the
    // reader has at the end of a story arc.
    "options": [
      {
        "text": "the text to render for this option. eg 'venture down the dark passage'",
        "arc": "the name of the story arc to navigate to. This will match the story-arc key at the very root of the JSON document"
      }
    ]
  },
  ...
}
```

Los campos en el fichero deben ser:

```json
{
  "intro": {
    "title": "The Little Blue Gopher",
    "paragraphs": [
      "Once upon a time, ...",
      "One of his friends ...",
      "On the other hand, ..."
    ],
    "options": [
      {
        "text": "That story about the Sticky Bandits isn't real, it is from Home Alone 2! Let's head to New York.",
        "chapter": "new-york"
      },
      {
        "text": "Gee, those bandits sound pretty real to me. Let's play it safe and try our luck in Denver.",
        "chapter": "denver"
      }
    ]
  }
```

Vamos a importar la historia en una *struct* usando el servicio online [JSON to Go](https://mholt.github.io/json-to-go/), que permite convertir un fichero JSON en la definición de una *struct* en Go.

El resultado es:

```go
type AutoGenerated struct {
    Intro struct {
        Title   string   `json:"title"`
        Story   []string `json:"story"`
        Options []struct {
            Text string `json:"text"`
            Arc  string `json:"arc"`
        } `json:"options"`
    } `json:"intro"`
}
```

## Paquete `cyoa`(*Choose your own adventure*)

Creamos un *package*  `cyoa` y pegamos la definición de la *struct* creada a partir del fichero JSON.

Sin embargo, en vez de tener un *struct* con varios *niveles*, disponer de *structs* separadas proporciona más flexibilidad, por lo que separamos la definición de la *struct* original en dos (también cambiamos algunos nombres):

```go
type Chapter struct {
  Title      string   `json:"title"`
  Paragraphs []string `json:"paragraphs"`
  Options    []Option `json:"options"`
}

type Option struct {
  Text    string `json:"text"`
  Chapter string `json:"chapter"`
}
```

La historia completa será un *map* en el que cada entrada contiene el nombre de un capítulo; algo del estilo `map[string]Chapter`.

Sin embargo, podemos *ocultar* los detalles de implementación definiendo un nuevo tipo de variable `Story`:

```go
type Story map[string]Chapter
```

De esta forma, a nivel conceptual, el desarrollador sólo se ocupa de manipular "historias", independientemente de cómo están implementadas, lo que puede facilitar el desarrollo.

## Servidor web

En general, una buena práctica es crear una carpeta `cmd` para cada uno de los comandos de la aplicación y después construir un binario que ejecutar... En nuestro caso, sólo va a tener uno, el componente encargado de publicar la historia en el navegador.

Creamos la carpeta:

```bash
mkdir -p cmd/cyoaweb
```

Dentro de esta carpeta construimos la aplicación.

Empezamos con la importación del fichero JSON que contiene la historia. El nombre del fichero lo pasamos como parámetro desde la CLI al programa.

Usamos el paquete `flag` y especificamos el fichero `gopher.json` como valor por defecto.

Mostramos el nombre del fichero desde el que cargamos la historia y a continuación abrimos el fichero.

Como el fichero está en formato JSON, usamos `json.NewDecoder()`. `json.NewDecoder` toma como parámetro un *Reader* (en este aso, el puntero al fichero); lo *descodificamos* y si no tenemos un error, lo mostramos por pantalla (para ver que todo funciona).

> Jon usa la forma *compacta*:
>
> ```go
> if err := d.Decode(&story); err != nil {
>      panic(err)
> }
> ```
>
> en la que asigna el resultado de *decodificar* el fichero y comprobar si hay un error en la misma  línea.
>
> Esto es equivalente a:
>
> ```go
> err = d.Decode(&story);
> if err != nil {
>    panic(err)
> }
> ```

Usamos `%+v` en `fmt.Printf()` para mostrar la *struct* con los valores obtenidos desde el fichero JSON (con `+`, también se muestran los nombres de los campos de la *struct*):

```bash
$ go run cmd/cyoaweb/main.go 
Using the story from file gopher.json.
map[debate:{Title:The Great Debate Paragraphs:[] Options:[{Text:Clearly that man in the fox outfit was the winner. Chapter:} {Text:I don't think
...
```

En este punto ya hemos comprobado que somos capaces de cargar el fichero con la historia y *parsearlo* en una *Struct*.

## *Refactor*

Una vez esta parte de la aplicación funcionando, *refactorizamos* para mejorar la calidad del código, de manera que sea más sencillo de mantener.

En general, la parte que vamos a tener que modificar más frecuentemente es la parte de *descodificación* y *parseo* del fichero JSON, por lo que quizás tiene más sentido que forme parte del *type* `Story`.

Movemos el código:

```go
  d := json.NewDecoder(f)
  var story cyoa.Story
  err = d.Decode(&story)
  if err != nil {
    panic(err)
  }
```

al fichero `story.go` y lo convertimos en una función `JsonStory`.

Esta función toma como argumento un `io.Reader` y devolverá una `Story` (y un `error`):

```go
// Punto de partida!!
func JsonStory(r io.Reader) (Story, error) {
  d := json.NewDecoder(r)
  var story cyoa.Story
  err = d.Decode(&story)
  if err != nil {
    panic(err)
  }
}
```

Cambiamos `f` por `r` como parámetro para `json.NewDecoder` y hacemos que en caso de error, en vez de lanzar un *panic*, devolvemos `nil` (no hay `Story`) y el error que se haya producido.

Si no se ha producido un error, devolvemos la `Story` (y `nil`).

Volviendo a `cmd/cyoaweb/main.go`, el bloque que hemos convertido en la función `JsonStory` se convierte en:

```go
  story, err := cyoa.JsonStory(f)
  if err != nil {
    panic(err)
  }
```

En realidad no hemos *reducido* el número de líneas de código, pero quizás tiene más sentido esta nueva organización, lo que puede simplificar el mantenimiento del código.

## Plantilla HTML y servidor web

Generamos una plantilla en HTML para mostrar nuestra historia.

Vamos a usar el paquete [`html/template`](https://pkg.go.dev/html/template) para poblar una plantilla de HTML desde Go.

> Usamos la plantilla que proporciona VSCode para HTML5, aunque no necesitamos la mayoría de campos *extra*.

```html
<!DOCTYPE html>
<html>
<head>
    <meta charset='utf-8'>
    <meta http-equiv='X-UA-Compatible' content='IE=edge'>
    <title>Choose Your Own Adventure</title>
    <meta name='viewport' content='width=device-width, initial-scale=1'>
</head>
<body>
    <h1>{{.Title}}</h1>

    {{range .Paragraphs}}
    <p>{{.}}</p>
    {{end}}

    <ul>
      {{range .Options}}
      <li><a href="/{{.Chapter}}">{{.Text}}</a><li>
      {{end}}
    </ul>
</body>
</html>
```

Podríamos guardar la plantilla como un fichero, y cargarlo desde la aplicación. Pero tratándose de una plantilla tan sencilla, lo más fácil es almacenar la plantilla en una variable, en el fichero `story.go`.

Creamos la variable `defaultHandlerTemplate` antes de la función `JsonStory` y pegamos el contenido de la plantilla usando (*backtics*) `` `...` `` (lo que nos permite tener cadenas de múltiples líneas).

Una vez tenemos la plantilla, necesitamos algo para *dibujar* (*render*) esta plantilla; tenemos que crear *algo* que pueda ser usado en las aplicaciones web para gestionar peticiones HTTP; para ello, revisando el paquete `http` tenemos `HandlerFunc` y `Handler` (que es un *interface*).

En este caso usaremos el *Handler* (más adelante veremos porqué).

Creamos un nuevo *handler* que toma una *Story* como argumento y devuelve un `http.Handler`:

> El autor no entra en detalles del porqué de esta "construcción", pero me temo que es uno de esos *patrones* que "tienen sentido" cuando sabes los tipos de objetos con los que estás trabajando.
>
> Al final, estamos consiguiendo lo que pretendíamos, que era crear una función a la que pasamos un `Story` y nos devuelve un `http.Handler` que puede usar un servidor web.
>
> Por lo que entiendo, creamos el objeto *handler* (un *struct*) basado en un *Story* y le asignamos el método *ServeHTTP*; la función `NewHandler` actúa como *constructor* del *handler* asociado a la *Story*.

```go
// Constructor de "handler"
func NewHandler(s Story) http.Handler {
  return handler{s}
}
// handler (con minúscula, no exportado)
type handler struct {
  s Story
}
// el método del handler para construir el `http.Handler` a partir de la `Story`
func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  ...
}
```

En `ServeHTTP`, *parseamos* la plantilla:

```go
func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  tpl := template.Must(template.New("").Parse(defaultHandlerTemplate))
  fmt.Printf("%+v", tpl)
}
```

Usamos [`template.Must()`](https://pkg.go.dev/html/template), que finaliza con *panic* si se produce un error al *parsear* la plantilla HTML. Esto es lo que necesitamos en este caso, porque si la plantilla no es válida, no se mostrará la historia y por tanto la aplicación no funcionará.

Muchas veces, en vez de validar la plantilla en el *handler*, se realiza la validación de la plantilla cuando se inicaliza el código (o se carga la plantilla).

Podríamos definir la variable `tpl *template.Template` al principio del código (como una variable global, pero no exportada) y usar la función [`init()`](https://go.dev/doc/effective_go#init) (que se ejecuta cuando se inicializa el código) para validar la plantilla.

Movemos la declaración de la función `tpl` al principio del fichero `story.go` y realizamos la validación de la plantilla en la función `init` del *package*:

```go
func init() {
  tpl = template.Must(template.New("").Parse(defaultHandlerTemplate))
}

var tpl *template.Template
```

De esta forma, el *handler* queda *vacío* (temporalmente):

```go
func NewHandler(s Story) http.Handler {
  return handler{s}
}

type handler struct {
  s Story
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
```

### Validando nuestro progreso

Para validar la plantilla, deberíamos pasar una *Story* al *handler*; la única *Story* que sabemos **seguro** que existirá es la `intro` (así lo hemos decidido, que la primera *story* en cualquier fichero sea `intro`).

En la función, elegimos pasar la *story* `intro` usando [template.Execute](https://pkg.go.dev/html/template#Template.Execute):

```go
err := tpl.Execute(w, h.s["intro"])
if err != nil {
  panic(err) // Porque estamos en desarrollo, pero no es una buena idea
}
```

Añadimos una nueva *flag* para especificar el puerto en el que escuchará el servidor web de *Choose your own adventure* y lanzamos el servidor.

En `cmd/cyoaweb/main.go`:

```go
port := flag.Int("port", "3000", "Port where the CYOA server listens")
```

Finalmente, eliminamos la línea que muestra el contenido del fichero por pantalla y los sustituimos por la llamada al *NewHandler*:

```go
h := cyoa.NewHandler(story)
fmt.Printf("Starting CYOA server on port %d\n", *port)
log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), h))
```

## Navegando por nuestra aventura

El siguiente paso es comprobar a dónde tenemos que saltar -a qué capítulo- al pulsar sobre las opciones que se muestran tras cada capítulo.

Empezamos averiguando en qué capítulo estamos; esto lo conseguimos a través del campo `r.URL.Path` en la *request* del *handler* (eliminamos espacios adicionales usando `strings.TrimSpace()` aunque no debería haber ninguno, sólo como medida de precaución).

A continuación, si el `path` está vacío o es `/`, asumimos que estamos empezando la *aventura*, por lo que redirigimos hacia el primer capítulo, que siempre se llama **intro**:

```go
path = "/intro"
```

Para el resto de casos, eliminamos la `/` inicial mediante:

```go
path = path[1:]
```

Si todo ha funcionado correctamente, ahora `path` contiene el nombre de un *chapter*:

```go
if chapter, ok := h.s[path], ok {
   ...
}
```

`h.s[path]` usa el valor contenido en `path` como *clave* del *map*; si lo encuentra, devuelve su valor, por lo que `chapter` contendrá el texto del capítulo.

Aunque no hace falta asignar el segundo valor a ninguna variable (indica si se ha encontrado algo o no), lo asignamos a `ok` porque de esta manera podemos validar si se ha encontrado el valor de `path` en el *map*.

Ejecutamos la plantilla -insertamos los datos procedentes del capítulo (si lo hemos encontrado)- y verificamos si se ha producido un error.

Como ahora tenemos el contenido del mapa en la variable *chapter*, eliminamos el `h.s["intro"]`.

```go
  if chapter, ok := h.s[path]; ok {
    err := tpl.Execute(w, chapter)
    if err != nil {
      log.Printf("%v", err)
      http.Error(w, "Something went wrong...", http.StatusInternalServerError)
    }
    return
  }
  http.Error(w, "Chapter not found", http.StatusNotFound)
```

Si se ha producido un error, lo mostramos en los logs (con `log.Printf()`) pero al usuario sólo le decimos que *algo ha salido mal* para evitar proporcionar demasiada información a un posible atacante. Como no sabemos qué es lo que ha pasado, devolvemos el código de error HTTP `http.StatusInternalServerError`.

Si no se ha encontrado el capítulo, mostramos el mensaje al usuario y el código HTTP `http.StatusNotFound`.

Podríamos convertir los mensajes de error en constantes, pero no es especialmente necesario para una aplicación como ésta.

## Plantilla personalizadas

Uno de los problemas que tiene la aplicación tal y como está ahora es que no es extensible, ni fácilmente modificable. La plantilla que usamos, por ejemplo, está *hardcodeada* en el código; así que una de las primeras modificaciones que podemos hacer es la de permitir personalizar la plantilla aplicada.

Para ello, definimos una variable `t *template.Template` en el *handler* (en el fichero `story.go`)

```go
type handler struct {
  s Story
  t *template.Template
}
```

También modficamos la función `NewHandler` para que el usuario pueda pasar una plantilla; si no lo hace, usaremos la plantilla por defecto:

```go
func NewHandler(s Story, t *template.Template) http.Handler {
  if t == nil {
    t = tpl
  }
  return handler{s, t}
}
```

También tenemos que actualizar la llamada desde `/cmd/cyoaweb/main.go`, aunque de momento, pasamos `nil`:

```go
h := cyoa.NewHandler(story, nil)
```

## Parametrizando las funciones

Si queremos ir ampliando las opciones de personalización de las funciones -para que se pueda construir el *path* de forma personalizada, o analizar de alguna otra forma, etc- tenemos que ir incluyendo cada vez más y más opciones... Y lo que es peor, si no se personalizan, tendremos que pasar `nil` o algo por el estilo.

Podemos usar [variadic functions](https://gobyexample.com/variadic-functions), que permiten pasar un número variable de parámetros a las funciones... Pero esta solución tampoco es ideal porque se pueden dar escenarios en los que tengamos dependencias entre algunos parámetros y no entre otros... Y entonces tendremos que controlar si se han proporcionado todos los parámetros relacionados, etc...

Lo ideal en este caso, sería poder extender las opciones disponibles en las funciones sin excesiva sobrecarga.

Al parecer, existe un patrón llamado *functional options* (David Cheney) que permite extender las funciones (Jon se refiere a [Functional options for friendly APIs](https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis)).

La idea es crear un nuevo `type`, en nuestro caso `HandlerOption` que es una función a la que pasamos un puntero (para que lo pueda modificar) de tipo `handler`:

```go
type HandlerOption func(h *handler)
```

A continuación, para extender la función, usaremos:

```go
func NewHandler(s Story, opts ...HandlerOption) http.Handler {
}
```

Ahora, lo que podemos hacer es crear una función que lo que hace es modificar el *handler* (aunque el *handler* nunca está expuesto para que el usuario lo pueda modificar directamente):

```go
func WithTemplate(t *template.Template) HandlerOption {
  return func(h *handler) {
    h.t = t
  }
}
```

Ahora, el constructor del *handler*:

```go
func NewHandler(s Story, opts ...HandlerOption) http.Handler {
  h := handler{s, tpl}
  for _, opt := range opts {
    opt(&h)
  }
  return h
}
```

Empezamos usando `s` (no es opcional, necesitamos una *story*) y usando la plantilla a partir de la variable global `tpl`.

Después, recorremos las opciones que se hayan pasado y pasamos la referencia al *handler*.

Con estas modificaciones, la llamada al *constructor* sería, si queremos usar todos las opciones por defecto:

```go
h := cyoa.NewHandler(story)
```

Y si queremos pasar una plantilla personalizada (en este caso, `tpl`)

```go
// Este template no usa ninguna variable, sólo la cadena Hello! en todos los casos
tpl := template.Must(template.New("").Parse("Hello!"))
h := cyoa.NewHandler(story, cyoa.WithTemplate(tpl))
```

Para que se use la plantilla "como parámetro" y no la variable global que estábamos usando hasta ahora, hay que modificar la función `ServeHTTP` en `story.go` y cambiar `tpl.Execute` por `h.t.Execute`:

```go
  if chapter, ok := h.s[path]; ok {
    err := h.t.Execute(w, chapter)
```
