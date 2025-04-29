package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

const monitoramentos = 3

func main() {
	exibeIntroducao()

	for {
		exibeMenu()

		switch leComando() {
		case 1:
			iniciarMonitoramento()
		case 2:
			fmt.Println("Exibindo logs...")
		case 0:
			fmt.Println("Saindo do programa")
			os.Exit(0)
		default:
			fmt.Println("Não conheço este comando")
			os.Exit(-1)
		}
	}
}

func exibeIntroducao() {
	nome := "Pedro"
	versao := 1.1
	fmt.Println("Olá, sr.", nome)
	fmt.Println("Este programa está na versão", versao)
}

func exibeMenu() {
	fmt.Println("1- Iniciar Monitoramento")
	fmt.Println("2- Exibir logs")
	fmt.Println("0- Sair do Programa")
}

func leComando() (comando int) {
	fmt.Scan(&comando)
	fmt.Println("")
	return
}

func iniciarMonitoramento() {
	fmt.Println("Monitorando...")
	sites := []string{
		"https://httpbin.org/status/200",
		"https://www.alura.com.br",
		"https://www.caelum.com.br",
	}

	for range monitoramentos {
		for _, site := range sites {
			testaSite(site)
		}
		time.Sleep(time.Second * 5)
	}
}

func testaSite(url string) {
	fmt.Println("Testando site: ", url)
	res, _ := http.Get(url)

	if res.StatusCode == http.StatusOK {
		fmt.Println("Site:", url, "foi carregado com sucesso!")
	} else {
		fmt.Println("Site:", url, "esta com problemas. Status Code:", res.StatusCode)
	}
}
