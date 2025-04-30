package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
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
			imprimeLogs()
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

	sites := leSitesDoArquivo()

	for range monitoramentos {
		for _, site := range sites {
			testaSite(site)
		}
		time.Sleep(time.Second * 5)
	}
}

func testaSite(url string) {
	fmt.Println("Testando site: ", url)
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
	}

	if res.StatusCode == http.StatusOK {
		fmt.Println("Site:", url, "foi carregado com sucesso!")
		registraLog(url, true)
	} else {
		fmt.Println("Site:", url, "esta com problemas. Status Code:", res.StatusCode)
		registraLog(url, false)
	}
}

func leSitesDoArquivo() []string {
	arquivo, err := os.Open("sites.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer arquivo.Close()

	sites := make([]string, 0, 10)
	leitor := bufio.NewReader(arquivo)
	for {
		linha, err := leitor.ReadString('\n')
		if err == io.EOF {
			break
		}
		sites = append(sites, strings.TrimSpace(linha))
	}

	return sites
}

func registraLog(site string, status bool) {
	arquivo, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
	}
	defer arquivo.Close()

	arquivo.WriteString(time.Now().Format("02/01/2006 15:04:05") + " - " + site + "- online: " + strconv.FormatBool(status) + "\n")
}

func imprimeLogs() {
	arquivo, err := ioutil.ReadFile("log.txt")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(arquivo))
}
