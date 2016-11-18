package ricevimenti_test

import (
	"testing"

	"github.com/remogatto/prettytest"
	"github.com/remogatto/ricevimenti"
)

// Start of setup

type testSuite struct {
	prettytest.Suite
}

var (
	config *ricevimenti.Config

	CSV_1 string = `
Data invio,Indirizzo IP,Cognome,Nome,Classe,Sezione,Indirizzo,AGOSTONI Elena,ANTOCI Francesca,BALSAMO Consiglia,COSENZA Ernst
2016-11-14 08:16:34,79.45.232.40,Grassi,Giulia,2,C,Economico sociale,Sì,No,No,No
2016-11-14 08:27:17,87.15.239.29,Chiapparino,Alice,2,A,Classico,No,No,No,No
2016-11-14 08:34:36,151.43.120.68,Zappi,Piero,5,A,Economico sociale,No,No,No,No
2016-11-14 08:36:37,89.96.138.74,Duiz,Simone,2,A,Classico,No,No,No,No
2016-11-14 08:54:47,79.54.50.194,Pecorari,Simone,3,B,Economico sociale,No,No,No,No
2016-11-14 10:37:10,2.224.178.39,Cavalieri,Noemi,1,E,Scienze umane,Sì,No,No,No
2016-11-14 14:24:31,79.54.48.214,Sivini,Giada,1,D,Scienze umane,No,No,Sì,No
`
	CSV_2 string = `
Data invio,Indirizzo IP,Cognome,Nome,Classe,Sezione,Indirizzo,AGOSTONI Elena,ANTOCI Francesca,BIANCHI Brigitta,BALSAMO Consiglia,COSENZA Ernst
2016-11-14 08:16:34,79.45.232.40,Carciotti,Sara,4,M,Musicale,Sì,No,sì,No,No
2016-11-14 08:27:17,87.15.239.29,Bugliano,Andrea,5,B,Scienze umane,No,No,No,Sì,Sì
`
	CSV_con_cognome_nomi_banditi = `
Data invio,Indirizzo IP,Cognome,Nome,Classe,Sezione,Indirizzo,AGOSTONI Elena,ANTOCI Francesca,BIANCHI Brigitta,BALSAMO Consiglia,COSENZA Ernst
2016-11-14 08:16:34,79.45.232.40,Carciotti,Sara,4,M,Musicale,Sì,No,sì,No,No
2016-11-14 08:27:17,87.15.239.29,Bugliano,Andrea,5,B,Scienze umane,No,No,No,Sì,Sì
2016-11-14 08:27:17,87.15.239.29,La Maiala,Sapida,5,B,Scienze umane,Sì,No,No,Sì,Sì
`
	CSV_con_duplicati string = `
Data invio,Indirizzo IP,Cognome,Nome,Classe,Sezione,Indirizzo,AGOSTONI Elena,ANTOCI Francesca,BIANCHI Brigitta,BALSAMO Consiglia,COSENZA Ernst
2016-11-14 08:16:34,79.45.232.40,Carciotti,Sara,4,M,Musicale,Sì,No,sì,No,No
2016-11-14 08:16:34,79.45.232.40,Carciotti,Sara,4,M,Musicale,Sì,No,sì,No,No
2016-11-14 08:27:17,87.15.239.29,Bugliano,Andrea,5,B,Scienze umane,No,No,No,No,Sì
2016-11-14 14:24:31,79.54.48.214,Sivini,Giada,1,D,Scienze umane,No,No,Sì,No,No
2016-11-14 08:54:47,79.54.50.194,Pecorari,Simone,3,B,Economico sociale,No,No,No,No,No
2016-11-14 14:24:31,79.54.48.214,siVini,giada,1,D,Scienze umane,No,No,Sì,No,No
`
	RISULTATO_HTML = `<h1>AGOSTONI Elena</h1>
<table>
<tr>
<th>Posizione</th>
<th>Genitore/i dell'alunno</th>
<th>Classe</th>
<th>Sezione</th>
<th>Indirizzo</th>
</tr>
<tr>
<td>1</td>
<td>GRASSI GIULIA</td>
<td>2</td>
<td>C</td>
<td>ECONOMICO SOCIALE</td>
</tr>
<tr>
<td>2</td>
<td>CAVALIERI NOEMI</td>
<td>1</td>
<td>E</td>
<td>SCIENZE UMANE</td>
</tr>
</table>`
)

func TestRunner(t *testing.T) {
	prettytest.Run(
		t,
		new(testSuite),
	)
}

// End of setup

// Your tests start here

func (t *testSuite) Before() {
	config = new(ricevimenti.Config)
	config.Banditi = []string{"LA MAIALA SAPIDA", "PIERINO PAPERINO"}
}

func (t *testSuite) TestAPIBase() {
	r := ricevimenti.NuoviRicevimenti(config, CSV_1)
	alunni := r.PrenotazioniDocente("AGOSTONI Elena")
	t.Not(t.Nil(alunni))
	t.True(alunni[0].Nome == "GIULIA")
	t.True(alunni[0].Cognome == "GRASSI")
}

func (t *testSuite) TestListaDocenti() {
	r := ricevimenti.NuoviRicevimenti(config, CSV_1)
	docenti := r.ListaDocenti()
	t.Equal(len(docenti), 4)
}

func (t *testSuite) TestFormattazioni() {
	r := ricevimenti.NuoviRicevimenti(config, CSV_1)
	alunni := r.PrenotazioniDocente("BALSAMO Consiglia")
	t.Equal("SIVINI", alunni[0].Cognome)
}

func (t *testSuite) TestElenchi() {
	r := ricevimenti.NuoviRicevimenti(config, CSV_1)
	alunni := r.PrenotazioniDocente("AGOSTONI Elena")
	t.Not(t.Nil(alunni))
	t.Equal("GRASSI", alunni[0].Cognome)
	t.Equal("CAVALIERI", alunni[1].Cognome)

	alunni = r.PrenotazioniDocente("BALSAMO Consiglia")
	t.Not(t.Nil(alunni))
	t.Equal("SIVINI", alunni[0].Cognome)
}

func (t *testSuite) TestCSVMultipli() {
	r := ricevimenti.NuoviRicevimenti(config, CSV_1, CSV_2)
	alunni := r.PrenotazioniDocente("AGOSTONI Elena")
	t.Not(t.Nil(alunni))
	t.Equal(3, len(alunni))
	t.Equal("GRASSI", alunni[0].Cognome)
	t.Equal("CAVALIERI", alunni[1].Cognome)
	t.Equal("CARCIOTTI", alunni[2].Cognome)

	alunni = r.PrenotazioniDocente("BALSAMO Consiglia")
	t.Equal(2, len(alunni))
	t.Equal("SIVINI", alunni[0].Cognome)
	t.Equal("BUGLIANO", alunni[1].Cognome)

}

func (t *testSuite) TestCSVConDuplicati() {
	r := ricevimenti.NuoviRicevimenti(config, CSV_con_duplicati)
	alunni := r.PrenotazioniDocente("AGOSTONI Elena")
	t.Not(t.Nil(alunni))
	t.Equal(1, len(alunni))
}

func (t *testSuite) TestCognomeNomeBanditi() {
	r := ricevimenti.NuoviRicevimenti(config, CSV_con_cognome_nomi_banditi)
	alunni := r.PrenotazioniDocente("AGOSTONI Elena")
	t.Not(t.Nil(alunni))
	t.Equal(1, len(alunni))
	t.Equal("CARCIOTTI", alunni[0].Cognome)
}

func (t *testSuite) TestRisultatoHTML() {
	r := ricevimenti.NuoviRicevimenti(config, CSV_1)
	markdown := r.GeneraHTML("AGOSTONI Elena")
	t.Equal(RISULTATO_HTML, markdown)
}
