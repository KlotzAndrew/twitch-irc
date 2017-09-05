package main

import "fmt"
import "net"
import "net/textproto"
import "strings"
import "bufio"
import "time"
import "os"
import "gopkg.in/yaml.v2"
import "io/ioutil"

type Config struct {
  Oauth string `yaml:"oauth"`
  Nickname string `yaml:"nickname"`
  Channel string `yaml:"channel"`
}

func readChannel(config Config) {
  conn, err := net.Dial("tcp", "irc.chat.twitch.tv:6667")
  if err != nil { panic(err) }
  defer conn.Close()

  conn.Write([]byte("PASS " + config.Oauth + "\r\n"))
  conn.Write([]byte("NICK " + config.Nickname + "\r\n"))
  conn.Write([]byte("JOIN " + config.Channel + "\r\n"))

  reader := bufio.NewReader(conn)
  tp := textproto.NewReader(reader)

  for {
    msg, err := tp.ReadLine()
    if err != nil { panic(err) }

    fmt.Println(msg)

    msgParts := strings.Split(msg, " ")

    if msgParts[0] == "PING" {
      conn.Write([]byte("PONG " + msgParts[1]))
      continue
    }
  }
}


func main() {
  var config Config
  yamlFile, err := ioutil.ReadFile("config.yml")

  err = yaml.Unmarshal(yamlFile, &config)
  if err != nil { panic(err) }

  conn, err := net.Dial("tcp", "irc.chat.twitch.tv:6667")
  if err != nil { panic(err) }
  defer conn.Close()

  conn.Write([]byte("PASS " + config.Oauth + "\r\n"))
  conn.Write([]byte("NICK " + config.Nickname + "\r\n"))
  conn.Write([]byte("JOIN " + config.Channel + "\r\n"))

  go readChannel(config)

  time.Sleep(time.Second * 3)

  fmt.Println("sending message...")
  scanner := bufio.NewScanner(os.Stdin)
  for scanner.Scan() {
    conn.Write([]byte("PRIVMSG #newtonee :" + scanner.Text() + " \r\n"))
  }
}
