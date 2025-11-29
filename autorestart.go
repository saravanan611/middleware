package apigate

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

/* restart the program on every day mention time to cleare the in-memory/catch-memory for good pratice and free the resource */
func AutoRestart(lRestartHour, lRestartMinute int) {

	go func() {

		lCurrentTime := time.Now()
		if !(lCurrentTime.Hour() == lRestartHour && lCurrentTime.Minute() == lRestartMinute) {
			lNextExecutionTime := time.Date(lCurrentTime.Year(), lCurrentTime.Month(), lCurrentTime.Day(), lRestartHour, lRestartMinute, 0, 0, lCurrentTime.Location())
			if lNextExecutionTime.Before(lCurrentTime) {
				lNextExecutionTime = lNextExecutionTime.Add(time.Duration(24 * time.Hour))
			}
			fmt.Println("Current Time:", lCurrentTime)
			fmt.Println("Next Execution Time:", lNextExecutionTime, lNextExecutionTime.Sub(lCurrentTime))
			durationUntilNextExecution := lNextExecutionTime.Sub(lCurrentTime)
			time.Sleep(durationUntilNextExecution)
		}

		log.Println("program going to restart within a minute...")
		time.Sleep(1 * time.Minute)

		lErr := restart()

		if lErr != nil {
			log.Println("Error restarting program:", lErr)
		}
		os.Exit(0)
	}()
}

func restart() (lErr error) {
	execPath, lErr := os.Executable()
	if lErr != nil {
		return lErr
	}

	cmd := exec.Command(execPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	lErr = cmd.Start()
	if lErr != nil {
		return lErr
	}
	return
}

/* this function treager in end of the program  like close glogal db connection
use defer to call this in func main() , this func is also auto restart your code where your program will panic
*/

func TreagerOnEnd(pEndFunc ...func()) {

	for lIdx, lFunc := range pEndFunc {
		log.Printf("going to execuate the end process %d \n", lIdx)
		lFunc()
	}

	if lErr := recover(); lErr != nil {
		log.Println("Caught panic:", lErr)

		if lErr := restart(); lErr != nil {
			log.Println("Error restarting program:", lErr)
		}
	}

}
