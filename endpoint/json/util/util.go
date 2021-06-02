package util

import (
  "os"
  "strings"
  "encoding/json"
  "io/ioutil"
  "net"
  "time"
  "fmt"
)

//формат времени для лога и JSON
const TIME_LAYOUT = "Jan 2 2006 15.04.05 -0700 MST"
//токены нужны для кодирования данных при отправке клиентам
var ENCODING_UNENCODED_TOKENS = []string{"%", ":", "[", "]", ",", "\""}
var ENCODING_ENCODED_TOKENS = []string{"%25", "%3A", "%5B", "%5D", "%2C", "%22"}
var DECODING_UNENCODED_TOKENS = []string{":", "[", "]", ",", "\"", "%"}
var DECODING_ENCODED_TOKENS = []string{"%3A", "%5B", "%5D", "%2C", "%22", "%25"}

//структура для имени пользователя и данных подключения клиента
type Client struct {
  //соединение
  Connection net.Conn
  //имя
  Username string
  //конфиг
  Properties Properties
}
//закрываем соединение и убираем клиента
func (client *Client) Close(doSendMessage bool) {
  if (doSendMessage) {
    SendClientMessage("disconnect", "", client, false, client.Properties)
  }
  client.Connection.Close();
  clients = removeEntry(client, clients);
}

//регистрируем соединение и кешируем его
func (client *Client) Register() {
  clients = append(clients, client);
}

//структура лога
type Action struct {
  //команды
  Command string      `json:"command"`
  //различные сообщения
  Content string      `json:"content"`
  //имя
  Username string     `json:"username"`
  //ip
  IP string           `json:"ip"`
  //отметка времени активности
  Timestamp string    `json:"timestamp"`
}

//основная конфигурация
type Properties struct {
  Hostname string
  Port string
  JSONEndpointPort string
  HasEnteredTheLobbyMessage string
  HasLeftTheLobbyMessage string
  ReceivedAMessage string

  //Аккаунты
  FirstAccount string
  SecondAccount string
  ThirdAccount string
  FourthAccount string
  FifthAccount string
  SixthAccount string
  SevenAccount string

  //Приоритеты аккаунтов
  FirstAccountPriority string
  SecondAccountPriority string
  ThirdAccountPriority string
  FourthAccountPriority string
  FifthAccountPriority string
  SixthAccountPriority string
  SevenAccountPriority string

  //Максимальный приритет подключившихся пользователей
  MaxUsers string
  //Лог
  LogFile string
}

//все действия (чат, подключение / отключение)
//произошедшие во время работы сервера
var actions = []Action{}
//кэшированные свойства конфигурации
var config = Properties{}
//список клиентов (статический)
var clients []*Client

//загружаем свойства конфигурации из файла "config.json"
func LoadConfig() Properties {
  if (config.Port != "") {
    return config;
  }
  pwd, _ := os.Getwd()

  payload, err := ioutil.ReadFile(pwd + "/config.json")
  CheckForError(err, "Невозможно открыть конфиг")

  var dat map[string]interface{}
  err = json.Unmarshal(payload, &dat)
  CheckForError(err, "Кажется JSON нужно закинуть в root папку (либо он сломался)")

//тащим конфиг прямо в структуре Properties
  var rtn = Properties {
    Hostname: dat["Hostname"].(string),
    Port: dat["Port"].(string),
    JSONEndpointPort: dat["JSONEndpointPort"].(string),
    HasEnteredTheLobbyMessage: dat["HasEnteredTheLobbyMessage"].(string),
    HasLeftTheLobbyMessage: dat["HasLeftTheLobbyMessage"].(string),
    ReceivedAMessage: dat["ReceivedAMessage"].(string),
    FirstAccount: dat["FirstAccount"].(string),
    SecondAccount: dat["SecondAccount"].(string),
    ThirdAccount: dat["ThirdAccount"].(string),
    FourthAccount: dat["FourthAccount"].(string),
    FifthAccount: dat["FifthAccount"].(string),
    SixthAccount: dat["SixthAccount"].(string),
    SevenAccount: dat["SevenAccount"].(string),
    FirstAccountPriority: dat["FirstAccountPriority"].(string),
    SecondAccountPriority: dat["SecondAccountPriority"].(string),
    ThirdAccountPriority: dat["ThirdAccountPriority"].(string),
    FourthAccountPriority: dat["FourthAccountPriority"].(string),
    FifthAccountPriority: dat["FifthAccountPriority"].(string),
    SixthAccountPriority: dat["SixthAccountPriority"].(string),
    SevenAccountPriority: dat["SevenAccountPriority"].(string),
    MaxUsers: dat["MaxUsers"].(string),
    LogFile: dat["LogFile"].(string),
  }
  config = rtn;
  return rtn;
}

//удалить запись о клиенте из сохраненных клиентов
func removeEntry(client *Client, arr []*Client) []*Client {
  rtn := arr
  index := -1
  for i, value := range arr {
    if (value == client) {
      index = i;
      break;
    }
  }

  if (index >= 0) {
    //у нас есть совпадение, создаем новый массив без совпадения
    rtn = make([]*Client, len(arr)-1)
    copy(rtn, arr[:index])
    copy(rtn[index:], arr[index+1:])
  }

  return rtn;
}

//отправил сообщение всем клиентам (кроме отправителя)
func SendClientMessage(messageType string, message string, client *Client, thisClientOnly bool, props Properties) {

  if (thisClientOnly) {
    //это сообщение предназначено только для указанного клиента
    message = fmt.Sprintf("/%v", messageType);
    fmt.Fprintln(client.Connection, message)

   } else if (client.Username != "") {
    //это сообщение предназначено для всех, кроме предоставленного клиента
    LogAction(messageType, message, client, props);

    //создать полезную нагрузку для отправки клиентам
    payload := fmt.Sprintf("/%v [%v] %v", messageType, client.Username, message);

    for _, _client := range clients {

      //написать сообщение клиенту
      if ((thisClientOnly && _client.Username == client.Username) ||
          (!thisClientOnly && _client.Username != "")) {

        fmt.Fprintln(_client.Connection, payload)
      }
    }
  }
}

//чекаем ошибочки
func CheckForError(err error, message string) {
  if err != nil {
      println(message + ": ", err.Error())
      os.Exit(1)
  }
}

//двойные кавычки одинарные кавычки
func EncodeCSV(value string) (string) {
  return strings.Replace(value, "\"", "\"\"", -1)
}

//простая кодировка http-ish для обработки специальных символов
func Encode(value string) (string) {
  return replace(ENCODING_UNENCODED_TOKENS, ENCODING_ENCODED_TOKENS, value)
}

//декодирование http-ish для обработки специальных символов
func Decode(value string) (string) {
  return replace(DECODING_ENCODED_TOKENS, DECODING_UNENCODED_TOKENS, value)
}

//замена токенов from на токены to (оба массива должны быть одинаковой длины)
func replace(fromTokens []string, toTokens []string, value string) (string) {
  for i:=0; i<len(fromTokens); i++ {
      value = strings.Replace(value, fromTokens[i], toTokens[i], -1)
  }
  return value;
}

//записываем действие в лог и минилог
func LogAction(action string, message string, client *Client, props Properties) {
  ip := client.Connection.RemoteAddr().String()
  timestamp := time.Now().Format(TIME_LAYOUT)

  actions = append(actions, Action {
    Command: action,
    Content: message,
    Username: client.Username,
    IP: ip,
    Timestamp: timestamp,
  })

  if (props.LogFile != "") {
    if (message == "") {
      message = "N/A"
    }
    fmt.Printf("Мини лог: %s, %s, %s\n", action, message, client.Username);

    logMessage := fmt.Sprintf("\"%s\", \"%s\", \"%s\", \"%s\", \"%s\"\n",
      EncodeCSV(client.Username), EncodeCSV(action), EncodeCSV(message),
        EncodeCSV(timestamp), EncodeCSV(ip))

    f, err := os.OpenFile(props.LogFile, os.O_APPEND|os.O_WRONLY, 0600)
    if (err != nil) {
      err = ioutil.WriteFile(props.LogFile, []byte{}, 0600)
      f, err = os.OpenFile(props.LogFile, os.O_APPEND|os.O_WRONLY, 0600)
      CheckForError(err, "Невозможно создать файл")
    }

    defer f.Close()
    _, err = f.WriteString(logMessage)
    CheckForError(err, "Невозможно записать в лог")
  }
}

func QueryMessages(actionType string, search string, username string) ([]Action) {

  isMatch := func(action Action) (bool) {
    if (actionType != "" && action.Command != actionType) {
      return false;
    }
    if (search != "" && !strings.Contains(action.Content, search)) {
      return false;
    }
    if (username != "" && action.Username != username) {
      return false;
    }
    return true;
  }

  rtn := make([]Action, 0, len(actions))

  //узнаем, какие элементы соответствуют критериям поиска, и добавляем их к тому, что мы будем возвращать
  for _, value := range actions {
    if (isMatch(value)) {
      rtn = append(rtn, value)
    }
  }

  return rtn;
}
