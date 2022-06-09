package main

import (
	"sync"
)

const N = 4            //nro de processos
type Processo struct { //estrutura do processo
	id           int
	estado       int
	cont         int
	contReqLocal int
	mu           sync.Mutex
}

type Requisicao struct { //requisicao de acesso à SC
	idReq   int
	contReq int
	ACK     bool
}

type inputChan [N]chan Processo //canal para cada processo
const channelBufferSize = 1     //tamanho do buffer de cada canal de entrada

func (p *Processo) EntrarSC(cont int) {
	p.estado = 1
	p.cont++
	p.contReqLocal = p.cont
	// se o processo conseguir acessar a SC...
	p.estado = 2
	p.mu.Lock() // espera a trava liberar e pega recurso.
	//...
	p.mu.Unlock() // libera a trava
	p.estado = 0
}

func (p *Processo) ResponderReq(r *Requisicao, idReq int, contReq int, ACK bool) {
	if contReq > p.cont {
		contReq = contReq + 1
	}

	switch p.estado {
	case 0:
		ACK = true
	case 1:
		if p.contReqLocal < contReq {
			//quer entrar na SC, entao se o pedido do outro processo for anterior ao meu, devo liberar ele
			ACK = true
		} else if p.contReqLocal == contReq && idReq < p.id {
			//se o outro processo quer no mesmo instante que esse e seu id (idReq)
			// e menor que o id, este da prioridade ao outro
		} else {
			//se n respondeu ate agr, nao libera o ACK
		}
	case 2:
		//mesmo caso do ultimo else acima, avisa que saiu da SC
		p.mu.Unlock()
		break
	}
}

func main() {

	var inCh inputChan // cada processo tem um canal de entrada, chamado inCh[i]
	for i := 0; i < N; i++ {
		inCh[i] = make(chan Processo, channelBufferSize) // criando cada um dos canais
	}

}
