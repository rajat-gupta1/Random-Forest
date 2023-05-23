import matplotlib.pyplot as plt
import numpy as np

def divide(arr):
    val = arr[0]
    for i in range(len(arr)):
        arr[i] /= val
        arr[i] = 1 / arr[i]
        arr[i] = round(arr[i], 2)

def generateGraphs():
    f = open("time.txt", "r")
    steal = []
    balance = []

    threads = {1, 2, 4, 6, 8, 12}
   
    for j in threads:
        time = 0
        for i in range(5):
            time += float(f.readline().strip('\n'))
        steal.append(time)

    balance.append(steal[0])

    threads = {2, 4, 6, 8, 12}

    for j in threads:
        time = 0
        for i in range(5):
            time += float(f.readline().strip('\n'))
        balance.append(time)

   
    f.close()

    divide(steal)
    divide(balance)

    xpoints = [1, 2, 4, 6, 8, 12]
    fig1 = plt.figure("Figure 1")
    plt.title("SpeedUp Graph Steal/Balance")

    ypoints = np.array(steal)
    plt.plot(xpoints, ypoints, marker = 'o', label = "STEAL")

    ypoints = np.array(balance)
    plt.plot(xpoints, ypoints, marker = 'o', label = "BALANCE")

    plt.xlabel("Number of threads")
    plt.ylabel("Speed Up")

    plt.legend()
    plt.savefig("speedup.png")


if __name__ == "__main__":
    generateGraphs()