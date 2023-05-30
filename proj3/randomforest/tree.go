package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"proj3/concurrent"
	"sort"
	"strconv"
	"time"
)

const usage = "Usage: run_mode tree_num tree_depth threads threshold thresholdBalance\n" +
	"run_mode = (s) - serial, (bal) - WorkBalancing, (stl) - WorkStealing \n" +
	"tree_num = number of trees in the forest \n" +
	"tree_depth     = Max depth of the tree\n" +
	"threads = Runs the parallel version of the program with the specified number of threads.\n" +
	"threshold = The number of items that a goroutine in the pool can grab from the executor in one time period\n" +
	"thresholdBalance = The threshold used to know when to perform balancing\n"

// Node to store the attributes related to a decision Tree
type DNode struct {

	// If the node is a leaf node
	nodeType string

	// The class predicted by the node
	predictedClass float64
	X              [][]float64

	// Which attribute is the node testing on
	testAttribute int
	testValue     float64
	children      []DNode
}

// The tree class to implement decision Tree
type Tree struct {
	maxDepth int
	root     DNode
}

// Struct to help in ArgSort
type Slice struct {
	sort.Float64Slice
	idx []int
}

// Func to help in ArgSort
func (s Slice) Swap(i, j int) {
	s.Float64Slice.Swap(i, j)
	s.idx[i], s.idx[j] = s.idx[j], s.idx[i]
}

// Method of class Tree to fit the decision tree
func (T *Tree) fit(X [][]float64, y []float64) {
	T.root = T.recursiveBuildTree(X, y, 0)
}

// Method of class Tree to build the decision tree
func (T *Tree) recursiveBuildTree(X [][]float64, y []float64, currDepth int) DNode {
	node := DNode{}
	node.X = X

	// Finding frequency of each element
	m := findFreq(y, 1.0)

	rows := len(X)
	if rows == 0 {
		// Sending a random number in case we want to predict the parent class
		node.predictedClass = -0.123
		return node
	}
	cols := len(X[0])

	ctr := true

	// Checking if all Xs are the same
	for i := 0; i < rows; i++ {
		for j := 1; j < cols; j++ {
			if X[i][0] != X[i][j] {
				ctr = false
				break
			}
		}
		if ctr == false {
			break
		}
	}

	// To check if the current node should be a leaf node
	if currDepth == T.maxDepth || len(m) <= 1 || ctr {
		node.nodeType = "leaf"
		if len(y) != 0 {
			node.predictedClass = T.classPredict(y)
		} else {
			node.predictedClass = -0.123
		}
	} else {

		// Deciding the attribute and on which point to divide the attribute
		A, val := T.importance(X, y)
		node.testAttribute, node.testValue = A, val
		node = T.nodeChildren(X, y, A, val, currDepth, node)
	}
	return node
}

// Function to build the children of the dtree
func (T *Tree) nodeChildren(X [][]float64, y []float64, A int, val float64, currDepth int, node DNode) DNode {
	var newX1 [][]float64
	var newY1 []float64
	var newX2 [][]float64
	var newY2 []float64

	n := len(y)

	// Splitting the data basis the attribute value
	for i := 0; i < n; i++ {
		if X[i][A] < val {
			newX1 = append(newX1, X[i][:])
			newY1 = append(newY1, y[i])
		} else {
			newX2 = append(newX2, X[i][:])
			newY2 = append(newY2, y[i])
		}
	}

	nodeLeft := T.recursiveBuildTree(newX1, newY1, currDepth+1)
	nodeRight := T.recursiveBuildTree(newX2, newY2, currDepth+1)

	// Meaning that there was not enough data for the left/right child
	if nodeLeft.predictedClass == -0.123 && nodeLeft.nodeType == "leaf" {
		nodeLeft.predictedClass = T.classPredict(y)
	} else if nodeRight.predictedClass == -0.123 && nodeRight.nodeType == "leaf" {
		nodeRight.predictedClass = T.classPredict(y)
	}
	node.children = append(node.children, nodeLeft)
	node.children = append(node.children, nodeRight)

	return node

}

// Function to predict the class of the given dataset
func (T *Tree) classPredict(y []float64) float64 {
	return mostFrequent(y)
}

// Function to calculate the importance and return the best attribute and the split value
func (T *Tree) importance(X [][]float64, y []float64) (int, float64) {
	minEntropy := math.Inf(2)
	minI := 1
	minVal := math.Inf(2)

	n := len(X[0])

	for i := 0; i < n; i++ {

		thisEntropy, val := T.importanceCont(ColSliceSingle(X, i), y)
		if thisEntropy < minEntropy {
			minEntropy = thisEntropy
			minI = i
			minVal = val
		}
	}
	return minI, minVal
}

// Helper function to importance
func (T *Tree) importanceCont(X []float64, y []float64) (float64, float64) {

	// Creating newX and newY, such that the elements in the newX are sorted
	newX := Slice{
		Float64Slice: sort.Float64Slice(X),
		idx:          make([]int, len(X)),
	}
	for i := range newX.idx {
		newX.idx[i] = i
	}

	sort.Sort(newX)
	argsort := newX.idx

	var newY []float64
	n := len(y)
	for i := 0; i < n; i++ {
		newY = append(newY, y[argsort[i]])
	}

	minEntropy := math.Inf(2)
	val := float64(n - 1)

	for i := 0; i < n-1; i++ {
		if X[i] == X[i+1] {
			continue
		}

		// Returning the split with least entropy
		thisEntropy := float64(i+1)/float64(n)*T.entropy(newY[:i+1]) + float64(n-i-1)/float64(n)*T.entropy(newY[i+1:])
		if thisEntropy < minEntropy {
			minEntropy = thisEntropy
			val = (X[i] + X[i+1]) / 2
		}
	}
	return minEntropy, val
}

// Function to predict the y Value for any unseen data
func (T *Tree) predict(X [][]float64) []float64 {

	var yPred []float64

	for i := 0; i < len(X); i++ {
		thisNode := T.root
		for thisNode.nodeType != "leaf" {
			if X[i][thisNode.testAttribute] < thisNode.testValue {
				thisNode = thisNode.children[0]
			} else {
				thisNode = thisNode.children[1]
			}
		}
		yPred = append(yPred, thisNode.predictedClass)
	}

	return yPred
}

// Function to calculate the Entropy. The node is split wherever the entropy is least
func (T *Tree) entropy(y []float64) float64 {

	n := float64(len(y))
	m := findFreq(y, n)
	sum := 0.0
	for _, value := range m {
		sum -= math.Log2(value) * value
	}
	return sum
}

// Function to create a map with frequency of each element in the list
func findFreq(y []float64, div float64) map[float64]float64 {
	m := make(map[float64]float64)
	for _, i := range y {
		_, ok := m[i]
		if ok {
			m[i] += 1.0 / div
		} else {
			m[i] = 1.0 / div
		}
	}
	return m
}

// Function to return the element occuring most frequently in a given array
func mostFrequent(arr []float64) float64 {
	m := map[float64]int{}
	var maxCnt int
	var freq float64
	for _, a := range arr {
		m[a]++
		if m[a] > maxCnt {
			maxCnt = m[a]
			freq = a
		}
	}
	return freq
}

// Creating a struct to hold all variables to be passed for parallelising the algorithm
type IntervalTask struct {
	XTrain   [][]float64
	XTest    [][]float64
	yTrain   []float64
	cols     int
	sqrtCols int
	rows     int
	i        int
}

// The function to perform computation for each thread
func calculateIntervals(XTrain [][]float64, XTest [][]float64, yTrain []float64, cols int, sqrtCols int, rows int, i int) []float64 {
	thisCols := rand.Perm(cols - 1)[:sqrtCols]
	var XTrainTemp [][]float64
	var XTestTemp [][]float64

	XTrainTemp = ColSlice2(XTrain, thisCols)

	XTestTemp = ColSlice2(XTest, thisCols)

	tree := Tree{maxDepth: i}
	tree.fit(XTrainTemp, yTrain)

	return tree.predict(XTestTemp)
}

// Creating a callable for our Executor
func NewIntervalTask(XTrain [][]float64, XTest [][]float64, yTrain []float64, cols int, sqrtCols int, rows int, i int) concurrent.Callable {
	return &IntervalTask{XTrain, XTest, yTrain, cols, sqrtCols, rows, i}
}

// Defining the Call function for the Executor
func (task *IntervalTask) Call() interface{} {

	yPred := calculateIntervals(task.XTrain, task.XTest, task.yTrain, task.cols, task.sqrtCols, task.rows, task.i)

	return yPred

}

// Creating a function to find a slice of columns
func ColSlice(Arr [][]float64, from int, to int) [][]float64 {
	var ArrRet [][]float64
	rows := len(Arr)
	for i := 0; i < rows; i++ {
		var tempData []float64
		for j := from; j < to; j++ {
			tempData = append(tempData, Arr[i][j])
		}
		ArrRet = append(ArrRet, tempData)
	}
	return ArrRet
}

// Creating a function to find a slice of columns
func ColSlice2(Arr [][]float64, lst []int) [][]float64 {
	var ArrRet [][]float64
	rows := len(Arr)
	for i := 0; i < rows; i++ {
		var tempData []float64
		for _, j := range lst {
			tempData = append(tempData, Arr[i][j])
		}
		ArrRet = append(ArrRet, tempData)
	}
	return ArrRet
}

// Creating a function to find a slice of columns
func ColSliceSingle(Arr [][]float64, col int) []float64 {
	var ArrRet []float64
	rows := len(Arr)
	for i := 0; i < rows; i++ {
		ArrRet = append(ArrRet, Arr[i][col])
	}
	return ArrRet
}

// Function to print accuracy of the Random Forest
func Accuracy (yPred [][]float64, yTest []float64) {

	var yPred2 []float64
	for j := 0; j < len(yPred[0]); j++ {
		yPred2 = append(yPred2, mostFrequent(ColSliceSingle(yPred, j)))
	}

	acc := 0
	for j := 0; j < len(yTest); j++ {
		if yPred2[j] == yTest[j] {
			acc++
		}
	}
	fmt.Printf("Accuracy: ")
	fmt.Println(float64(acc) / float64(len(yTest)))
}

// Function to split the data into train and test
func TrainTestSplit (rows int, cols int, data2 [][]float64) ([][]float64, [][]float64, []float64, []float64) {

	X := ColSlice(data2, 0, cols-1)
	y := ColSliceSingle(data2, cols-1)

	randList := rand.Perm(rows)
	var XTrain [][]float64
	var yTrain []float64
	var XTest [][]float64
	var yTest []float64

	for i := 0; i < rows*2/3; i++ {
		XTrain = append(XTrain, X[randList[i]][:])
		yTrain = append(yTrain, y[randList[i]])
	}

	for i := rows * 2 / 3; i < rows; i++ {
		XTest = append(XTest, X[randList[i]][:])
		yTest = append(yTest, y[randList[i]])
	}
	return XTrain, XTest, yTrain, yTest
}

// Function to read the data and preprocess it
func ReadPreProcess(str string) ([][]float64, int, int) {
	f, err := os.Open(str)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()

	if err != nil {
		log.Fatal(err)
	}

	var data2 [][]float64
	rows := len(data)
	cols := len(data[0])

	// Converting values to float and putting 0 wherever the data is missing
	for i := 0; i < rows; i++ {
		var tempData []float64
		for j := 0; j < cols; j++ {
			thisFloat, error1 := strconv.ParseFloat(data[i][j], 64)
			if error1 != nil {
				thisFloat = 0
			}
			tempData = append(tempData, thisFloat)
		}
		data2 = append(data2, tempData)
	}
	return data2, rows, cols
}

func main() {
	
	if len(os.Args) < 4 {
		fmt.Println(usage)
		return
	}

	implementationType := os.Args[1]
	trees, _ := strconv.Atoi(os.Args[2])
	i, _ := strconv.Atoi(os.Args[3])

	// Reading, preprocessing and splitting the data into train and test
	data, rows, cols := ReadPreProcess("./randomforest/arrhythmia.csv")
	XTrain, XTest, yTrain, yTest := TrainTestSplit(rows, cols, data)

	// To find the number of features to pass to each random forest
	sqrtCols := int(math.Round(math.Sqrt(float64(cols))))

	var yPred [][]float64

	strt := time.Now()

	if implementationType == "s" {
		for k := 0; k < trees; k++ {
			yPred = append(yPred, calculateIntervals(XTrain, XTest, yTrain, cols, sqrtCols, rows, i))
		}

	} else {

		threadCount, _ := strconv.Atoi(os.Args[4])
		threshold, _ := strconv.Atoi(os.Args[5])
		executor := concurrent.NewWorkStealingExecutor(threadCount, threshold)
		if implementationType == "bal" {
			thresholdBalance, _ := strconv.Atoi(os.Args[6])
			executor = concurrent.NewWorkBalancingExecutor(threadCount, threshold, thresholdBalance)
		}
		var futures []concurrent.Future
		for k := 0; k < trees; k++ {
			futures = append(futures, executor.Submit(NewIntervalTask(XTrain, XTest, yTrain, cols, sqrtCols, rows, i)))
		}

		for _, future := range futures {
			yPred = append(yPred, future.Get().([]float64))
		}
		executor.Shutdown()

	}

	// To check if the code is running properly. Commenting it because we only want the time in output
	Accuracy(yPred, yTest)

	end := time.Since(strt).Seconds()
	fmt.Printf("Time Taken: %.2fs\n", end)
}
