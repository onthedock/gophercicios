package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
)

const (
	QUIZ_FILE_NOT_FOUND         int = 1
	ERROR_READING_QUIZ_CSV_FILE int = 2
	ERROR_NO_ENTIENDO_RESPUESTA int = 3
)

type pregunta struct {
	enunciado string
	respuesta string
}

func parselines(lines [][]string) []pregunta {
	var res []pregunta
	res = make([]pregunta, len(lines))
	for i, problema := range lines {
		res[i] = pregunta{
			enunciado: problema[0],
			respuesta: strings.Trim(problema[1], " "),
		}
	}
	return res
}

func openFile(fileName string) *os.File {
	file, err := os.Open(fileName)
	if err != nil {
		exit("Error al abrir el fichero.\n[ERROR] %s\n", err, QUIZ_FILE_NOT_FOUND)
	}
	return file
}

func readFile(f *csv.Reader) [][]string {
	var lineas [][]string
	lineas, err := f.ReadAll()
	if err != nil {
		exit("Error al leer el archivo.\n[ERROR] %s\n", err, ERROR_READING_QUIZ_CSV_FILE)
	}
	return lineas
}

func getUserAnswer() string {
	r := bufio.NewReader(os.Stdin)
	respuestaUsuario, err := r.ReadString('\n')
	if err != nil {
		exit("Error al leer la respuesta.\n[ERROR] %s\n", err, ERROR_NO_ENTIENDO_RESPUESTA)
	}
	return strings.Trim(respuestaUsuario, "\n ")
}

func exit(msg string, err error, exitcode int) {
	fmt.Printf(msg, err.Error())
	os.Exit(exitcode)
}

func main() {
	csvFile := flag.String("csv", "problems.csv", "Fichero de problemas en formato 'enunciado,respuesta' (CSV)")
	flag.Parse()

	r := csv.NewReader(openFile(*csvFile))

	var lineas [][]string
	lineas = readFile(r)
	problemas := parselines(lineas)

	var puntuacion int = 0
	for i := range problemas {
		fmt.Printf("Pregunta %d:\t %s =\n", i+1, problemas[i].enunciado)

		if problemas[i].respuesta == getUserAnswer() {
			puntuacion++
		}
	}
	fmt.Printf("Ha respondido %d respuestas correctas de %d preguntas.\n", puntuacion, len(problemas))
}
