package elevator




func buttonPressedServiceOrder()

func serviceOrderAtFloor()

type elevatorState int;

const(
	idle elevatorState = 0;
	mooving = 1;
	doorOpen = 2;
	error_ = 3;
)

