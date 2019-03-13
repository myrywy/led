package main

import (
    "fmt"
    "os"
    "strconv"
    "errors"
)

var VALUE_PATH string = "/sys/class/gpio/gpio%d/value"

func readPinValue(pin_number int) (int64, error) {
    pin_value_path := fmt.Sprintf(VALUE_PATH, pin_number)
    fmt.Println("reading ", pin_value_path)
    f, err := os.OpenFile(pin_value_path, os.O_RDONLY, 0755)
    if err != nil {
        return 0, errors.New("Cannot access value file.")
    }
    defer f.Close()
    file_content := make([]byte, 1)
    bytes_read, err := f.Read(file_content)
    if err != nil || bytes_read != 1 {
        return 0, errors.New("Error reading value.")
    }
    led_state, err := strconv.ParseInt(string(file_content), 10, 8)
    if err != nil {
        return 0, errors.New("Value pseudo-file content invalid.")
    }
    return led_state, nil
}

func main() {
        value, err := readPinValue(21)
        if err != nil {
            fmt.Println(err)
            return
        }
        fmt.Println("Wartość: ", value)
}


