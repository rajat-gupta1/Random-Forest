#!/bin/bash
#
#SBATCH --mail-user=rajatgupta@uchicago.edu
#SBATCH --mail-type=ALL
#SBATCH --job-name=part6
#SBATCH --output=/home/rajatgupta/course/parallel/project-3-rajat-gupta1/proj3/out/%j.%N.stdout
#SBATCH --error=/home/rajatgupta/course/parallel/project-3-rajat-gupta1/proj3/out/%j.%N.stderr
#SBATCH --chdir=/home/rajatgupta/course/parallel/project-3-rajat-gupta1/proj3
#SBATCH --partition=debug
#SBATCH --nodes=1
#SBATCH --ntasks=1
#SBATCH --cpus-per-task=16
#SBATCH --mem-per-cpu=900
#SBATCH --exclusive
#SBATCH --time=03:00:00


module load golang/1.16.2

# Sequential
for i in {1..5}
do
    go run randomforest/tree.go s 200 4 >> time.txt
done

# Stealing
for j in 1 2 3 4 6
do
    for i in {1..5}
    do
        go run randomforest/tree.go stl 200 4 $(($j * 2)) 10 >> time.txt
    done
done

# Balancing
for j in 1 2 3 4 6
do
    for i in {1..5}
    do
        go run randomforest/tree.go bal 200 4 $(($j * 2)) 10 2 >> time.txt
    done
done