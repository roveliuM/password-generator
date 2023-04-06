package main

import (
    "bufio"
    "fmt"
    "log"
    "math/rand"
    "os"
    "strings"
    "sync"
    "time"

    "github.com/schollz/closestmatch"
    "github.com/schollz/gpt-2"
)

func main() {
    // Get user inputs
    var name, surname, year, phone, nicknames string
    fmt.Print("Enter your name: ")
    fmt.Scanln(&name)
    fmt.Print("Enter your surname: ")
    fmt.Scanln(&surname)
    fmt.Print("Enter your year of birth: ")
    fmt.Scanln(&year)
    fmt.Print("Enter your phone number: ")
    fmt.Scanln(&phone)
    fmt.Print("Enter your nicknames (separated by commas): ")
    fmt.Scanln(&nicknames)

    // Load GPT-2 model
    model, err := gpt2.Load("345M")
    if err != nil {
        log.Fatal(err)
    }

    // Generate passwords
    currentYear := strconv.Itoa(time.Now().Year())
    numPasswords := 10000
    passwordsPerRoutine := numPasswords / 4
    passwords := make([]string, numPasswords)

    var wg sync.WaitGroup
    wg.Add(4)

    for i := 0; i < 4; i++ {
        go func(start int) {
            defer wg.Done()
            for j := start; j < start+passwordsPerRoutine; j++ {
                var password string
                switch rand.Intn(7) {
                case 0:
                    password = name + surname + year
                case 1:
                    password = name + year + surname
                case 2:
                    password = surname + name + year
                case 3:
                    password = surname + year + name
                case 4:
                    password = year + name + surname
                case 5:
                    password = year + surname + name
                case 6:
                    password = name + surname + year[2:]
                }
                switch rand.Intn(6) {
                case 0:
                    password += name[0:2] + phone
                case 1:
                    password += surname[0:2] + phone
                case 2:
                    password += year[2:] + phone
                case 3:
                    password += phone + name[0:2]
                case 4:
                    password += phone + surname[0:2]
                case 5:
                    password += phone + year[2:]
                }
                switch rand.Intn(3) {
                case 0:
                    password += currentYear
                case 1:
                    password += currentYear[2:]
                case 2:
                    password += currentYear[1:]
                }
                if nicknames != "" {
                    nicknameList := strings.Split(nicknames, ",")
                    for _, nickname := range nicknameList {
                        password += strings.TrimSpace(nickname)
                    }
                }

                // Generate additional words using GPT-2
                words := strings.Split(password, "")
                for len(words) < 20 {
                    input := strings.Join(words, " ")
                    output, err := model.Generate(input, 1)
                    if err != nil {
                        log.Fatal(err)
                    }
                    words = append(words, strings.Split(output[0], " ")...)
                }

                // Use closestmatch to find the closest words to the input words
                cm := closestmatch.New(words, []int{2})
                words = cm.Closest(strings.Join(words, " "))

                password = strings.Join(words, "")

                passwords[j] = password
            }
        }(i * passwordsPerRoutine)
    }

    wg.Wait()

    // Write passwords to file
    file, err := os.Create("passwords.txt")
    if err != nil {
        fmt.Println("Error creating file:", err)
        return
    }
    defer file.Close()

    for _, password := range passwords {
        _, err := file.WriteString(password + "\n")
        if err != nil {
            fmt.Println("Error writing to file:", err)
            return
        }
    }

    fmt.Println("Passwords generated and saved to passwords.txt")
}
