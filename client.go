package main

import (
  "fmt"
  "os"
  "net"
  "bufio"
  "regexp"
  "strings"
  "./util"
)


var standardInputMessageRegex, _ = regexp.Compile(`^\/([^\s]*)\s*(.*)$`)
var chatServerResponseRegex, _ = regexp.Compile(`^\/([^\s]*)\s?(?:\[([^\]]*)\])?\s*(.*)$`)

//структура для команд сервера
type Command struct {

  Command, Username, Body string
}

//главная программа
func main() {

  username, password, properties := getConfig();

  conn, err := net.Dial("tcp", properties.Hostname + ":" + properties.Port)
  util.CheckForError(err, "В соединении отказано")
  defer conn.Close()

  //мы слушаем команды сервера и консоли пользователя
  go watchForConnectionInput(username, password, properties, conn)
  for true {
    watchForConsoleInput(conn)
  }
}

//разбираем аргументы, которые будут использоваться при подключении к серверу
func getConfig() (string, string, util.Properties) {
    fmt.Print("Введите имя: ")
    var username string
    fmt.Scanln(&username)
    fmt.Print("Введите пароль: ")
        var password string
    fmt.Scanln(&password)
    properties := util.LoadConfig()
    return username, password, properties
}

//продолжаем слушать ввод с консоли
//отправляем команду серверу
func watchForConsoleInput(conn net.Conn) {
  reader := bufio.NewReader(os.Stdin)

  for true {
    message, err := reader.ReadString('\n')
    util.CheckForError(err, "Потеряно соединение")

    message = strings.TrimSpace(message)
    if (message != "") {
      command := parseInput(message)

      if (command.Command == "") {
        //если команды нет, значит отправляем как сообщение
        sendCommand("message", message, conn);
      } else {
        switch command.Command {
          //команда отключения от сервера
        case "disconnect":
            sendCommand("disconnect", "", conn)

        default:
          fmt.Printf("Несуществующая команда \"%s\"\n", command.Command)
        }
      }
    }
  }
}

//слушаем команды от сервера
func watchForConnectionInput(username string, password string, properties util.Properties, conn net.Conn) {
  reader := bufio.NewReader(conn)

  for true {
    message, err := reader.ReadString('\n')
    util.CheckForError(err, "Потеряно соединение");
    message = strings.TrimSpace(message)
    if (message != "") {
      Command := parseCommand(message)
      switch Command.Command {

        // отправляем наше имя пользователя и пароль
        case "ready":
          sendCommand("user", username, conn)

        //пользователь подключился к серверу чата
        case "connect":
          fmt.Printf(properties.HasEnteredTheLobbyMessage + "\n", Command.Username)

        //и отключился
        case "disconnect":
          fmt.Printf(properties.HasLeftTheLobbyMessage + "\n", Command.Username)

        //разные сообщения при смене администратора
        case "readonlynewconnect":
                  if (Command.Username != username) {
          fmt.Printf("Подключился %v с более высоким приоритетом. У вас статус ReadOnly \n", Command.Username)
          }

        case "writeandreadifdisconnect":
          fmt.Printf("У вас статус Write&Read\n")

        case "ifdisconnect":
                  if (Command.Username != username) {
          fmt.Printf("Администратор был %v. Теперь у нас новый администратор. %v\n", Command.Username)
          }

        case "writeandread":
            fmt.Printf("У вас статус Write&Read\n")

        case "readonly":
            fmt.Printf("У вас статус ReadOnly \n")

        case "nopermission":
            fmt.Printf("Неверное имя пользователя \n")


        //пользователь отправил соо
        case "message":
          if (Command.Username != username) {
            fmt.Printf(properties.ReceivedAMessage + "\n", Command.Username, Command.Body)
          }

      }
    }
  }
}

//отправляем команду на сервер
//команды имеют форму /command
func sendCommand(command string, body string, conn net.Conn) {
  message := fmt.Sprintf("/%v %v\n", util.Encode(command), util.Encode(body));
  conn.Write([]byte(message))
}

//анализируем входное сообщение и возвращаем команду
func parseInput(message string) Command {
  res := standardInputMessageRegex.FindAllStringSubmatch(message, -1)
  if (len(res) == 1) {
    //это команда
    return Command {
      Command: res[0][1],
      Body: res[0][2],
    }
  } else {
    return Command {
      Body: util.Decode(message),
    }
  }
}

//ищем команды
func parseCommand(message string) Command {
  res := chatServerResponseRegex.FindAllStringSubmatch(message, -1)
  if (len(res) == 1) {
    //есть совпадение
    return Command {
      Command: util.Decode(res[0][1]),
      Username: util.Decode(res[0][2]),
      Body: util.Decode(res[0][3]),
    }
  } else {
    return Command{}
  }
}
