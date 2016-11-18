package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/remogatto/ricevimenti"
)

var (
	verbose *bool
	Usage   = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()

	}
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func logging(format string, message ...interface{}) {
	if *verbose {
		log.Printf(format, message...)
	}
}

func main() {
	data := make([]string, 0)

	outputDir := flag.String("o", "output", "Cartella di output")
	banditi := flag.String("b", "", "Elenco di Cognomi/Nomi banditi separati da virgola")
	verbose = flag.Bool("verbose", false, "Attiva i log")
	flag.Parse()

	if len(flag.Args()) > 0 {
		config := new(ricevimenti.Config)
		config.Banditi = strings.Split(*banditi, ",")

		for _, filename := range flag.Args() {
			logging("Leggo il file %s\n", filename)
			csvData, err := ioutil.ReadFile(filename)
			check(err)
			data = append(data, string(csvData))

		}
		logging("%s\n", "Elaboro i dati CSV...")
		r := ricevimenti.NuoviRicevimenti(config, data...)
		err := os.Mkdir(*outputDir, 0777)
		check(err)
		logging("Creo cartella %s...", *outputDir)
		for _, docente := range r.ListaDocenti() {
			base := strings.Join(strings.Split(docente, " "), "_")
			htmlFilename := fmt.Sprintf("%s.html", base)
			completePath := path.Join(*outputDir, htmlFilename)
			f, err := os.Create(completePath)
			check(err)
			defer f.Close()
			f.WriteString(r.GeneraHTML(docente))
			docFilename := filepath.Join(*outputDir, fmt.Sprintf("%s.doc", base))
			logging("Converto da %s a %s...", completePath, docFilename)
			_, err = exec.Command("pandoc", "-o", docFilename, completePath).CombinedOutput()
			check(err)
		}
	} else {
		Usage()
	}
}
