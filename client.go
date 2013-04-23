package ardrone

import (
	"fmt"
	"github.com/felixge/ardrone/commands"
	"github.com/felixge/ardrone/navdata"
	"log"
	"net"
	"sync"
	"time"
)

type State struct {
	Pitch     float64  // -1 = max back, 1 = max forward
	Roll      float64  // -1 = max left, 1 = max right
	Yaw       float64  // -1 = max counter clockwise, 1 = max clockwise
	Vertical  float64  // -1 = max down, 1 = max up
	Land      bool     // Must be true for landing
	Emergency bool     // Used to disable / trigger emergency mode
	Config    []KeyVal // Config values to send
}

type AnimationId int

const (
	PHI_M30_DEG AnimationId = iota
	PHI_30_DEG
	THETA_M30_DEG
	THETA_30_DEG
	THETA_20_DEG_YAW_200_DEG
	THETA_20_DEG_YAWM_200_DEG
	TURNAROUND
	TURNAROUND_GODOWN
	YAW_SHAKE
	YAW_DANCE
	PHI_DANCE
	THETA_DANCE
	VZ_DANCE
	WAVE
	PHI_THETA_MIXED
	DOUBLE_PHI_THETA_MIXED
	FLIP_AHEAD
	FLIP_BEHIND
	FLIP_LEFT
	FLIP_RIGHT
)

type KeyVal struct {
	Key   string
	Value string
}

type Client struct {
	Config      Config
	navdataConn *navdata.Conn
	commands    *commands.Sequence
	controlConn net.Conn

	stateLock sync.Mutex
	state     State

	Navdata chan *navdata.Navdata // @TODO: make read-only
}

type Config struct {
	Ip             string
	NavdataPort    int
	AtPort         int
	NavdataTimeout time.Duration
}

func DefaultConfig() Config {
	return Config{
		Ip:             "192.168.1.1",
		NavdataPort:    5554,
		AtPort:         5556,
		NavdataTimeout: 2000 * time.Millisecond,
	}
}

func Connect(config Config) (*Client, error){
	client := &Client{Config: config}
	return client, client.Connect()
}

func (client *Client) Connect() error {
	navdataAddr := addr(client.Config.Ip, client.Config.NavdataPort)
	navdataConn, err := navdata.Dial(navdataAddr)
	if err != nil {
		return err
	}

	client.navdataConn = navdataConn
	client.navdataConn.SetReadTimeout(client.Config.NavdataTimeout)

	controlAddr := addr(client.Config.Ip, client.Config.AtPort)
	controlConn, err := net.Dial("udp", controlAddr)

	if err != nil {
		return err
	}

	client.controlConn = controlConn
	client.commands = &commands.Sequence{}

	client.Navdata = make(chan *navdata.Navdata, 0)

	go client.sendLoop()
	go client.navdataLoop()

	// disable emergency mode (if on) and request demo navdata from drone.
	for {
		data := <-client.Navdata

		state := State{Land: true}
		// Sets emergency state if we are in an emergency (which disables it)
		state.Emergency = data.Header.State&navdata.STATE_EMERGENCY_LANDING != 0

		// Request demo navdata if we are not receiving it yet
		if data.Demo == nil {
			state.Config = []KeyVal{{Key: "general:navdata_demo", Value: "TRUE"}}
		} else {
			state.Config = []KeyVal{}
		}

		client.Apply(state)

		// Once emergency is disabled and full navdata is being sent, we are done
		if !state.Emergency && data.Demo != nil {
			break
		}
	}

	return nil
}

func (client *Client) Animate(id AnimationId, arg int) {
	val := fmt.Sprintf("%d,%d", id, arg)
	config := KeyVal{Key: "control:flight_anim", Value: val}
	client.ApplyFor(300 * time.Millisecond, State{Config: []KeyVal{config}})
}

// @TODO Implement error return value
func (client *Client) Takeoff() {
	for {
		// State's zero value will make the drone takeoff/hover
		client.Apply(State{})
		select {
		case data := <-client.Navdata:
			if data.Demo.ControlState == navdata.CONTROL_HOVERING {
				return
			}
		}
	}
}

// @TODO Implement error return value
func (client *Client) Land() {
	for {
		client.Apply(State{Land: true})
		select {
		case data := <-client.Navdata:
			if data.Demo.ControlState == navdata.CONTROL_LANDED {
				return
			}
		}
	}
}

// Apply sets the desired state of the drone. Internally this is turned into
// one or more commands that the sendLoop transmits to the drone every 30ms.
func (client *Client) Apply(state State) {
	client.stateLock.Lock()
	defer client.stateLock.Unlock()
	client.state = state
}

// Applies a given state for a certain duration, and resets the state to its
// zero value (hover) before returning.
func (client *Client) ApplyFor(duration time.Duration, state State) {
	client.Apply(state)
	time.Sleep(duration)
	// Set zero state (causes drone to hover)
	client.Apply(State{})
}

func (client *Client) sendLoop() {
	for {
		client.stateLock.Lock()
		client.commands.Add(&commands.Ref{
			Fly:       !client.state.Land,
			Emergency: client.state.Emergency,
		})
		client.commands.Add(&commands.Pcmd{
			Pitch:    client.state.Pitch,
			Roll:     client.state.Roll,
			Yaw:      client.state.Yaw,
			Vertical: client.state.Vertical,
		})
		for _, config := range client.state.Config {
			client.commands.Add(&commands.Config{Key: config.Key, Value: config.Value})
		}

		client.stateLock.Unlock()

		// @TODO: This interface for creating the commands / tracking the sequence
		// numbers is BS - need to figure out a way to make it better.
		message := client.commands.ReadMessage()
		// @TODO: Handle Write() errors
		client.controlConn.Write([]byte(message))
		time.Sleep(30 * time.Millisecond)
	}
}

func (client *Client) navdataLoop() {
	for {
		navdata, err := client.navdataConn.ReadNavdata()
		// @TODO figure out a better way to handle this, maybe an error channel
		if err != nil {
			log.Printf("error: %s\n", err)
			continue;
		}

		// non-blocking sent into Navdata channel
		select {
		case client.Navdata <- navdata:
		default:
		}
	}
}

func (client *Client) Send() {
	message := client.commands.ReadMessage()
	//fmt.Printf("message: %#v\n", message)
	client.controlConn.Write([]byte(message))
}

func addr(ip string, port int) string {
	return fmt.Sprintf("%s:%d", ip, port)
}
