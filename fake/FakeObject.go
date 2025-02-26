package fake

import "sync"

type FakeObj struct {
	mu sync.Mutex
	v  int
}

func (obj *FakeObj) SetVal(v int) { // méthode de classe avec passage par référence
	obj.mu.Lock() // bloque l'utilisation de l'objet
	obj.v = v
	obj.mu.Unlock() // débloque l'utilisation de l'objet
}

func (obj *FakeObj) GetVal() int {
	obj.mu.Lock()
	defer obj.mu.Unlock()
	return obj.v
}
