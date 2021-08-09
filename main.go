package main

import (
	"embed"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
    "flag"

	"github.com/gorilla/mux"
	"gopkg.in/olahol/melody.v1"

	"github.com/lirix360/ReadmangaGrabber/config"
	"github.com/lirix360/ReadmangaGrabber/data"
	"github.com/lirix360/ReadmangaGrabber/logger"
	"github.com/lirix360/ReadmangaGrabber/manga"
	"github.com/lirix360/ReadmangaGrabber/tools"
)

//go:embed index.html
var webUI embed.FS

func main() {
    
    //Агрументы командной строки
    addr := flag.String("addr","127.0.0.1","ip адрес сервера или домен")
    port := flag.String("port", "8888", "Порт сервера")
    flag.Parse()
    //Преобразование аргументов
    url := *addr+":"+*port
    //
    
	var err error

	logger.Log.Info("Запуск приложения!")

	r := mux.NewRouter()
    
    m := melody.New()

	r.HandleFunc("/saveConfig", config.SaveConfig)
	r.HandleFunc("/loadConfig", config.LoadConfig)

	r.HandleFunc("/getChaptersList", manga.GetChaptersList)
	r.HandleFunc("/downloadManga", manga.DownloadManga)

	r.HandleFunc("/closeApp", func(w http.ResponseWriter, r *http.Request) {
		logger.Log.Info("Закрытие приложения...")
		os.Exit(0)
	})

	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		err := m.HandleRequest(w, r)
		if err != nil {
			logger.Log.Error("Ошибка при обработке данных WS:", err)
		}
	})

	go func() {
		for {
			msgData := <-data.WSChan
			wsData, err := json.Marshal(msgData)
			if err != nil {
				logger.Log.Error("Ошибка при сериализации данных для отправки через WS:", err)
				continue
			}
			err = m.Broadcast(wsData)
			if err != nil {
				logger.Log.Error("Ошибка при отправке данных через WS:", err)
				continue
			}
		}
	}()

	r.PathPrefix("/").Handler(http.FileServer(http.FS(webUI)))

	srv := &http.Server{
		Handler:      r,
		Addr:         url,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	err = tools.OpenBrowser("http://"+url+"/")
	if err != nil {   
		logger.Log.Fatal("Ошибка при открытии браузера:", err)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
	go func() {
		<-signalChan
		data.WSChan <- data.WSData{Cmd: "closeApp"}
		logger.Log.Info("Закрытие приложения...")
		os.Exit(0)
	}()

	err = srv.ListenAndServe()
	if err != nil {
		logger.Log.Fatal("Ошибка при запуске веб-сервера:", err)
	}
}
