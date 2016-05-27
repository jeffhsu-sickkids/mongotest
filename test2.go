package main

import (
    "fmt"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "log"
)

type Person struct {
    Id    bson.ObjectId `bson:"_id,omitempty"`      // see http://godoc.org/labix.org/v2/mgo/bson#Marshal
    Name  string        `bson:"name"`               // define field name in the database if you want
    Phone string                                    // the name different from the struc, if not define,
}                                                   // it will use lowercase name.

var (
    IsDrop = true
)

func main() {
    session, err := mgo.Dial("localhost")
    if err != nil {
        panic(err)
    }
    defer session.Close()

    // optional. switch the session to a monotonic behavior.
    session.SetMode(mgo.Monotonic, true)

    // drop database
    if IsDrop {
        err = session.DB("moon").DropDatabase()
        if err != nil {
            panic(err)
        }
    }

    // collection People
    c := session.DB("moon").C("people")
    result := Person{}

    // Index
    index := mgo.Index{
        Key:        []string{"name", "phone"},
        Unique:     true,
        DropDups:   true,
        Background: true,
        Sparse:     true,
    }

    err = c.EnsureIndex(index)
    if err != nil {
        panic(err)
    }

    count1, err := c.Count()
    fmt.Println("<database drop | total peoples> :", count1)

    // create
    err = c.Insert(&Person{Name: "Ale", Phone: "+11 11 1111 1111"},
                    &Person{Name: "Ale", Phone: "+22 22 2222 2222"},
                        &Person{Name: "Cla", Phone: "+33 33 3333 3333"})
    if err != nil {
        log.Fatal(err)
    }

    count2, err := c.Count()
    fmt.Println("<database created | insert 3 peoples | total peoples> :", count2)

    // query one
    err = c.Find(bson.M{"name": "Ale"}).One(&result)        // see http://godoc.org/gopkg.in/mgo.v2#Collection.Find
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("<query one | name: Ale> Phone :", result.Phone)

    // query all [a]
    var results []Person
    err = c.Find(bson.M{"name": "Ale"}).Sort("-name").All(&results)    // see http://godoc.org/gopkg.in/mgo.v2#Query.All
     if err != nil {
        panic(err)
    }

    fmt.Println()
    fmt.Println("<query all [a] | results>")
    fmt.Println(results)

    // query all [b]
    iter := c.Find(nil).Limit(3).Iter()
    fmt.Println()
    fmt.Println("<query all [b] | results>")
    for iter.Next(&result) {
        fmt.Printf("results: %v\n", result.Name + " , " + result.Phone)
    }
    err = iter.Close()
    if err != nil {
        log.Fatal(err)
    }

    // update
    selector := bson.M{"name": "Cla"}
    updator := bson.M{"$set": bson.M{"phone": "+44 44 4444 4444"}}
    err = c.Update(selector, updator)
    if err != nil {
        panic(err)
    }

    err = c.Find(bson.M{"name": "Cla"}).One(&result)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println()
    fmt.Println("<updated | name : Cla> Phone :", result.Phone)

    // remove
    err = c.Remove(bson.M{"_id": result.Id})
    if err != nil {
        log.Fatal(err)
    }

    count3, err := c.Count()
    fmt.Println("<total peoples after delete> :", count3)


}
