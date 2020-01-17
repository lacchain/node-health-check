package main

import (
	"fmt"
	"time"

	"github.com/xin053/hsperfdata"
)

var previousOldSpacePercentage float64 = 0
var previousNFGC int64 = 0
var i = 0

func executeReadJavaProcess() <-chan bool {
	fmt.Println("Starting Verification")
	c := make(chan bool)
	go func() {
		restart := false
		for i := 0; ; i++ {
			time.Sleep(delayTimeMinutes * time.Minute)
			_, restartByProcessHealthy := analyzeJavaProcess() //if true ==> restart
			restart = restartByProcessHealthy
			c <- restart
		}
	}()
	return c
}

func analyzeJavaProcess() (interface{}, bool) {
	err, currentNFGC, edenPercentage, s0Percentage, s1Percentage, relativeOldSpacePercentage, absoluteOldSpacePercentage := getData()
	if err != nil {
		return err, false
	}

	hasNFGCIncremented := (currentNFGC > previousNFGC)
	hasOldSpaceIncremented := (int(relativeOldSpacePercentage) >= int(previousOldSpacePercentage))
	oldSpaceCondition := hasOldSpaceIncremented && hasNFGCIncremented && (relativeOldSpacePercentage > 95.0) && (absoluteOldSpacePercentage > 90)
	youngSpaceCondition := (edenPercentage < 2.0) && (s0Percentage < 2.0) && (s1Percentage < 2.0)
	restart := (oldSpaceCondition && youngSpaceCondition)
	previousOldSpacePercentage = relativeOldSpacePercentage
	previousNFGC = currentNFGC
	report()
	return nil, restart
}

func report() interface{} {
	javaProcessPath, err := getDataPath(processName)
	if err != nil {
		return err
	}

	data, err := readData(javaProcessPath)
	if err != nil {
		return err
	}

	edenUsed, edenSpaceCapacity, edenMaxSpaceCapacity := getEdenCapacities(data)
	edenPercentage := divide(float64(edenUsed), float64(edenSpaceCapacity)) * 100
	fmt.Println("Eden statistics\n",
		"eden used:", edenUsed, "\n",
		"eden capacity:", edenSpaceCapacity, "\n",
		"eden MaxSpaceCapacity", edenMaxSpaceCapacity, "\n",
		"eden used/edenCapacity (%)", fmt.Sprintf("%2f", edenPercentage))

	s0Used, s0SpaceCapacity, s0MaxSpaceCapacity := getS0Capacities(data)
	s0Percentage := divide(float64(s0Used), float64(s0SpaceCapacity)) * 100
	fmt.Println("s0 statistics\n",
		"s0 used:", s0Used, "\n",
		"s0 capacity:", s0SpaceCapacity, "\n",
		"s0 MaxSpaceCapacity", s0MaxSpaceCapacity, "\n",
		"s0 used/s0Capacity (%)", fmt.Sprintf("%2f", s0Percentage))

	s1Used, s1SpaceCapacity, s1MaxSpaceCapacity := getS1Capacities(data)
	s1Percentage := divide(float64(s1Used), float64(s1SpaceCapacity)) * 100
	fmt.Println("s1 statistics\n",
		"s1 used:", s1Used, "\n",
		"s1 capacity:", s1SpaceCapacity, "\n",
		"s1 MaxSpaceCapacity", s1MaxSpaceCapacity, "\n",
		"s1 used/s1Capacity (%)", fmt.Sprintf("%2f", s1Percentage))

	youngSpaceCapacity, youngSpaceMaxCapacity := getYoungCapacities(data)
	youngUsed := edenUsed + s0Used + s1Used
	fmt.Println("young statistics\n",
		"young used:", youngUsed, "\n",
		"young capacity:", youngSpaceCapacity, "\n",
		"young MaxSpaceCapacity", youngSpaceMaxCapacity)

	oldUsed, oldSpaceCapacity, oldSpaceMaxCapacity := getOldSpaceCapacities(data)
	relativeOldSpacePercentage := divide(float64(oldUsed), float64(oldSpaceCapacity)) * 100
	absoluteOldSpacePercentage := divide(float64(oldUsed), float64(oldSpaceMaxCapacity)) * 100
	fmt.Println("old statistics\n",
		"old used", oldUsed, "\n",
		"old capacity", oldSpaceCapacity, "\n",
		"old MaxSpaceCapacity", oldSpaceMaxCapacity, "\n",
		"oldUsed/oldCapacity(%)", fmt.Sprintf("%2f", relativeOldSpacePercentage)+"%", "\n",
		"oldUsed/maxOldCapacity(%)", absoluteOldSpacePercentage,
	)

	currentNFGC := getFGC(data)
	fmt.Println("Number of Full Garbage Collector:", currentNFGC)
	fmt.Println("*************************************************************************")
	return nil
}

func getData() (interface{}, int64, float64, float64, float64, float64, float64) {
	javaProcessPath, err := getDataPath(processName)
	if err != nil {
		return err, 0, 0, 0, 0, 0, 0
	}

	data, err := readData(javaProcessPath)
	if err != nil {
		return err, 0, 0, 0, 0, 0, 0
	}

	edenUsed, edenSpaceCapacity, _ := getEdenCapacities(data)
	edenPercentage := divide(float64(edenUsed), float64(edenSpaceCapacity)) * 100

	s0Used, s0SpaceCapacity, _ := getS0Capacities(data)
	s0Percentage := divide(float64(s0Used), float64(s0SpaceCapacity)) * 100

	s1Used, s1SpaceCapacity, _ := getS1Capacities(data)
	s1Percentage := divide(float64(s1Used), float64(s1SpaceCapacity)) * 100

	oldUsed, oldSpaceCapacity, oldSpaceMaxCapacity := getOldSpaceCapacities(data)

	relativeOldSpacePercentage := divide(float64(oldUsed), float64(oldSpaceCapacity)) * 100
	absoluteOldSpacePercentage := divide(float64(oldUsed), float64(oldSpaceMaxCapacity)) * 100

	currentNFGC := getFGC(data)

	return nil, currentNFGC, edenPercentage, s0Percentage, s1Percentage, relativeOldSpacePercentage, absoluteOldSpacePercentage
}

func getDataPath(processName string) (string, interface{}) {
	dataPathsByProcessName, err := hsperfdata.DataPathsByProcessName(processName)
	if err != nil {
		fmt.Println("GetDataPathByProcessName() error =", err)
		return "", err
	}

	for _, value := range dataPathsByProcessName {
		return value, nil
	}

	return "", nil
}

func readData(dataPath string) (map[string]interface{}, interface{}) {
	data, err := hsperfdata.ReadPerfData(dataPath, true)
	if err != nil {
		fmt.Println("readData() error =", err)
		return map[string]interface{}{}, err
	}

	return data, nil
}

func getOldSpaceCapacities(data map[string]interface{}) (int64, int64, int64) {
	oldUsed := data["sun.gc.generation.1.space.0.used"].(int64)
	oldSpaceCapacity := data["sun.gc.generation.1.space.0.capacity"].(int64)
	oldSpaceMaxCapacity := data["sun.gc.generation.1.space.0.maxCapacity"].(int64)
	return oldUsed, oldSpaceCapacity, oldSpaceMaxCapacity
}

func printAllData(data map[string]interface{}) {
	for key, value := range data {
		fmt.Println(key, " ====> ", value)
	}
}

func getFGC(data map[string]interface{}) int64 {
	return data["sun.gc.collector.1.invocations"].(int64)
}

func getYoungCapacities(data map[string]interface{}) (int64, int64) {
	youngSpaceMaxCapacity := data["sun.gc.generation.0.maxCapacity"].(int64)
	youngSpaceCapacity := data["sun.gc.generation.0.capacity"].(int64)
	return youngSpaceCapacity, youngSpaceMaxCapacity
}

func getEdenCapacities(data map[string]interface{}) (int64, int64, int64) {
	edenUsed := data["sun.gc.generation.0.space.0.used"].(int64)
	edenSpaceCapacity := data["sun.gc.generation.0.space.0.capacity"].(int64)
	edenMaxSpaceCapacity := data["sun.gc.generation.0.space.0.maxCapacity"].(int64)
	return edenUsed, edenSpaceCapacity, edenMaxSpaceCapacity
}

func getS0Capacities(data map[string]interface{}) (int64, int64, int64) {
	s0Used := data["sun.gc.generation.0.space.1.used"].(int64)
	s0SpaceCapacity := data["sun.gc.generation.0.space.1.capacity"].(int64)
	s0MaxSpaceCapacity := data["sun.gc.generation.0.space.1.maxCapacity"].(int64)
	return s0Used, s0SpaceCapacity, s0MaxSpaceCapacity
}

func getS1Capacities(data map[string]interface{}) (int64, int64, int64) {
	s1Used := data["sun.gc.generation.0.space.2.used"].(int64)
	s1SpaceCapacity := data["sun.gc.generation.0.space.2.capacity"].(int64)
	s1MaxSpaceCapacity := data["sun.gc.generation.0.space.2.maxCapacity"].(int64)
	return s1Used, s1SpaceCapacity, s1MaxSpaceCapacity
}
