package main

import (
    "fmt"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "encoding/json"
    "log"
)

type Patient struct {
    Name string     `json:"name"`
    Telecom []PhoneNumber   `json:"telecom"`
}

type PhoneNumber struct {
    Phone string    `json:"phone"`
}

func main() {
    jsonData := []byte(`
        {
            "name": "Jeff Hsu",
            "telecom" : [
             {
                "phone": "+1 647 555 5555"
             }
            ]
        }
        `)
    goPatient := Patient{}
    json.Unmarshal(jsonData, &goPatient)

    fmt.Println("converted json into go struct:", goPatient)

    session, err := mgo.Dial("localhost")
    if err != nil {
        panic(err)
    }
    defer session.Close()

    session.SetMode(mgo.Monotonic, true)

    c:= session.DB("resource").C("patient")
    err = c.Insert(&goPatient)

    if err != nil {
        log.Fatal(err)
    }


    result := Patient{}
    err = c.Find(bson.M{"name": "Jeff Hsu"}).One(&result)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(result)



}
