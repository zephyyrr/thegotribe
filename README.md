# TheGoTribe

TheGoTribe is go bindings for TheEyeTribe eye trackers.


## Usage
 ````go
 import (
 	"log"
 	tet "github.com/zephyyrr/thegotribe"

 func main() {
 	et, err := tet.Create()
 	if err != nil{
 		log.Fatal(err)
 	}
 	for frame := range et.Frames {
 		log.Println(frame.Average)
 	}
 }
 ````