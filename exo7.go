package main

/**
  SUJET : CRIBLE D'ÉRATOSTHÈNE
  AUTEUR : Alexandre Cazala
**/

import (
    "fmt"
    "time"
)

func main() {
  result_channel := primesNumber()

  // Affichage des nombres premiers inférieurs à 100
  for n := <-result_channel; n < 1000; n = <-result_channel {
    fmt.Println(n)
  }
  close(result_channel)
}

func primesNumber() chan int {
  result_channel := make(chan int)
  go func(){
      production_channel := make(chan int, 10)
      production_channel <- 2
      go Worker(production_channel, result_channel)
  }()
  return result_channel
}

func Worker(production_channel chan int, result_channel chan int) {
    // Premiere etape : On récupère notre valeur (qui est toujours différente de 0)
    to_filter_channel := make(chan int, 10)
    init_value := 0
    for ; init_value == 0 ; {
      select {
      case init_value = <-production_channel: // On vient de récupérer notre valeur
        result_channel <- init_value // notre première valeur est un nombre premier
        go Worker(to_filter_channel, result_channel) // On lance un nouveau worker
      default: // Au cas où, il n'y ait rien, on dort pour réduire la consommation CPU
        time.Sleep(100)
      }
    }
    cpt := init_value
    for {
      cpt = init_value + cpt
      select {
      case to_filter_channel <- cpt : // on push notre compteur au worker suivant
      default: // Si  on peut pas push au worker suivant alors on dort un peu pour économiser le CPU
        time.Sleep(100)
      }
    }
}
