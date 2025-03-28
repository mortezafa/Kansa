package main

//
//import (
//	dn "github.com/mitchellh/go-ps"
//	"log"
//	"time"
//)
//
//func main() {
//
//	time.Sleep(5 * time.Second)
//
//	allPro, _ := dn.Processes()
//	isAnkirunning := false
//	isAnkiActive := false
//
//	for _, pro := range allPro {
//
//		if pro.Executable() == "anki" {
//			isAnkirunning = true
//		}
//	}
//
//	if isWindowActive("anki") {
//		isAnkiActive = true
//	} else {
//		log.Printf("Anki not active lel")
//	}
//
//	if isAnkirunning && isAnkiActive {
//		for {
//			start := time.Now()
//
//			time.Sleep(2 * time.Second)
//
//			elapsed := time.Since(start)
//			log.Printf("Elapsed time: %s\n", elapsed)
//
//			if !isWindowActive("anki") {
//				log.Printf("exited Anki stopping timer!")
//				break
//			}
//
//		}
//	} else {
//		log.Printf("------------CANT START TIMER------------\n")
//		log.Printf("Anki is running: %v \nAnki is active: %v", isAnkirunning, isAnkiActive)
//	}
//
//}
