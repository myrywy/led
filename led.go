package main

import (
    "fmt"
    "os"
    "strconv"
)

var VALUE_PATH string = "/sys/class/gpio/gpio21/value"

func main() {
        f, err := os.OpenFile(VALUE_PATH, os.O_RDONLY, 0755)
        if err != nil {
            fmt.Println("Przykro mi, błąd przy otwieraniu pliku :(")
            fmt.Println(err)
            return
        }
        defer f.Close()
        value := make([]byte, 1)
        bytesRead, err := f.Read(value)
        if err != nil || bytesRead != 1 {
            fmt.Println("Sorry, blad odczytu pliku :c")
            fmt.Println(err)
            return
        }
        led_state, err := strconv.ParseInt(string(value), 10, 8)
        if err != nil {
            fmt.Println("Inwalidzka zawartość pliku :(")
            return
        }
        fmt.Println("Wartość: ", led_state)
}


