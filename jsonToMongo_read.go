package main

import (
    "fmt"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "encoding/json"
    "log"
    "io/ioutil"
)

type Patient struct {
    Identifiers []Identifier    `json:"identifier"`
    Names []Name     `json:"name"`
    Telecoms []Telecom   `json:"telecom"`
}

type Identifier struct {
    System string   `json:"system"`
    Value string    `json:"value"`
}

type Name struct {
    Family []string     `json:"family"`
    Given []string      `json:"given"`
}

type Telecom struct {
    System string    `json:"system"`
    Value string    `json:"value"`
    Use string  `json:"use"`
}

func getPatient() Patient {
    raw, err := ioutil.ReadFile("data.json")
    if err != nil {
        panic(err)
    }

    patient := Patient{}
    json.Unmarshal(raw, &patient)

    return patient

}

func insertDoc(p Patient) {
    session, err := mgo.Dial("localhost")
    if err != nil {
        panic(err)
        }
    defer session.Close()

    session.SetMode(mgo.Monotonic, true)

    c:= session.DB("resource").C("patient")
    err = c.Insert(&p)

    if err != nil {
        log.Fatal(err)
    }
}

func main() {
    goPatient := getPatient()

    fmt.Println("converted json into go struct:", goPatient)
    insertDoc(goPatient)

    // Checking by querying
    session, err := mgo.Dial("localhost")
    if err != nil {
        panic(err)
        }
    defer session.Close()

    c:= session.DB("resource").C("patient")

    result := Patient{}
    err = c.Find(bson.M{}).One(&result)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(result)

}
