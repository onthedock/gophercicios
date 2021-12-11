# Detalles del ejercicio 1

El ejercicio está dividido en dos partes para simplifircar el proceso de explicación así como para hacerlo más sencillo de solucionar. La segunda parte es más difícil que la primera, por lo que si te quedas encallado, no dudes en avanzar a otro problema y volver más adelante a la parte 2.

## Primera parte

Crea un programa que lea las preguntas de una prueba desde un fichero CSV y que muestre las preguntas al usuario, guardando un registro de cuántas preguntas acierta y cuántas falla. Independientemente de si la respuesta es correcta o no, se debe preguntar la siguiente pregunta a continuación.

El fichero CSV por defecto es `problems.csv`, pero el usuario puede personalizar el nombre del fichero a través de una opción (*flag*).

El fichero CSV tendrá el formato siguiente, con la primera columna con la pregunta y la segunda columna la respuesta a la pregunta:

```csv
5+5,10
7+3,10
1+1,2
8+3,11
1+2,3
8+6,14
3+1,4
1+4,5
5+1,6
2+3,5
3+3,6
2+4,6
5+2,7
```

Se puede asumir que las pruebas serán relativamente cortas (<100 preguntas) y que las respuestas serán una única palabra o número.

Al final de la prueba el programa debe mostrar el número total de preguntas acertadas y de cuántas preguntas en total consistía la prueba.
