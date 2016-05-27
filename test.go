package main

import (
        //"fmt"
        "log"
        "gopkg.in/mgo.v2"
        "gopkg.in/mgo.v2/bson"
)

type Person struct {
        Name string
        Phone string ",omitempty"
}

func main() {

        session, err := mgo.Dial("localhost")           // open an connection -> Dial function
        if err != nil {                                 //  if you have a
                panic(err)
        }
        defer session.Close()                           // session must close at the end

        session.SetMode(mgo.Monotonic, true)            // Optional. Switch the session to a monotonic behavior.

        c := session.DB("moon").C("people")
        err = c.Insert(&Person{Name: "Ale", Phone: "+11 11 1111 1111"})
        if err != nil {
                log.Fatal(err)
        }

        result := Person{}
        err = c.Find(bson.M{"name": "Ale"}).One(&result)
        if err != nil {
                log.Fatal(err)
        }

        //fmt.Println("Phone:", result.Phone)

        // data, err := bson.Marshal(&Person{Name:"Bob"})
        // if err != nil {
        //         panic(err)
        // }
        // fmt.Printf("%q", data)
}
