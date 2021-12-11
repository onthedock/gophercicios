package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	QUIZ_FILE_NOT_FOUND         int = 1
	ERROR_READING_QUIZ_CSV_FILE int = 2
	ERROR_NO_ENTIENDO_RESPUESTA int = 3
	PRUEBA_FINALIZADA           int = 0
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

func getUserAnswer() {

	r := bufio.NewReader(os.Stdin)
	respuestaUsuario, err := r.ReadString('\n')
	if err != nil {
		exit("Error al leer la respuesta.\n[ERROR] %s\n", err, ERROR_NO_ENTIENDO_RESPUESTA)
	}
	canalRespuestaUsuario <- strings.Trim(respuestaUsuario, "\n ")
	return
}

func exit(msg string, err error, exitcode int) {
	fmt.Printf(msg, err.Error())
	os.Exit(exitcode)
}

func printResult(puntuacion int, total int) {
	fmt.Printf("Ha respondido %d respuestas correctas de %d preguntas.\n", puntuacion, total)
	os.Exit(PRUEBA_FINALIZADA)
}

var canalRespuestaUsuario = make(chan string)

func main() {
	csvFile := flag.String("csv", "problems.csv", "Fichero de problemas en formato 'enunciado,respuesta' (CSV)")
	limite := flag.Int("limite", 30, "Tiempo limite para completar la prueba en segundos")
	flag.Parse()

	var fHandle *os.File = openFile(*csvFile)
	defer fHandle.Close()

	r := csv.NewReader(fHandle)

	var lineas [][]string
	lineas = readFile(r)
	problemas := parselines(lineas)

	var puntuacion int = 0

	temporizador := time.NewTimer(time.Duration(*limite) * time.Second)

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
