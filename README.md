# TCP сервер для отладки сетевой?сетевых бибилотек

TCP сервер принимает соединение, отображает поступающие в него данные, и, при необходимости генерирует исходящие данные.
Размер исходящих пакетов данных можно регулировать.

Написан для тестирования/отладки zifi бибилотеки

## Сборка:
 go get
 go build

## Закуск
Unix/MAC:
```
 ./tcp-server-zifilibtest
```
Windows:
```
tcp-server-zifilibtest
```

## Параметры запуска

```
  -h
      help.
      Выводит подсказку по пратаметрам
  -delay int
    	delay as seconds (default 2)
      Интервал в секундах между циклом приемом/передачей данных
  -host string
    	listen address (default "localhost")
      Адрес интерфейса сервера
  -port int
    	port number (default 3333)
      Номер порта
  -type string
    	type tcp/udp (default "tcp")
      Тип протокола (tcp/udp)
  -maxblocksize int
    	genetated block size. maximum. (default: 0)
      Максимальный размер блока генерируемых на отправку данных
      Если минимвльный != максимальному то случаным образом генерируется пакет длиной в мах-мин
  -minblocksize int
    	genetated block size. minimum. (default: 0)
      Минимальный размер блока генерируемых данных
      Если минимвльный != максимальному то случаным образом генерируется пакет длиной в мах-мин
  -onlyrecevedata
    	Receve data only mode (default: false)
      Режим "только принемать данные"
  -onlysenddata
    	Send data only mode (default: false)
      Режим "только отправлять данные"
  -randomdatasend
    	Generate random data to send (default: false)
      Выходные данные генерируются случайным набором байт

```
