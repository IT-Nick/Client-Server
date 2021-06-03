package main

import (
  "net"
  "fmt"
  "bufio"
  "strings"
  "regexp"
  "./util"
  "./endpoint/json"
  "time"
  "strconv"
)																																																																																							//ПАТехецкий далбаеб//ПАТехецкий далбаеб//ПАТехецкий далбаеб

//хранит в себе наших пользователей и их приоритеты
var m = make(map[string]int)

//хранит в себе подключившихся пользователей
var users = make([]int, 0)

//главная программа
func main() {

  // запускаем сервер
  properties := util.LoadConfig()
  psock, err := net.Listen("tcp", ":" + properties.Port)
  util.CheckForError(err, "Не удалось создать сервер")

  fmt.Printf("Сервер запущен на порту %v...\n", properties.Port)

	//переводим приоритеты (string) из конфига в (int) для работы с картой
	firstP, _ := strconv.Atoi(properties.FirstAccountPriority)
	SecondP, _ := strconv.Atoi(properties.SecondAccountPriority)
	ThirdP, _ := strconv.Atoi(properties.ThirdAccountPriority)
	FourthP, _ := strconv.Atoi(properties.FourthAccountPriority)
	FifthP, _ := strconv.Atoi(properties.FifthAccountPriority)
	SixthP, _ := strconv.Atoi(properties.SixthAccountPriority)
	SevenP, _ := strconv.Atoi(properties.SevenAccountPriority)

 	//добавляем приоритеты в карту
	m[properties.FirstAccount] = firstP
	m[properties.SecondAccount] = SecondP
	m[properties.ThirdAccount] = ThirdP
	m[properties.FourthAccount] = FourthP
	m[properties.FifthAccount] = FifthP
	m[properties.SixthAccount] = SixthP
	m[properties.SevenAccount] = SevenP

	//Ставим максимальный приоритет подключившихся пользователей
	maxusers, _ := strconv.Atoi(properties.MaxUsers)
	users = append(users, maxusers)


  // запускаем JSON-конфиг
  go json.Start();
  for {
    // принимаем соединения
    conn, err := psock.Accept()
    util.CheckForError(err, "Невозможно принять соединение")

		// отслеживаем детали клиента
    client := util.Client{Connection: conn, Properties: properties}
    client.Register();

		//реализуем неблокирующую обработку клиентского запроса
    channel := make(chan string)
    go waitForInput(channel, &client)
    go handleInput(channel, &client, properties)

		util.SendClientMessage("ready", properties.Port, &client, true, properties)
  }
}

//находим приоритет который нужно удалить из среза при дисконекте
func findMe(users []int, priority int) int {
	for i, n := range users {
			 if (priority == n) {
				 return i
			 }
	}
	return len(users)
}

//удаляем приоритет из среза при дисконекте
func removeIndex(users []int, thisindexremove int) []int {
    return append(users[:thisindexremove], users[thisindexremove+1:]...)
}

//Находим наименьший приоритет
func add(min int) int {
	for _, v := range users {
			 if (v < min) {
				 min = v
			 }
	}
	return min
}


// ждем ввода от клиента и сигнализируем каналу
func waitForInput(out chan string, client *util.Client) {
  defer close(out)

  reader := bufio.NewReader(client.Connection)
  for {
    line, err := reader.ReadBytes('\n')
    if err != nil {
      // соединение было закрыто, удаляем клиента
      client.Close(true);
      return
    }
    out <- string(line)
  }
}

// прослушиваем обновления канала для клиента и обрабатываем сообщение
func handleInput(in <-chan string, client *util.Client, props util.Properties) {
	min := users[0]
  for {
  	message := <- in
    if (message != "") {
      message = strings.TrimSpace(message)
      action, body := getAction(message)
      if (action != "") {
        switch action {

          //пользователь предоставил свое имя
          case "user":
		client.Username = body
		priority := m[client.Username]//priority - приоритет текущего клиента

		sum := add(min)
		users = append(users, priority)

			fmt.Println("Наша очередь (первое 8 - лимит): ", users)
			fmt.Println("Администратор под приоритетом: ", sum)

				//проверка на подключение приоритетного клиента
				if priority <= sum {
					util.SendClientMessage("writeandread", body, client, true, props)
					util.SendClientMessage("readonlynewconnect", body, client, false, props)
				} else {
					util.SendClientMessage("readonly", body, client, true, props)
				}

					//проверка на несуществующее имя
					if _, ok := m[body]; ok {
            					util.SendClientMessage("connect", "", client, false, props)
					} else {
						util.SendClientMessage("nopermission", body, client, true, props)
						time.Sleep(1 * time.Second)
						client.Close(false);
					}

	  // пользователь отправил сообщение
          case "message":
		sum := add(min)
		priority := m[client.Username]

		//проверка на наличие Write&Read
		if priority <= sum {
          		util.SendClientMessage("message", body, client, false, props)
		} else {
			util.SendClientMessage("readonly", body, client, true, props)
		}

          // пользователь отключается
          case "disconnect":
		priority := m[client.Username]

		sum := add(min)
		//находим клиента из среза
		thisindexremove := findMe(users, priority)
		//и удаляем
		users = removeIndex(users, thisindexremove)

			fmt.Println("Наша очередь (первое 8 - лимит): ", users)
			fmt.Println("Администратор под приоритетом: ", sum)

				if priority <= sum {
					util.SendClientMessage("ifdisconnect", body, client, false, props)
				}
		
            			client.Close(false);
          default:
            util.SendClientMessage("unrecognized", action, client, true, props)
        }
      }
    }
  }
}

// разбираем содержимое сообщения и возвращаем отдельные значения
func getAction(message string) (string, string) {
  actionRegex, _ := regexp.Compile(`^\/([^\s]*)\s*(.*)$`)
  res := actionRegex.FindAllStringSubmatch(message, -1)
  if (len(res) == 1) {
    return res[0][1], res[0][2]
  }
  return "", ""
}
