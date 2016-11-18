package ricevimenti

import (
	"encoding/csv"
	"fmt"
	"log"
	"strings"
)

const (
	DATA_INVIO = iota
	IP
	COGNOME
	NOME
	CLASSE
	SEZIONE
	INDIRIZZO
	PRIMO_DOCENTE
)

const Intestazione = `
<tr>
<th>Posizione</th>
<th>Genitore/i dell'alunno</th>
<th>Classe</th>
<th>Sezione</th>
<th>Indirizzo</th>
</tr>`

type Config struct {
	Banditi []string
}

type Alunno struct {
	Nome      string
	Cognome   string
	Classe    string
	Sezione   string
	Indirizzo string
}

type Ricevimenti struct {
	Config *Config

	docenti    map[string][]*Alunno
	id_docenti map[int]string
}

func alunniDuplicati(alunno1 *Alunno, alunno2 *Alunno) bool {
	if (alunno1.Nome == alunno2.Nome) && (alunno1.Cognome == alunno2.Cognome) && (alunno1.Classe == alunno2.Classe) && (alunno1.Sezione == alunno2.Sezione) && (alunno1.Indirizzo == alunno2.Indirizzo) {
		return true
	}
	return false
}

func (r *Ricevimenti) alunnoBandito(alunno *Alunno) bool {
	for _, a := range r.Config.Banditi {
		if a == fmt.Sprintf("%s %s", alunno.Cognome, alunno.Nome) {
			return true
		}
	}
	return false
}

func (r *Ricevimenti) listaContiene(docente string, alunno *Alunno) bool {
	for _, a := range r.docenti[docente] {
		if r.alunnoBandito(alunno) {
			return true
		}
		if alunniDuplicati(a, alunno) {
			return true
		}
	}
	return false
}

func NuoviRicevimenti(config *Config, data ...string) *Ricevimenti {

	r := new(Ricevimenti)

	r.Config = config

	r.id_docenti = make(map[int]string)
	r.docenti = make(map[string][]*Alunno)

	for _, csvData := range data {
		reader := csv.NewReader(strings.NewReader(csvData))
		records, err := reader.ReadAll()
		if err != nil {
			log.Fatal(err)
		}

		n_campi := len(records[0])
		n_docenti := n_campi - PRIMO_DOCENTE
		n_records := len(records)

		for i := 0; i < n_docenti; i++ {
			if len(r.docenti[records[0][i+PRIMO_DOCENTE]]) == 0 {
				r.docenti[records[0][i+PRIMO_DOCENTE]] = make([]*Alunno, 0)
			}
			r.id_docenti[i] = records[0][i+PRIMO_DOCENTE]
		}

		for i := 1; i < n_records; i++ {
			alunno := new(Alunno)
			alunno.Nome = strings.ToUpper(records[i][NOME])
			alunno.Cognome = strings.ToUpper(records[i][COGNOME])
			alunno.Classe = records[i][CLASSE]
			alunno.Sezione = records[i][SEZIONE]
			alunno.Indirizzo = strings.ToUpper(records[i][INDIRIZZO])

			for t := 0; t < n_docenti; t++ {
				docente := r.id_docenti[t]
				if records[i][t+PRIMO_DOCENTE] == "Sì" && !r.listaContiene(docente, alunno) {
					r.docenti[docente] = append(r.docenti[docente], alunno)
				}
			}
		}
	}
	return r
}

func (r *Ricevimenti) ListaDocenti() []string {
	docenti := make([]string, 0)
	for docente, _ := range r.docenti {
		docenti = append(docenti, docente)
	}
	return docenti
}

func (r *Ricevimenti) PrenotazioniDocente(cognomeNome string) []*Alunno {
	return r.docenti[cognomeNome]
}

func (r *Ricevimenti) GeneraHTML(cognomeNome string) string {
	alunni := r.PrenotazioniDocente(cognomeNome)
	lista := ""
	for pos, alunno := range alunni {
		lista += "<tr>\n"
		lista += fmt.Sprintf("<td>%d</td>\n", pos+1)
		lista += fmt.Sprintf("<td>%s %s</td>\n", alunno.Cognome, alunno.Nome)
		lista += fmt.Sprintf("<td>%s</td>\n", alunno.Classe)
		lista += fmt.Sprintf("<td>%s</td>\n", alunno.Sezione)
		lista += fmt.Sprintf("<td>%s</td>\n", alunno.Indirizzo)
		lista += "</tr>\n"
	}
	risultato := fmt.Sprintf("<h1>%s</h1>\n<table>%s\n%s</table>", cognomeNome, Intestazione, lista)
	return risultato
}
