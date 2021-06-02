# TCP Client-Server на языке GO
Домашнее задание по метрологии заключалось в написании клиент-серверного приложения с администратором, имеющим права "wr" и клиентами с правами "r". Писал на языке [GOLANG](https://golang.org/), в среде разработки [ATOM](https://atom.io/)

## Реализованные возможности

1. [Редактирование конфига](#Редактирование-конфига)
2. [Шифрование сообщений](#Шифрование-сообщений)
3. [Очередь](#Очередь)
    
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
____
# Очередь

## Подсветка кода


**Внимание:** Если в тексте таблицы нужно использовать символ "вертикальная черта - `|`", то в место него необходимо написать замену на комбинацию HTML-кода* `&#124;`, это нужно для того, что бы таблица не потеряла ориентации.    
*) - Можно использовать ASCII и/или UTF коды.

**Пример:**
```
| Обозначение | Описание | Пример регулярного выражения|
|----:|:----:|:----------|
| literal | Строка содержит символьный литерал literal | foo |
| re1&#124;re2 | Строка содержит регулярные выражения `rel` или `re2` | foo&#124;bar |
```
**Результат:**

| Обозначение | Описание | Пример регулярного выражения|
|----:|:----:|:----------|
| literal | Строка содержит символьный литерал literal | foo |
| re1&#124;re2 | Строка содержит регулярные выражения `rel` или `re2` | foo&#124;bar |

[:arrow_up:Оглавление](#Оглавление) 
____
