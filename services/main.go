package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func main() {
	// Définir le répertoire à surveiller
	watchDir := "../my-vue-project/src"

	// Créer un nouveau watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Ajouter le répertoire à surveiller
	err = filepath.Walk(watchDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return watcher.Add(path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// Fonction pour compiler les fichiers Vue.js
	compileVue := func() {
		cmd := exec.Command("npm", "run", "serve")
		cmd.Dir = watchDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Println("Erreur de compilation:", err)
		} else {
			fmt.Println("Compilation réussie.")
		}
	}

	// Fonction pour ouvrir l'URL dans le navigateur
	openBrowser := func(url string) {
		var err error

		switch runtime.GOOS {
		case "linux":
			err = exec.Command("xdg-open", url).Start()
		case "windows":
			err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
		case "darwin":
			err = exec.Command("open", url).Start()
		default:
			err = fmt.Errorf("unsupported platform")
		}

		if err != nil {
			log.Println("Erreur lors de l'ouverture du navigateur:", err)
		}
	}

	// Lancer une première compilation et ouvrir l'URL
	compileVue()
	openBrowser("http://localhost:8080")

	// Surveiller les événements de changement de fichier
	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					fmt.Println("Fichier modifié:", event.Name)
					compileVue()
					openBrowser("http://localhost:8080")
				}
			case err := <-watcher.Errors:
				log.Println("Erreur:", err)
			}
		}
	}()

	// Garder le programme en cours d'exécution
	<-done
}
