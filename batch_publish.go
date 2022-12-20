package main

import (
	"fmt"
	"sync"
)

// To be used in upcoming pub/sub project. How do I batch-publish messages to pub/sub topic to minimize network calls?
// Publish to one endpoint/GRPC method. That method can be thought of as the "gateway".
//The protobuf contains a "Repeated" field for 0 or more messages. On the server, type check and  route to appropriate fn.

type Storage struct {
	fakeDB      []*Student
	fakeDBMutex sync.Mutex
}
type Student struct {
	id   int
	name string
}

func store(student *Student, stg *Storage) {
	fmt.Println(student.name)

	defer stg.fakeDBMutex.Unlock()
	stg.fakeDBMutex.Lock()

	stg.fakeDB = append(stg.fakeDB, student)
}

func batchStore(students []*Student, stg *Storage) {
	for _, s := range students {
		fmt.Println(s.name)
		stg.fakeDBMutex.Lock()
		stg.fakeDB = append(stg.fakeDB, s)
		stg.fakeDBMutex.Unlock()
	}

}

// Mimicking a network call for subscriber to push to pub/sub. In the real world, storage wouldn't need to be passed
//Server will handle storage. This batch publishes a protobuf message.
func pushToFakeServer(a interface{}, s *Storage) {
	publishGateway(a, s)
}

func publishGateway(a interface{}, stg *Storage) {
	switch obj := a.(type) {
	case *Student:
		store(obj, stg)
	case []*Student:
		batchStore(obj, stg)
	}
}

func generateFakeData() []*Student {
	v := []string{"A", "B", "C", "D", "AA", "BB", "CC", "DD"}
	var res []*Student
	for i := 0; i < 8; i++ {
		s := &Student{
			id:   i,
			name: v[i],
		}
		res = append(res, s)
	}
	return res
}

func main() {
	//fmt.Println("Hello, world!")
	var studentStorage []*Student
	st := &Storage{
		fakeDB: studentStorage,
	}

	c := &Student{
		id:   1,
		name: "Charles",
	}

	d := &Student{
		id:   2,
		name: "Edwin",
	}

	collection := []*Student{c, d}
	//publish(c)
	//batchPublish(collection)
	data := generateFakeData()
	data = append(data, collection...)

	pushToFakeServer(data, st)
}
