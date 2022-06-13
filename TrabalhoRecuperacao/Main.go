package main

import (
	"fmt"
	"strconv"
)

type ComponentesProcesso struct {
	id            int
	estado        int
	cont          int
	contReqLocal  int
	filaPendentes []Requisicao
}
type Requisicao struct {
	ACK     bool
	id      int
	contReq int
}

const N = 10 // numero de processos
type inputChan [N]chan Requisicao

var mutex = make(chan struct{}, 1)
var inCh inputChan

func main() {
	var naoTermina = make(chan struct{})
	mutex <- struct{}{}
	for i := 0; i < N; i++ {
		inCh[i] = make(chan Requisicao, 1)
		var pendentes = make([]Requisicao, N)
		p := ComponentesProcesso{id: i, estado: 0, cont: 0, contReqLocal: 0, filaPendentes: pendentes}
		go processo(p)
		go processoRespondedor(p)
	}
	naoTermina <- struct{}{}
}
func recurso(p ComponentesProcesso) {
	fmt.Println("____________________________________________________________")

	fmt.Println("Processo: " + strconv.Itoa(p.id) + " esta no recurso")

	fmt.Println("____________________________________________________________")
}

func processo(p ComponentesProcesso) {
	fmt.Println("processo executando: ", p.id)
	for {
		<-mutex //entra secao critica
		fmt.Println("processo esta na secao critica: ", p.id)
		p.estado = 1
		p.cont++
		p.contReqLocal = p.cont
		for i := 0; i < N; i++ {
			if i != p.id {
				fmt.Println("processo:  " + strconv.Itoa(p.id) + " pedindo acesso para processo: " + strconv.Itoa(i))
				r := Requisicao{ACK: true, id: p.id, contReq: p.contReqLocal}
				inCh[i] <- r

			}

		}
		retorno := true

		contador := 0
		for contador < (N - 1) {
			resposta := <-inCh[p.id]
			fmt.Println("processo:  " + strconv.Itoa(p.id) + " recebendo resposta do processo: " + strconv.Itoa(resposta.id) + " resposta: " + strconv.FormatBool(resposta.ACK))
			if resposta.ACK == false {
				retorno = false
			}
			contador++
		}

		if retorno == true {
			fmt.Println("processo: " + strconv.Itoa(p.id) + " pode acessar recurso")
			p.estado = 2
			recurso(p)

		}
		mutex <- struct{}{} //sai secao critica
		fmt.Println("processo saiu da secao critica: ", p.id)
		p.estado = 0

		contadorDependentes := 0
		for contadorDependentes < len(p.filaPendentes) {
			resposta := Requisicao{ACK: true, id: p.id, contReq: p.contReqLocal}
			fmt.Println("processo: " + strconv.Itoa(p.id) + " respondendo pendente: " + strconv.Itoa(p.filaPendentes[contadorDependentes].id))
			inCh[p.filaPendentes[contadorDependentes].id] <- resposta
			contadorDependentes++
		}

	}
}
func processoRespondedor(p ComponentesProcesso) {
	fmt.Println("processo respondedor executando: ", p.id)
	for {
		pedido := <-inCh[p.id]
		fmt.Println("processo respondedor: " + strconv.Itoa(p.id) + " recebeu requisicao  do processo  " + strconv.Itoa(pedido.id))
		if pedido.contReq > p.cont {
			p.cont = pedido.contReq + 1
		}

		switch p.estado {
		case 0:
			fmt.Println("processo respondedor: " + strconv.Itoa(p.id) + " respondeu requisicao com ok para o processo  " + strconv.Itoa(pedido.id))
			resposta := Requisicao{ACK: true, id: p.id, contReq: p.contReqLocal}
			inCh[pedido.id] <- resposta

		case 1:
			if p.contReqLocal < pedido.contReq {
				fmt.Println("processo" + strconv.Itoa(p.id) + " respondeu requisicao com ok pois pedido é anterior para o processo " + strconv.Itoa(pedido.id))
				resposta := Requisicao{ACK: true, id: p.id, contReq: p.contReqLocal}
				inCh[pedido.id] <- resposta
			} else if p.contReqLocal == pedido.contReq && pedido.id < p.id {
				fmt.Println("processo" + strconv.Itoa(p.id) + " respondeu requisicao com ok pois pedido tem id menor para o processo " + strconv.Itoa(pedido.id))
				resposta := Requisicao{ACK: true, id: p.id, contReq: p.contReqLocal}
				inCh[pedido.id] <- resposta
			} else {
				fmt.Println("processo" + strconv.Itoa(p.id) + " recebeu requisicao competiu e perdeu para o processo  " + strconv.Itoa(pedido.id))
				p.filaPendentes = append(p.filaPendentes, pedido)
			}
		case 2:
			fmt.Println("processo" + strconv.Itoa(p.id) + " esta na secao critica adicionou pendente  " + strconv.Itoa(pedido.id))
			p.filaPendentes = append(p.filaPendentes, pedido)
			fmt.Println("processo: " + strconv.Itoa(pedido.id) + " foi adicionado como pendente do processo: " + strconv.Itoa(p.id))
		}

	}
}
