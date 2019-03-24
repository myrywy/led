package led

import (
    "os"
    "fmt"
    "strconv"
    "errors"
)

type BinaryIo interface {
    readPinValue() (int64, error) 
    setPinValue(value int) error 
}

type Pin int

var VALUE_PATH string = "/sys/class/gpio/gpio%d/value"

func (pin_number *Pin) readPinValue() (int64, error) {
    pin_value_path := fmt.Sprintf(VALUE_PATH, pin_number)
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

func (pin_number *Pin) setPinValue(value int) error {
    pin_value_path := fmt.Sprintf(VALUE_PATH, pin_number)
    f, err := os.OpenFile(pin_value_path, os.O_WRONLY, 0755)
    if err != nil {
        return errors.New("Cannot access value file.")
    }
    defer f.Close()
    file_content := []byte(fmt.Sprintf("%d", value))
    _, err = f.Write(file_content)
    if err != nil {
        return errors.New("Error writing value.")
    }
    return nil
}
