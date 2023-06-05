places, rides, group_nm = map(int, input().split())
# Gets groups as list of integers
groups = list(map(int, [input() for _ in range(group_nm)]))

# Initialize dicts
profits = {}
groups_after = {}

for i in range(group_nm):
    # Starting values
    current_index = i
    profits[i] = 0

    while True:
        # Group that's about to ride
        next_grp = groups[current_index]

        # Go out of loop when no more places are available
        if profits[i] + next_grp > places:
            break

        # Increase profits by the number of people in group
        profits[i] += next_grp

        # Increment the index
        current_index += 1

        # Reset the index if we reached the end of the list
        current_index = 0 if current_index == group_nm else current_index

        # We passed through the whole list yet there are more places are available
        if current_index == i:
            break

    # Once done, we want to save the index of the group that needs to ride next
    groups_after[i] = current_index

# Initialize total and reset current index
total = 0
current_index = 0

# Sum up the profits
for i in range(rides):
    total += profits[current_index]
    current_index = groups_after[current_index]

# Solution, yay!
print(total)

# A thank you to CodinGame user Tobou who hinted me towards the "caching" idea when optimizing the solution