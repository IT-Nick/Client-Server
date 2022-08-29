# TCP Client-Server на языке GO
Клиент-Серверное приложение с администратором, имеющим права "wr" и клиентами с правами "r". 
## Реализованные возможности

1. [Редактирование конфига](#Редактирование-конфига)
2. [Шифрование сообщений](#Шифрование-сообщений)
3. [Очередь](#Очередь)
4. [Неблокирующая обработка клиентского запроса](#неблокирующая-обработка-клиентского-запроса)
5. [Установка и запуск](#Установка-и-запуск)
    
## Редактирование конфига
В конфиге можно менять [ИМЯ], [ПРИОРИТЕТ], а так же другие параметры как [ВЫВОДИМЫЕ СООБЩЕНИЯ] и [АДРЕС И ПОРТ]
```JSON
{
  "Port": "5555",
  "JSONEndpointPort": "8080",
  "Hostname": "localhost",
  "HasEnteredTheLobbyMessage": "[%s] вошел в лобби",
  "HasLeftTheLobbyMessage": "[%s] вышел из лобби",
  "ReceivedAMessage": "[%s] сказал: %s",
  "FirstAccount": "Anton",
  "SecondAccount": "Kirill",
  "ThirdAccount": "Stepan",
  "FourthAccount": "Анатолий",
  "FifthAccount": "Алексей",
  "SixthAccount": "Добавьте_сюда_новый_аккаунт",
  "SevenAccount": "Добавьте_сюда_новый_аккаунт",
  "FirstAccountPriority": "1",
  "SecondAccountPriority": "2",
  "ThirdAccountPriority": "3",
  "FourthAccountPriority": "4",
  "FifthAccountPriority": "5",
  "SixthAccountPriority": "6",
  "SevenAccountPriority": "7",
  "MaxUsers": "8",
  "LogFile": "./log.csv"
}

```

## Шифрование сообщений

Были созданы токены, которые нужны для кодирования данных при отправке клиентам
```GO
var ENCODING_UNENCODED_TOKENS = []string{"%", ":", "[", "]", ",", "\""}
var ENCODING_ENCODED_TOKENS = []string{"%25", "%3A", "%5B", "%5D", "%2C", "%22"}
var DECODING_UNENCODED_TOKENS = []string{":", "[", "]", ",", "\"", "%"}
var DECODING_ENCODED_TOKENS = []string{"%3A", "%5B", "%5D", "%2C", "%22", "%25"}

func EncodeCSV(value string) (string) {
  return strings.Replace(value, "\"", "\"\"", -1)
}

//кодируем
func Encode(value string) (string) {
  return replace(ENCODING_UNENCODED_TOKENS, ENCODING_ENCODED_TOKENS, value)
}

//декодируем
func Decode(value string) (string) {
  return replace(DECODING_ENCODED_TOKENS, DECODING_UNENCODED_TOKENS, value)
}

//заменяем
func replace(fromTokens []string, toTokens []string, value string) (string) {
  for i:=0; i<len(fromTokens); i++ {
      value = strings.Replace(value, fromTokens[i], toTokens[i], -1)
  }
  return value;
}
```

# Очередь
Очередь реализована через срез в который вносится приоритет пользователя при подключении (сам приоритет берется из мапы) и соответственно при дисконекте приоритет удаляется из среза. Только приоритетный пользователь имеет статус Write&Read. Так же при заходе приоритетного клиента у всех статус сменится на ReadOnly, а при его выходе приоритетным станет следующий клиент с наивысшим приоритетом. (т.е. твоя позиция в очереди не зависит от того, когда ты зашел, а только от твоего приоритета и приоритета других присоединившихся пользователей).
```GO
//хранит в себе наших пользователей и их приоритеты
var m = make(map[string]int)

//хранит в себе подключившихся пользователей
var users = make([]int, 0)
```
# Неблокирующая обработка клиентского запроса
Для реализации неблокирующего запроса создаем две горутины (одна читает, другая пишет) и там уже внутри функций реализуем "чат"
```GO
channel := make(chan string)
go waitForInput(channel, &client)
go handleInput(channel, &client, properties)
```
# Установка и запуск
Для запуска программ нужно закинуть config.json в корневую папку вместе с программами и запустить сначала "server" а затем "client" и все.
____
Для тестирования/редактирования кода нужно скачать язык [go](https://golang.org/) и установить его следуя официальной инструкции:

    Extract the archive you downloaded into /usr/local, creating a Go tree in /usr/local/go.
    Important: This step will remove a previous installation at /usr/local/go, if any, prior to extracting. Please back up any data before proceeding.

    For example, run the following as root or through sudo:

    rm -rf /usr/local/go && tar -C /usr/local -xzf go1.16.5.linux-amd64.tar.gz

    Add /usr/local/go/bin to the PATH environment variable.
    You can do this by adding the following line to your $HOME/.profile or /etc/profile (for a system-wide installation):

    export PATH=$PATH:/usr/local/go/bin

    Note: Changes made to a profile file may not apply until the next time you log into your computer. To apply the changes immediately, just run the shell             commands directly or execute them from the profile using a command such as source $HOME/.profile.

    Verify that you've installed Go by opening a command prompt and typing the following command:

    $ go version

    Confirm that the command prints the installed version of Go.
    
Для компиляции программы используем команды: 
```
GOOS=linux GOARCH=amd64 GO111MODULE=off go build server.go
GOOS=linux GOARCH=amd64 GO111MODULE=off go build client.go
```
Либо без ручного указания платформы: 
```
GO111MODULE=off go build server.go
GO111MODULE=off go build client.go
```
____
## P.S.
Эта штука работает. Правда я до конца не разобрался с файлами и папками в Golang, поэтому у меня дублируется util.go и возможно конфиг придется кидать в корневую папку пользователя linux (на винде все норм).
