package lfu

type lfu struct {

	maxBytes int

	onEvicted func(key string, value interface{})

	usedBytes int


}
