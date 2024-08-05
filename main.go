package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// Constants for the problem and algorithm parameters
const (
	maxShiftTime     = 720.0 // 12 hours in minutes
	costPerDriver    = 500.0
	tabuListSize     = 10
	maxIterations    = 100
	initialTabuValue = 1000.0
	neighborhoodSize = 10
)

// Load represents a delivery task with pickup and dropoff locations
type Load struct {
	id      int
	pickup  [2]float64
	dropoff [2]float64
}

// Solution represents a set of routes and their associated cost
type Solution struct {
	routes [][]int
	cost   float64
}

// Global variables to store problem data and precomputed distances
var (
	loads            []Load
	distanceMatrix   [][]float64
	deliveryDistance []float64
)

func main() {
	// Check if a data file path is provided
	if len(os.Args) < 2 {
		fmt.Println("Please provide a data file path.")
		return
	}

	dataFile := os.Args[1]
	// Read loads from the provided file
	if err := readLoads(dataFile); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	// Initialize distance matrices
	initializeMatrices()
	// Run the tabu search algorithm
	bestSolution := tabuSearch()
	// Print the best solution found
	printSolution(bestSolution)
}

// readLoads reads load data from the specified file
func readLoads(filename string) error {
	// Open and read the file
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "loadNumber") {
			continue // Skip header line
		}
		// Parse load data and add to loads slice
		parts := strings.Fields(line)
		id, _ := strconv.Atoi(parts[0])
		pickup := parseCoordinates(parts[1])
		dropoff := parseCoordinates(parts[2])
		loads = append(loads, Load{id, pickup, dropoff})
	}
	return scanner.Err()
}

// parseCoordinates converts a string coordinate to a float64 pair
func parseCoordinates(coord string) [2]float64 {
	coord = strings.Trim(coord, "()")
	parts := strings.Split(coord, ",")
	x, _ := strconv.ParseFloat(parts[0], 64)
	y, _ := strconv.ParseFloat(parts[1], 64)
	return [2]float64{x, y}
}

// initializeMatrices precomputes distance matrices for efficiency
func initializeMatrices() {
	totalLoads := len(loads)
	deliveryDistance = make([]float64, totalLoads)
	distanceMatrix = make([][]float64, totalLoads+1)

	for i := range distanceMatrix {
		distanceMatrix[i] = make([]float64, totalLoads+1)
	}

	// Calculate distances between loads and depot
	for i, load := range loads {
		deliveryDistance[i] = euclideanDistance(load.pickup, load.dropoff)
		distanceMatrix[0][i+1] = euclideanDistance([2]float64{0, 0}, load.pickup)
		distanceMatrix[i+1][0] = euclideanDistance(load.dropoff, [2]float64{0, 0})
		for j, otherLoad := range loads {
			if i != j {
				distanceMatrix[i+1][j+1] = euclideanDistance(load.dropoff, otherLoad.pickup)
			}
		}
	}
}

// euclideanDistance calculates the Euclidean distance between two points
func euclideanDistance(a, b [2]float64) float64 {
	return math.Sqrt(math.Pow(a[0]-b[0], 2) + math.Pow(a[1]-b[1], 2))
}

// tabuSearch implements the Tabu Search algorithm
func tabuSearch() Solution {
	rand.Seed(time.Now().UnixNano())

	// Initialize a random initial solution
	currentSolution := generateInitialSolution()
	bestSolution := currentSolution

	tabuList := make(map[string]float64)
	tabuCounter := make(map[string]int)

	// Main loop of the Tabu Search algorithm
	for iteration := 0; iteration < maxIterations; iteration++ {
		neighbors := generateNeighborhood(currentSolution)
		bestNeighbor := Solution{cost: math.Inf(1)}

		// Find the best non-tabu neighbor
		for _, neighbor := range neighbors {
			if tabuValue, ok := tabuList[neighborKey(neighbor)]; ok && tabuValue > 0 {
				continue
			}
			if neighbor.cost < bestNeighbor.cost {
				bestNeighbor = neighbor
			}
		}

		// Update best solution if necessary
		if bestNeighbor.cost < bestSolution.cost {
			bestSolution = bestNeighbor
		}

		// Update tabu list
		updateTabuList(tabuList, tabuCounter, bestNeighbor)

		currentSolution = bestNeighbor
	}

	return bestSolution
}

// generateInitialSolution creates a random initial solution
func generateInitialSolution() Solution {
	var solution Solution
	remainingLoads := make([]int, len(loads))
	for i := range remainingLoads {
		remainingLoads[i] = i + 1
	}

	// Create routes until all loads are assigned
	for len(remainingLoads) > 0 {
		var route []int
		currentNode := 0
		routeTime := 0.0

		// Build a single route
		for len(remainingLoads) > 0 {
			nextNode := selectNextNode(currentNode, remainingLoads, routeTime)
			if nextNode == 0 {
				break
			}
			route = append(route, nextNode)
			routeTime += distanceMatrix[currentNode][nextNode] + deliveryDistance[nextNode-1]
			currentNode = nextNode
			// Remove the selected load from remainingLoads
			for i, load := range remainingLoads {
				if load == nextNode {
					remainingLoads = append(remainingLoads[:i], remainingLoads[i+1:]...)
					break
				}
			}
		}

		if len(route) > 0 {
			solution.routes = append(solution.routes, route)
		}
	}

	solution.cost = calculateCost(solution)
	return solution
}

// generateNeighborhood creates a set of neighbor solutions
func generateNeighborhood(solution Solution) []Solution {
	var neighbors []Solution

	for i := 0; i < neighborhoodSize; i++ {
		neighbor := swapRandomRoutes(solution)
		neighbor.cost = calculateCost(neighbor)
		neighbors = append(neighbors, neighbor)
	}

	return neighbors
}

// swapRandomRoutes creates a new solution by swapping two random routes
func swapRandomRoutes(solution Solution) Solution {
	// Clone solution and swap routes
	var newSolution Solution
	newSolution.routes = make([][]int, len(solution.routes))
	copy(newSolution.routes, solution.routes)

	if len(newSolution.routes) < 2 {
		return newSolution
	}

	i, j := rand.Intn(len(newSolution.routes)), rand.Intn(len(newSolution.routes))
	for i == j {
		j = rand.Intn(len(newSolution.routes))
	}

	newSolution.routes[i], newSolution.routes[j] = newSolution.routes[j], newSolution.routes[i]

	return newSolution
}

// selectNextNode chooses the next load to add to a route
func selectNextNode(currentNode int, remainingLoads []int, routeTime float64) int {
	var probabilities []float64
	var sum float64

	// Calculate probabilities for each remaining load
	for _, load := range remainingLoads {
		if routeTime+distanceMatrix[currentNode][load]+deliveryDistance[load-1]+distanceMatrix[load][0] > maxShiftTime {
			probabilities = append(probabilities, 0)
		} else {
			probability := 1.0 / distanceMatrix[currentNode][load]
			probabilities = append(probabilities, probability)
			sum += probability
		}
	}

	if sum == 0 {
		return 0
	}

	// Select a load based on the calculated probabilities
	randomValue := rand.Float64() * sum
	for i, probability := range probabilities {
		randomValue -= probability
		if randomValue <= 0 {
			return remainingLoads[i]
		}
	}

	return 0
}

// updateTabuList manages the tabu list, adding new entries and removing old ones
func updateTabuList(tabuList map[string]float64, tabuCounter map[string]int, solution Solution) {
	key := neighborKey(solution)
	if len(tabuList) >= tabuListSize {
		for k := range tabuList {
			if tabuCounter[k] > 0 {
				tabuCounter[k]--
			}
			if tabuCounter[k] == 0 {
				delete(tabuList, k)
				delete(tabuCounter, k)
			}
		}
	}
	tabuList[key] = initialTabuValue
	tabuCounter[key] = tabuListSize
}

// neighborKey generates a unique key for a solution
func neighborKey(solution Solution) string {
	var sb strings.Builder
	for _, route := range solution.routes {
		sb.WriteString(fmt.Sprintf("%v-", route))
	}
	return sb.String()
}

// calculateCost computes the total cost of a solution
func calculateCost(solution Solution) float64 {
	totalDistance := 0.0
	for _, route := range solution.routes {
		previousNode := 0
		for _, node := range route {
			totalDistance += distanceMatrix[previousNode][node] + deliveryDistance[node-1]
			previousNode = node
		}
		totalDistance += distanceMatrix[previousNode][0]
	}
	return totalDistance + float64(len(solution.routes))*costPerDriver
}

// printSolution outputs the solution in the required format
func printSolution(solution Solution) {
	for _, route := range solution.routes {
		fmt.Printf("[%s]\n", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(route)), ","), "[]"))
	}
}
