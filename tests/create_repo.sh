#!/bin/bash

#This file is exclusively for testing purposes. It creates a dummy repository with the specified branches and commits.

# Create a new directory for the repository
mkdir dummy
cd dummy || exit

# Initialize a new Git repository
git init

# Create 3 commits in the master branch
for i in {1..3}; do
  touch master${i}.txt
  git add master${i}.txt
  git commit -m "master commit ${i}"
done

# Create branch_A and add 3 commits
git checkout -b branch_A
for i in {1..3}; do
  touch branch_A${i}.txt
  git add branch_A${i}.txt
  git commit -m "branch_A commit ${i}"
done

# Switch back to master and create branch_B, then add 3 commits
git checkout master
git checkout -b branch_B
for i in {1..3}; do
  touch branch_B${i}.txt
  git add branch_B${i}.txt
  git commit -m "branch_B commit ${i}"
done

echo "Repository created with the specified branches and commits."
