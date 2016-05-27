package main

import (
    //"encoding/json"
	//"math/big"

	"github.com/iq4health/kani"
    "fmt"

    "log"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "reflect"
)

type Person struct {
        Name string
        Phone string
}

func initi() {
	// ObjectId or $oid
	kani.RegisterType("ObjectId",
		func(val interface{}) kani.Node {
			if id, ok := val.(bson.ObjectId); ok {
				node := kani.NewObj()
				node.Set("$oid", kani.NewStr(id.Hex()))
				return node
			}
			return kani.NewNull()
		})
}

func main(){
    initi()
    session, err := mgo.Dial("localhost")           // open an connection -> Dial function
    if err != nil {                                 //  if you have a
            panic(err)
    }
    defer session.Close()                           // session must close at the end

    session.SetMode(mgo.Monotonic, true)            // Optional. Switch the session to a monotonic behavior.

    c := session.DB("moon").C("people")
    // err = c.Insert(&Person{Name: "Aless", Phone: "+11 11 1111 1111"})
    // if err != nil {
    //         log.Fatal(err)
    // }

    var result map[string]interface{}
    err = c.Find(bson.M{}).One(&result)
    if err != nil {
            log.Fatal(err)
    }

    fmt.Println(result)

    for _, m := range result {
        fmt.Println(reflect.TypeOf(m))
    }

    n2, f := kani.Fy(result)
    fmt.Println(f)
    fmt.Println(n2.Defy())

}
