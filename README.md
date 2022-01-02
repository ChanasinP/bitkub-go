# Bitkub API Golang Client

## 📝 Table of Contents

- [Installing](#installing)
- [Usage](#usage)
- [Authors](#authors)

## 🏁 Install <a name = "installing"></a>

```
go get github.com/ChanasinP/bitkub-go
```

## 🎈 Usage <a name="usage"></a>

```
package main

import (
	"log"

	"github.com/ChanasinP/bitkub-go"
)

func main() {
	api := bitkub.NewBitkub("API_KEY", "API_SECRET")

	status, err := api.GetServerStatus()
	if err != nil {
		log.Panic(err)
	}

	log.Println("Server status")
	for _, s := range status {
		log.Printf("Name: %+s, Status: %s, Message: %s", s.Name, s.Status, s.Message)
	}
}

```

## ✍️ Authors <a name = "authors"></a>

- [@ChanasinP](https://github.com/ChanasinP) - Idea & Initial work
