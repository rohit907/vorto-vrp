# Tabu Search for Vehicle Routing Problem (VRP)

This Go application implements the Tabu Search algorithm to solve the Vehicle Routing Problem (VRP). The VRP involves finding the optimal routes for a fleet of vehicles to deliver goods to a set of locations, minimizing the total cost, including travel distance and the number of drivers.

## Features

- Reads load data from a file.
- Precomputes distance matrices for efficiency.
- Implements Tabu Search to find near-optimal routes.
- Generates and evaluates neighboring solutions.
- Manages a Tabu list to avoid cycling back to previously visited solutions.
- Outputs the best routes found along with their costs.

## Prerequisites

- Go 1.18 or later

## Installation

1. **Clone the repository:**

   ```bash
   git clone https://github.com/rohitb907/vorto-vrp.git
   cd vorto-vrp
   ```

2. **Build the application:**
    ```bash
    go build -o vorto-vrp
    ```

    
3. **Run the application:**
    ```bash
    go run main.go problem20.txt
    ```

    
**Data File Format**

The data file should contain load information in the following format:
```bash
loadNumber pickup dropoff
1 (15.564823759233205,50.91649123557091) (19.516838671539674,72.5683671095093)
2 (15,25) (35,45)
...
```
Where each line (after the header) represents a load with:

loadNumber An integer LoadID.

Pickup coordinates (x,y) in the format (x,y).

Dropoff coordinates (x,y) in the format (x,y).

**Example Output**

The output will list the routes and their costs in the following format:
```bash
   [1,2,3]
   [4,5]
```

**Run the complete test evaluation**
 ```bash
    python3 evaluateShared.py --cmd "go run main.go" --problemDir Training
 ```
    

**Example Output**


![ouput](https://github.com/user-attachments/assets/ad9f9f8e-feca-4ddf-b4b2-544552fa37a8)




**Reference**


https://www.researchgate.net/publication/221704721_A_Tabu_Search_Heuristic_for_the_Vehicle_Routing_Problem
https://www.geeksforgeeks.org/what-is-tabu-search/
