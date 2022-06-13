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

var recurso = make(chan string, 1)
var mutex = make(chan struct{}, 1)
var inCh inputChan

func main() {
	mutex <- struct{}{}
	for i := 0; i < N; i++ {
		inCh[i] = make(chan Requisicao, 1)
		var pendentes = make([]Requisicao, N)
		p := ComponentesProcesso{id: i, estado: 0, cont: 0, contReqLocal: 0, filaPendentes: pendentes}
		go processo(p)
		go processoRespondedor(p)
	}
	retorno := <-recurso
	fmt.Println(retorno)

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
				r := Requisicao{ACK: true, id: p.id, contReq: p.contReqLocal}
				inCh[i] <- r
			}
		}
		retorno := true
		for i := 0; i < N-1; i++ {
			resposta := <-inCh[p.id]
			if resposta.ACK == false {
				retorno = false
			}
		}
		if retorno == true {
			p.estado = 2
			mensagem := "Processo" + strconv.Itoa(p.id) + "Esta acessando o recurso"
			recurso <- mensagem //acessa recurso
		}
		mutex <- struct{}{} //sai secao critica
		fmt.Println("processo saiu da secao critica: ", p.id)

		p.estado = 0
		for i := range p.filaPendentes {
			resposta := Requisicao{ACK: true, id: p.id, contReq: p.contReqLocal}
			inCh[p.filaPendentes[i].id] <- resposta
		}
		p.filaPendentes = make([]Requisicao, N)

	}
}
func processoRespondedor(p ComponentesProcesso) {
	fmt.Println("processo respondedor executando: ", p.id)
	for {
		pedido := <-inCh[p.id]
		fmt.Println("processo" + strconv.Itoa(p.id) + "recebeu requisicao  do processo  " + strconv.Itoa(pedido.id))
		if pedido.contReq > p.cont {
			p.cont = pedido.contReq + 1
		}

		switch p.estado {
		case 0:
			resposta := Requisicao{ACK: true, id: p.id, contReq: p.contReqLocal}
			inCh[pedido.id] <- resposta
			return

		case 1:
			if p.contReqLocal < pedido.contReq {
				resposta := Requisicao{ACK: true, id: p.id, contReq: p.contReqLocal}
				inCh[pedido.id] <- resposta
				return
			} else if p.contReqLocal == pedido.contReq && pedido.id < p.id {
				resposta := Requisicao{ACK: true, id: p.id, contReq: p.contReqLocal}
				inCh[pedido.id] <- resposta
				return
			}
			p.filaPendentes = append(p.filaPendentes, pedido)

		case 2:
			p.filaPendentes = append(p.filaPendentes, pedido)
		}

	}
}
