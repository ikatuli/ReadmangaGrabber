# Manga Grabber (v2) [![Build status](https://api.travis-ci.com/lirix360/ReadmangaGrabber.svg?branch=master)](https://travis-ci.com/github/lirix360/ReadmangaGrabber) [![Go Report Card](https://img.shields.io/badge/go%20report-A+-brightgreen.svg?style=flat)](https://goreportcard.com/report/github.com/lirix360/readmangagrabber)

Утилита для скачивания манги с сайтов ReadManga, MintManga, SelfManga и MangaLib.

## Возможности

* Скачивание целой манги / указанного списка глав из манги
* Создание PDF файлов для скачанных глав
* Создание CBZ файлов для скачанных глав

**Возможности скачивания платной манги нет и не будет!**

![Интерфейс](https://lirix360.github.io/ReadmangaGrabber/screenshot.png?raw=true)

## Аргументы командной строки

 * addr string  
        ip адрес сервера или домен (default "127.0.0.1")
 * conf string__
        Путь к конфигурационному файлу
 * port string  
        Порт сервера (default "8888")

<<<<<<< HEAD
        
## Примеры

### Запуск с изменением сетивого адреса и порта 

```
./grabber_linux_x64 --addr "192.168.42.46"
```


```
./grabber_linux_x64  -port "8080" --addr "my.localdomain"
```

### Запуск с конфигурационного файла

```
./grabber_linux_x64 -conf ~/builds/grabber_config.json
```
