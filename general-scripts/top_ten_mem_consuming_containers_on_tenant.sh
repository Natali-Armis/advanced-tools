#!/bin/bash

function human_readable() {
    local bytes=$1
    local kib=$((bytes/1024))
    local mib=$((kib/1024))
    local gib=$((mib/1024))
    if (( gib > 0 )); then
        echo "${gib}GiB"
    elif (( mib > 0 )); then
        echo "${mib}MiB"
    elif (( kib > 0 )); then
        echo "${kib}KiB"
    else
        echo "${bytes}B"
    fi
}

total_machine_memory=$(grep MemTotal /proc/meminfo | awk '{print $2 * 1024}')
container_ids=$(docker ps -q)
container_names=$(docker ps --format "{{.ID}} {{.Names}}")

declare -A container_memory_usage
declare -A container_memory_limit
declare -A container_names_map
total_memory_usage=0

while read -r id name; do
    container_names_map["$id"]=$name
done <<< "$container_names"

for container_id in $container_ids; do
    full_container_id=$(docker inspect --format '{{.Id}}' "$container_id")
    cgroup_path=$(find /sys/fs/cgroup -name "docker-$full_container_id.scope" 2>/dev/null)
    if [ -z "$cgroup_path" ]; then
        echo "Cgroup path not found for container ID $container_id"
        continue
    fi
    memory_usage=$(cat "$cgroup_path/memory.current")
    memory_limit=$(docker inspect --format '{{.HostConfig.Memory}}' "$container_id")
    if [ "$memory_limit" -eq 0 ]; then
        memory_limit="unlimited"
    else
        memory_limit=$(human_readable $memory_limit)
    fi

    container_memory_usage["$container_id"]=$memory_usage
    container_memory_limit["$container_id"]=$memory_limit
    total_memory_usage=$((total_memory_usage + memory_usage))
done

top_containers=$(for id in "${!container_memory_usage[@]}"; do
    echo "$id ${container_memory_usage[$id]}"
done | sort -k2 -n -r | head -10)

echo "Top 10 containers by memory usage:"
echo "Container ID        Container Name        Memory Usage        Memory Limit"
echo "-------------------------------------------------------------------------"
echo "$top_containers" | while read -r line; do
    container_id=$(echo $line | awk '{print $1}')
    memory_usage=$(echo $line | awk '{print $2}')
    human_memory=$(human_readable $memory_usage)
    container_name=${container_names_map[$container_id]}
    memory_limit=${container_memory_limit[$container_id]}
    echo "$container_id      $container_name      $human_memory      $memory_limit"
done

total_human_memory=$(human_readable $total_memory_usage)
total_memory_percentage=$(awk "BEGIN {printf \"%.2f\", ($total_memory_usage/$total_machine_memory)*100}")
total_machine_memory_human=$(human_readable $total_machine_memory)

echo "-------------------------------------------------------------------------"
echo "Total memory usage of all containers: $total_human_memory"
echo "Total memory of the machine: $total_machine_memory_human"
echo "Percentage of memory used by containers: $total_memory_percentage%"