# TCP Client-Server на языке GO
Домашнее задание по метрологии заключалось в написании клиент-серверного приложения с администратором, имеющим права "wr" и клиентами с правами "r". Писал на языке [GOLANG](https://golang.org/), в среде разработки [ATOM](https://atom.io/)

## Реализованные возможности

1. [Редактирование конфига](#Редактирование-конфига)
2. [Шифрование сообщений](#Шифрование-сообщений)
3. [Очередь](#Очередь)
4. [Неблокирующая обработка клиентского запроса](#неблокирующая-обработка-клиентского-запроса)
    
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
____
## P.S.
Эта штука работает. Правда я до конца не разобрался с файлами и папками в Golang, поэтому у меня дублируется util.go и возможно конфиг придется кидать в корневую папку пользователя linux (на винде все норм).
