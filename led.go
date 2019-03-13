package main

import (
    "fmt"
    "os"
    "strconv"
)

var VALUE_PATH string = "/sys/class/gpio/gpio%i/value"

func readPinValue(pin_number int) -> int, error {
    pin_value_path := fmt.Spritf(VALUE_PATH, pin_number)
    f, err := os.OpenFile(pin_value_path, os.O_RDONLY, 0755)
    if err != nil {
        return 0, error("Cannot access value file.")
    }
    defer f.Close()
    file_content := make([]byte, 1)
    bytes_read, err := f.Read(file_content)
    if err != nil || bytes_read != 1 {
        return 0, "Error reading value."
    }
    led_state, err := strconv.ParseInt(string(file_content), 10, 8)
    if err != nil {
        return 0, "Value pseudo-file content invalid."
    }
    return led_state, nil
}

func main() {
        fmt.Println("Wartość: ", readPinValue(21))
}


