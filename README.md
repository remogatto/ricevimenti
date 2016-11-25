# Generatore di liste di prenotazione

A partire dal corrente anno scolastico, presso il Liceo "Carducci-Dante" di Trieste, è possibile prenotare i colloqui con i docenti nelle giornate di ricevimento pomeridiano, mediante modulo web realizzato con [LimeSurvey](https://www.limesurvey.org/).

Il modulo consente di scegliere il nome dei docenti attraverso una semplice interfaccia utente. Una volta raccolti i dati, un algoritmo produce le liste di prenotazione.

L’algoritmo è costruito in modo  da ottimizzare il flusso dei colloqui, evitando le situazioni di contemporaneità degli stessi. La posizione di ciascun genitore viene scelta in base a due criteri, nel seguente ordine:

1. Data e ora di invio della prenotazione
2. Eventuale contemporaneità dei colloqui

Il software è sviluppato in [Go](https://golang.org/), un linguaggio di programmazione opensource, compilato e particolarmente efficiente, sviluppato da Google.
