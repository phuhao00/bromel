package robot

import "sync"

type Mgr struct {
	sync.RWMutex
	robots       []*Robot
	onlineRobots map[uint64]*Robot
}

var robotManager *Mgr

func InitRobotManager(num int32) *Mgr {
	robotManager = &Mgr{
		robots:       make([]*Robot, 0, num),
		onlineRobots: make(map[uint64]*Robot),
	}
	return robotManager
}

func (mgr *Mgr) addRobot(server string) *Robot {
	robot := newRobot(server)
	mgr.robots = append(mgr.robots, robot)
	return robot
}

func (mgr *Mgr) addOnlineRobot(robot *Robot) {
	mgr.Lock()
	mgr.onlineRobots[robot.robotID] = robot
	mgr.Unlock()
}
