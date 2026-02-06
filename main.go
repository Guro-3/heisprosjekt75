package main

import (
    "Driver-go/elevio"
    "fmt"
)

func main(){
    // 1. initialisere server
     //2. lage kanaler som go rutinene kan bruke
    go PollButtons()
    go PollFloorSensor()
    
    //en dør go funksjon, som starter timer sedenr door timou tilbake 
        // den må få et start dør event, starte ny timer med ønsket dyration
        // tid ute sender ut på kanal timout

    // 4. kjører gotonearesFloor elller noe i den duren
    // 5. staarte er for select loop:
    // hvis det kommer noe på kanal button: kall buttonPressedServiceOrder
    //hvis det kommer noe på kanal floor reached: kall serviceOrderAtFloor
    // hvis det kommer noe på kanal start timer blir start dør timer kalt, og det kommer noe på stopptimer blir doortimout kaldt
    
    // til senere.....
    //trenger vel melding inn melding ut kanal?
    // og en assignd order, stateheartbeat kanal?


    // go routine for recive and send

    
    
}

