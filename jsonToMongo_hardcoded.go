package main

import (
    "fmt"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "encoding/json"
    "log"
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

func main() {
    jsonData := []byte(`
        {
        "identifier": [
        {
            "system": "urn:oid:1.2.36.146.595.217.0.1",
            "value": "12345"
        }
        ],
        "name": [
        {
            "family": [
                "Chalmers"
            ],
            "given": [
                "Peter",
                "James"
            ]
        }
        ],
        "telecom": [
            {
            "system": "phone",
            "value": "(03) 5555 6473",
            "use": "work"
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
err = c.Find(bson.M{}).One(&result)
if err != nil {
    log.Fatal(err)
}

fmt.Println(result)



}
