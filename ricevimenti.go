package ricevimenti

import (
	"encoding/csv"
	"fmt"
	"log"
	"sort"
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

const IntestazioneAlunno = `
<tr>
<th>Posizione</th>
<th>Docente</th>
</tr>`

type Config struct {
	Banditi             []string
	MassimoPrenotazioni int
}

type Alunno struct {
	Nome      string
	Cognome   string
	Classe    string
	Sezione   string
	Indirizzo string
}

type Posizione struct {
	Pos     int
	Docente string
}

type PerPosizione []*Posizione

func (p PerPosizione) Len() int              { return len(p) }
func (p PerPosizione) Swap(i, j int)         { p[i], p[j] = p[j], p[i] }
func (p PerPosizione) Less(i, j int) bool    { return p[i].Pos < p[j].Pos }
func (p PerPosizione) ToArray() []*Posizione { return ([]*Posizione)(p) }

type Ricevimenti struct {
	Config *Config

	docenti         map[string][]*Alunno
	id_docenti      map[int]string
	PosizioniAlunni map[string]PerPosizione
}

func alunniDuplicati(alunno1 *Alunno, alunno2 *Alunno) bool {
	if (alunno1.Nome == alunno2.Nome) && (alunno1.Cognome == alunno2.Cognome) && (alunno1.Classe == alunno2.Classe) && (alunno1.Sezione == alunno2.Sezione) && (alunno1.Indirizzo == alunno2.Indirizzo) {
		return true
	}
	return false
}

func (a *Alunno) cognomeNome() string {
	return fmt.Sprintf("%s %s", a.Cognome, a.Nome)
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
	if r.alunnoBandito(alunno) {
		return true
	}
	for _, a := range r.docenti[docente] {
		if a != nil {
			if alunniDuplicati(a, alunno) {
				return true
			}
		}
	}
	return false
}

func (r *Ricevimenti) posizioneDisponibile(pos int, alunno *Alunno) bool {
	posizioni := r.PosizioniAlunni[alunno.cognomeNome()]
	for _, p := range posizioni.ToArray() {
		if p.Pos == pos {
			return false
		}
	}
	return true
}

func (r *Ricevimenti) PosizioniAlunno(cognomeNome string) []*Posizione {
	return r.PosizioniAlunni[cognomeNome]
}

func (r *Ricevimenti) inserisciAlunno(docente string, alunno *Alunno) {

	for pos, a := range r.docenti[docente] {
		if (a == nil) && r.posizioneDisponibile(pos, alunno) {
			r.docenti[docente][pos] = alunno
			posizione := new(Posizione)
			posizione.Pos = pos
			posizione.Docente = docente
			r.PosizioniAlunni[alunno.cognomeNome()] = append(r.PosizioniAlunni[alunno.cognomeNome()], posizione)
			return
		}
	}
}

func NuoviRicevimenti(config *Config, data ...string) *Ricevimenti {

	r := new(Ricevimenti)

	r.Config = config

	r.id_docenti = make(map[int]string)
	r.docenti = make(map[string][]*Alunno)
	r.PosizioniAlunni = make(map[string]PerPosizione)

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
				r.docenti[records[0][i+PRIMO_DOCENTE]] = make([]*Alunno, config.MassimoPrenotazioni)
			}
			r.id_docenti[i] = strings.TrimLeft(records[0][i+PRIMO_DOCENTE], " ")
		}

		for i := 1; i < n_records; i++ {
			alunno := new(Alunno)
			alunno.Nome = strings.Trim(strings.ToUpper(records[i][NOME]), " ")
			alunno.Cognome = strings.Trim(strings.ToUpper(records[i][COGNOME]), " ")
			alunno.Classe = records[i][CLASSE]
			alunno.Sezione = records[i][SEZIONE]
			alunno.Indirizzo = strings.ToUpper(records[i][INDIRIZZO])
			posizione := new(Posizione)
			if !r.alunnoBandito(alunno) {
				posizione.Pos = -1
				r.PosizioniAlunni[alunno.cognomeNome()] = append(r.PosizioniAlunni[alunno.cognomeNome()], posizione)
			}
			for t := 0; t < n_docenti; t++ {
				docente := r.id_docenti[t]
				if (records[i][t+PRIMO_DOCENTE] == "Sì") && (!r.listaContiene(docente, alunno)) {
					r.inserisciAlunno(docente, alunno)
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

// Grazie al sistema automatico di prenotazione, si è potuto evitare
// la sovrapposizione dei colloqui dello stesso genitore con più
// docenti. Per cui gli spazi creatisi devono essere riempiti con le
// prenotazioni giunte via libretto.
func (r *Ricevimenti) GeneraHTML(cognomeNome string) string {

	nota := "Grazie al sistema automatico di prenotazione, è stato possibile evitare la sovrapposizione dei colloqui dello stesso genitore con più docenti. Per cui, gli spazi eventualmente creatisi, saranno riempiti con le prenotazioni giunte via libretto."

	alunni := r.PrenotazioniDocente(cognomeNome)
	lista := ""
	for pos, alunno := range alunni {
		if alunno != nil {
			lista += "<tr>\n"
			lista += fmt.Sprintf("<td>%d</td>\n", pos+1)
			lista += fmt.Sprintf("<td>%s %s</td>\n", alunno.Cognome, alunno.Nome)
			lista += fmt.Sprintf("<td>%s</td>\n", alunno.Classe)
			lista += fmt.Sprintf("<td>%s</td>\n", alunno.Sezione)
			lista += fmt.Sprintf("<td>%s</td>\n", alunno.Indirizzo)
			lista += "</tr>\n"
		} else if pos < r.Config.MassimoPrenotazioni {
			lista += "<tr>\n"
			lista += fmt.Sprintf("<td>%d</td>\n", pos+1)
			lista += fmt.Sprintf("<td>%s</td>\n", "DISPONIBILE")
			lista += fmt.Sprintf("<td>%s</td>\n", "")
			lista += fmt.Sprintf("<td>%s</td>\n", "")
			lista += fmt.Sprintf("<td>%s</td>\n", "")
			lista += "</tr>\n"

		}
	}
	risultato := fmt.Sprintf("<h1>%s</h1>\n<h2>Nota</h2>\n<p>%s</p><table>%s\n%s</table>", cognomeNome, nota, Intestazione, lista)
	return risultato
}

func (r *Ricevimenti) GeneraHTMLAlunno(cognomeNome string) string {
	posizioni := r.PosizioniAlunno(cognomeNome)
	sort.Sort(PerPosizione(posizioni))
	lista := ""
	for _, pos := range posizioni {
		if pos.Pos == -1 {
			continue
		}
		lista += "<tr>\n"
		lista += fmt.Sprintf("<td>%d</td>\n", pos.Pos+1)
		lista += fmt.Sprintf("<td>%s</td>\n", pos.Docente)
		lista += "</tr>\n"
	}

	risultato := fmt.Sprintf("<h1>%s</h1>\n<table>%s\n%s</table>", cognomeNome, IntestazioneAlunno, lista)
	return risultato
}
