package elevator

type elevatorState int;

const NUMFloors = 4
const numCabButtons = 1
const numHallButtons = 1




const(
	idle elevatorState = 0;
	mooving = 1;
	doorOpen = 2;
	error_ = 3;
)

type Elevator struct{
	currentFloor int
	lastFloor int
	cabOrderMatrix  [NUMFloors][numCabButtons]bool
	HallorderMatrix  [NUMFloors][numHallButtons]bool
	dir Driver.MotorDirection
	state elevatorState
}

